// Copyright 2017 Jörg Kost All rights reserved.
// joerg.kost@gmx.com, jk@ip-clear.de
// Use of this source code is governed by Apache 2.0
// license that can be found in the LICENSE.MD file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/ipcjk/mlxsh/libhost"
	"github.com/ipcjk/mlxsh/netironDevice"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var cliWriteTimeout, cliReadTimeout time.Duration
var cliHostname, cliPassword, cliUsername, cliEnablePassword string
var cliSpeedMode bool
var debug, version, quiet bool
var cliMaxParallel int
var cliScriptFile, cliConfigFile, cliRouterFile, cliLabel string
var selectedHosts []libhost.HostConfig

type chanHost struct {
	hostName string
	message  string
	err      error
}

func init() {
	flag.StringVar(&cliScriptFile, "script", "", "script file to to execute, if no file is found, its used as a direct command")
	flag.StringVar(&cliConfigFile, "config", "", "Configuration file to insert, its used as a direct command")
	flag.StringVar(&cliLabel, "label", "", "label-selection for run commands on a group of routers, e.g. 'location=munich,environment=prod'")
	flag.StringVar(&cliHostname, "hostname", "", "Router hostname")
	flag.StringVar(&cliPassword, "password", "", "user password")
	flag.StringVar(&cliUsername, "username", "", "username")
	flag.StringVar(&cliEnablePassword, "enable", "", "enable password")
	flag.IntVar(&cliMaxParallel, "c", 2, "concurrent working threads / connections to the routers")
	flag.DurationVar(&cliReadTimeout, "readtimeout", time.Second*15, "timeout for reading poll on cli select")
	flag.DurationVar(&cliWriteTimeout, "writetimeout", time.Millisecond*0, "timeout to stall after a write to cli")
	flag.BoolVar(&debug, "debug", false, "Enable debug for read / write")
	flag.BoolVar(&cliSpeedMode, "speedmode", false, "Enable speed mode write, will ignore any output from the cli while writing")
	flag.BoolVar(&quiet, "q", false, "quiet mode, no output except error on connecting & co")
	flag.BoolVar(&version, "version", false, "prints version and exit")

	if os.Getenv("JK") == "1" {
		log.Println("Developer configuration active")
		flag.StringVar(&cliRouterFile, "routerdb", "config_jk.yaml", "Input file in yaml for username,password and host configuration if not specified on command-line")
	} else {
		flag.StringVar(&cliRouterFile, "routerdb", "mlxsh.yaml", "Input file in yaml for username,password and host configuration if not specified on command-line")
	}

	flag.Parse()

	if version {
		log.Println("mlxsh 0.2b (C) 2017 by Jörg Kost, jk@ip-clear.de")
		os.Exit(0)
	}

	if cliHostname == "" && cliLabel == "" {
		log.Println("No host/router or selector given, abort...")
		os.Exit(0)
	} else if cliHostname != "" && cliLabel != "" {
		log.Println("Cant run in targetHost-mode and Groupselector")
		os.Exit(0)
	}

	if cliRouterFile != "" {
		file, err := os.Open(cliRouterFile)
		if err != nil {
			log.Fatal(err)
		}

		selectedHosts, err = libhost.LoadMatchesFromYAML(file, cliLabel, cliHostname)
		if err != nil {
			log.Fatal(err)
		}

	}
	/*  Hostname on cli but did not found in list */
	if cliHostname != "" && len(selectedHosts) == 0 {
		selectedHosts = append(selectedHosts, libhost.HostConfig{Hostname: cliHostname, Username: cliUsername, Password: cliPassword, EnablePassword: cliEnablePassword, SpeedMode: cliSpeedMode, SSHPort: 22})
	}

	if len(selectedHosts) == 0 {
		log.Fatal("Could not find any target host for this labels")
	}

	/* Possible overwrite settings from CliParameters */
	for x := range selectedHosts {
		selectedHosts[x].ApplyCliSettings(cliScriptFile, cliConfigFile, cliWriteTimeout, cliReadTimeout)
	}
}

func main() {
	hostChannel := make(chan chanHost, 1)
	var wg sync.WaitGroup
	var semaphore = make(chan struct{}, cliMaxParallel)

	// worker
	for x := range selectedHosts {
		wg.Add(1)
		go func(x int) {
			semaphore <- struct{}{}

			var err error
			var buffer = new(bytes.Buffer)

			router := netironDevice.NetironDevice(
				netironDevice.NetironConfig{HostConfig: selectedHosts[x], Debug: debug, W: buffer})

			defer func() {
				if router != nil {
					router.CloseConnection()
				}
				hostChannel <- chanHost{message: buffer.String(), hostName: selectedHosts[x].Hostname, err: err}
				wg.Done()
				<-semaphore
			}()

			if router == nil {
				err = fmt.Errorf("Cant instance object")
				return
			}

			if err = router.ConnectPrivilegedMode(); err != nil {
				return
			}

			if _, err = router.SkipPageDisplayMode(); err != nil {
				return
			}

			if err = router.GetPromptMode(); err != nil {
				return
			}

			if selectedHosts[x].Filename != "" {
				var input io.Reader
				file, err := os.Open(selectedHosts[x].Filename)
				defer file.Close()

				if err != nil && os.IsNotExist(err) {
					command := strings.Replace(selectedHosts[x].Filename, ";", "\n", -1)
					input = strings.NewReader(command)
					if debug {
						log.Printf("Cant open file: %s, will read from command line argument\n", err)
					}
				} else if err != nil {
					log.Printf("Cant open file: %s\n", err)
				} else {
					input = file
				}

				/* Execution Mode starts here */
				if selectedHosts[x].ExecMode {
					if err := router.RunCommands(input); err != nil {
						return
					}
				} else {

					/* Configuration Mode starts here */
					if err = router.ConfigureTerminalMode(); err != nil {
						return
					}
					if err := router.PasteConfiguration(input); err != nil {
						return
					}
					if err := router.WriteConfiguration(); err != nil {
						return
					}

				}
			}
		}(x)
	}

	// closer
	go func() {
		wg.Wait()
		close(hostChannel)
	}()

	// printer
	for elems := range hostChannel {
		fmt.Printf("╔══════════════════════════════════════════════════════════════════════╗\n")
		if elems.err != nil {
			fmt.Printf("║%-70s║\n", elems.hostName)
			fmt.Printf("║%-70s║\n", "No success:")
			fmt.Printf("║%-70s║\n", elems.err)
			if elems.message != "" {
				fmt.Println(elems.message)
			}
			fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
		} else {
			if quiet == false {
				fmt.Printf("║%-70s║\n", elems.hostName)
				fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
				fmt.Println(elems.message)
			}
		}

	}

}
