/*
 * Copyright (c) 2020. Liero AB
 *
 * Permission is hereby granted, free of charge, to any person obtaining
 * a copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software
 * is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
 * OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
 * IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
 * CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
 * TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"git.liero.se/opentelco/go-swpx/shared"
)

type Config struct {
	Key       string
	Region    string
	Providers []ProviderConfig
	Drivers   []DriverConfig
}

type ProviderConfig struct {
	Weight int
}

type DriverConfig struct {
	OIDs          []string
	RegexMatchOID string
}

// func ParseConfig(hclText string) (*Config, error) {
// 	result := &Config{}

// 	hclParseTree, err := hcl.Parse(hclText)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err := hcl.Unmarshal(b, result)
// 	if err != nil {

// 	}
// 	// if rawRegion := hclParseTree.Get("region", false); rawRegion != nil {
// 	// 	result.Region = rawRegion.Value.(string)
// 	// }

// 	log.Printf("%+v\n", result)

// 	return result, nil
// }

func FindAndParse() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	home = path.Join(home, shared.DEFAULT_CORE_PATH)

	var ff = func(pathX string, infoX os.FileInfo, errX error) error {

		// first thing to do, check error. and decide what to do about it
		if errX != nil {
			log.Fatal(errX)
		}

		fmt.Printf("pathX: %v\n", pathX)

		// find out if it's a dir or file, if file, print info
		if infoX.IsDir() {
			fmt.Printf("is dir.\n")
		} else {
			fmt.Printf("  dir: 「%v」\n", filepath.Dir(pathX))
			fmt.Printf("  file name 「%v」\n", infoX.Name())
			fmt.Printf("  extenion: 「%v」\n", filepath.Ext(pathX))
		}

		return nil
	}

	err = filepath.Walk(home, ff)

	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", home, err)
	}

}
