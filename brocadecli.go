// Copyright 2017 Jörg Kost All rights reserved.
// joerg.kost@gmx.com, jk@premium-datacenter.de
// Use of this source code is governed by a GPLv2-style
// license that can be found in the LICENSE file.

package main

import (
	device "./device"
	"flag"
	"os"
	"log"
	"io"
	"time"
)


var passWord, userName, fileName, hostName, enable string
var readTimeout, writeTimeout time.Duration
var debug, speedMode bool

func init() {
	flag.StringVar(&fileName, "filename", "", "Configuration file to insert")
	flag.StringVar(&hostName, "hostname", "rt1", "Router hostname")
	flag.StringVar(&passWord, "password", "password", "user password")
	flag.StringVar(&userName, "username", "username", "username")
	flag.StringVar(&enable, "enable", "enablepassword", "enable password")
	flag.DurationVar(&readTimeout, "readtimeout", time.Second*15, "timeout for reading poll on cli select")
	flag.DurationVar(&writeTimeout, "writetimeout", time.Millisecond*0, "timeout to stall after a write to cli")
	flag.BoolVar(&debug, "debug", false,  "Enable debug for read / write")
	flag.BoolVar(&speedMode, "speedmode", false,  "Enable speed mode write, will ignore any output from the cli while writing")

}

func main() {
	flag.Parse()
	router := device.Brocade(device.DEVICE_MLX, hostName, 22, enable, userName, passWord,
		readTimeout,writeTimeout, debug, speedMode)

	router.ConnectPrivilegedMode()
	/* router.ExecPrivilegedMode("show ip route ... longer") */
	router.ConfigureTerminalMode()

	if fileName != "" {
		file, err := os.Open(fileName)
		defer file.Close()
		if err != nil {
			log.Printf("Cant open file: %s", err)
		} else {
			log.Println("START PROGRAMMING FROM CONFIGFILE")
			router.PasteConfiguration(io.Reader(file))
			log.Println("\nEND\n")
		}
	}

	router.CloseConnection()
}
