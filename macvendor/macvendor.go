/*
 * File: macvendor.go
 * Project: macvendor
 * File Created: Sunday, 14th February 2021 3:31:40 pm
 * Author: Mathias Ehrlin (mathias.ehrlin@vx.se)
 * -----
 * Last Modified: Sunday, 14th February 2021 4:09:23 pm
 * Modified By: Mathias Ehrlin (mathias.ehrlin@vx.se>)
 * -----
 * Copyright - 2021 VX Service Delivery AB
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * -----
 */

package macvendor

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type vendor struct {
}

func ReadDataFile() error {
	file, err := os.Open("vendor_list.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		r := strings.Split(scanner.Text(), "")
		var key []string
		var value []string
		for ix, v := range r {
			if ix <= 5 {
				key = append(key, v)
			} else {
				value = append(value, v)
			}
		}
		fmt.Println(strings.ToLower(strings.Join(key, "")), " => ", strings.TrimSpace(strings.Join(value, "")))

	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
