package main

import (
	"math"
	"regexp"
	"strconv"
	"strings"

	"git.liero.se/opentelco/go-dnc/models/pb/metric"
	"git.liero.se/opentelco/go-dnc/models/pb/snmpc"
	"git.liero.se/opentelco/go-swpx/proto/go/networkelement"
	"git.liero.se/opentelco/go-swpx/shared/oids"
	"github.com/araddon/dateparse"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var reFindIndexinOID = regexp.MustCompile(`(\d+)$`) // used to get the last number of the oid

type discoveryItem struct {
	Index       int64
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

func populateDiscoveryMap(logger hclog.Logger, task *snmpc.Task, discoveryMap map[int]*discoveryItem) {
	for _, m := range task.Metrics {
		index, _ := strconv.Atoi(reFindIndexinOID.FindString(m.Oid))
		switch m.GetName() {
		case "ifIndex":
			if val, ok := discoveryMap[index]; ok {
				i, _ := strconv.Atoi(m.GetValue())
				val.Index = int64(i)
			} else {
				i, _ := strconv.Atoi(m.GetValue())
				discoveryMap[index] = &discoveryItem{
					Index: int64(i),
				}
			}
		case "ifAlias":
			if val, ok := discoveryMap[index]; ok {
				val.Alias = m.GetValue()
			} else {
				discoveryMap[index] = &discoveryItem{
					Alias: m.GetValue(),
				}
			}

		case "ifXEntry":

			if val, ok := discoveryMap[index]; ok {
				val.Descr = m.GetValue()
			} else {
				discoveryMap[index] = &discoveryItem{
					Descr: m.GetValue(),
				}
			}

		case "ifDescr":

			if val, ok := discoveryMap[index]; ok {
				val.Descr = m.GetValue()
			} else {
				discoveryMap[index] = &discoveryItem{
					Descr: m.GetValue(),
				}
			}

		case "ifType":

			i, _ := strconv.Atoi(m.GetValue())

			if val, ok := discoveryMap[index]; ok {
				val.ifType = i
			} else {
				discoveryMap[index] = &discoveryItem{
					ifType: i,
				}
			}

		case "ifMtu":
			i, _ := strconv.Atoi(m.GetValue())
			if val, ok := discoveryMap[index]; ok {
				val.mtu = i
			} else {
				discoveryMap[index] = &discoveryItem{
					mtu: i,
				}
			}

		case "ifPhysAddress":
			if val, ok := discoveryMap[index]; ok {
				val.physAddress = m.GetValue()
			} else {
				discoveryMap[index] = &discoveryItem{
					physAddress: m.GetValue(),
				}
			}

		case "ifAdminStatus":
			i, _ := strconv.Atoi(m.GetValue())
			if val, ok := discoveryMap[index]; ok {
				val.adminStatus = i
			} else {
				discoveryMap[index] = &discoveryItem{
					adminStatus: i,
				}
			}

		case "ifOperStatus":
			i, _ := strconv.Atoi(m.GetValue())
			if val, ok := discoveryMap[index]; ok {
				val.operStatus = i
			} else {
				discoveryMap[index] = &discoveryItem{
					operStatus: i,
				}
			}

		case "ifLastChange":
			ts := dateparse.MustParse(m.GetValue())
			tspb := timestamppb.New(ts)
			if val, ok := discoveryMap[index]; ok {
				val.lastChange = tspb
			} else {
				discoveryMap[index] = &discoveryItem{
					lastChange: tspb,
				}
			}
		case "ifHighSpeed":
			i, _ := strconv.Atoi(m.GetValue())

			if val, ok := discoveryMap[index]; ok {
				val.highSpeed = i
			} else {
				discoveryMap[index] = &discoveryItem{
					highSpeed: i,
				}
			}
		}
	}
}

func getIfXEntryInformation(m *metric.Metric, elementInterface *networkelement.Interface) {

	i, _ := strconv.Atoi(m.GetValue())
	switch {
	case strings.HasPrefix(m.Oid, oids.IfOutUcastPkts):
		elementInterface.Stats.Output.Unicast = int64(i)

	case strings.HasPrefix(m.Oid, oids.IfInBroadcastPkts):
		elementInterface.Stats.Input.Broadcast = int64(i)

	case strings.HasPrefix(m.Oid, oids.IfInMulticastPkts):
		elementInterface.Stats.Input.Multicast = int64(i)
	case strings.HasPrefix(m.Oid, oids.IfOutBroadcastPkts):
		elementInterface.Stats.Output.Broadcast = int64(i)

	case strings.HasPrefix(m.Oid, oids.IfOutMulticastPkts):
		elementInterface.Stats.Output.Multicast = int64(i)

	case strings.HasPrefix(m.Oid, oids.IfAlias):
		elementInterface.Alias = m.GetValue()

	case strings.HasPrefix(m.Oid, oids.IfHighSpeed):
		elementInterface.Speed = int64(i)

	}

}

func getSystemInformation(m *metric.Metric, ne *networkelement.Element) {
	switch m.Oid {
	case oids.SysContact:
		ne.Contact = m.GetValue()
	case oids.SysDescr:
		ne.Version = m.GetValue()
	case oids.SysLocation:
		ne.Location = m.GetValue()
	case oids.SysName:
		ne.Sysname = m.GetValue()
	// case oids.SysORLastChange:
	// case oids.SysObjectID:
	case oids.SysUpTime:
		ne.Uptime = m.GetValue()
	}
}

func itemToInterface(v *discoveryItem) *networkelement.Interface {
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

// convert uW(int) to dB(float64)
// rounds to 2 nearest decimals
func convertToDb(uw int64) float64 {
	v := 10 * math.Log10(float64(uw)/1000)
	f := math.Round(v*100) / 100
	if math.IsInf(f, 0) || math.IsNaN(f) {
		return -40
	}
	return f
}
