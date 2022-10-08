package lexer

import (
	"fmt"
	"testing"
)

var (
	testVrpPortConfigExample string = `
#
interface GigabitEthernet0/0/1
 description uprn-200001093619-01#1#
 undo enable snmp trap updown
 port link-type hybrid
 port hybrid pvid vlan 902
 undo port hybrid vlan 1
 port hybrid tagged vlan 296
 port hybrid untagged vlan 902
 loopback-detect recovery-time 255
 loopback-detect packet vlan 902
 loopback-detect enable
 loopback-detect action quitvlan
 stp disable
 stp bpdu-filter enable
 traffic-secure inbound acl name NO_V6
 traffic-policy KBPS_1000000 inbound
 bpdu disable
 undo lldp enable
 qos schedule-profile PQ
 arp anti-attack check user-bind enable
 ip source check user-bind enable
 user-bind ip sticky-mac
 multicast-suppression packets 5
 broadcast-suppression packets 5
 unicast-suppression block outbound
 dhcp snooping enable
 dhcp snooping check dhcp-request enable
 dhcp snooping check dhcp-chaddr enable
 dhcp snooping check dhcp-rate enable
 dhcp snooping check dhcp-rate 5
 dhcpv6 option37 rebuild enable
 dhcpv6 option18 rebuild enable
 dhcpv6 option37 format user-defined "uprn-200001093619-01#1"
 nd snooping enable dhcpv6 only
 nd snooping check na enable
 nd snooping check rs enable
 dhcp option82 rebuild enable
 dhcp option82 remote-id format user-defined "uprn-200001093619-01#1"
 #
return

#
interface GigabitEthernet0/0/1
 description uprn-200001093619-01#1#
 undo enable snmp trap updown
 port link-type hybrid
 port hybrid pvid vlan 902
 undo port hybrid vlan 1
 port hybrid tagged vlan 296
 port hybrid untagged vlan 902
 loopback-detect recovery-time 255
 loopback-detect packet vlan 902
 loopback-detect enable
 loopback-detect action quitvlan
 stp disable
 stp bpdu-filter enable
 traffic-secure inbound acl name NO_V6
 traffic-policy KBPS_1000000 inbound
 bpdu disable
 undo lldp enable
 qos schedule-profile PQ
 arp anti-attack check user-bind enable
 ip source check user-bind enable
 user-bind ip sticky-mac
 multicast-suppression packets 5
 broadcast-suppression packets 5
 unicast-suppression block outbound
 dhcp snooping enable
 dhcp snooping check dhcp-request enable
 dhcp snooping check dhcp-chaddr enable
 dhcp snooping check dhcp-rate enable
 dhcp snooping check dhcp-rate 5
 dhcpv6 option37 rebuild enable
 dhcpv6 option18 rebuild enable
 dhcpv6 option37 format user-defined "uprn-second-one-01#1"
 nd snooping enable dhcpv6 only
 nd snooping check na enable
 nd snooping check rs enable
 dhcp option82 rebuild enable
 dhcp option82 remote-id format user-defined "uprn-second-one-01#1"
 #
return`
)

var testVrpPortConfigExampleShort string = `
#
interface GigabitEthernet0/0/1
 description uprn-200001093619-01#1#
 dhcp option82 remote-id format user-defined "uprn-200001093619-01#1"
#
 `

var testVrpMainConfig string = `
#
!Software Version V200R013C00SPC500
#
sysname sait-hea-a2
#
FTP acl 2000
#
info-center filter-id bymodule-alias SHELL DISPLAY_CMDRECORD
info-center filter-id bymodule-alias SRM TXPOWER_EXCEEDMAJOR
info-center filter-id bymodule-alias SRM RXPOWER_EXCEEDMINOR
info-center filter-id bymodule-alias SRM RXPOWER_RESUME
info-center filter-id bymodule-alias SRM OPTPWRABNORMAL
info-center filter-id bymodule-alias SRM OPTPWRRESUME
info-center filter-id bymodule-alias SECE DAI_DROP_PACKET
info-center filter-id bymodule-alias IFPDT IF_STATE
info-center filter-id bymodule-alias SW_SNPG IGMPV2_PKT
info-center loghost source Vlanif55
info-center timestamp debugging date precision-time tenth-second
info-center timestamp log date precision-time millisecond
#
file prompt quiet
#
dns resolve
dns server 192.168.50.1
#
vlan batch 55 296 700 800 999
#
mac-address aging-time 3600
#
stp bpdu-filter default
#
authentication-profile name default_authen_profile
authentication-profile name dot1x_authen_profile
authentication-profile name dot1xmac_authen_profile
authentication-profile name mac_authen_profile
authentication-profile name multi_authen_profile
authentication-profile name portal_authen_profile
#`

func Test_lex(t *testing.T) {

	l := lex("ok", testVrpPortConfigExample, "")
	l.options = lexOptions{
		emitComment: true,
		breakOK:     true,
		continueOK:  true,
	}
	items := make([]item, 0)
	for {

		item := l.nextItem()
		fmt.Println(item)
		items = append(items, item)
		if item.typ == itemEOF || item.typ == itemError {
			break
		}
	}

}

func Test_parse(t *testing.T) {

	l := lex("ok", testVrpPortConfigExample, "")
	l.options = lexOptions{
		emitComment: true,
		breakOK:     true,
		continueOK:  true,
	}

	tree := NewTree(l)
	tree.parse()

}
