package main

import rl "github.com/chzyer/readline"

/* defaultGetCompletion
is a parameter list of the default parameters for the get command. This
list will be added to every completer */
var defaultGetCompletion = rl.PcItem("get",
	rl.PcItem("filter"),
	rl.PcItem("hosts"),
	rl.PcItem("selhosts"),
	rl.PcItem("allhosts"))

/* defaultSetCompletion
is a parameter list of the default parameters for the get command. This
list will be added to every completer */
var defaultSetCompletion = rl.PcItem("set",
	rl.PcItem("filter"),
	rl.PcItem("complete",
		rl.PcItem("junos"),
		rl.PcItem("netiron"),
		rl.PcItem("vdx")),
)

/* cliNetironCompleter
is an autocompletion tree for the Netiron command line
*/
var cliNetironCompleter = rl.NewPrefixCompleter(
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
	defaultGetCompletion,
	defaultSetCompletion,
)

/* cliJunOSCompleter
is an autocompletion tree for the Netiron command line
*/
var cliJunOSCompleter = rl.NewPrefixCompleter(
	rl.PcItem("show",
		rl.PcItem("arp", rl.PcItem("no-resolve")),
		rl.PcItem("bfd"),
		rl.PcItem("bgp",
			rl.PcItem("group"),
			rl.PcItem("neighbor"),
			rl.PcItem("summary"),
		),
		rl.PcItem("chassis",
			rl.PcItem("alarms"),
			rl.PcItem("environment"),
			rl.PcItem("firmware"),
			rl.PcItem("hardware"),
			rl.PcItem("location"),
			rl.PcItem("pic"),
			rl.PcItem("pic-mode"),
			rl.PcItem("routing-engine"),
		),
		rl.PcItem("ethernet-switching",
			rl.PcItem("filters"),
			rl.PcItem("interfaces"),
			rl.PcItem("next-hops"),
			rl.PcItem("statistics"),
			rl.PcItem("table"),
		),
		rl.PcItem("firewall",
			rl.PcItem("application"),
			rl.PcItem("counter"),
			rl.PcItem("filter"),
			rl.PcItem("log"),
			rl.PcItem("terse"),
		),
		rl.PcItem("igmp"),
		rl.PcItem("interfaces"),
		rl.PcItem("ipv6",
			rl.PcItem("neighbors"),
			rl.PcItem("router-advertisement"),
		),
		rl.PcItem("lacp",
			rl.PcItem("interfaces"),
			rl.PcItem("statistics"),
			rl.PcItem("timeouts"),
		),
		rl.PcItem("lldp",
			rl.PcItem("detail"),
			rl.PcItem("local-information"),
			rl.PcItem("neighbors"),
			rl.PcItem("statistics"),
		),
		rl.PcItem("show"),
		rl.PcItem("mpls",
			rl.PcItem("interface"),
			rl.PcItem("lsp"),
			rl.PcItem("path"),
		),
		rl.PcItem("ospf",
			rl.PcItem("database"),
			rl.PcItem("interface"),
			rl.PcItem("log"),
			rl.PcItem("neighbor",
				rl.PcItem("area"),
				rl.PcItem("brief"),
				rl.PcItem("detail"),
				rl.PcItem("extensive"),
				rl.PcItem("instance"),
				rl.PcItem("interface"),
			),
			rl.PcItem("overview"),
			rl.PcItem("route"),
			rl.PcItem("statistics"),
		),
		rl.PcItem("ospf3",
			rl.PcItem("database"),
			rl.PcItem("interface"),
			rl.PcItem("log"),
			rl.PcItem("neighbor",
				rl.PcItem("area"),
				rl.PcItem("brief"),
				rl.PcItem("detail"),
				rl.PcItem("extensive"),
				rl.PcItem("instance"),
				rl.PcItem("interface"),
			),
			rl.PcItem("overview"),
			rl.PcItem("route"),
			rl.PcItem("statistics"),
		),
		rl.PcItem("route",
			rl.PcItem("advertising-protocol"),
			rl.PcItem("best"),
			rl.PcItem("brief"),
			rl.PcItem("detail"),
			rl.PcItem("instance"),
			rl.PcItem("martians"),
			rl.PcItem("next-hop"),
			rl.PcItem("protocol",
				rl.PcItem("arp"),
				rl.PcItem("bgp"),
				rl.PcItem("direct"),
				rl.PcItem("isis"),
				rl.PcItem("local"),
				rl.PcItem("mpls"),
				rl.PcItem("ospf"),
				rl.PcItem("ospf2"),
				rl.PcItem("ospf3"),
				rl.PcItem("static"),
			),
		),
		rl.PcItem("sflow"),
		rl.PcItem("snmp"),
		rl.PcItem("spanning-tree"),
		rl.PcItem("system",
			rl.PcItem("alarms",
				rl.PcItem("arp"),
				rl.PcItem("commit"),
				rl.PcItem("configurations"),
				rl.PcItem("connections"),
				rl.PcItem("license"),
				rl.PcItem("login"),
				rl.PcItem("memory"),
				rl.PcItem("processes"),
				rl.PcItem("reboot"),
				rl.PcItem("services"),
				rl.PcItem("software"),
				rl.PcItem("storage"),
				rl.PcItem("uptime"),
				rl.PcItem("users"),
			)),
		rl.PcItem("version"),
		rl.PcItem("virtual-chassis"),
		rl.PcItem("vlans"),
		rl.PcItem("vrrp"),
	),
	defaultGetCompletion,
	defaultSetCompletion,
)

