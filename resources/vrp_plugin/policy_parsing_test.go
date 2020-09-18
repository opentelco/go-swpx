package main

import (
	"git.liero.se/opentelco/go-swpx/proto/traffic_policy"
	"testing"
)

func Test_ParsePolicyShaping(t *testing.T) {
	data := `<liero-test-a1>display current-configuration interface GigabitEthernet0/0/1 | include traffic-policy|shaping
 traffic-policy KBPS_100000 inbound
 qos queue 0 shaping cir 100000 pir 100000 cbs 2500000 pbs 2500000
<liero-test-a1>`

	policy, err := parsePolicy(data)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedPolicy := "KBPS_100000"
	gotPolicy := policy.Inbound

	if gotPolicy != expectedPolicy {
		t.Errorf("expected %s policy, got %s", expectedPolicy, gotPolicy)
	}

	expectedCbs := float64(2500000)
	gotCbs := policy.Qos.Shaping.Cbs

	if expectedCbs != gotCbs {
		t.Errorf("expected %f cbs, got %f", gotCbs, expectedCbs)
	}
}

func Test_ParseQOSTable(t *testing.T) {
	data := `<liero-test-a1>display qos queue statistics interface GigabitEthernet 0/0/1
------------------------------------------------------------
  Queue ID          : 0
  CIR(kbps)         : 100,000
  PIR(kbps)         : 100,000
  Passed Packets    : 48,140
  Passed Rate(pps)  : 0
  Passed Bytes      : 9,940,984
  Passed Rate(bps)  : 24
  Dropped Packets   : 0
  Dropped Rate(pps) : 0
  Dropped Bytes     : 0
  Dropped Rate(bps) : 0
------------------------------------------------------------
  Queue ID          : 1
  CIR(kbps)         : 0
  PIR(kbps)         : 1,000,000
  Passed Packets    : 0
  Passed Rate(pps)  : 0
  Passed Bytes      : 0
  Passed Rate(bps)  : 0
  Dropped Packets   : 0
  Dropped Rate(pps) : 0
  Dropped Bytes     : 0
  Dropped Rate(bps) : 0
------------------------------------------------------------
  Queue ID          : 2
  CIR(kbps)         : 0
  PIR(kbps)         : 1,000,000
  Passed Packets    : 0
  Passed Rate(pps)  : 0
  Passed Bytes      : 0
  Passed Rate(bps)  : 0
  Dropped Packets   : 0
  Dropped Rate(pps) : 0
  Dropped Bytes     : 0
  Dropped Rate(bps) : 0
------------------------------------------------------------
  Queue ID          : 3
  CIR(kbps)         : 0
  PIR(kbps)         : 1,000,000
  Passed Packets    : 0
  Passed Rate(pps)  : 0
  Passed Bytes      : 0
  Passed Rate(bps)  : 0
  Dropped Packets   : 0
  Dropped Rate(pps) : 0
  Dropped Bytes     : 0
  Dropped Rate(bps) : 0
------------------------------------------------------------
  Queue ID          : 4
  CIR(kbps)         : 0
  PIR(kbps)         : 1,000,000
  Passed Packets    : 0
  Passed Rate(pps)  : 0
  Passed Bytes      : 0
  Passed Rate(bps)  : 0
  Dropped Packets   : 0
  Dropped Rate(pps) : 0
  Dropped Bytes     : 0
  Dropped Rate(bps) : 0
------------------------------------------------------------
  Queue ID          : 5
  CIR(kbps)         : 0
  PIR(kbps)         : 1,000,000
  Passed Packets    : 0
  Passed Rate(pps)  : 0
  Passed Bytes      : 0
  Passed Rate(bps)  : 0
  Dropped Packets   : 0
  Dropped Rate(pps) : 0
  Dropped Bytes     : 0
  Dropped Rate(bps) : 0
------------------------------------------------------------
  Queue ID          : 6
  CIR(kbps)         : 0
  PIR(kbps)         : 1,000,000
  Passed Packets    : 506,611
  Passed Rate(pps)  : 0
  Passed Bytes      : 36,391,040
  Passed Rate(bps)  : 120
  Dropped Packets   : 0
  Dropped Rate(pps) : 0
  Dropped Bytes     : 0
  Dropped Rate(bps) : 0
------------------------------------------------------------
  Queue ID          : 7
  CIR(kbps)         : 0
  PIR(kbps)         : 1,000,000
  Passed Packets    : 204
  Passed Rate(pps)  : 0
  Passed Bytes      : 13,056
  Passed Rate(bps)  : 0
  Dropped Packets   : 0
  Dropped Rate(pps) : 0
  Dropped Bytes     : 0
  Dropped Rate(bps) : 0
------------------------------------------------------------
<liero-test-a1>`

	val, err := parseQueueStatistics(data)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := 8
	got := len(val.QueueStatistics)

	if got != expected {
		t.Errorf("expected %d stats, got %d", expected, got)
	}

}

