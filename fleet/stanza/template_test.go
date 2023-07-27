package stanza

import (
	"testing"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/stanzapb"
	"github.com/stretchr/testify/assert"
)

func Test_Template(t *testing.T) {
	st := &stanzapb.Stanza{
		Id:       "test",
		Template: `Hostname: {{ .Hostname }}`,
	}

	d := &devicepb.Device{
		Hostname: "test",
	}

	x, err := FromTemplate(st.Id, st.Content, d)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Hostname: test", x)
}
