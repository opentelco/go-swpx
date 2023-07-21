package device

import (
	"testing"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
)

func Test_getChanges(t *testing.T) {
	type args struct {
		deviceA *devicepb.Device
		deviceB *devicepb.Device
	}
	tests := []struct {
		name string
		args args
		want []*devicepb.Change
	}{
		{
			name: "no changes",
			args: args{
				deviceA: &devicepb.Device{
					Hostname: "test",
				},
				deviceB: &devicepb.Device{
					Hostname: "test",
				},
			},
			want: []*devicepb.Change{},
		},
		{
			name: "domain changed",
			args: args{
				deviceA: &devicepb.Device{
					Hostname: "test",
				},
				deviceB: &devicepb.Device{
					Hostname: "test",
					Domain:   "test",
				},
			},
			want: []*devicepb.Change{
				{
					Field:    "domain",
					OldValue: "",
					NewValue: "test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getChanges(tt.args.deviceA, tt.args.deviceB); len(got) != len(tt.want) {
				t.Errorf("getChanges() = %v, want %v", got, tt.want)
			}
		})
	}
}
