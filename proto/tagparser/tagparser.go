/*
 * Copyright (c) 2020. (vx-oss-app) VX Service Delivery AB
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by lunne, Last modified: 2020-06-09, 09:37
 */

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

func init() {

}

var replaceFile = true

var rootCmd = &cobra.Command{
	Use:  "go-tagparser [filePath]",
	Args: cobra.MinimumNArgs(0),
	Long: fmt.Sprintf(`
 go-tagparser <path>
 simple injection of BSON tags. It will search the selected folder for *.pb.go and inject bson into the file

 it injects a bson tag that is the same as the json tag
 Id                   string   protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"
 - will become -
 Id                   string   protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" bson:"id,omitempty"


 if bson already exists as a tag it will not touch row.



 `),
	Run: func(cmd *cobra.Command, args []string) {
		path := "../proto/go.pkg/oss/"
		if len(args) > 0 {
			path = args[0]
		}
		fmt.Println("using path:", path)
		walk(path)
	},
}

func walk(path string) {
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.Contains(info.Name(), "pb.go") {
				err := parseFile(path)
				if err != nil {
					log.Fatal(err)
				}
			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}
}

func parseFile(fname string) error {
	f, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}

	idTag := regexp.MustCompile("bson:\"id\"")
	tag := regexp.MustCompile("json:\".*\"")
	rows := bytes.Split(f, []byte("\n"))
	newfile := make([][]byte, 0)
	var writefile bool
	for _, r := range rows {
		// if the file has the json tag
		if bytes.Contains(r, []byte("json:")) && !bytes.Contains(r, []byte("bson:")) {
			writefile = true
			newrow := make([]byte, 0)
			locs := tag.FindIndex(r)

			new := bytes.Replace(r[locs[0]:locs[1]], []byte("json"), []byte(" bson"), -1)
			new = bytes.Replace(new, []byte(",omitempty"), []byte(""), -1)
			p1 := r[0:locs[1]]
			p2 := new
			p3 := r[locs[1]:]

			newrow = append(newrow, p1...)
			newrow = append(newrow, p2...)
			newrow = append(newrow, p3...)

			newrow = idTag.ReplaceAll(newrow, []byte("bson:\"_id\"")) // replace id with _id
			newfile = append(newfile, newrow)

		} else { // no tag, just insert the line as it is
			newfile = append(newfile, r)
		}

	}

	if writefile {
		// if replaceFile is disabled
		if !replaceFile {
			d, s := filepath.Split(fname)
			s = fmt.Sprintf("test_%s", s)
			fname = filepath.Join(d, s)
		}
		output := bytes.Join(newfile, []byte("\n")) // join the array with new lines
		err = ioutil.WriteFile(fname, output, 0644)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("wrote tag changes to %s\n", fname)
	}

	return nil
}

func main() {
	// var fname string
	// fname = "../proto/go.pkg/oss/nbp/nbp.pb.go"
	// fname = "../proto/go.pkg/oss/inventory/connection/connection.pb.go"
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
