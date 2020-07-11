/*
Copyright 2020 Adam Yi. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Auto update challenge_list and container_bundle BUILD files for CTF challenges

package main

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bazelbuild/buildtools/build"
	"github.com/bmatcuk/doublestar"
)

var GCR_PREFIX string

var mode string

func checkIfImage(path string) (isImage bool, isChallenge bool) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	f, err := build.ParseBuild(path, data)
	if err != nil {
		panic(err)
	}
	rules := f.Rules("")
	for _, r := range rules {
		if r.AttrString("name") == "image" {
			isImage = true
		}
		if r.Kind() == "ctf_challenge" {
			isChallenge = true
		}
	}
	return
}

func updateChalList(path string, chals []string) {
	log.Printf("Updating %v", path)
	file, err := os.OpenFile(path, os.O_RDWR, 0000)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	f, err := build.ParseBuild(path, data)
	if err != nil {
		panic(err)
	}
	rules := f.Rules("challenges_list")
	if len(rules) != 1 {
		panic("there is no or more than one challenges_list, i'm confused")
	}
	r := rules[0]
	list := &build.ListExpr{}
	for _, chal := range chals {
		list.List = append(list.List, &build.StringExpr{Value: chal})
	}
	build.SortStringList(list)
	r.SetAttr("deps", list)
	out := build.Format(f)
	switch mode {
	case "fix":
		file.Seek(0, io.SeekStart)
		file.Write(out)
		file.Truncate(int64(len(out)))
	case "check":
		if bytes.Compare(data, out) != 0 {
			log.Printf("%v not up to date", path)
			os.Exit(1)
		}
	case "stdout":
		os.Stdout.Write(out)
	}
}

func updateContainersList(path string, containers []string) {
	log.Printf("Updating %v", path)
	log.Println(containers)
	file, err := os.OpenFile(path, os.O_RDWR, 0000)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	f, err := build.ParseBuild(path, data)
	if err != nil {
		panic(err)
	}

	rules := f.Rules("container_bundle")
	if len(rules) != 1 {
		panic("There's no or more than 1 container_bundle, i'm confused")
	}
	r := rules[0]
	m := make(map[string]string)
	mc := make(map[string][]build.Comment)
	for _, c := range containers {
		m[GCR_PREFIX+c+":latest"] = "//" + c + ":image"
	}
	if r.AttrString("name") == "all_containers" {
		images, ok := r.Attr("images").(*build.DictExpr)
		if !ok {
			panic("images not DictExpr")
		}
		for _, e := range images.List {
			kv, ok := e.(*build.KeyValueExpr)
			if !ok {
				panic("shouldn't happen but unable to cast to KeyValueExpr")
			}
			cos := kv.Comment().Suffix
			for _, co := range cos {
				if strings.Index(co.Token, "ctflark: keep") >= 0 || strings.Index(co.Token, "ctflark:keep") >= 0 {
					m[kv.Key.(*build.StringExpr).Value] = kv.Value.(*build.StringExpr).Value
					mc[kv.Key.(*build.StringExpr).Value] = cos
					break
				}
			}
		}
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	dict := &build.DictExpr{}

	for _, k := range keys {
		kve := &build.KeyValueExpr{Key: &build.StringExpr{Value: k}, Value: &build.StringExpr{Value: m[k]}, Comments: build.Comments{Suffix: mc[k]}}
		dict.List = append(dict.List, kve)
	}
	r.SetAttr("images", dict)
	out := build.Format(f)
	switch mode {
	case "fix":
		file.Seek(0, io.SeekStart)
		file.Write(out)
		file.Truncate(int64(len(out)))
	case "check":
		if bytes.Compare(data, out) != 0 {
			log.Printf("%v not up to date", path)
			os.Exit(1)
		}
	case "stdout":
		os.Stdout.Write(out)
	}
}

func main() {
	flag.StringVar(&mode, "mode", "fix", "fix/check/stdout")
	flag.Parse()
	checkdirs := [...]string{"infra", "challenges"}
	builtimages := make([]string, 0)
	workspace := os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	log.Printf("Running ctflark for %v", workspace)
	for _, cd := range checkdirs {
		matches, err := doublestar.Glob(workspace + "/" + cd + "/**/BUILD.bazel")
		if err != nil {
			panic(err)
		}
		matches2, err := doublestar.Glob(workspace + "/" + cd + "/**/BUILD")
		if err != nil {
			panic(err)
		}
		matches = append(matches, matches2...)
		builtchals := make([]string, 0)
		for _, m := range matches {
			pkg, err := filepath.Rel(workspace, m)
			if err != nil {
				panic(err)
			}
			pkg = pkg[0:strings.LastIndex(pkg, "/")]
			// log.Printf("%v %v", pkg,m)
			isImage, isChal := checkIfImage(m)
			if isImage {
				builtimages = append(builtimages, pkg)
			}
			if isChal {
				builtchals = append(builtchals, "//"+pkg+":challenge")
			}
		}
		log.Printf("%v challenge_list: %v", cd, builtchals)
		updateChalList(workspace+"/"+cd+"/BUILD", builtchals)
	}
	log.Printf("Built Images: %v", builtimages)

	updateContainersList(workspace+"/BUILD", builtimages)

	log.Println("All done! Have a great day!")
}
