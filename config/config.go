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
	"io/ioutil"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

const (
	DefaultMongoTimeout = time.Second * 5
)

var (
	ErrResourceConfigurationMissing = fmt.Errorf("resource plugin was not found or loaded")
	ErrProviderConfigurationMissing = fmt.Errorf("provider plugin was not found or loaded")
)

func loadFile(fname string) ([]byte, error) {
	content, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func parseConfig(src []byte, filename string, cfg interface{}) error {
	var diags hcl.Diagnostics
	if len(src) == 0 {
		return fmt.Errorf("no byte array supplied")
	}

	file, diags := hclsyntax.ParseConfig(src, filename, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return fmt.Errorf("config parse: %w", diags)
	}

	diags = gohcl.DecodeBody(file.Body, nil, cfg)
	if diags.HasErrors() {
		return fmt.Errorf("config parse: %w", diags)
	}

	return nil
}

func LoadConfig(fname string, cfg interface{}) error {
	content, err := loadFile(fname)
	if err != nil {
		return err
	}

	return parseConfig(content, fname, cfg)
}
