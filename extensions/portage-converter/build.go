// Copyright © 2019-2020 Ettore Di Giacinto <mudler@gentoo.org>,
//                       Daniele Rondina <geaaru@sabayonlinux.org>
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, see <http://www.gnu.org/licenses/>.

package main

import (
	"errors"
	"io/ioutil"

	compiler "github.com/mudler/luet/pkg/compiler"
	pkg "github.com/mudler/luet/pkg/package"
	st "github.com/mudler/luet/pkg/spectooling"

	"gopkg.in/yaml.v2"
)

type LuetCompilationSpecSanitized struct {
	Steps      []string `json:"steps,omitempty" yaml:"steps,omitempty"`
	Env        []string `json:"env,omitempty" yaml:"env,omitempty"`
	Prelude    []string `json:"prelude,omitempty" yaml:"prelude,omitempty"` // Are run inside the image which will be our builder
	Image      string   `json:"image,omitempty" yaml:"image,omitempty"`
	Seed       string   `json:"seed,omitempty" yaml:"seed,omitempty"`
	PackageDir string   `json:"package_dir,omitempty" yaml:"package_dir,omitempty"`

	Retrieve []string `json:"retrieve,omitempty" yaml:"retrieve,omitempty"`

	Unpack   bool     `json:"unpack,omitempty" yaml:"unpack,omitempty"`
	Includes []string `json:"includes,omitempty" yaml:"includes,omitempty"`
	Excludes []string `json:"excludes,omitempty" yaml:"excludes,omitempty"`

	PackageRequires  []*st.DefaultPackageSanitized `json:"requires,omitempty" yaml:"requires,omitempty"`
	PackageConflicts []*st.DefaultPackageSanitized `json:"conflicts,omitempty" yaml:"conflicts,omitempty"`
	Provides         []*st.DefaultPackageSanitized `json:"provides,omitempty" yaml:"provides,omitempty"`
}

func NewLuetCompilationSpecSanitized(s *compiler.LuetCompilationSpec) *LuetCompilationSpecSanitized {
	return &LuetCompilationSpecSanitized{
		Steps:      s.Steps,
		Env:        s.Env,
		Prelude:    s.Prelude,
		Image:      s.Image,
		Seed:       s.Seed,
		PackageDir: s.PackageDir,
		Retrieve:   s.Retrieve,
		Unpack:     s.Unpack,
		Includes:   s.Includes,
		Excludes:   s.Excludes,
	}
}

func NewLuetCompilationSpecSanitizedFromYaml(data []byte) (*LuetCompilationSpecSanitized, error) {
	ans := &LuetCompilationSpecSanitized{}
	if err := yaml.Unmarshal(data, ans); err != nil {
		return nil, err
	}
	return ans, nil
}

func NewLuetCompilationSpecSanitizedFromFile(file string) (*LuetCompilationSpecSanitized, error) {
	if file == "" {
		return nil, errors.New("Invalid file path")
	}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	ans, err := NewLuetCompilationSpecSanitizedFromYaml(content)
	if err != nil {
		return nil, err
	}

	return ans, nil
}

func (s *LuetCompilationSpecSanitized) Yaml() ([]byte, error) {
	return yaml.Marshal(s)
}

func (s *LuetCompilationSpecSanitized) WriteBuildDefinition(path string) error {
	data, err := s.Yaml()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}

func (s *LuetCompilationSpecSanitized) Clone() (*LuetCompilationSpecSanitized, error) {
	data, err := s.Yaml()
	if err != nil {
		return nil, err
	}

	return NewLuetCompilationSpecSanitizedFromYaml(data)
}

func (s *LuetCompilationSpecSanitized) AddRequires(pkgs []*pkg.DefaultPackage) {
	for _, pkg := range pkgs {
		p := st.NewDefaultPackageSanitized(pkg)
		s.PackageRequires = append(s.PackageRequires, p)
	}
}

func (s *LuetCompilationSpecSanitized) AddConflicts(pkgs []*pkg.DefaultPackage) {
	for _, pkg := range pkgs {
		p := st.NewDefaultPackageSanitized(pkg)
		s.PackageConflicts = append(s.PackageConflicts, p)
	}
}
