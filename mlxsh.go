// Copyright 2017 Jörg Kost All rights reserved.
// joerg.kost@gmx.com, jk@ip-clear.de
// Use of this source code is governed by Apache 2.0
// license that can be found in the LICENSE.MD file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/ipcjk/mlxsh/slxDevice"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/ipcjk/mlxsh/junosDevice"
	"github.com/ipcjk/mlxsh/libhost"
	"github.com/ipcjk/mlxsh/netironDevice"
	"github.com/ipcjk/mlxsh/routerDevice"
	"github.com/ipcjk/mlxsh/vdxDevice"
	"github.com/mattn/go-isatty"
)

var cliWriteTimeout, cliReadTimeout time.Duration
var cliHostname, cliPassword, cliUsername, cliEnablePassword string
var debug, version, quiet, cliHostCheck, cliSpeedMode bool
var outputIsTerminal, cliNoColor, shellMode bool
var cliMaxParallel int
var cliScriptFile, cliConfigFile, cliRouterFile, cliLabel, cliType, cliKeyFile, cliHostFile string
var selectedHosts, allHosts []libhost.HostConfig

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
	flag.StringVar(&cliType, "clitype", "mlxe", "Router type")
	flag.StringVar(&cliKeyFile, "i", "", "Path to a ssh private key (in openssh2-format) that will be used for connections ")
	flag.StringVar(&cliHostFile, "sf", "", "Path to the known-hosts-file (in openssh2-format) that will be used for validating hostkeys, defaults to .ssh/known_hosts ")
	flag.IntVar(&cliMaxParallel, "c", 20, "concurrent working threads")
	flag.DurationVar(&cliReadTimeout, "readtimeout", time.Second*30, "timeout for reading poll on cli select")
	flag.DurationVar(&cliWriteTimeout, "writetimeout", time.Millisecond*0, "timeout to stall after a write to cli")
	flag.BoolVar(&shellMode, "shell", false, "Run in libreadline command line prompt mode")
	flag.BoolVar(&debug, "debug", false, "Enable debug for read / write")
	flag.BoolVar(&cliHostCheck, "s", false, "Enable strict hostkey checking for ssh connections")
	flag.BoolVar(&cliSpeedMode, "speedmode", false, "Enable speed mode write, will ignore any output from the cli while writing")
	flag.BoolVar(&quiet, "q", false, "quiet mode, no output except error on connecting & co")
	flag.BoolVar(&version, "version", false, "prints version and exit")
	flag.BoolVar(&cliNoColor, "nocolor", false, "Disable color printing when output line is a terminal")

	if os.Getenv("JK") == "1" {
		log.Println("Developer configuration active")
		flag.StringVar(&cliRouterFile, "routerdb", "config_jk.yaml", "Input file in yaml for username,password and host configuration if not specified on command-line")
	} else {
		flag.StringVar(&cliRouterFile, "routerdb", "mlxsh.yaml", "Input file in yaml for username,password and host configuration if not specified on command-line")
	}

	flag.Parse()

	if version {
		log.Println("mlxsh 0.5 (C) 2018 by Jörg Kost, jk@ip-clear.de")
		os.Exit(0)
	}

	if isatty.IsTerminal(os.Stdout.Fd()) {
		outputIsTerminal = true
	}

	if !outputIsTerminal && shellMode {
		log.Println("Cant run in shellmode without terminal")
		os.Exit(0)
	}

	if cliHostname == "" && cliLabel == "" && !shellMode {
		log.Println("No host/router or selector given, abort...")
		os.Exit(0)
	} else if cliHostname != "" && cliLabel != "" && shellMode == false {
		log.Println("Cant run in targetHost-mode or groupselection")
		os.Exit(0)
	}

	if cliHostFile == "" {
		cliHostFile = getUserKnownHostsFile()
	}

	if cliRouterFile != "" {
		file, err := os.Open(cliRouterFile)
		if err != nil {
			log.Fatal(err)
		}

		selectedHosts, allHosts, err = libhost.LoadMatchesFromYAML(file, cliLabel, cliHostname)
		if err != nil {
			log.Fatal(err)
		}
	}

	/* Setup done for shellmode, rest setup is only done for one-shot */
	if shellMode {
		return
	}

	/*  Hostname on cli but did not found in list */
	if cliHostname != "" && len(selectedHosts) == 0 {
		selectedHosts = append(selectedHosts, libhost.HostConfig{Hostname: cliHostname, Username: cliUsername, Password: cliPassword, EnablePassword: cliEnablePassword, DeviceType: cliType, SpeedMode: cliSpeedMode, SSHPort: 22})
	}

	if len(selectedHosts) == 0 {
		log.Fatal("Could not find any target host for this labels")
	}
}

