package main

import "io"

/*RouterInt is the minimum interface that a router module should have */
type RouterInt interface {
	Close()
	ConfigureTerminalMode() error
	Connect() error
	PasteConfiguration(io.Reader) error
	RunCommands(io.Reader) error
	CommitConfiguration() error
}
