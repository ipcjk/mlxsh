// Copyright 2017 Jörg Kost All rights reserved.
// joerg.kost@gmx.com, jk@ip-clear.de
// Use of this source code is governed by a GPLv2-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"github.com/ipcjk/brocadecli/device"
	"github.com/ipcjk/brocadecli/libhost"
	"gopkg.in/yaml.v1"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var cli libhost.HostEntry
var selectedHosts []libhost.HostEntry
var debug, version bool
var scriptFile, configFile, routerFile, label string

func init() {
	flag.StringVar(&scriptFile, "script", "", "script file to to execute, if no file is found, its used as a direct command")
	flag.StringVar(&configFile, "config", "", "Configuration file to insert, its used as a direct command")
	flag.StringVar(&label, "label", "", "selector for run commands on a group of routers")
	flag.StringVar(&cli.Hostname, "hostname", "", "Router hostname")
	flag.StringVar(&cli.Password, "password", "", "user password")
	flag.StringVar(&cli.Username, "username", "", "username")
	flag.StringVar(&cli.EnablePassword, "enable", "", "enable password")
	flag.DurationVar(&cli.ReadTimeout, "readtimeout", time.Second*15, "timeout for reading poll on cli select")
	flag.DurationVar(&cli.WriteTimeout, "writetimeout", time.Millisecond*0, "timeout to stall after a write to cli")
	flag.BoolVar(&debug, "debug", false, "Enable debug for read / write")
	flag.BoolVar(&cli.SpeedMode, "speedmode", false, "Enable speed mode write, will ignore any output from the cli while writing")
	flag.BoolVar(&version, "version", false, "prints version and exit")

	if version {
		log.Println("brocadecli 0.x (C) 2017 by Jörg Kost, jk@ip-clear.de")
		os.Exit(0)
	}

	if os.Getenv("JK") != "" {
		log.Println("Developer configuration active")
		flag.StringVar(&routerFile, "routerdb", "config_jk.yaml", "Input file in yaml for username,password and host configuration if not specified on command-line")
	} else {
		flag.StringVar(&routerFile, "routerdb", "broconfig.yaml", "Input file in yaml for username,password and host configuration if not specified on command-line")
	}

	flag.Parse()

	if routerFile != "" {
		loadMergeConfig()
	}

}

func main() {
	for _, selectHost := range selectedHosts {
		var err error

		router := device.Brocade(selectHost.DeviceType, selectHost.Hostname, selectHost.SSHPort, selectHost.EnablePassword, selectHost.Username, selectHost.Password,
			selectHost.ReadTimeout, selectHost.WriteTimeout, debug, selectHost.SpeedMode)

		if err = router.ConnectPrivilegedMode(); err != nil {
			log.Fatal(err)
		}

		if _, err = router.SkipPageDisplayMode(); err != nil {
			log.Fatal(err)
		}

		if err = router.GetPromptMode(); err != nil {
			log.Fatal(err)
		}

		if selectHost.Filename != "" {
			var input io.Reader
			file, err := os.Open(selectHost.Filename)
			defer file.Close()

			if err != nil && os.IsNotExist(err) {
				input = strings.NewReader(selectHost.Filename)
				if debug {
					log.Printf("Cant open file: %s, will read from command line argument\n", err)
				}
			} else if err != nil {
				log.Printf("Cant open file: %s\n", err)
			} else {
				input = file
			}

			if selectHost.ExecMode == true {
				router.RunCommandsFromReader(input)
			} else {
				router.ConfigureTerminalMode()
				router.PasteConfiguration(input)
				router.WriteConfiguration()
			}

		}

		router.CloseConnection()
	}

}

func loadMergeConfig() {
	var hostsConfig []libhost.HostEntry

	if cli.Hostname == "" && label == "" {
		log.Fatal("No host/router or selector given, abort...")
	} else if cli.Hostname != "" && label != "" {
		log.Fatal("Cant run in targetHost-mode and Groupselector")
	}

	source, err := ioutil.ReadFile(routerFile)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(source, &hostsConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for _, Host := range hostsConfig {
		/* Code for targetHost-mode */
		if Host.Hostname == cli.Hostname {
			if debug {
				log.Println("Overwrite cli settings for " + cli.Hostname + " from " + configFile)
			}

			/* Better copy structure?
			but then we will copy the booleans?
			*/

			if Host.Password != "" {
				cli.Password = Host.Password
			}

			if Host.EnablePassword != "" {
				cli.EnablePassword = Host.EnablePassword
			}

			if Host.Username != "" {
				cli.Username = Host.Username
			}

			if Host.Filename != "" {
				cli.Filename = Host.Filename
			}

			if Host.ExecMode == true {
				cli.ExecMode = true
			}

			if Host.SpeedMode == true {
				cli.SpeedMode = true
			}

			if Host.ConfigFile != "" {
				cli.Filename = Host.ConfigFile
				cli.ExecMode = false
			}

			if Host.ScriptFile != "" {
				cli.Filename = Host.ScriptFile
				cli.ExecMode = true
			}

			if scriptFile != "" {
				cli.Filename = scriptFile
				cli.ExecMode = true
			}

			if configFile != "" {
				cli.Filename = configFile
				cli.ExecMode = false
			}

			if Host.SSHPort == 0 {
				cli.SSHPort = 22
			} else {
				cli.SSHPort = Host.SSHPort
			}

			selectedHosts = append(selectedHosts, cli)
			break

		} else if Host.MatchLabels(label) {
			var newHost = Host

			if configFile != "" {
				newHost.Filename = configFile
				newHost.ExecMode = false
			}

			if scriptFile != "" {
				newHost.Filename = scriptFile
				newHost.ExecMode = true
			}

			if newHost.SSHPort == 0 {
				newHost.SSHPort = 22
			} else {
				newHost.SSHPort = Host.SSHPort
			}

			if newHost.ReadTimeout == 0 {
				newHost.ReadTimeout = cli.ReadTimeout
			}

			if newHost.WriteTimeout == 0 {
				newHost.WriteTimeout = cli.WriteTimeout
			}

			selectedHosts = append(selectedHosts, newHost)
		}
	}
}
