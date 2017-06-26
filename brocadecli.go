// Copyright 2017 Jörg Kost All rights reserved.
// joerg.kost@gmx.com, jk@ip-clear.de
// Use of this source code is governed by a GPLv2-style
// license that can be found in the LICENSE file.

package main

import (
	device "./device"
	"flag"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type HostEntry struct {
	Hostname        string `yaml:"Hostname"`
	Username        string `yaml:"Username"`
	Password        string `yaml:"Password"`
	EnablePassword  string `yaml:"EnablePassword"`
	DeviceType      string `yaml:"DeviceType"`
	KeyFile         string `yaml:"KeyFile"`
	StrictHostCheck bool   `yaml:"StrictHostCheck"`
	Filename        string `yaml:"FileName"`
	ExecMode        string `yaml:"ExecMode"`
	SpeedMode       string `yaml:"SpeedMode"`
	ReadTimeout time  `yaml:Readtimeout`
	WriteTimeout time  `yaml:Writetimeout`
}

var Hosts []HostEntry

var passWord, userName, fileName, hostName, enable, logDir string
var readTimeout, writeTimeout time.Duration
var debug, speedMode, execMode bool
var outputFile, configFile string

func init() {
	flag.StringVar(&fileName, "filename", "", "Configuration file to insert")
	flag.StringVar(&hostName, "hostname", "rt1", "Router hostname")
	flag.StringVar(&passWord, "password", "password", "user password")
	flag.StringVar(&userName, "username", "username", "username")
	flag.StringVar(&enable, "enable", "enablepassword", "enable password")
	flag.DurationVar(&readTimeout, "readtimeout", time.Second*15, "timeout for reading poll on cli select")
	flag.DurationVar(&writeTimeout, "writetimeout", time.Millisecond*0, "timeout to stall after a write to cli")
	flag.BoolVar(&debug, "debug", false, "Enable debug for read / write")
	flag.BoolVar(&speedMode, "speedmode", false, "Enable speed mode write, will ignore any output from the cli while writing")
	flag.BoolVar(&execMode, "execmode", false, "Exec commands / input from filename instead of paste configuration")
	flag.StringVar(&logDir, "logdir", "", "Record session into logDir, automatically gzip")
	flag.StringVar(&outputFile, "outputfile", "", "Output file, else stdout")

	if os.Getenv("JK") != "" {
		flag.StringVar(&configFile, "configfile", "config_jk.yaml", "Input file in yaml for username,password and host configuration if not specified on command-line")
	} else {
		flag.StringVar(&configFile, "configfile", "config.yaml", "Input file in yaml for username,password and host configuration if not specified on command-line")
	}

	flag.Parse()
	if configFile != "" {
		loadConfig()
	}

}

func main() {
	var err error
	router := device.Brocade(device.DEVICE_MLX, hostName, 22, enable, userName, passWord,
		readTimeout, writeTimeout, debug, speedMode)

	if err = router.ConnectPrivilegedMode(); err != nil {
		log.Fatal(err)
	}

	if _, err = router.SkipPageDisplayMode(); err != nil {
		log.Fatal(err)
	}

	if err = router.GetPromptMode(); err != nil {
		log.Fatal(err)
	}

	if fileName != "" {
		file, err := os.Open(fileName)
		defer file.Close()
		if err != nil {
			log.Printf("Cant open file: %s", err)
		} else {
			if execMode == true {
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
	source, err := ioutil.ReadFile(configFile)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(source, &Hosts)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for _, Host := range Hosts {
		if Host.Hostname == hostName {
			if debug {
				log.Println("Overwrite cli settings for " + hostName + " from " + configFile)
			}

			if Host.Password != "" {
				passWord = Host.Password
			}

			if Host.EnablePassword != "" {
				enable = Host.EnablePassword
			}

			if Host.Username != "" {
				userName = Host.Username
			}

			if Host.Filename != "" {
				fileName = Host.Filename
			}

			if Host.ExecMode == "True" {
				execMode = true
			}
			if Host.SpeedMode == "True" {
				speedMode = true
			}
		}
	}
}
