// Copyright 2017 Jörg Kost All rights reserved.
// joerg.kost@gmx.com, jk@premium-datacenter.de
// Use of this source code is governed by a GPLv2-style
// license that can be found in the LICENSE file.

package main

import (
	device "./device"
	"flag"
	"log"
	"os"
	"time"
)

var passWord, userName, fileName, hostName, enable, logDir string
var readTimeout, writeTimeout time.Duration
var debug, speedMode, execMode bool
var outputFile string

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
}

func main() {
	flag.Parse()
	router := device.Brocade(device.DEVICE_MLX, hostName, 22, enable, userName, passWord,
		readTimeout, writeTimeout, debug, speedMode)

	router.ConnectPrivilegedMode()
	router.SkipPageDisplayMode()

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
