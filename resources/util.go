package resources

import (
	"git.liero.se/opentelco/go-dnc/models/protobuf/metric"
	"git.liero.se/opentelco/go-dnc/models/protobuf/transport"
	"git.liero.se/opentelco/go-swpx/proto/networkelement"
	"git.liero.se/opentelco/go-swpx/shared/oids"
	"google.golang.org/protobuf/types/known/timestamppb"
	"regexp"
	"strconv"
	"strings"
)

var ReFindIndexinOID = regexp.MustCompile("(\\d+)$") // used to get the last number of the oid

type DiscoveryItem struct {
	Index       int
	Descr       string
	Alias       string
	ifType      int
	mtu         int
	physAddress string
	adminStatus int
	operStatus  int
	lastChange  *timestamppb.Timestamp
	highSpeed   int
}

func PopulateDiscoveryMap(task *transport.Message_Snmpc, discoveryMap map[int]*DiscoveryItem) {
	for _, m := range task.Snmpc.Metrics {
		index, _ := strconv.Atoi(ReFindIndexinOID.FindString(m.Oid))
		switch m.GetName() {
		case "ifIndex":
			if val, ok := discoveryMap[index]; ok {
				val.Index = int(m.GetIntValue())
			} else {
				discoveryMap[index] = &DiscoveryItem{
					Index: int(m.GetIntValue()),
				}
			}
		case "ifAlias":
			if val, ok := discoveryMap[index]; ok {
			} else {
				val.Alias = m.GetStringValue()
				discoveryMap[index] = &DiscoveryItem{
					Alias: m.GetStringValue(),
				}
			}
		case "ifDescr":
			if val, ok := discoveryMap[index]; ok {
				val.Descr = m.GetStringValue()
			} else {
				discoveryMap[index] = &DiscoveryItem{
					Descr: m.GetStringValue(),
				}
			}
		case "ifType":
			if val, ok := discoveryMap[index]; ok {
				val.ifType = int(m.GetIntValue())
			} else {
				discoveryMap[index] = &DiscoveryItem{
					ifType: int(m.GetIntValue()),
				}
			}
		case "ifMtu":
			if val, ok := discoveryMap[index]; ok {
				val.mtu = int(m.GetIntValue())
			} else {
				discoveryMap[index] = &DiscoveryItem{
					mtu: int(m.GetIntValue()),
				}
			}
		case "ifPhysAddress":
			if val, ok := discoveryMap[index]; ok {
				val.physAddress = m.GetHwaddrValue()
			} else {
				discoveryMap[index] = &DiscoveryItem{
					physAddress: m.GetHwaddrValue(),
				}
			}
		case "ifAdminStatus":
			if val, ok := discoveryMap[index]; ok {
				val.adminStatus = int(m.GetIntValue())
			} else {
				discoveryMap[index] = &DiscoveryItem{
					adminStatus: int(m.GetIntValue()),
				}
			}
		case "ifOperStatus":
			if val, ok := discoveryMap[index]; ok {
				val.operStatus = int(m.GetIntValue())
			} else {
				discoveryMap[index] = &DiscoveryItem{
					operStatus: int(m.GetIntValue()),
				}
			}
		case "ifLastChange":
			if val, ok := discoveryMap[index]; ok {
				val.lastChange = m.GetTimestampValue()
			} else {
				discoveryMap[index] = &DiscoveryItem{
					lastChange: m.GetTimestampValue(),
				}
			}
		case "ifHighSpeed":
			if val, ok := discoveryMap[index]; ok {
				val.highSpeed = int(m.GetIntValue())
			} else {
				discoveryMap[index] = &DiscoveryItem{
					highSpeed: int(m.GetIntValue()),
				}
			}
		}
	}
}

