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
