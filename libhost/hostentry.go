package libhost

import (
	"fmt"
	"gopkg.in/yaml.v1"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

/*
HostEntry represents the host structure,
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
	Labels          map[string]string `yaml:"Labels"`
	Password        string            `yaml:"Password"`
	ReadTimeout     time.Duration     `yaml:"Readtimeout"`
	ScriptFile      string            `yaml:"ScriptFile"`
	SpeedMode       bool              `yaml:"SpeedMode"`
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

/* LoadFromYAML reads a yaml configuration reader source
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

/* ApplyCliSettings overwrites given cli parameters/set defaults
 */

func (h *HostConfig) ApplyCliSettings(scriptFile, configFile string, writeTimeout time.Duration, readTimeout time.Duration) {

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

}

/* LoadMatchFromYAML reads the yaml configuration reader source
and returns a slice of hosts that matches the given labels
*/
func LoadMatchesFromYAML(r io.Reader, label, hostname string) ([]HostConfig, error) {
	var hostsConfig []HostConfig
	var hostsMatch []HostConfig

	hostsConfig, err := LoadAllFromYAML(r)
	if err != nil {
		return []HostConfig{}, fmt.Errorf("Cant load from yaml source: %s", err)
	}

	for _, Host := range hostsConfig {

		if hostname != "" && hostname == Host.Hostname {
			hostsMatch = append(hostsMatch, Host)
			return hostsMatch, nil
		} else if Host.MatchLabels(label) {
			hostsMatch = append(hostsMatch, Host)
		}
	}

	return hostsMatch, nil
}
