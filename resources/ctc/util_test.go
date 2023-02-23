package main

import (
	"fmt"
	"testing"
)

func Test_ConvertToDb(t *testing.T) {
	x := convertToDb(13032)
	fmt.Println(x)

}