func Test_ParsePolicyStatistics(t *testing.T) {
	data := `<liero-test-a1>display traffic policy statistics interface GigabitEthernet0/0/1 inbound verbose classifier-base

 Interface: GigabitEthernet0/0/1
 Traffic policy inbound: KBPS_100000
 Rule number: 28
 Current status: success
 Statistics interval: 300
---------------------------------------------------------------------
 Classifier: COS operator or
 Behavior: KBPS_100000
 Board : 0
---------------------------------------------------------------------
 Matched          |      Packets:                             0
                  |      Bytes:                               -
                  |      Rate(pps):                           0
                  |      Rate(bps):                           -
---------------------------------------------------------------------
   Passed         |      Packets:                             0
                  |      Bytes:                               -
                  |      Rate(pps):                           0
                  |      Rate(bps):                           -
---------------------------------------------------------------------
   Dropped        |      Packets:                             0
                  |      Bytes:                               -
                  |      Rate(pps):                           0
                  |      Rate(bps):                           -
---------------------------------------------------------------------
     Filter       |      Packets:                             0
                  |      Bytes:                               -
---------------------------------------------------------------------
     Car          |      Packets:                             0
                  |      Bytes:                               -
---------------------------------------------------------------------
 Classifier: DISCARD operator or
 Behavior: DISCARD
 Board : 0
---------------------------------------------------------------------
 Matched          |      Packets:                             0
                  |      Bytes:                               -
                  |      Rate(pps):                           0
                  |      Rate(bps):                           -
---------------------------------------------------------------------
   Passed         |      Packets:                             0
                  |      Bytes:                               -
                  |      Rate(pps):                           0
                  |      Rate(bps):                           -
---------------------------------------------------------------------
   Dropped        |      Packets:                             0
                  |      Bytes:                               -
                  |      Rate(pps):                           0
                  |      Rate(bps):                           -
---------------------------------------------------------------------
     Filter       |      Packets:                             0
                  |      Bytes:                               -
---------------------------------------------------------------------
     Car          |      Packets:                             0
                  |      Bytes:                               -
---------------------------------------------------------------------
 Classifier: IPTV operator and
 Behavior: IPTV
 Board : 0
---------------------------------------------------------------------
 Matched          |      Packets:                             0
                  |      Bytes:                               -
                  |      Rate(pps):                           0
                  |      Rate(bps):                           -
---------------------------------------------------------------------
   Passed         |      Packets:                             0
                  |      Bytes:                               -
                  |      Rate(pps):                           0
                  |      Rate(bps):                           -
---------------------------------------------------------------------
   Dropped        |      Packets:                             0
                  |      Bytes:                               -
                  |      Rate(pps):                           0
                  |      Rate(bps):                           -
---------------------------------------------------------------------
     Filter       |      Packets:                             0
                  |      Bytes:                               -
---------------------------------------------------------------------
     Car          |      Packets:                             0
                  |      Bytes:                               -
---------------------------------------------------------------------
 Classifier: VOIP operator or
 Behavior: VOIP
 Board : 0
---------------------------------------------------------------------
 Matched          |      Packets:                             0
                  |      Bytes:                               -
                  |      Rate(pps):                           0
                  |      Rate(bps):                           -
---------------------------------------------------------------------
   Passed         |      Packets:                             0
                  |      Bytes:                               -
                  |      Rate(pps):                           0
                  |      Rate(bps):                           -
---------------------------------------------------------------------
   Dropped        |      Packets:                             0
                  |      Bytes:                               -
                  |      Rate(pps):                           0
                  |      Rate(bps):                           -
---------------------------------------------------------------------
     Filter       |      Packets:                             0
                  |      Bytes:                               -
---------------------------------------------------------------------
     Car          |      Packets:                             0
                  |      Bytes:                               -
---------------------------------------------------------------------
 Classifier: OTHER operator and
 Behavior: KBPS_100000
 Board : 0
---------------------------------------------------------------------
 Matched          |      Packets:                       163,143
                  |      Bytes:                               -
                  |      Rate(pps):                           0
                  |      Rate(bps):                           -
---------------------------------------------------------------------
   Passed         |      Packets:                       163,143
                  |      Bytes:                               -
                  |      Rate(pps):                           0
                  |      Rate(bps):                           -
---------------------------------------------------------------------
   Dropped        |      Packets:                             0
                  |      Bytes:                               -
                  |      Rate(pps):                           0
                  |      Rate(bps):                           -
---------------------------------------------------------------------
     Filter       |      Packets:                             0
                  |      Bytes:                               -
---------------------------------------------------------------------
     Car          |      Packets:                             0
                  |      Bytes:                               -
---------------------------------------------------------------------
<liero-test-a1>`

	policy := &traffic_policy.ConfiguredTrafficPolicy{}
	err := parsePolicyStatistics(policy, data)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedPolicy := "KBPS_100000"
	gotPolicy := policy.InboundStatistics.TrafficPolicy

	expectedClassifiers := 5
	gotClassifiers := len(policy.InboundStatistics.Classifiers)

	if gotPolicy != expectedPolicy {
		t.Errorf("expected %s policy, got %s", expectedPolicy, gotPolicy)
	}

	if gotClassifiers != expectedClassifiers {
		t.Errorf("expected %d classifiers, got %d", expectedClassifiers, gotClassifiers)
	}

}
