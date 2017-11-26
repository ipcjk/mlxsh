package main

import "io"

type RouterInt interface {
	CloseConnection()
	ConfigureTerminalMode() error
	ConnectPrivilegedMode() error
	PasteConfiguration(io.Reader) error
	RunCommands(io.Reader) error
	WriteConfiguration() error
}