var cliVDXCompleter = rl.NewPrefixCompleter(
	rl.PcItem("show access-list"),
	rl.PcItem("show access-list-log buffer"),
	rl.PcItem("show access-list-log buffer config"),
	rl.PcItem("show ag"),
	rl.PcItem("show ag map"),
	rl.PcItem("show ag nport-utilization"),
	rl.PcItem("show ag pg"),
	rl.PcItem("show arp"),
	rl.PcItem("show arp access-list"),
	rl.PcItem("show bare-metal"),
	rl.PcItem("show bfd"),
	rl.PcItem("show bfd neighbors"),
	rl.PcItem("show bfd neighbors application"),
	rl.PcItem("show bfd neighbors dest-ip"),
	rl.PcItem("show bfd neighbors details"),
	rl.PcItem("show bfd neighbors interface"),
	rl.PcItem("show bfd neighbors vrf"),
	rl.PcItem("show bgp evpn dampened-routes"),
	rl.PcItem("show bgp evpn interface port-channel"),
	rl.PcItem("show bgp evpn interface tunnel"),
	rl.PcItem("show bgp evpn l2route detail"),
	rl.PcItem("show bgp evpn l2route next-hop"),
	rl.PcItem("show bgp evpn l2route summary"),
	rl.PcItem("show bgp evpn l2route type arp"),
	rl.PcItem("show bgp evpn l2route type auto-discovery"),
	rl.PcItem("show bgp evpn l2route type ethernet-segment"),
	rl.PcItem("show bgp evpn l2route type inclusive-multicast"),
	rl.PcItem("show bgp evpn l2route type mac"),
	rl.PcItem("show bgp evpn l2route type nd"),
	rl.PcItem("show bgp evpn l2route unreachable"),
	rl.PcItem("show bgp evpn l3vni"),
	rl.PcItem("show bgp evpn neighbors"),
	rl.PcItem("show bgp evpn neighbors advertised-routes detail"),
	rl.PcItem("show bgp evpn neighbors advertised-routes type"),
	rl.PcItem("show bgp evpn neighbors routes best"),
	rl.PcItem("show bgp evpn neighbors routes detail"),
	rl.PcItem("show bgp evpn neighbors routes not-installed-best"),
	rl.PcItem("show bgp evpn neighbors routes type"),
	rl.PcItem("show bgp evpn neighbors routes unreachable"),
	rl.PcItem("show bgp evpn neighbors routes-summary"),
	rl.PcItem("show bgp evpn routes"),
	rl.PcItem("show bgp evpn routes best"),
	rl.PcItem("show bgp evpn routes detail"),
	rl.PcItem("show bgp evpn routes local"),
	rl.PcItem("show bgp evpn routes next-hop"),
	rl.PcItem("show bgp evpn routes no-best"),
	rl.PcItem("show bgp evpn routes not-installed-best"),
	rl.PcItem("show bgp evpn routes rd"),
	rl.PcItem("show bgp evpn routes rd type"),
	rl.PcItem("show bgp evpn routes summary"),
	rl.PcItem("show bgp evpn routes type arp"),
	rl.PcItem("show bgp evpn routes type auto-discovery"),
	rl.PcItem("show bgp evpn routes type ethernet-segment"),
	rl.PcItem("show bgp evpn routes type inclusive-multicast"),
	rl.PcItem("show bgp evpn routes type ipv4-prefix"),
	rl.PcItem("show bgp evpn routes type ipv6-prefix"),
	rl.PcItem("show bgp evpn routes type mac"),
	rl.PcItem("show bgp evpn routes type nd"),
	rl.PcItem("show bgp evpn routes unreachable"),
	rl.PcItem("show bgp evpn summary"),
	rl.PcItem("show bpdu-drop"),
	rl.PcItem("show capture packet interface"),
	rl.PcItem("show cee maps"),
	rl.PcItem("show cert-util ldapca"),
	rl.PcItem("show cert-util sshkey"),
	rl.PcItem("show cert-util syslogca"),
	rl.PcItem("show chassis"),
	rl.PcItem("show cipherset"),
	rl.PcItem("show class-maps"),
	rl.PcItem("show cli"),
	rl.PcItem("show cli history"),
	rl.PcItem("show clock"),
	rl.PcItem("show config snapshot"),
	rl.PcItem("show copy-support status"),
	rl.PcItem("show crypto ca"),
	rl.PcItem("show crypto key"),
	rl.PcItem("show dadstatus"),
	rl.PcItem("show debug dhcp packet"),
	rl.PcItem("show debug dhcp packet buffer"),
	rl.PcItem("show debug ip bgp all"),
	rl.PcItem("show debug ip igmp"),
	rl.PcItem("show debug ip pim"),
	rl.PcItem("show debug ipv6 packet"),
	rl.PcItem("show debug lacp"),
	rl.PcItem("show debug lldp"),
	rl.PcItem("show debug spanning-tree"),
	rl.PcItem("show debug udld"),
	rl.PcItem("show debug vrrp"),
	rl.PcItem("show defaults threshold"),
	rl.PcItem("show default-vlan"),
	rl.PcItem("show dpod"),
	rl.PcItem("show diag burninerrshow"),
	rl.PcItem("show diag burninerrshowerrLog"),
	rl.PcItem("show diag burninstatus"),
	rl.PcItem("show diag post results"),
	rl.PcItem("show diag setcycle"),
	rl.PcItem("show diag status"),
	rl.PcItem("show dot1x"),
	rl.PcItem("show dot1x all"),
	rl.PcItem("show dot1x diagnostics interface"),
	rl.PcItem("show dot1x interface"),
	rl.PcItem("show dot1x session-info interface"),
	rl.PcItem("show dot1x statistics interface"),
	rl.PcItem("show dpod"),
	rl.PcItem("show dport-test"),
	rl.PcItem("show edge-loop-detection detail"),
	rl.PcItem("show edge-loop-detection globals"),
	rl.PcItem("show edge-loop-detection interface"),
	rl.PcItem("show edge-loop-detection rbridge-id"),
	rl.PcItem("show environment fan"),
	rl.PcItem("show environment history"),
	rl.PcItem("show environment power"),
	rl.PcItem("show environment sensor"),
	rl.PcItem("show environment temp"),
	rl.PcItem("show event-handler activations"),
	rl.PcItem("show fabric all"),
	rl.PcItem("show fabric ecmp group"),
	rl.PcItem("show fabric ecmp load-balance"),
	rl.PcItem("show fabric isl"),
	rl.PcItem("show fabric islports"),
	rl.PcItem("show fabric login-policy"),
	rl.PcItem("show fabric port-channel"),
	rl.PcItem("show fabric route linkinfo"),
	rl.PcItem("show fabric route multicast"),
	rl.PcItem("show fabric route neighbor-state"),
	rl.PcItem("show fabric route pathinfo"),
	rl.PcItem("show fabric route topology"),
	rl.PcItem("show fabric trunk"),
	rl.PcItem("show fcoe devices"),
	rl.PcItem("show fcoe fabric-map"),
	rl.PcItem("show fcoe fcf-group"),
	rl.PcItem("show fcoe fcoe-enodes"),
	rl.PcItem("show fcoe fcoe-map"),
	rl.PcItem("show fcoe fcport-group"),
	rl.PcItem("show fcoe interface"),
	rl.PcItem("show fcoe login"),
	rl.PcItem("show fcsp auth-secret dh-chap"),
	rl.PcItem("show fibrechannel login"),
	rl.PcItem("show file"),
	rl.PcItem("show fips"),
	rl.PcItem("show firmwaredownloadhistory"),
	rl.PcItem("show firmwaredownloadstatus"),
	rl.PcItem("show global-running-config"),
	rl.PcItem("show ha"),
	rl.PcItem("show hardware port-group"),
	rl.PcItem("show hardware connector-group"),
	rl.PcItem("show hardware-profile"),
	rl.PcItem("show history"),
	rl.PcItem("show http server status"),
	rl.PcItem("show interface"),
	rl.PcItem("show interface description"),
	rl.PcItem("show interface FibreChannel"),
	rl.PcItem("show interface management"),
	rl.PcItem("show interface stats"),
	rl.PcItem("show interface status"),
	rl.PcItem("show interface trunk"),
	rl.PcItem("show inventory"),
	rl.PcItem("show ip anycast-gateway"),
	rl.PcItem("show ip arp inspection"),
	rl.PcItem("show ip arp inspection interfaces"),
	rl.PcItem("show ip arp inspection statistics"),
	rl.PcItem("show ip arp suppression-cache"),
	rl.PcItem("show ip arp suppression-statistics"),
	rl.PcItem("show ip arp suppression-status"),
	rl.PcItem("show ip as-path-list"),
	rl.PcItem("show ip bgp"),
	rl.PcItem("show ip bgp attribute-entries"),
	rl.PcItem("show ip bgp dampened-paths"),
	rl.PcItem("show ip bgp filtered-routes"),
	rl.PcItem("show ip bgp flap-statistics"),
	rl.PcItem("show ip bgp neighbors"),
	rl.PcItem("show ip bgp neighbors advertised-routes"),
	rl.PcItem("show ip bgp neighbors flap-statistics"),
	rl.PcItem("show ip bgp neighbors received"),
	rl.PcItem("show ip bgp neighbors received-routes"),
	rl.PcItem("show ip bgp neighbors routes"),
	rl.PcItem("show ip bgp neighbors routes-summary"),
	rl.PcItem("show ip bgp peer-group"),
	rl.PcItem("show ip bgp rbridge-id"),
	rl.PcItem("show ip bgp routes"),
	rl.PcItem("show ip bgp routes community"),
	rl.PcItem("show ip bgp summary"),
	rl.PcItem("show ip community-list"),
	rl.PcItem("show ip dhcp relay address interface"),
	rl.PcItem("show ip dhcp relay address rbridge-id"),
	rl.PcItem("show ip dhcp relay gateway"),
	rl.PcItem("show ip dhcp relay statistics"),
	rl.PcItem("show ip dns"),
	rl.PcItem("show ip extcommunity-list"),
	rl.PcItem("show ip fabric-virtual-gateway"),
	rl.PcItem("show ip igmp groups"),
	rl.PcItem("show ip igmp interface"),
	rl.PcItem("show ip igmp snooping"),
	rl.PcItem("show ip igmp static-groups"),
	rl.PcItem("show ip igmp statistics interface"),
	rl.PcItem("show ip interface"),
	rl.PcItem("show ip interface loopback"),
	rl.PcItem("show ip interface ve"),
	rl.PcItem("show ip next-hop"),
	rl.PcItem("show ip ospf"),
	rl.PcItem("show ip ospf area"),
	rl.PcItem("show ip ospf border-routers"),
	rl.PcItem("show ip ospf config"),
	rl.PcItem("show ip ospf database"),
	rl.PcItem("show ip ospf filtered-lsa area"),
	rl.PcItem("show ip ospf interface"),
	rl.PcItem("show ip ospf neighbor"),
	rl.PcItem("show ip ospf redistribute route"),
	rl.PcItem("show ip ospf routes"),
	rl.PcItem("show ip ospf summary"),
	rl.PcItem("show ip ospf traffic"),
	rl.PcItem("show ip ospf virtual"),
	rl.PcItem("show ip pim bsr"),
	rl.PcItem("show ip pim group"),
	rl.PcItem("show ip pim mcache"),
	rl.PcItem("show ip pim neighbor"),
	rl.PcItem("show ip pim rpf"),
	rl.PcItem("show ip pim rp-hash"),
	rl.PcItem("show ip pim rp-map"),
	rl.PcItem("show ip pim rp-set"),
	rl.PcItem("show ip pim traffic"),
	rl.PcItem("show ip pim-sparse"),
	rl.PcItem("show ip prefix-list"),
	rl.PcItem("show ip route"),
	rl.PcItem("show ip route import"),
	rl.PcItem("show ip route system-summary"),
	rl.PcItem("show ip route-map"),
	rl.PcItem("show ip static route"),
	rl.PcItem("show ipv6 anycast-gateway"),
	rl.PcItem("show ipv6 bgp attribute-entries"),
	rl.PcItem("show ipv6 bgp dampened-paths"),
	rl.PcItem("show ipv6 bgp filtered-routes"),
	rl.PcItem("show ipv6 bgp filtered-routes detail"),
	rl.PcItem("show ipv6 bgp flap-statistics"),
	rl.PcItem("show ipv6 bgp neighbors"),
	rl.PcItem("show ipv6 bgp neighbors advertised-routes"),
	rl.PcItem("show ipv6 bgp neighbors flap-statistics"),
	rl.PcItem("show ipv6 bgp neighbors last-packet-with-error"),
	rl.PcItem("show ipv6 bgp neighbors rbridge-id"),
	rl.PcItem("show ipv6 bgp neighbors received"),
	rl.PcItem("show ipv6 bgp neighbors received-routes"),
	rl.PcItem("show ipv6 bgp neighbors rib-out-routes"),
	rl.PcItem("show ipv6 bgp neighbors routes"),
	rl.PcItem("show ipv6 bgp neighbors routes-summary"),
	rl.PcItem("show ipv6 bgp peer-group"),
	rl.PcItem("show ipv6 bgp rbridge-id"),
	rl.PcItem("show ipv6 bgp routes"),
	rl.PcItem("show ipv6 bgp routes community"),
	rl.PcItem("show ipv6 bgp summary"),
	rl.PcItem("show ipv6 counters interface"),
	rl.PcItem("show ipv6 dhcp relay address interface"),
	rl.PcItem("show ipv6 dhcp relay address rbridge-id"),
	rl.PcItem("show ipv6 dhcp relay statistics"),
	rl.PcItem("show ipv6 fabric-virtual-gateway"),
	rl.PcItem("show ipv6 interface"),
	rl.PcItem("show ipv6 mld groups"),
	rl.PcItem("show ipv6 mld interface"),
	rl.PcItem("show ipv6 mld snooping"),
	rl.PcItem("show ipv6 mld statistics"),
	rl.PcItem("show ipv6 nd interface"),
	rl.PcItem("show ipv6 nd suppression-cache"),
	rl.PcItem("show ipv6 nd suppression-statistics"),
	rl.PcItem("show ipv6 nd suppression-status"),
	rl.PcItem("show ipv6 neighbor"),
	rl.PcItem("show ipv6 ospf"),
	rl.PcItem("show ipv6 ospf area"),
	rl.PcItem("show ipv6 ospf database"),
	rl.PcItem("show ipv6 ospf interface"),
	rl.PcItem("show ipv6 ospf memory"),
	rl.PcItem("show ipv6 ospf neighbor"),
	rl.PcItem("show ipv6 ospf redistribute route"),
	rl.PcItem("show ipv6 ospf routes"),
	rl.PcItem("show ipv6 ospf spf"),
	rl.PcItem("show ipv6 ospf summary"),
	rl.PcItem("show ipv6 ospf virtual-links"),
	rl.PcItem("show ipv6 ospf virtual-neighbor"),
	rl.PcItem("show ipv6 prefix-list"),
	rl.PcItem("show ipv6 raguard"),
	rl.PcItem("show ipv6 route"),
	rl.PcItem("show ipv6 static route"),
	rl.PcItem("show ipv6 vrrp"),
	rl.PcItem("show lacp"),
	rl.PcItem("show lacp sys-id"),
	rl.PcItem("show license"),
	rl.PcItem("show license id"),
	rl.PcItem("show linecard"),
	rl.PcItem("show lldp interface"),
	rl.PcItem("show lldp neighbors"),
	rl.PcItem("show lldp statistics"),
	rl.PcItem("show logging auditlog"),
	rl.PcItem("show logging raslog"),
	rl.PcItem("show mac-address-table"),
	rl.PcItem("show mac-address-table consistency-check"),
	rl.PcItem("show mac-address-table mac-move"),
	rl.PcItem("show maps dashboard"),
	rl.PcItem("show maps group"),
	rl.PcItem("show maps policy"),
	rl.PcItem("show media"),
	rl.PcItem("show media interface"),
	rl.PcItem("show media linecard"),
	rl.PcItem("show media optical-monitoring interface"),
	rl.PcItem("show media tunable-optic-sfpp"),
	rl.PcItem("show mgmt-ip-service"),
	rl.PcItem("show mm"),
	rl.PcItem("show monitor"),
	rl.PcItem("show name-server brief"),
	rl.PcItem("show name-server detail"),
	rl.PcItem("show name-server nodefind"),
	rl.PcItem("show name-server zonemember"),
	rl.PcItem("show nas statistics"),
	rl.PcItem("show netconf client-capabilities"),
	rl.PcItem("show netconf-state capabilities"),
	rl.PcItem("show netconf-state datastores"),
	rl.PcItem("show netconf-state schemas"),
	rl.PcItem("show netconf-state sessions"),
	rl.PcItem("show netconf-state statistics"),
	rl.PcItem("show notification stream"),
	rl.PcItem("show nsx-controller"),
	rl.PcItem("show ntp status"),
	rl.PcItem("show openflow"),
	rl.PcItem("show openflow controller"),
	rl.PcItem("show openflow flow"),
	rl.PcItem("show openflow group"),
	rl.PcItem("show openflow interface"),
	rl.PcItem("show openflow meter"),
	rl.PcItem("show openflow queues"),
	rl.PcItem("show openflow resources"),
	rl.PcItem("show overlapping-vlan-resource usage"),
	rl.PcItem("show overlay-gateway"),
	rl.PcItem("show ovsdb-server"),
	rl.PcItem("show policymap"),
	rl.PcItem("show port port-channel"),
	rl.PcItem("show port-channel"),
	rl.PcItem("show port-channel-redundancy-group"),
	rl.PcItem("show port-profile"),
	rl.PcItem("show port-profile domain"),
	rl.PcItem("show port-profile interface"),
	rl.PcItem("show port-profile name"),
	rl.PcItem("show port-security"),
	rl.PcItem("show port-security addresses"),
	rl.PcItem("show port-security interface"),
	rl.PcItem("show port-security oui interface"),
	rl.PcItem("show port-security sticky interface"),
	rl.PcItem("show process cpu"),
	rl.PcItem("show process info"),
	rl.PcItem("show process memory"),
	rl.PcItem("show prom-access"),
	rl.PcItem("show qos flowcontrol interface"),
	rl.PcItem("show qos interface"),
	rl.PcItem("show qos maps"),
	rl.PcItem("show qos maps dscp-cos"),
	rl.PcItem("show qos maps dscp-mutation"),
	rl.PcItem("show qos maps dscp-traffic-class"),
	rl.PcItem("show qos queue interface"),
	rl.PcItem("show qos rcv-queue interface"),
	rl.PcItem("show qos rcv-queue multicast"),
	rl.PcItem("show qos red profiles"),
	rl.PcItem("show qos red statistics interface"),
	rl.PcItem("show qos tx-queue interface"),
	rl.PcItem("show rbridge-id"),
	rl.PcItem("show rbridge-running config"),
	rl.PcItem("show rbridge-local-running-config"),
	rl.PcItem("show redundancy"),
	rl.PcItem("show rmon"),
	rl.PcItem("show rmon history"),
	rl.PcItem("show route-map"),
	rl.PcItem("show route-map interface"),
	rl.PcItem("show running reserved-vlan"),
	rl.PcItem("show running-config"),
	rl.PcItem("show running-config aaa"),
	rl.PcItem("show running-config aaa accounting"),
	rl.PcItem("show running-config ag"),
	rl.PcItem("show running-config banner"),
	rl.PcItem("show running-config cee-map"),
	rl.PcItem("show running-config class-map"),
	rl.PcItem("show running-config diag post"),
	rl.PcItem("show running-config dot1x"),
	rl.PcItem("show running-config dpod"),
	rl.PcItem("show running-config event-handler"),
	rl.PcItem("show running-config fabric route mcast"),
	rl.PcItem("show running-config fcoe"),
	rl.PcItem("show running-config fcsp auth"),
	rl.PcItem("show running-config hardware"),
	rl.PcItem("show running-config interface fcoe"),
	rl.PcItem("show running-config interface FibreChannel"),
	rl.PcItem("show running-config interface fortygigabitethernet"),
	rl.PcItem("show running-config interface fortygigabitethernet cee"),
	rl.PcItem("show running-config interface fortygigabitethernet channel-group"),
	rl.PcItem("show running-config interface fortygigabitethernet description"),
	rl.PcItem("show running-config interface fortygigabitethernet dot1x"),
	rl.PcItem("show running-config interface fortygigabitethernet fabric"),
	rl.PcItem("show running-config interface fortygigabitethernet fcoeport"),
	rl.PcItem("show running-config interface fortygigabitethernet lacp"),
	rl.PcItem("show running-config interface fortygigabitethernet lldp"),
	rl.PcItem("show running-config interface fortygigabitethernet mac"),
	rl.PcItem("show running-config interface fortygigabitethernet mtu"),
	rl.PcItem("show running-config interface fortygigabitethernet qos"),
	rl.PcItem("show running-config interface fortygigabitethernet rmon"),
	rl.PcItem("show running-config interface fortygigabitethernet sflow"),
	rl.PcItem("show running-config interface fortygigabitethernet shutdown"),
	rl.PcItem("show running-config interface fortygigabitethernet switchport"),
	rl.PcItem("show running-config interface fortygigabitethernet udld"),
	rl.PcItem("show running-config interface fortygigabitethernet vlan"),
	rl.PcItem("show running-config interface gigabitethernet"),
	rl.PcItem("show running-config interface gigabitethernet bpdu-drop"),
	rl.PcItem("show running-config interface gigabitethernet description"),
	rl.PcItem("show running-config interface gigabitethernet dot1x"),
	rl.PcItem("show running-config interface gigabitethernet lacp"),
	rl.PcItem("show running-config interface gigabitethernet lldp"),
	rl.PcItem("show running-config interface gigabitethernet mac"),
	rl.PcItem("show running-config interface gigabitethernet mtu"),
	rl.PcItem("show running-config interface gigabitethernet port-profile-port"),
	rl.PcItem("show running-config interface gigabitethernet priority-tag"),
	rl.PcItem("show running-config interface gigabitethernet qos"),
	rl.PcItem("show running-config interface gigabitethernet rmon"),
	rl.PcItem("show running-config interface gigabitethernet sflow"),
	rl.PcItem("show running-config interface gigabitethernet shutdown"),
	rl.PcItem("show running-config interface gigabitethernet switchport"),
	rl.PcItem("show running-config interface gigabitethernet udld"),
	rl.PcItem("show running-config interface gigabitethernet vlan"),
	rl.PcItem("show running-config interface management"),
	rl.PcItem("show running-config interface port-channel"),
	rl.PcItem("show running-config interface tengigabitethernet"),
	rl.PcItem("show running-config interface tengigabitethernet bpdu-drop"),
	rl.PcItem("show running-config interface tengigabitethernet cee"),
	rl.PcItem("show running-config interface tengigabitethernet description"),
	rl.PcItem("show running-config interface tengigabitethernet dot1x"),
	rl.PcItem("show running-config interface tengigabitethernet fabric"),
	rl.PcItem("show running-config interface tengigabitethernet fcoeport"),
	rl.PcItem("show running-config interface tengigabitethernet lacp"),
	rl.PcItem("show running-config interface tengigabitethernet lldp"),
	rl.PcItem("show running-config interface tengigabitethernet mac"),
	rl.PcItem("show running-config interface tengigabitethernet mtu"),
	rl.PcItem("show running-config interface tengigabitethernet qos"),
	rl.PcItem("show running-config interface tengigabitethernet rmon"),
	rl.PcItem("show running-config interface tengigabitethernet sflow"),
	rl.PcItem("show running-config interface tengigabitethernet shutdown"),
	rl.PcItem("show running-config interface tengigabitethernet udld"),
	rl.PcItem("show running-config interface tengigabitethernet vlan"),
	rl.PcItem("show running-config interface vlan"),
	rl.PcItem("show running-config interface vlan ip"),
	rl.PcItem("show running-config ip access-list"),
	rl.PcItem("show running-config ipv6 access-list"),
	rl.PcItem("show running-config ip dns"),
	rl.PcItem("show running-config ip igmp"),
	rl.PcItem("show running-config ip route"),
	rl.PcItem("show running-config ldap-server"),
	rl.PcItem("show running-config line"),
	rl.PcItem("show running-config logging"),
	rl.PcItem("show running-config logging auditlog class"),
	rl.PcItem("show running-config logging raslog"),
	rl.PcItem("show running-config logging syslog-client"),
	rl.PcItem("show running-config logging syslog-facility"),
	rl.PcItem("show running-config logging syslog-server"),
	rl.PcItem("show running-config mac access-list"),
	rl.PcItem("show running-config mac-address-table"),
	rl.PcItem("show running-config monitor"),
	rl.PcItem("show running-config nas server-ip"),
	rl.PcItem("show running-config ntp"),
	rl.PcItem("show running-config ntp authentication-key"),
	rl.PcItem("show running-config openflow-controller"),
	rl.PcItem("show running-config overlay-gateway"),
	rl.PcItem("show running-config ovsdb-server"),
	rl.PcItem("show running-config password-attributes"),
	rl.PcItem("show running-config police-priority-map"),
	rl.PcItem("show running-config policy-map"),
	rl.PcItem("show running-config port-profile"),
	rl.PcItem("show running-config port-profile activate"),
	rl.PcItem("show running-config port-profile fcoe-profile"),
	rl.PcItem("show running-config port-profile qos-profile"),
	rl.PcItem("show running-config port-profile security-profile"),
	rl.PcItem("show running-config port-profile static"),
	rl.PcItem("show running-config port-profile vlan-profile"),
	rl.PcItem("show running-config port-profile-domain"),
	rl.PcItem("show running-config preprovision"),
	rl.PcItem("show running-config protocol cdp"),
	rl.PcItem("show running-config protocol edge"),
	rl.PcItem("show running-config protocol lldp"),
	rl.PcItem("show running-config protocol spanning-tree mstp"),
	rl.PcItem("show running-config protocol spanning-tree pvst"),
	rl.PcItem("show running-config protocol spanning-tree rpvst"),
	rl.PcItem("show running-config protocol spanning-tree rstp"),
	rl.PcItem("show running-config protocol spanning-tree stp"),
	rl.PcItem("show running-config protocol udld"),
	rl.PcItem("show running-config radius-server"),
	rl.PcItem("show running-config rbridge-id"),
	rl.PcItem("show running-config rbridge-id crypto"),
	rl.PcItem("show running-config rbridge-id event-handler"),
	rl.PcItem("show running-config rbridge-id hardware-profile"),
	rl.PcItem("show running-config rbridge-id interface"),
	rl.PcItem("show running-config rbridge-id linecard"),
	rl.PcItem("show running-config rbridge-id maps"),
	rl.PcItem("show running-config rbridge-id openflow"),
	rl.PcItem("show running-config rbridge-id ssh"),
	rl.PcItem("show running-config rmon"),
	rl.PcItem("show running-config role"),
	rl.PcItem("show running-config route-map"),
	rl.PcItem("show running-config rule"),
	rl.PcItem("show running-config secpolicy"),
	rl.PcItem("show running-config sflow"),
	rl.PcItem("show running-config sflow-policy"),
	rl.PcItem("show running-config sflow-profile"),
	rl.PcItem("show running-config snmp-server"),
	rl.PcItem("show running-config snmp-server context"),
	rl.PcItem("show running-config snmp-server engineid"),
	rl.PcItem("show running-config snmp-server mib community-map"),
	rl.PcItem("show running-config ssh"),
	rl.PcItem("show running-config ssh server"),
	rl.PcItem("show running-config ssh server key-exchange"),
	rl.PcItem("show running-config support"),
	rl.PcItem("show running-config support autoupload-param"),
	rl.PcItem("show running-config support support-param"),
	rl.PcItem("show running-config switch-attributes"),
	rl.PcItem("show running-config system-monitor"),
	rl.PcItem("show running-config system-monitor-mail"),
	rl.PcItem("show running-config tacacs-server"),
	rl.PcItem("show running-config telnet server"),
	rl.PcItem("show running-config threshold-monitor"),
	rl.PcItem("show running-config threshold-monitor interface"),
	rl.PcItem("show running-config threshold-monitor security"),
	rl.PcItem("show running-config threshold-monitor sfp"),
	rl.PcItem("show running-config username"),
	rl.PcItem("show running-config vcs"),
	rl.PcItem("show running-config vlag-commit-mode"),
	rl.PcItem("show running-config zoning"),
	rl.PcItem("show running-config zoning defined-configuration"),
	rl.PcItem("show running-config zoning enabled-configuration"),
	rl.PcItem("show secpolicy"),
	rl.PcItem("show sflow"),
	rl.PcItem("show sflow-profile"),
	rl.PcItem("show sfm"),
	rl.PcItem("show sfp"),
	rl.PcItem("show slots"),
	rl.PcItem("show span path"),
	rl.PcItem("show spanning-tree"),
	rl.PcItem("show spanning-tree brief"),
	rl.PcItem("show spanning-tree interface"),
	rl.PcItem("show spanning-tree mst brief"),
	rl.PcItem("show spanning-tree mst detail"),
	rl.PcItem("show spanning-tree mst instance"),
	rl.PcItem("show spanning-tree mst interface"),
	rl.PcItem("show ssh server status"),
	rl.PcItem("show ssh server rekey-interval status"),
	rl.PcItem("show startup-config"),
	rl.PcItem("show startup-db"),
	rl.PcItem("show statistics access-list"),
	rl.PcItem("show storm-control"),
	rl.PcItem("show support"),
	rl.PcItem("show system"),
	rl.PcItem("show system internal asic counter blk"),
	rl.PcItem("show system internal asic counter drop-reason"),
	rl.PcItem("show system internal asic counter interface"),
	rl.PcItem("show system internal asic counter mem blk"),
	rl.PcItem("show system internal bgp evpn interface"),
	rl.PcItem("show system internal bgp evpn l2route type"),
	rl.PcItem("show system internal bgp evpn l3vni"),
	rl.PcItem("show system internal bgp evpn neighbor"),
	rl.PcItem("show system internal bgp evpn routes type"),
	rl.PcItem("show system internal bgp evpn variables"),
	rl.PcItem("show system internal bgp evpn vlan-db"),
	rl.PcItem("show system internal bgp ipv4 config"),
	rl.PcItem("show system internal bgp ipv4 neighbor"),
	rl.PcItem("show system internal bgp ipv4 network"),
	rl.PcItem("show system internal bgp ipv4 nexthop"),
	rl.PcItem("show system internal bgp ipv4 tcpdump"),
	rl.PcItem("show system internal bgp ipv4 variables"),
	rl.PcItem("show system internal bgp ipv6 neighbor"),
	rl.PcItem("show system internal bgp ipv6 network"),
	rl.PcItem("show system internal bgp ipv6 nexthop"),
	rl.PcItem("show system internal bgp ipv6 tcpdump"),
	rl.PcItem("show system internal bgp ipv6 variables"),
	rl.PcItem("show system internal dcm"),
	rl.PcItem("show system internal nas"),
	rl.PcItem("show system internal nsm"),
	rl.PcItem("show system internal nsx"),
	rl.PcItem("show system internal ofagt"),
	rl.PcItem("show system internal ovsdb"),
	rl.PcItem("show system monitor"),
	rl.PcItem("show system pstat interface"),
	rl.PcItem("show telnet server status"),
	rl.PcItem("show threshold monitor"),
	rl.PcItem("show track summary"),
	rl.PcItem("show tunnel"),
	rl.PcItem("show tunnel nsx service-node"),
	rl.PcItem("show tunnel status"),
	rl.PcItem("show udld"),
	rl.PcItem("show udld interface"),
	rl.PcItem("show udld statistics"),
	rl.PcItem("show users"),
	rl.PcItem("show vcs"),
	rl.PcItem("show version"),
	rl.PcItem("show virtual-fabric status"),
	rl.PcItem("show vlag-partner-info"),
	rl.PcItem("show vlan"),
	rl.PcItem("show vlan brief"),
	rl.PcItem("show vlan classifier"),
	rl.PcItem("show vlan private-vlan"),
	rl.PcItem("show vlan rspan-vlan"),
	rl.PcItem("show vnetwork"),
	rl.PcItem("show vrf"),
	rl.PcItem("show vrrp"),
	rl.PcItem("show zoning enabled-configuration"),
	rl.PcItem("show zoning operation-info"),
)