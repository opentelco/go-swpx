package main_test

import (
	"fmt"
	"math"
	"testing"
)

func Test_uWConvertion(t *testing.T) {
	rx := 256
	fRX := float64(rx)
	fmt.Printf("The value is: %.2f", math.Round((10*math.Log10(fRX/1000)*100)/100))
}
