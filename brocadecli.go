// Copyright 2017 Jörg Kost All rights reserved.
// joerg.kost@gmx.com, jk@ip-clear.de
// Use of this source code is governed by a GPLv2-style
// license that can be found in the LICENSE file.

package main

import (
	device "./device"
	libhost "./libhost"
	"flag"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var targetHost libhost.HostEntry
var debug bool
var logDir, outputFile, configFile string

func init() {
	flag.StringVar(&targetHost.Filename, "filename", "", "Configuration file to insert")
	flag.StringVar(&targetHost.Hostname, "hostname", "rt1", "Router hostname")
	flag.StringVar(&targetHost.Password, "password", "password", "user password")
	flag.StringVar(&targetHost.Username, "username", "username", "username")
	flag.StringVar(&targetHost.EnablePassword, "enable", "enablepassword", "enable password")
	flag.DurationVar(&targetHost.ReadTimeout, "readtimeout", time.Second*15, "timeout for reading poll on cli select")
	flag.DurationVar(&targetHost.WriteTimeout, "writetimeout", time.Millisecond*0, "timeout to stall after a write to cli")
	flag.BoolVar(&debug, "debug", false, "Enable debug for read / write")
	flag.BoolVar(&targetHost.SpeedMode, "speedmode", false, "Enable speed mode write, will ignore any output from the cli while writing")
	flag.BoolVar(&targetHost.ExecMode, "execmode", false, "Exec commands / input from filename instead of paste configuration")
	flag.StringVar(&logDir, "logdir", "", "Record session into logDir, automatically gzip")
	flag.StringVar(&outputFile, "outputfile", "", "Output file, else stdout")

	if os.Getenv("JK") != "" {
		log.Println("Developer configuration active")
		flag.StringVar(&configFile, "configfile", "config_jk.yaml", "Input file in yaml for username,password and host configuration if not specified on command-line")
	} else {
		flag.StringVar(&configFile, "configfile", "broconfig.yaml", "Input file in yaml for username,password and host configuration if not specified on command-line")
	}

	flag.Parse()
	if configFile != "" {
		loadConfig()
	}

}

func main() {
	var err error
	router := device.Brocade(device.DEVICE_MLX, targetHost.Hostname, 22, targetHost.EnablePassword, targetHost.Username, targetHost.Password,
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
	/* router.ExecPrivilegedMode("clear ip bgp neighbor ... soft") */
	router.CloseConnection()
}

func loadConfig() {
	var hostsConfig []libhost.HostEntry

	source, err := ioutil.ReadFile(configFile)
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

		}
	}
}
