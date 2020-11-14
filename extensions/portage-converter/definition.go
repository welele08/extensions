/*
Copyright (C) 2020  Daniele Rondina <geaaru@sabayonlinux.org>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.

*/
package main

import (
	"errors"
	"io/ioutil"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type PortageConverterSpecs struct {
	SkippedResolutions PortageConverterSkips `json:"skipped_resolutions,omitempty" yaml:"skipped_resolutions,omitempty"`

	IncludeFiles  []string                   `json:"include_files,omitempty" yaml:"include_files,omitempty"`
	Artefacts     []PortageConverterArtefact `json:"artefacts,omitempty" yaml:"artefacts,omitempty"`
	BuildTmplFile string                     `json:"build_template_file" yaml:"build_template_file"`
}

type PortageConverterSkips struct {
	Packages   []PortageConverterPkg `json:"packages,omitempty" yaml:"packages,omitempty"`
	Categories []string              `json:"categories,omitempty" yaml:"categories,omitempty"`
}

type PortageConverterPkg struct {
	Name     string `json:"name" yaml:"name"`
	Category string `json:"category" yaml:"category"`
}

type PortageConverterArtefact struct {
	Tree     string   `json:"tree" yaml:"tree"`
	Packages []string `json:"packages" yaml:"packages"`
}

type PortageConverterInclude struct {
	SkippedResolutions PortageConverterSkips      `json:"skipped_resolutions,omitempty" yaml:"skipped_resolutions,omitempty"`
	Artefacts          []PortageConverterArtefact `json:"artefacts,omitempty" yaml:"artefacts,omitempty"`
}

func SpecsFromYaml(data []byte) (*PortageConverterSpecs, error) {
	ans := &PortageConverterSpecs{}
	if err := yaml.Unmarshal(data, ans); err != nil {
		return nil, err
	}
	return ans, nil
}

func IncludeFromYaml(data []byte) (*PortageConverterInclude, error) {
	ans := &PortageConverterInclude{}
	if err := yaml.Unmarshal(data, ans); err != nil {
		return nil, err
	}
	return ans, nil
}

func LoadSpecsFile(file string) (*PortageConverterSpecs, error) {

	if file == "" {
		return nil, errors.New("Invalid file path")
	}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	ans, err := SpecsFromYaml(content)
	if err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(path.Dir(file))
	if err != nil {
		return nil, err
	}

	if len(ans.IncludeFiles) > 0 {

		for _, include := range ans.IncludeFiles {

			if include[0:1] != "/" {
				include = filepath.Join(absPath, include)
			}
			content, err := ioutil.ReadFile(include)
			if err != nil {
				return nil, err
			}

			data, err := IncludeFromYaml(content)
			if err != nil {
				return nil, err
			}

			if len(data.SkippedResolutions.Packages) > 0 {
				ans.SkippedResolutions.Packages = append(ans.SkippedResolutions.Packages,
					data.SkippedResolutions.Packages...)
			}

			if len(data.SkippedResolutions.Categories) > 0 {
				ans.SkippedResolutions.Categories = append(ans.SkippedResolutions.Categories,
					data.SkippedResolutions.Categories...)
			}

		}
	}

	if ans.BuildTmplFile != "" && ans.BuildTmplFile[0:1] != "/" {
		// Convert in abs path
		ans.BuildTmplFile = filepath.Join(absPath, ans.BuildTmplFile)
	}

	return ans, nil
}

func (s *PortageConverterSpecs) GetArtefacts() []PortageConverterArtefact {
	return s.Artefacts
}

func (a *PortageConverterArtefact) GetPackages() []string { return a.Packages }
func (a *PortageConverterArtefact) GetTree() string       { return a.Tree }
