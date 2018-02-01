package main

import "io"

/*RouterInt is the minimum interface that a router module should have */
type RouterInt interface {
	Close()
	CommitConfiguration() error
	ConfigureTerminalMode() error
	Connect() error
	PasteConfiguration(io.Reader) error
	RunCommands(io.Reader) error
}
