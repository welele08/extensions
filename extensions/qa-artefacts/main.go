package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Package struct {
	Name     string `yaml:"name"`
	Category string `yaml:"category"`
	Version  string `yaml:"version"`
}

type Spec struct {
	Package Package `yaml:"package"`
}

type Metadata struct {
	Files []string `yaml:"files"`
	Spec  Spec     `yaml:"compilespec"`
}

var Cache map[string]Package = map[string]Package{}

func (p Package) String() string {
	return fmt.Sprintf("%s/%s-%s", p.Category, p.Name, p.Version)
}

func checkExistingFiles(yamlpath string) {
	dat, err := ioutil.ReadFile(yamlpath)
	if err != nil {
		log.Println(err)
		return
	}

	var Meta Metadata
	if err := yaml.Unmarshal(dat, &Meta); err != nil {
		log.Println(err)
		return
	}

	for _, f := range Meta.Files {
		if pack, ok := Cache[f]; ok {
			fmt.Println(Meta.Spec.Package, "ships", f, "already owned by", pack)
			continue
		}
		Cache[f] = Meta.Spec.Package
	}
}

func CheckFileConflict(path string) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		log.Panicf("failed reading directory: %s", err)
	}

	for _, e := range entries {
		if strings.Contains(e.Name(), "yaml") && e.Name() != "repository.yaml" {
			checkExistingFiles(filepath.Join(path, e.Name()))
		}
	}
}

func main() {
	CheckFileConflict(os.Args[1])
}
