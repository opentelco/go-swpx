package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
)

func Test_Diff(t *testing.T) {
	a, _ := os.ReadFile("first.cfg")
	b, _ := os.ReadFile("second.cfg")
	aString := string(a)
	bString := string(b)

	edits := myers.ComputeEdits(span.URIFromPath(""), aString, bString)
	diff := fmt.Sprint(gotextdiff.ToUnified("previous-config", "new-config", aString, edits))
	fmt.Println(diff)

}
