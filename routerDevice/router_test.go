package router_test

import (
	"bytes"
	"github.com/ipcjk/mlxsh/libhost"
	"github.com/ipcjk/mlxsh/routerDevice"
	"regexp"
	"testing"
)

func TestCreatePseudoRouter(t *testing.T) {
	var configureErrors = `(?i)(invalid command|unknown|Warning|Error|not found)`

	var pseudoRouter = &router.Router{
		PromptModes:        make(map[string]string),
		ErrorMatches:       regexp.MustCompile(configureErrors),
		PromptDetect:       `[@?\.\d\w-]+> ?$`,
		PromptReadTriggers: []string{">"},
		PromptReplacements: map[string][]string{
			"SSHConfigPrompt":    {">", "#"},
			"SSHConfigPromptPre": {">", "#"}},
	}

	pseudoRouter.Close()

}

func TestDetectPrompt(t *testing.T) {
	var rtc router.RunTimeConfig

	var configureErrors = `(?i)(invalid command|unknown|Warning|Error|not found)`
	var pseudoRouter = router.Router{
		PromptModes:        make(map[string]string),
		ErrorMatches:       regexp.MustCompile(configureErrors),
		PromptDetect:       `[@?\.\d\w-]+> ?$`,
		PromptReadTriggers: []string{">"},
		PromptReplacements: map[string][]string{
			"SSHConfigPrompt":    {">", "#"},
			"SSHConfigPromptPre": {">", "#"}},
	}

	if err := pseudoRouter.DetectPrompt(rtc, "joerg@core-10>"); err != nil {
		t.Errorf("Cant detect prompt! :%s", err)
	}

}

func CreatePseudoRunTimeConfig(t *testing.T) {
	var buffer = new(bytes.Buffer)

	var rtc = &router.RunTimeConfig{
		HostConfig: libhost.HostConfig{Hostname: "core-10", Username: "joerg", Password: "foobar"}, W: buffer}

	router.GenerateDefaults(rtc)

	if rtc.SSHPort != 22 {
		t.Error("Wrong default SSH-Port")
	}

	if rtc.ReadTimeout == 0 {
		t.Error("ReadTimeout of 0 is dangerous!")
	}

}
