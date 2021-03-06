package main

/* this looks like a unique packages, but it actually belongs
to the main package and uses the following main global variables:

Hosts loaded and merged from the YAML-file:
selectedHosts
allHosts

Command prompt filter:
cliLabel
cliScriptFile
cliConfigFile

Sets a decision tree for the rl.guess function
cliCompleter

*/

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	rl "github.com/chzyer/readline"
	"github.com/ipcjk/mlxsh/libhost"
)

func filterInput(r rune) (rune, bool) {
	switch r {
	case rl.CharCtrlZ:
		return r, false
	}
	return r, true
}

func printHosts(x []libhost.HostConfig) {
	for _, v := range x {
		fmt.Printf("\x1b[32m%s\x1b[0m", v.Hostname)

		/* read and sort labels */
		labels := make([]string, 0)
		for k := range v.Labels {
			labels = append(labels, k)
		}
		sort.Strings(labels)

		for _, k := range labels {
			fmt.Printf(" \x1b[34m%s=%s\x1b[0m", k, v.Labels[k])
		}
		fmt.Printf("\n")
	}
}

func printAllHosts() {
	printHosts(allHosts)
}

func printSelectedHosts() {
	fmt.Println("Hosts matched:")
	if len(selectedHosts) == 0 {
		fmt.Println("None")
		return
	}
	printHosts(selectedHosts)
}

/* loadAutoCompletion
will load an autocompletion tree for the router / switches with the
most matches (counting DeviceType-field from yaml)
*/
func loadAutoCompletion(l *rl.Instance) {
	var defaultCliCompletionName = "netiron"
	var newCompletionName = ""
	var countCompletionName = -1

	if len(selectedHosts) == 0 {
		return
	}

	/* Try to load cliCompletion for the type of hosts with the most matches */
	var countDeviceTypes = make(map[string]int)
	if len(selectedHosts) > 0 {
		for x := range selectedHosts {
			switch strings.ToLower(selectedHosts[x].DeviceType) {
			case "vdx", "slx":
				countDeviceTypes["vdx"]++
			case "mlx", "cer", "mlxe", "xmr", "iron", "turobiron", "icx", "fcs":
				countDeviceTypes["netiron"]++
			case "juniper", "junos", "mx", "ex", "j":
				countDeviceTypes["junos"]++
			}
		}
	}

	for key := range countDeviceTypes {
		if countDeviceTypes[key] > countCompletionName {
			newCompletionName, countCompletionName = key, countDeviceTypes[key]
		}
	}

	if newCompletionName == "" {
		newCompletionName = defaultCliCompletionName
	}

	loadAutoCompletionNamed(l, newCompletionName)
}

func loadAutoCompletionNamed(l *rl.Instance, newCompletionName string) {
	/* Currently hardcoded ;( */
	switch newCompletionName {
	case "netiron", "mlx", "brocade", "cer", "tuborion", "mlxe":
		l.Config.AutoComplete = cliNetironCompleter
	case "junos", "juniper":
		l.Config.AutoComplete = cliJunOSCompleter
	case "vdx", "slx":
		l.Config.AutoComplete = cliVDXCompleter
	default:
		l.Config.AutoComplete = cliNetironCompleter
		newCompletionName = "netiron"
	}
	fmt.Println("Set", newCompletionName, "as default command line autocompletion tree")
}

/* setFilter executes a filter set on allHosts and will also load a pre-defined
auto completion tree */
func setFilter(label string) {
	var err error
	selectedHosts, err = libhost.LoadMatchesFromSlice(allHosts, label)
	if err != nil {
		fmt.Printf("Cant find any matches for the filter")
	}
	cliLabel = label

	printSelectedHosts()
}

/* Command line mode: Read commands or config statements from shell, execute on targets */
func runShellMode() {
	l, err := rl.NewEx(&rl.Config{
		Prompt:              "\033[31mmlxsh>\033[0m ",
		HistoryFile:         getUserHistoryFile(),
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}

	/* Set defaults from flags */
	setFilter(cliLabel)
	loadAutoCompletion(l)

	defer l.Close()

	for {
		line, err := l.Readline()
		if err == rl.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "clear "):
			switch line[6:] {
			default:
				prepareRunCmd("run", line[6:])
			}
		case strings.HasPrefix(line, "ls "):
			switch line[3:] {
			case "hosts", "selhosts":
				printSelectedHosts()
			case "allhosts":
				printAllHosts()
			case "filter":
				fmt.Println(cliLabel)
			default:
				printSelectedHosts()
			}
		case strings.HasPrefix(line, "request "):
			switch line[8:] {
			default:
				prepareRunCmd("run", line[8:])
			}
		case strings.HasPrefix(line, "run "):
			switch line[4:] {
			default:
				prepareRunCmd("run", line[4:])
			}
		case strings.HasPrefix(line, "mset filter "):
			setFilter(line[12:])
			loadAutoCompletion(l)
		case strings.HasPrefix(line, "mset complete "):
			loadAutoCompletionNamed(l, line[14:])
		case strings.HasPrefix(line, "show "):
			switch line[5:] {
			default:
				prepareRunCmd("run", "show "+line[5:])
			}
		case line == "exit":
			goto exit
		case line == "":
		default:
			fmt.Println("command not found or param missing:", strconv.Quote(line))
		}
	}
exit:
	/* close all connections? */
	/* cleanup some things? */
}

func prepareRunCmd(confOrRun, line string) {
	if confOrRun == "run" {
		cliScriptFile = line
		cliConfigFile = ""
	} else {
		cliConfigFile = line
		cliScriptFile = ""
	}
	run()
}
