package resources

import (
	"fmt"
	"testing"
)

func Test_ConvertToDb(t *testing.T) {
	x := ConvertToDb(13032)
	fmt.Println(x)

}
