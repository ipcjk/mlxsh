package libhost

import (
	"strings"
	"time"
)

/*
HostEntry represents the host structure,
that is being imported from the yaml configuration
 */
type HostEntry struct {
	Hostname        string            `yaml:"Hostname"`
	Username        string            `yaml:"Username"`
	Password        string            `yaml:"Password"`
	EnablePassword  string            `yaml:"EnablePassword"`
	SSHPort         int               `yaml:"SSHPort"`
	DeviceType      string            `yaml:"DeviceType"`
	KeyFile         string            `yaml:"KeyFile"`
	StrictHostCheck bool              `yaml:"StrictHostCheck"`
	Filename        string            `yaml:"FileName"`
	ScriptFile      string            `yaml:"ScriptFile"`
	ConfigFile      string            `yaml:"ConfigFile"`
	ExecMode        bool              `yaml:"ExecMode"`
	SpeedMode       bool              `yaml:"SpeedMode"`
	ReadTimeout     time.Duration     `yaml:"Readtimeout"`
	WriteTimeout    time.Duration     `yaml:"Writetimeout"`
	Labels          map[string]string `yaml:"Labels"`
}

/*
MatchLabels checks the given labels on the command line and
returns true or false
 */
func (h HostEntry) MatchLabels(userLabels string) bool {

	userLabels += ","

	selectArgs := strings.Split(userLabels, ",")
	for _, label := range selectArgs {
		partArgs := strings.Split(label, "=")

		if len(partArgs) < 2 {
			continue
		}

		if selectorValue, ok := h.Labels[partArgs[0]]; ok {
			if selectorValue != partArgs[1] {
				return false
			}
		} else {
			return false
		}

	}

	return true
}
