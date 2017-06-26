package libhost

import "time"

type HostEntry struct {
	Hostname        string `yaml:"Hostname"`
	Username        string `yaml:"Username"`
	Password        string `yaml:"Password"`
	EnablePassword  string `yaml:"EnablePassword"`
	DeviceType      string `yaml:"DeviceType"`
	KeyFile         string `yaml:"KeyFile"`
	StrictHostCheck bool   `yaml:"StrictHostCheck"`
	Filename        string `yaml:"FileName"`
	ExecMode        bool `yaml:"ExecMode"`
	SpeedMode       bool `yaml:"SpeedMode"`
	ReadTimeout time.Duration  `yaml:Readtimeout`
	WriteTimeout time.Duration  `yaml:Writetimeout`
}
