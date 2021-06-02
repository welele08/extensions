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

package devkit

import (
	"fmt"

	specs "github.com/Luet-lab/extensions/extensions/repo-devkit/pkg/specs"
	. "github.com/mudler/luet/pkg/logger"
	luet_pkg "github.com/mudler/luet/pkg/package"
)

type RepoList struct {
	*RepoKnife
}

func NewRepoList(s *specs.LuetRDConfig,
	backend, path string, opts map[string]string) (*RepoList, error) {

	knife, err := NewRepoKnife(s, backend, path, opts)
	if err != nil {
		return nil, err
	}

	ans := &RepoList{
		RepoKnife: knife,
	}

	return ans, nil
}

func (c *RepoList) ListPkgsAvailable() ([]*luet_pkg.DefaultPackage, error) {
	ans := []*luet_pkg.DefaultPackage{}

	err := c.RepoKnife.Analyze()
	if err != nil {
		return nil, err
	}

	for _, art := range c.MetaMap {
		ans = append(ans, art.CompileSpec.Package)
	}

	return ans, nil
}

func (c *RepoList) ListPkgsMissing() ([]*luet_pkg.DefaultPackage, error) {
	ans := []*luet_pkg.DefaultPackage{}
	mPkgs := make(map[string]bool, 0)

	err := c.RepoKnife.Analyze()
	if err != nil {
		return nil, err
	}

	// Creating map with available packages
	for _, art := range c.MetaMap {
		mPkgs[art.CompileSpec.Package.HumanReadableString()] = true
		DebugC(fmt.Sprintf("Find %s",
			art.CompileSpec.Package.HumanReadableString()))
	}

	for _, p := range c.ReciperRuntime.GetDatabase().World() {

		DebugC(fmt.Sprintf("Checking %s", p.HumanReadableString()))
		if _, ok := mPkgs[p.HumanReadableString()]; !ok {

			if c.Specs.List.ToIgnore(p.(*luet_pkg.DefaultPackage)) {
				Debug("Ignoring package %s", p.HumanReadableString())
				continue
			} else {
				ans = append(ans, p.(*luet_pkg.DefaultPackage))
			}
		}
	}

	return ans, nil
}
