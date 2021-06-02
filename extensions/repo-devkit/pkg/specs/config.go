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
	"fmt"
	"io/ioutil"

	. "github.com/mudler/luet/pkg/logger"
	luet_pkg "github.com/mudler/luet/pkg/package"
	luet_version "github.com/mudler/luet/pkg/versioner"
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
func (c *LuetRDConfig) GetList() *LuetRDCList       { return &c.List }

func (c *LuetRDCCleaner) HasExcludes() bool {
	return len(c.Excludes) > 0
}

func (c *LuetRDCList) HasFilters() bool {
	return len(c.ExcludePkgs) > 0
}

func (c *LuetPackage) GetName() string     { return c.Name }
func (c *LuetPackage) GetCategory() string { return c.Category }
func (c *LuetPackage) GetVersion() string  { return c.Version }
func (p *LuetPackage) HumanReadableString() string {
	return fmt.Sprintf("%s/%s-%s", p.Category, p.Name, p.Version)
}

func (c *LuetRDCList) ToIgnore(pkg *luet_pkg.DefaultPackage) bool {
	ans := false

	if c.HasFilters() {
		pSelector, err := luet_version.ParseVersion(pkg.GetVersion())
		if err != nil {
			Warning(fmt.Sprintf(
				"Error on create package selector for package %s: %s",
				pkg.HumanReadableString(), err.Error()))
			return true
		}

		for _, f := range c.ExcludePkgs {
			if f.GetName() != pkg.GetName() ||
				f.GetCategory() != pkg.GetCategory() {
				continue
			}

			selector, err := luet_version.ParseVersion(f.GetVersion())
			if err != nil {
				Warning(fmt.Sprintf(
					"Error on create version selector for package %s: %s",
					f.HumanReadableString(), err.Error()))
				continue
			}

			admit, err := luet_version.PackageAdmit(selector, pSelector)
			if err != nil {
				Warning(fmt.Sprintf("Error on check package %s: %s",
					f.HumanReadableString(), err.Error()))
				continue
			}

			if admit {
				ans = true
			}

		}
	}

	return ans
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
