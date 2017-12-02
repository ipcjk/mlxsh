package main

import "io"

type RouterInt interface {
	Close()
	ConfigureTerminalMode() error
	Connect() error
	PasteConfiguration(io.Reader) error
	RunCommands(io.Reader) error
	CommitConfiguration() error
}
