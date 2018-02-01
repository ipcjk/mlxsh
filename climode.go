package main

/* this looks like a unique packages, but it actually belongs
to the main package and uses the following main global variables:

Hosts loaded and merged from the YAML-file:
selectedHosts
allHosts

Command prompt filter:
cliLabel

Sets a decision tree for the rl.guess function
cliCompleter

*/

import (
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"

	rl "github.com/chzyer/readline"
	"github.com/ipcjk/mlxsh/libhost"
)

var cliCompleter = rl.NewPrefixCompleter(
	rl.PcItem("show",
		rl.PcItem("access-list"),
		rl.PcItem("acl-policy"),
		rl.PcItem("arp"),
		rl.PcItem("bfd",
			rl.PcItem("applications"),
			rl.PcItem("mpls"),
			rl.PcItem("neighbors"),
			rl.PcItem("neighbors bgp"),
			rl.PcItem("neighbors details"),
			rl.PcItem("neighbors interface"),
			rl.PcItem("neighbors isis"),
			rl.PcItem("neighbors ospf"),
			rl.PcItem("neighbors ospf"),
			rl.PcItem("neighbors static"),
			rl.PcItem("neighbors static")),
		rl.PcItem("chassis"),
		rl.PcItem("configuration "),
		rl.PcItem("cpu histogram"),
		rl.PcItem("interface ethernet"),
		rl.PcItem("interfaces tunnel"),
		rl.PcItem("ip",
			rl.PcItem("bgp", rl.PcItem("attribute-entries"),
				rl.PcItem("config"),
				rl.PcItem("dampened-paths"),
				rl.PcItem("filtered-routes "),
				rl.PcItem("flap-statistics "),
				rl.PcItem("ipv6"),
				rl.PcItem("neighbors "),
				rl.PcItem("neighbors advertised-routes"),
				rl.PcItem("neighbors flap-statistics "),
				rl.PcItem("neighbors last-packet-with-error "),
				rl.PcItem("neighbors received "),
				rl.PcItem("neighbors received-routes"),
				rl.PcItem("neighbors rib-out-routes"),
				rl.PcItem("routes community "),
				rl.PcItem("neighbors routes "),
				rl.PcItem("neighbors routes-summary"),
				rl.PcItem("peer-group"),
				rl.PcItem("routes"),
				rl.PcItem("summary"),
				rl.PcItem("vrf neighbors "),
				rl.PcItem("vrf routes "),
				rl.PcItem("vrf ")),
			rl.PcItem("interface"),
			rl.PcItem("mbgp ipv6"),
			rl.PcItem("multicast"),
			rl.PcItem("multicast vpls"),
			rl.PcItem("ospf"),
			rl.PcItem("route"),
			rl.PcItem("static-arp"),
			rl.PcItem("vrrp"),
			rl.PcItem("vrrp-extended")),
		rl.PcItem("ipsec",
			rl.PcItem("egress-config"),
			rl.PcItem("egress-spi-table"),
			rl.PcItem("error-count"),
			rl.PcItem("ingress-config"),
			rl.PcItem("ingress-spi-table"),
			rl.PcItem("policy"),
			rl.PcItem("profile"),
			rl.PcItem("proposal"),
			rl.PcItem("sa"),
			rl.PcItem("statistics")),
		rl.PcItem("ip-tunnels"),
		rl.PcItem("ipv6",
			rl.PcItem("access-list bindings "),
			rl.PcItem("access-list receive accounting "),
			rl.PcItem("bgp"),
			rl.PcItem("bgp neighbors"),
			rl.PcItem("bgp routes "),
			rl.PcItem("bgp summary"),
			rl.PcItem("dhcp-relay interface"),
			rl.PcItem("dhcp-relay options"),
			rl.PcItem("interface tunnel"),
			rl.PcItem("ospf interface "),
			rl.PcItem("vrrp"),
			rl.PcItem("vrrp-extended")),
		rl.PcItem("isis"),
		rl.PcItem("license"),
		rl.PcItem("log"),
		rl.PcItem("module"),
		rl.PcItem("mpls",
			rl.PcItem("autobw-threshold-table "),
			rl.PcItem("bypass-lsp"),
			rl.PcItem("config"),
			rl.PcItem("forwarding"),
			rl.PcItem("interface"),
			rl.PcItem("label-range"),
			rl.PcItem("ldp"),
			rl.PcItem("ldp database"),
			rl.PcItem("ldp fec"),
			rl.PcItem("ldp interface"),
			rl.PcItem("ldp neighbor"),
			rl.PcItem("ldp path"),
			rl.PcItem("ldp peer"),
			rl.PcItem("ldp session "),
			rl.PcItem("ldp statistics"),
			rl.PcItem("ldp tunnel "),
			rl.PcItem("lsp"),
			rl.PcItem("lsp_pmp_xc "),
			rl.PcItem("path"),
			rl.PcItem("policy "),
			rl.PcItem("route "),
			rl.PcItem("rsvp",
				rl.PcItem("interface"),
				rl.PcItem("neighbor"),
				rl.PcItem("session"),
				rl.PcItem("session backup"),
				rl.PcItem("session brief"),
				rl.PcItem("session bypass"),
				rl.PcItem("session destination"),
				rl.PcItem("session detail"),
				rl.PcItem("session detour"),
				rl.PcItem("session down"),
				rl.PcItem("session extensive"),
				rl.PcItem("session (ingress/egress)"),
				rl.PcItem("session (interface)"),
				rl.PcItem("session name"),
				rl.PcItem("session pmp"),
				rl.PcItem("session pp"),
				rl.PcItem("session ppend"),
				rl.PcItem("session transit"),
				rl.PcItem("session up"),
				rl.PcItem("session wide"),
				rl.PcItem("statistics")),
			rl.PcItem("static-lsp"),
			rl.PcItem("statistics",
				rl.PcItem("pe"),
				rl.PcItem("bypass-lsp"),
				rl.PcItem("label"),
				rl.PcItem("ldp transit"),
				rl.PcItem("ldp tunnel "),
				rl.PcItem("lsp"),
				rl.PcItem("oam"),
				rl.PcItem("vll"),
				rl.PcItem("vll-local"),
				rl.PcItem("vpls"),
				rl.PcItem("vrf")),
			rl.PcItem("summary"),
			rl.PcItem("ted database"),
			rl.PcItem("ted path"),
			rl.PcItem("vll"),
			rl.PcItem("vll-local"),
			rl.PcItem("vpls")),
		rl.PcItem("openflow",
			rl.PcItem("controller"),
			rl.PcItem("flows"),
			rl.PcItem("groups"),
			rl.PcItem("interface"),
			rl.PcItem("meters"),
			rl.PcItem("queues")),
		rl.PcItem("rate-limit",
			rl.PcItem("counters bum-drop"),
			rl.PcItem("detail"),
			rl.PcItem("interface"),
			rl.PcItem("ipv6 hoplimit-expired-to-cpu"),
			rl.PcItem("option-pkt-to-cpu"),
			rl.PcItem("ttl-expired-to-cpu")),
		rl.PcItem("route-map"),
		rl.PcItem("running-config"),
		rl.PcItem("sflow statistics "),
		rl.PcItem("spanning-tree "),
		rl.PcItem("statistics "),
		rl.PcItem("terminal"),
		rl.PcItem("version"),
		rl.PcItem("vlan"),
	),
	rl.PcItem("get",
		rl.PcItem("filter"),
		rl.PcItem("hosts"),
		rl.PcItem("selhosts"),
		rl.PcItem("allhosts"),
	),
	rl.PcItem("set",
		rl.PcItem("filter"),
	),
)

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case rl.CharCtrlZ:
		return r, false
	}
	return r, true
}