func GetIfXEntryInformation(m *metric.Metric, elementInterface *networkelement.Interface) {

	switch {
	case strings.HasPrefix(m.Oid, oids.IfOutUcastPkts):
		elementInterface.Stats.Output.Unicast = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfInBroadcastPkts):
		elementInterface.Stats.Input.Broadcast = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfInMulticastPkts):
		elementInterface.Stats.Input.Multicast = m.GetIntValue()
	case strings.HasPrefix(m.Oid, oids.IfOutBroadcastPkts):
		elementInterface.Stats.Output.Broadcast = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfOutMulticastPkts):
		elementInterface.Stats.Output.Multicast = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfAlias):
		elementInterface.Alias = m.GetStringValue()

	case strings.HasPrefix(m.Oid, oids.IfHighSpeed):
		elementInterface.Speed = m.GetIntValue()

	}

}

func GetIfEntryInformation(m *metric.Metric, elementInterface *networkelement.Interface) {
	switch {
	case strings.HasPrefix(m.Oid, oids.IfOutOctets):
		elementInterface.Stats.Output.Bytes = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfInOctets):
		elementInterface.Stats.Input.Bytes = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfInUcastPkts):
		elementInterface.Stats.Input.Unicast = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfInErrors):
		elementInterface.Stats.Input.Errors = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfOutErrors):
		elementInterface.Stats.Output.Errors = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfDescr):
		elementInterface.Description = m.GetStringValue()

	case strings.HasPrefix(m.Oid, oids.IfType):
		elementInterface.Type = networkelement.InterfaceType(m.GetIntValue())

	case strings.HasPrefix(m.Oid, oids.IfMtu):
		elementInterface.Mtu = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfLastChange):
		elementInterface.LastChanged = m.GetTimestampValue()

	case strings.HasPrefix(m.Oid, oids.IfPhysAddress):
		elementInterface.Hwaddress = m.GetStringValue()

	case strings.HasPrefix(m.Oid, oids.IfOperStatus):
		elementInterface.AdminStatus = networkelement.InterfaceStatus(m.GetIntValue())

	case strings.HasPrefix(m.Oid, oids.IfAdminStatus):
		elementInterface.OperationalStatus = networkelement.InterfaceStatus(m.GetIntValue())

	}
}

func GetHuaweiInformation(m *metric.Metric, elementInterface *networkelement.Interface) {
	switch {
	case strings.HasPrefix(m.Oid, oids.HuaIfEtherStatInCRCPkts):
		elementInterface.Stats.Input.CrcErrors = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.HuaIfEtherStatInPausePkts):
		elementInterface.Stats.Input.Pauses = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.HuaIfEthIfStatReset):
		elementInterface.Stats.Resets = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.HuaIfEtherStatOutPausePkts):
		elementInterface.Stats.Output.Pauses = m.GetIntValue()
	}
}

func GetSystemInformation(m *metric.Metric, ne *networkelement.Element) {
	switch m.Oid {
	case oids.SysContact:
		ne.Contact = m.GetStringValue()
	case oids.SysDescr:
		ne.Version = m.GetStringValue()
	case oids.SysLocation:
		ne.Location = m.GetStringValue()
	case oids.SysName:
		ne.Sysname = m.GetStringValue()
	// case oids.SysORLastChange:
	// case oids.SysObjectID:
	case oids.SysUpTime:
		ne.Uptime = m.GetStringValue()
	}
}

func ItemToInterface(v *DiscoveryItem) *networkelement.Interface {
	iface := &networkelement.Interface{
		AggregatedId:      "",
		Index:             int64(v.Index),
		Alias:             v.Alias,
		Description:       v.Descr,
		Hwaddress:         v.physAddress,
		Type:              networkelement.InterfaceType(v.ifType),
		AdminStatus:       networkelement.InterfaceStatus(v.adminStatus),
		OperationalStatus: networkelement.InterfaceStatus(v.operStatus),
		LastChanged:       v.lastChange,
		Speed:             int64(v.highSpeed),
		Mtu:               int64(v.mtu),
	}
	return iface
}
