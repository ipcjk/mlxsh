// Copyright 2017 Jörg Kost All rights reserved.
// joerg.kost@gmx.com, jk@ip-clear.de
// Use of this source code is governed by a GPLv2-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/ipcjk/brocadecli/device"
	"github.com/ipcjk/brocadecli/libhost"
	"flag"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var targetHost libhost.HostEntry
var debug, version bool
var scriptFile, configFile string
var logDir, outputFile, routerFile string

func init() {
	flag.StringVar(&scriptFile, "script", "", "script file to to execute")
	flag.StringVar(&configFile, "config", "", "Configuration file to insert")
	flag.StringVar(&targetHost.Hostname, "hostname", "", "Router hostname")
	flag.StringVar(&targetHost.Password, "password", "", "user password")
	flag.StringVar(&targetHost.Username, "username", "", "username")
	flag.StringVar(&targetHost.EnablePassword, "enable", "", "enable password")
	flag.DurationVar(&targetHost.ReadTimeout, "readtimeout", time.Second*15, "timeout for reading poll on cli select")
	flag.DurationVar(&targetHost.WriteTimeout, "writetimeout", time.Millisecond*0, "timeout to stall after a write to cli")
	flag.BoolVar(&debug, "debug", false, "Enable debug for read / write")
	flag.BoolVar(&targetHost.SpeedMode, "speedmode", false, "Enable speed mode write, will ignore any output from the cli while writing")
	flag.StringVar(&logDir, "logdir", "", "Record session into logDir, automatically gzip")
	flag.StringVar(&outputFile, "outputfile", "", "Output file, else stdout")
	flag.BoolVar(&version, "version", false, "prints version and exit")

	if version {
		log.Println("brocadecli 0.1 (C) 2017 by Jörg Kost, jk@ip-clear.de")
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
	var err error
	router := device.Brocade(device.DEVICE_MLX, targetHost.Hostname, targetHost.SSHPort, targetHost.EnablePassword, targetHost.Username, targetHost.Password,
		targetHost.ReadTimeout, targetHost.WriteTimeout, debug, targetHost.SpeedMode)

	if err = router.ConnectPrivilegedMode(); err != nil {
		log.Fatal(err)
	}

	if _, err = router.SkipPageDisplayMode(); err != nil {
		log.Fatal(err)
	}

	if err = router.GetPromptMode(); err != nil {
		log.Fatal(err)
	}

	if targetHost.Filename != "" {
		file, err := os.Open(targetHost.Filename)
		defer file.Close()
		if err != nil {
			log.Printf("Cant open file: %s", err)
		} else {
			if targetHost.ExecMode == true {
				router.RunCommandsFromReader(file)
			} else {
				router.ConfigureTerminalMode()
				router.PasteConfiguration(file)
				router.WriteConfiguration()
			}
		}
	}
	/* router.ExecPrivilegedMode("show ip route ... longer") */
	router.CloseConnection()
}

func loadMergeConfig() {
	var hostsConfig []libhost.HostEntry

	if targetHost.Hostname == "" {
		log.Fatal("No host/router given, abort...")
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
		if Host.Hostname == targetHost.Hostname {
			if debug {
				log.Println("Overwrite cli settings for " + targetHost.Hostname + " from " + configFile)
			}

			/* Better copy structure?
			but then we will copy the booleans?
			*/

			if Host.Password != "" {
				targetHost.Password = Host.Password
			}

			if Host.EnablePassword != "" {
				targetHost.EnablePassword = Host.EnablePassword
			}

			if Host.Username != "" {
				targetHost.Username = Host.Username
			}

			if Host.Filename != "" {
				targetHost.Filename = Host.Filename
			}

			if Host.ExecMode == true {
				targetHost.ExecMode = true
			}

			if Host.SpeedMode == true {
				targetHost.SpeedMode = true
			}

			if Host.ConfigFile != "" {
				targetHost.Filename = Host.ConfigFile
				targetHost.ExecMode = false
			}

			if Host.ScriptFile != "" {
				targetHost.Filename = Host.ScriptFile
				targetHost.ExecMode = true
			}

			if scriptFile != "" {
				targetHost.Filename = scriptFile
				targetHost.ExecMode = true
			}

			if configFile != "" {
				targetHost.Filename = configFile
				targetHost.ExecMode = false
			}

			if Host.SSHPort == 0 {
				targetHost.SSHPort = 22
			} else {
				targetHost.SSHPort = Host.SSHPort
			}

			break
		}
	}
}