func main() {
	if shellMode {
		runShellMode()
	} else {
		run()
	}
}

func applyCliSettings() {
	/* Possible overwrite settings from CliParameters */
	for x := range selectedHosts {
		selectedHosts[x].ApplyCliSettings(cliScriptFile, cliConfigFile, cliWriteTimeout, cliReadTimeout, cliHostCheck, cliKeyFile, cliHostFile)
	}
}

/* Config or Exec-Statements running from command line parameter or file input */
func run() {
	hostChannel := make(chan chanHost, 1)
	var wg sync.WaitGroup
	var semaphore = make(chan struct{}, cliMaxParallel)

	applyCliSettings()

	// worker
	for x := range selectedHosts {
		wg.Add(1)
		go func(x int) {
			semaphore <- struct{}{}

			var err error
			var buffer = new(bytes.Buffer)
			var singleRouter RouterInt

			switch strings.ToLower(selectedHosts[x].DeviceType) {
			case "vdx":
				singleRouter = RouterInt(vdxDevice.VdxDevice(router.RunTimeConfig{HostConfig: selectedHosts[x], Debug: debug, W: buffer}))
			case "slx":
				singleRouter = RouterInt(slxDevice.SlxDevice(router.RunTimeConfig{HostConfig: selectedHosts[x], Debug: debug, W: buffer}))
			case "mlx", "cer", "mlxe", "xmr", "iron", "turobiron", "icx", "fcs":
				singleRouter = RouterInt(netironDevice.NetironDevice(
					router.RunTimeConfig{HostConfig: selectedHosts[x], Debug: debug, W: buffer}))
			case "juniper", "junos", "mx", "ex", "j":
				singleRouter = RouterInt(junosDevice.JunosDevice(router.RunTimeConfig{HostConfig: selectedHosts[x], Debug: debug, W: buffer}))
			default:
				/* Default always to Netiron for compatible  */
				singleRouter = RouterInt(netironDevice.NetironDevice(
					router.RunTimeConfig{HostConfig: selectedHosts[x], Debug: debug, W: buffer}))
			}

			defer func() {
				if singleRouter != nil {
					singleRouter.Close()
				}
				hostChannel <- chanHost{message: buffer.String(), hostName: selectedHosts[x].Hostname, err: err}
				wg.Done()
				<-semaphore
			}()

			if singleRouter == nil {
				err = fmt.Errorf("can't instance router object")
				return
			}

			if err = singleRouter.Connect(); err != nil {
				return
			}

			if selectedHosts[x].Filename != "" {
				var input io.Reader
				var file *os.File

				file, err = os.Open(selectedHosts[x].Filename)
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
					if err = singleRouter.RunCommands(input); err != nil {
						return
					}
				} else {

					/* Configuration Mode starts here */
					if err = singleRouter.ConfigureTerminalMode(); err != nil {
						return
					}
					if err = singleRouter.PasteConfiguration(input); err != nil {
						return
					}
					if err = singleRouter.CommitConfiguration(); err != nil {
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
		state := "OK"

		if elems.err != nil {
			state = "err"
		}

		if !outputIsTerminal || cliNoColor {
			fmt.Printf("%s: [%-20s]", state, elems.hostName)
		} else if state == "err" {
			fmt.Printf("\x1b[31m%s: [%-20s]\x1b[0m", state, elems.hostName)
		} else {
			fmt.Printf("\x1b[32m%s: [%-20s]\x1b[0m", state, elems.hostName)
		}

		if state == "err" {
			fmt.Printf(" errors: %s, messages: %s", elems.message, elems.err)
		} else if !quiet {
			fmt.Printf(" %s", elems.message)
		}

		fmt.Printf("\n")
	}
}

func getUserKnownHostsFile() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH") + `\.ssh\known_hosts`
		if home == "" {
			home = os.Getenv("USERPROFILE") + `\.ssh\known_hosts`
		}
		return home
	}
	return os.Getenv("HOME") + "/.ssh/known_hosts"
}

func getUserHistoryFile() string {
	var historyFile = "/.mlxsh_history"
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH") + historyFile
		if home == "" {
			home = os.Getenv("USERPROFILE") + historyFile
		}
		return home
	}
	return os.Getenv("HOME") + historyFile
}