func printHosts(x []libhost.HostConfig) {
	for _, v := range x {
		fmt.Printf("\x1b[32m%s", v.Hostname)

		/* read and sort labels */
		labels := make([]string, 0)
		for k := range v.Labels {
			labels = append(labels, k)
		}
		sort.Strings(labels)

		for _, k := range labels {
			fmt.Printf(" \x1b[34m%s=%s", k, v.Labels[k])
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
func runCliMode() {

	l, err := rl.NewEx(&rl.Config{
		Prompt:              "\033[31mmlxsh>\033[0m ",
		HistoryFile:         getUserHistoryFile(),
		AutoComplete:        cliCompleter,
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	//log.SetOutput(l.Stderr())
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
		case strings.HasPrefix(line, "show"):
			switch line[5:] {
			default:
				prepareRunCmd("run", "show "+line[5:])
			}
		case strings.HasPrefix(line, "get"):
			switch line[4:] {
			case "hosts":
			case "selhosts":
				printSelectedHosts()
			case "allhosts":
				printAllHosts()
			case "filter":
				fmt.Println(cliLabel)
			}
		case strings.HasPrefix(line, "set filter"):
			setFilter(line[11:])
		case line == "help":
			usage(l.Stderr())
		case line == "bye":
			goto exit
		case line == "":
		default:
			log.Println("you said:", strconv.Quote(line))
		}
	}
exit:
	/* close all connections? */
	/* cleanup some things? */
}

func usage(w io.Writer) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, cliCompleter.Tree("    "))
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
