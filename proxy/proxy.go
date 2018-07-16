// Copyright 2017 Qian Qiao
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package proxy provides tools proxying requests from different backend
// servers
package proxy

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/qqiao/cli"
)

var configPath string

// Target is for each proxy target configuration
type Target struct {
	Path   string `json:"path"`
	Target string `json:"target"`
}

// Config encapsulates proxy configuration
type Config struct {
	Targets []Target `json:"targets"`
}

var config Config

// NewComponent returns command line component that can be used to form
// complex commandline application
func NewComponent() *cli.Component {
	comp := &cli.Component{
		UsageLine: "proxy [-config CONFIG_FILE]",
		Short:     "proxies requests into different server backends",
		Run:       proxy,
	}

	comp.FlagSet().StringVar(&configPath, "config", "proxy_config.json",
		"configuration for proxy backends")

	return comp
}

func proxy(ctx context.Context, comp *cli.Component, args []string) {
	flagSet := comp.FlagSet()

	if flag.ErrHelp == flagSet.Parse(args) {
		return
	}

	if !filepath.IsAbs(configPath) {
		cwd, err := os.Getwd()

		if nil != err {
			log.Printf("Unable to determine current working directory. Error: %v",
				err)
			os.Exit(1)
		}

		configPath = filepath.Join(cwd, configPath)
	}

	jsonFile, err := os.Open(configPath)
	if nil != err {
		log.Printf("Unable to open '%s'. Error: %v", configPath, err)
		os.Exit(1)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if nil != err {
		log.Printf("Unable to read '%s'. Error: %v", configPath, err)
		os.Exit(1)
	}

	if err := json.Unmarshal(byteValue, &config); nil != err {
		log.Printf("Unable to parse '%s'. Error: %v", configPath, err)
		os.Exit(1)
	}

	log.Printf("%v", config)
}
