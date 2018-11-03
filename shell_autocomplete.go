package main

import rl "github.com/chzyer/readline"

/* defaultGetCompletion
is a parameter list of the default parameters for the get command. This
list will be added to every completer */
var defaultGetCompletion = rl.PcItem("ls",
	rl.PcItem("filter"),
	rl.PcItem("hosts"),
	rl.PcItem("selhosts"),
	rl.PcItem("allhosts"),
)

/* defaultSetCompletion
is a parameter list of the default parameters for the get command. This
list will be added to every completer */
var defaultSetCompletion = rl.PcItem("mset",
	rl.PcItem("filter"),
	rl.PcItem("complete",
		rl.PcItem("junos"),
		rl.PcItem("netiron"),
		rl.PcItem("vdx")))

/* cliNetironCompleter
is an autocompletion tree for the Netiron command line
*/
var cliNetironCompleter = rl.NewPrefixCompleter(
	rl.PcItem("clear access-list receive accounting"),
	rl.PcItem("clear arp-guard-statistics"),
	rl.PcItem("clear arp vrf"),
	rl.PcItem("clear arp-inspection-statistics"),
	rl.PcItem("clear rate-limit arp"),
	rl.PcItem("clear bm histogram"),
	rl.PcItem("clear cpu histogram sequence"),
	rl.PcItem("clear dot1x-mka statistics"),
	rl.PcItem("clear ikev2 statistics"),
	rl.PcItem("clear ikev2 sa"),
	rl.PcItem("clear ip bgp dampening"),
	rl.PcItem("clear ip bgp flag-statistics"),
	rl.PcItem("clear ip bgp local routes"),
	rl.PcItem("clear ip bgp neighbor"),
	rl.PcItem("clear ip bgp routes"),
	rl.PcItem("clear ip bgp traffic"),
	rl.PcItem("clear ip bgp vrf"),
	rl.PcItem("clear ip ospf"),
	rl.PcItem("clear ip rip local routes"),
	rl.PcItem("clear ip rip routes"),
	rl.PcItem("clear ip vrrp statistics"),
	rl.PcItem("clear ip vrrp-extended statistics"),
	rl.PcItem("clear ipsec error-count"),
	rl.PcItem("clear ipsec sa"),
	rl.PcItem("clear ipsec statistics"),
	rl.PcItem("clear ipsec statistics tunnel"),
	rl.PcItem("clear ipv6 bgp dampening"),
	rl.PcItem("clear ipv6 bgp flap-statistics"),
	rl.PcItem("clear ipv6 bgp local routes"),
	rl.PcItem("clear ipv6 bgp neighbor"),
	rl.PcItem("clear ipv6 bgp routes"),
	rl.PcItem("clear ipv6 bgp traffic"),
	rl.PcItem("clear ipv6 ospf"),
	rl.PcItem("clear ipv6 rip route"),
	rl.PcItem("clear ipv6 vrrp statistics"),
	rl.PcItem("clear ipv6 vrrp-extended statistics"),
	rl.PcItem("clear isis shortcut"),
	rl.PcItem("clear mac-address vpls"),
	rl.PcItem("clear macsec statistics"),
	rl.PcItem("clear memory histogram"),
	rl.PcItem("clear metro mp-vlp-queue"),
	rl.PcItem("clear mmrp statistics"),
	rl.PcItem("clear mpls auto-bandwidth-samples"),
	rl.PcItem("clear mpls ldp neighbor"),
	rl.PcItem("clear mpls ldp statistics"),
	rl.PcItem("clear mpls rsvp statistics session"),
	rl.PcItem("clear mpls statistics"),
	rl.PcItem("clear mvrp statistics"),
	rl.PcItem("clear openflow"),
	rl.PcItem("clear pki counters"),
	rl.PcItem("clear pki crl"),
	rl.PcItem("clear rate-limit arp"),
	rl.PcItem("clear rate-limit counters bum-drop"),
	rl.PcItem("clear rate-limit counters ip-option-pkt-to-cpu"),
	rl.PcItem("clear rate-limit counters ipv6-hoplimit-expired-to-cpu"),
	rl.PcItem("clear rate-limit counters ip-ttl-expired-to-cpu"),
	rl.PcItem("clear statistics openflow"),
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
		rl.PcItem("configuration"),
		rl.PcItem("cpu histogram"),
		rl.PcItem("interface ethernet"),
		rl.PcItem("interfaces tunnel"),
		rl.PcItem("ip",
			rl.PcItem("bgp", rl.PcItem("attribute-entries"),
				rl.PcItem("config"),
				rl.PcItem("dampened-paths"),
				rl.PcItem("filtered-routes"),
				rl.PcItem("flap-statistics"),
				rl.PcItem("ipv6"),
				rl.PcItem("neighbors"),
				rl.PcItem("neighbors advertised-routes"),
				rl.PcItem("neighbors flap-statistics"),
				rl.PcItem("neighbors last-packet-with-error"),
				rl.PcItem("neighbors received"),
				rl.PcItem("neighbors received-routes"),
				rl.PcItem("neighbors rib-out-routes"),
				rl.PcItem("routes community"),
				rl.PcItem("neighbors routes"),
				rl.PcItem("neighbors routes-summary"),
				rl.PcItem("peer-group"),
				rl.PcItem("routes"),
				rl.PcItem("summary"),
				rl.PcItem("vrf neighbors"),
				rl.PcItem("vrf routes"),
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
			rl.PcItem("access-list bindings"),
			rl.PcItem("access-list receive accounting"),
			rl.PcItem("bgp"),
			rl.PcItem("bgp neighbors"),
			rl.PcItem("bgp routes"),
			rl.PcItem("bgp summary"),
			rl.PcItem("dhcp-relay interface"),
			rl.PcItem("dhcp-relay options"),
			rl.PcItem("interface tunnel"),
			rl.PcItem("ospf interface"),
			rl.PcItem("vrrp"),
			rl.PcItem("vrrp-extended")),
		rl.PcItem("isis"),
		rl.PcItem("license"),
		rl.PcItem("log"),
		rl.PcItem("module"),
		rl.PcItem("mpls",
			rl.PcItem("autobw-threshold-table"),
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
			rl.PcItem("ldp session"),
			rl.PcItem("ldp statistics"),
			rl.PcItem("ldp tunnel"),
			rl.PcItem("lsp"),
			rl.PcItem("lsp_pmp_xc"),
			rl.PcItem("path"),
			rl.PcItem("policy"),
			rl.PcItem("route"),
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
				rl.PcItem("ldp tunnel"),
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
		rl.PcItem("sflow statistics"),
		rl.PcItem("spanning-tree"),
		rl.PcItem("statistics"),
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
	rl.PcItem("clear ag nport-utilization"),
	rl.PcItem("clear arp"),
	rl.PcItem("clear bfd counters"),
	rl.PcItem("clear bgp evpn l2routes"),
	rl.PcItem("clear bgp evpn local routes"),
	rl.PcItem("clear bgp evpn mac-route dampening"),
	rl.PcItem("clear bgp evpn neighbor"),
	rl.PcItem("clear bgp evpn routes"),
	rl.PcItem("clear counters"),
	rl.PcItem("clear counters (IP)"),
	rl.PcItem("clear counters (MAC)"),
	rl.PcItem("clear counters access-list"),
	rl.PcItem("clear counters interface"),
	rl.PcItem("clear counters slot-id"),
	rl.PcItem("clear counters storm-control"),
	rl.PcItem("clear dot1x statistics"),
	rl.PcItem("clear dot1x statistics interface"),
	rl.PcItem("clear edge-loop-detection"),
	rl.PcItem("clear ip arp inspection statistics"),
	rl.PcItem("clear ip arp suppression-cache"),
	rl.PcItem("clear ip arp suppression-statistics"),
	rl.PcItem("clear ip bgp dampening"),
	rl.PcItem("clear ip bgp flap-statistics"),
	rl.PcItem("clear ip bgp local routes"),
	rl.PcItem("clear ip bgp neighbor"),
	rl.PcItem("clear ip bgp routes"),
	rl.PcItem("clear ip bgp traffic"),
	rl.PcItem("clear ip dhcp relay statistics"),
	rl.PcItem("clear ip fabric-virtual-gateway"),
	rl.PcItem("clear ip igmp groups"),
	rl.PcItem("clear ip igmp statistics interface"),
	rl.PcItem("clear ip ospf"),
	rl.PcItem("clear ip pim mcache"),
	rl.PcItem("clear ip pim rp-map"),
	rl.PcItem("clear ip pim traffic"),
	rl.PcItem("clear ip route"),
	rl.PcItem("clear ipv6 bgp dampening"),
	rl.PcItem("clear ipv6 bgp flap-statistics"),
	rl.PcItem("clear ipv6 bgp local routes"),
	rl.PcItem("clear ipv6 bgp neighbor"),
	rl.PcItem("clear ipv6 bgp routes"),
	rl.PcItem("clear ipv6 bgp traffic"),
	rl.PcItem("clear ipv6 counters"),
	rl.PcItem("clear ipv6 dhcp relay statistics"),
	rl.PcItem("clear ipv6 fabric-virtual-gateway"),
	rl.PcItem("clear ipv6 mld groups"),
	rl.PcItem("clear ipv6 mld statistics"),
	rl.PcItem("clear ipv6 nd suppression-cache"),
	rl.PcItem("clear ipv6 nd suppression-statistics"),
	rl.PcItem("clear ipv6 neighbor"),
	rl.PcItem("clear ipv6 ospf"),
	rl.PcItem("clear ipv6 route"),
	rl.PcItem("clear ipv6 vrrp statistics"),
	rl.PcItem("clear lacp"),
	rl.PcItem("clear lacp counters"),
	rl.PcItem("clear lldp neighbors"),
	rl.PcItem("clear lldp statistics"),
	rl.PcItem("clear logging auditlog"),
	rl.PcItem("clear logging raslog"),
	rl.PcItem("clear mac-address-table conversational"),
	rl.PcItem("clear mac-address-table dynamic"),
	rl.PcItem("clear maps dashboard"),
	rl.PcItem("clear nas statistics"),
	rl.PcItem("clear openflow"),
	rl.PcItem("clear overlay-gateway"),
	rl.PcItem("clear policy-map-counters"),
	rl.PcItem("clear sessions"),
	rl.PcItem("clear sflow statistics"),
	rl.PcItem("clear spanning-tree counter"),
	rl.PcItem("clear spanning-tree detected-protocols"),
	rl.PcItem("clear statistics openflow"),
	rl.PcItem("clear support"),
	rl.PcItem("clear udld statistics"),
	rl.PcItem("clear vrrp statistics"),
	rl.PcItem("show access-list accounting"),
	rl.PcItem("show access-list bindings"),
	rl.PcItem("show access-list receive accounting"),
	rl.PcItem("show arp"),
	rl.PcItem("show arp-guard-access-list"),
	rl.PcItem("show arp-guard port-bindings"),
	rl.PcItem("show arp-guard statistics ethernet"),
	rl.PcItem("show bfd"),
	rl.PcItem("show bfd applications"),
	rl.PcItem("show bfd mpls"),
	rl.PcItem("show bfd neighbors"),
	rl.PcItem("show bfd neighbors bgp"),
	rl.PcItem("show bfd neighbors details"),
	rl.PcItem("show bfd neighbors interface"),
	rl.PcItem("show bfd neighbors isis"),
	rl.PcItem("show bfd neighbors ospf"),
	rl.PcItem("show bfd neighbors ospf6"),
	rl.PcItem("show bfd neighbors static"),
	rl.PcItem("show bfd neighbors static6"),
	rl.PcItem("show bip slot"),
	rl.PcItem("show cam-detail-eth"),
	rl.PcItem("show cam-detail-ip"),
	rl.PcItem("show cam ifl"),
	rl.PcItem("show cam ipvpn"),
	rl.PcItem("show cam uda"),
	rl.PcItem("show configuration"),
	rl.PcItem("show cpu histogram"),
	rl.PcItem("show cpu histogram sequence"),
	rl.PcItem("show dot1x-mka group"),
	rl.PcItem("show dot1x-mka config"),
	rl.PcItem("show dot1x-mka sessions brief"),
	rl.PcItem("show dot1x-mka sessions ethernet"),
	rl.PcItem("show dot1x-mka statistics"),
	rl.PcItem("show egress-truncate"),
	rl.PcItem("show ikev2 policy"),
	rl.PcItem("show ikev2 profile"),
	rl.PcItem("show ikev2 proposal"),
	rl.PcItem("show ikev2 sa"),
	rl.PcItem("show ikev2 session"),
	rl.PcItem("show ikev2 statistics"),
	rl.PcItem("show interface ethernet"),
	rl.PcItem("show interfaces tunnel"),
	rl.PcItem("show ip allow-src-multicast"),
	rl.PcItem("show ip bgp neighbors"),
	rl.PcItem("show ip bgp summary"),
	rl.PcItem("show ip http client"),
	rl.PcItem("show ip interface"),
	rl.PcItem("show ip ospf"),
	rl.PcItem("show ip route"),
	rl.PcItem("show ip static-arp"),
	rl.PcItem("show ip vrrp"),
	rl.PcItem("show ip vrrp-extended"),
	rl.PcItem("show ipsec egress-config"),
	rl.PcItem("show ipsec egress-spi-table"),
	rl.PcItem("show ipsec error-count"),
	rl.PcItem("show ipsec ingress-config"),
	rl.PcItem("show ipsec ingress-spi-table"),
	rl.PcItem("show ipsec policy"),
	rl.PcItem("show ipsec profile"),
	rl.PcItem("show ipsec proposal"),
	rl.PcItem("show ipsec sa"),
	rl.PcItem("show ipsec statistics"),
	rl.PcItem("show ip-tunnels"),
	rl.PcItem("show ipv6 access-list bindings"),
	rl.PcItem("show ipv6 access-list receive accounting"),
	rl.PcItem("show ipv6 bgp neighbors"),
	rl.PcItem("show ipv6 bgp summary"),
	rl.PcItem("show ipv6 dhcp-relay interface"),
	rl.PcItem("show ipv6 dhcp-relay options"),
	rl.PcItem("show ipv6 interface tunnel"),
	rl.PcItem("show ipv6 ospf interface"),
	rl.PcItem("show ipv6 vrrp"),
	rl.PcItem("show ipv6 vrrp-extended"),
	rl.PcItem("show isis"),
	rl.PcItem("show isis shortcut"),
	rl.PcItem("show macsec ethernet"),
	rl.PcItem("show macsec statistics ethernet"),
	rl.PcItem("show memory histogram"),
	rl.PcItem("show metro mp-vlp-queue"),
	rl.PcItem("show mmrp"),
	rl.PcItem("show mmrp attributes"),
	rl.PcItem("show mmrp config"),
	rl.PcItem("show mmrp statistics"),
	rl.PcItem("show mpls autobw-threshold-table"),
	rl.PcItem("show mpls bypass-lsp"),
	rl.PcItem("show mpls config"),
	rl.PcItem("show mpls forwarding"),
	rl.PcItem("show mpls interface"),
	rl.PcItem("show mpls label-range"),
	rl.PcItem("show mpls ldp"),
	rl.PcItem("show mpls ldp database"),
	rl.PcItem("show mpls ldp fec"),
	rl.PcItem("show mpls ldp interface"),
	rl.PcItem("show mpls ldp neighbor"),
	rl.PcItem("show mpls ldp path"),
	rl.PcItem("show mpls ldp peer"),
	rl.PcItem("show mpls ldp session"),
	rl.PcItem("show mpls ldp statistics"),
	rl.PcItem("show mpls ldp tunnel"),
	rl.PcItem("show mpls lsp"),
	rl.PcItem("show mpls lsp_p2mp_xc"),
	rl.PcItem("show mpls path"),
	rl.PcItem("show mpls policy"),
	rl.PcItem("show mpls route"),
	rl.PcItem("show mpls rsvp interface"),
	rl.PcItem("show mpls rsvp neighbor"),
	rl.PcItem("show mpls rsvp session"),
	rl.PcItem("show mpls rsvp session backup"),
	rl.PcItem("show mpls rsvp session brief"),
	rl.PcItem("show mpls rsvp session bypass"),
	rl.PcItem("show mpls rsvp session destination"),
	rl.PcItem("show mpls rsvp session detail"),
	rl.PcItem("show mpls rsvp session detour"),
	rl.PcItem("show mpls rsvp session down"),
	rl.PcItem("show mpls rsvp session extensive"),
	rl.PcItem("show mpls rsvp session"),
	rl.PcItem("show mpls rsvp session"),
	rl.PcItem("show mpls rsvp session name"),
	rl.PcItem("show mpls rsvp session p2mp"),
	rl.PcItem("show mpls rsvp session p2p"),
	rl.PcItem("show mpls rsvp session ppend"),
	rl.PcItem("show mpls rsvp session transit"),
	rl.PcItem("show mpls rsvp session up"),
	rl.PcItem("show mpls rsvp session wide"),
	rl.PcItem("show mpls rsvp statistics"),
	rl.PcItem("show mpls static-lsp"),
	rl.PcItem("show mpls statistics 6pe"),
	rl.PcItem("show mpls statistics bypass-lsp"),
	rl.PcItem("show mpls statistics label"),
	rl.PcItem("show mpls statistics ldp transit"),
	rl.PcItem("show mpls statistics ldp tunnel"),
	rl.PcItem("show mpls statistics lsp"),
	rl.PcItem("show mpls statistics oam"),
	rl.PcItem("show mpls statistics vll"),
	rl.PcItem("show mpls statistics vll-local"),
	rl.PcItem("show mpls statistics vpls"),
	rl.PcItem("show mpls statistics vrf"),
	rl.PcItem("show mpls summary"),
	rl.PcItem("show mpls ted database"),
	rl.PcItem("show mpls ted path"),
	rl.PcItem("show mpls vll"),
	rl.PcItem("show mpls vll-local"),
	rl.PcItem("show mpls vpls"),
	rl.PcItem("show mstp"),
	rl.PcItem("show mvrp"),
	rl.PcItem("show mvrp attributes"),
	rl.PcItem("show mvrp config"),
	rl.PcItem("show mvrp statistics"),
	rl.PcItem("show nht-table ipsec-based"),
	rl.PcItem("show openflow"),
	rl.PcItem("show openflow controller"),
	rl.PcItem("show openflow flows"),
	rl.PcItem("show openflow groups"),
	rl.PcItem("show openflow interface"),
	rl.PcItem("show openflow meters"),
	rl.PcItem("show openflow queues"),
	rl.PcItem("show pim interface"),
	rl.PcItem("show pim multicast-filter"),
	rl.PcItem("show pki certificates"),
	rl.PcItem("show pki counters"),
	rl.PcItem("show pki crls"),
	rl.PcItem("show pki enrollment-profile"),
	rl.PcItem("show pki entity"),
	rl.PcItem("show pki key mypubkey"),
	rl.PcItem("show pki trustpoint"),
	rl.PcItem("show rate-limit counters bum-drop"),
	rl.PcItem("show rate-limit detail"),
	rl.PcItem("show rate-limit interface"),
	rl.PcItem("show rate-limit ipv6 hoplimit-expired-to-cpu"),
	rl.PcItem("show rate-limit option-pkt-to-cpu"),
	rl.PcItem("show rate-limit ttl-expired-to-cpu"),
	rl.PcItem("show rmon alarm"),
	rl.PcItem("show rmon statistics"),
	rl.PcItem("show route-map"),
	rl.PcItem("show rstp"),
	rl.PcItem("show running-config"),
	rl.PcItem("show sflow statistics"),
	rl.PcItem("show spanning-tree"),
	rl.PcItem("show statistics"),
	rl.PcItem("show sysmon config"),
	rl.PcItem("show sysmon results brief"),
	rl.PcItem("show sysmon results detail"),
	rl.PcItem("show sysmon schedule"),
	rl.PcItem("show telemetry"),
	rl.PcItem("show terminal"),
	rl.PcItem("show tm-voq-stat queue-drops"),
	rl.PcItem("show vlan"),
	rl.PcItem("show vlan tvf-lag-lb"),
	defaultGetCompletion,
	defaultSetCompletion,
)