/*
Copyright (C) 2020-2021  Daniele Rondina <geaaru@sabayonlinux.org>

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
package specs

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func NewLuetRDConfig() *LuetRDConfig {
	return &LuetRDConfig{
		Cleaner: LuetRDCCleaner{
			Excludes: []string{},
		},
	}
}

func (c *LuetRDConfig) GetCleaner() *LuetRDCCleaner { return &c.Cleaner }

func (c *LuetRDCCleaner) HasExcludes() bool {
	return len(c.Excludes) > 0
}

func SpecsFromYaml(data []byte) (*LuetRDConfig, error) {
	ans := NewLuetRDConfig()
	if err := yaml.Unmarshal(data, ans); err != nil {
		return nil, err
	}
	return ans, nil
}

func LoadSpecsFile(file string) (*LuetRDConfig, error) {
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

	return ans, nil
}
