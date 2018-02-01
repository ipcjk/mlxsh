package libhost

import (
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v1"
)

/*
HostConfig represents the host structure,
that is being imported from the yaml configuration
*/
type HostConfig struct {
	ConfigFile      string            `yaml:"ConfigFile"`
	DeviceType      string            `yaml:"DeviceType"`
	EnablePassword  string            `yaml:"EnablePassword"`
	ExecMode        bool              `yaml:"ExecMode"`
	Filename        string            `yaml:"FileName"`
	Hostname        string            `yaml:"Hostname"`
	KeyFile         string            `yaml:"KeyFile"`
	KnownHosts      string            `yaml:"KnownHosts"`
	Labels          map[string]string `yaml:"Labels"`
	Password        string            `yaml:"Password"`
	ReadTimeout     time.Duration     `yaml:"Readtimeout"`
	ScriptFile      string            `yaml:"ScriptFile"`
	SpeedMode       bool              `yaml:"SpeedMode"`
	SSHIP           string            `yaml:"SSHIP"`
	SSHPort         int               `yaml:"SSHPort"`
	StrictHostCheck bool              `yaml:"StrictHostCheck"`
	Username        string            `yaml:"Username"`
	WriteTimeout    time.Duration     `yaml:"Writetimeout"`
}

/*
MatchLabels checks the given labels on the command line and
returns true or false
*/
func (h HostConfig) MatchLabels(userLabels string) bool {

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

/*LoadAllFromYAML reads a yaml configuration reader source
and returns a slice of hosts
*/
func LoadAllFromYAML(r io.Reader) ([]HostConfig, error) {

	var hostsConfig []HostConfig

	source, err := ioutil.ReadAll(r)
	if err != nil {
		return []HostConfig{}, fmt.Errorf("Cant read from yaml source: %s", err)
	}

	err = yaml.Unmarshal(source, &hostsConfig)
	if err != nil {
		return []HostConfig{}, fmt.Errorf("Cant parse  yaml source: %s", err)
	}

	return hostsConfig, nil
}

/*ApplyCliSettings overwrites given cli parameters/set defaults */
func (h *HostConfig) ApplyCliSettings(scriptFile, configFile string, writeTimeout time.Duration, readTimeout time.Duration, HostCheck bool, KeyFile string, HostFile string) {

	if configFile != "" {
		h.Filename = configFile
		h.ConfigFile = configFile
		h.ExecMode = false
	}

	if scriptFile != "" {
		h.Filename = scriptFile
		h.ScriptFile = scriptFile
		h.ExecMode = true
	}

	if writeTimeout != 0 {
		h.WriteTimeout = writeTimeout
	}

	if readTimeout != 0 {
		h.ReadTimeout = readTimeout
	}

	if HostCheck {
		h.StrictHostCheck = true
	}

	if KeyFile != "" {
		h.KeyFile = KeyFile
	}

	if HostFile != "" {
		h.KnownHosts = HostFile
	}

}

/*LoadMatchesFromSlice reads the allHosts slice
and returns a slice of hosts that matches the given labels
*/
func LoadMatchesFromSlice(allHosts []HostConfig, label string) ([]HostConfig, error) {
	var hostsMatch []HostConfig

	for _, Host := range allHosts {
		if label != "" && Host.MatchLabels(label) {
			hostsMatch = append(hostsMatch, Host)
		} else if label == "" {
			hostsMatch = append(hostsMatch, Host)
		}
	}

	sort.Slice(hostsMatch, func(i, j int) bool {
		return hostsMatch[i].Hostname < hostsMatch[j].Hostname
	})

	return hostsMatch, nil
}

/*LoadMatchesFromYAML reads the yaml configuration reader source
and returns a slice of hosts that matches the given labels and also
a slice with all hosts from the yaml file
*/
func LoadMatchesFromYAML(r io.Reader, label, hostname string) ([]HostConfig, []HostConfig, error) {

	var allHosts []HostConfig
	var hostsMatch []HostConfig

	allHosts, err := LoadAllFromYAML(r)
	if err != nil {
		return []HostConfig{}, []HostConfig{}, fmt.Errorf("Cant load from yaml source: %s", err)
	}

	for _, Host := range allHosts {
		if hostname != "" && hostname == Host.Hostname {
			hostsMatch = append(hostsMatch, Host)
			return nil, hostsMatch, nil
		} else if label != "" && Host.MatchLabels(label) {
			hostsMatch = append(hostsMatch, Host)
		}
	}

	sort.Slice(hostsMatch, func(i, j int) bool {
		return hostsMatch[i].Hostname < hostsMatch[j].Hostname
	})

	sort.Slice(allHosts, func(i, j int) bool {
		return allHosts[i].Hostname < allHosts[j].Hostname
	})

	return hostsMatch, allHosts, nil
}
