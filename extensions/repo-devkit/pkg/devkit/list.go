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
	"errors"
	"fmt"

	specs "github.com/Luet-lab/extensions/extensions/repo-devkit/pkg/specs"

	"github.com/Luet-lab/luet-portage-converter/pkg/converter"
	. "github.com/mudler/luet/pkg/config"
	. "github.com/mudler/luet/pkg/logger"
	luet_pkg "github.com/mudler/luet/pkg/package"
	luet_tree "github.com/mudler/luet/pkg/tree"
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

func (c *RepoList) ListPkgsMissingByDeps(treePaths []string, withResolve bool) ([]*luet_pkg.DefaultPackage, error) {
	ans := []*luet_pkg.DefaultPackage{}
	reciperBuild := luet_tree.NewCompilerRecipe(luet_pkg.NewInMemoryDatabase(false))

	list, err := c.ListPkgsMissing()
	if err != nil {
		return list, err
	}

	pc := converter.NewPortageConverter("", "repoman")

	// Load ReciperBuildtime
	for _, t := range treePaths {
		if c.Verbose {
			InfoC(fmt.Sprintf(":evergreen_tree: Loading tree %s...", t))
		} else {
			DebugC(fmt.Sprintf(":evergreen_tree: Loading tree %s...", t))
		}
		err := reciperBuild.Load(t)
		if err != nil {
			return ans, errors.New("Error on load tree" + err.Error())
		}
	}

	// Using local load of the three to reduce log verbosity.
	pc.ReciperBuild = reciperBuild
	pc.TreePaths = treePaths

	// Create a map of the packages
	mMissings := make(map[string]*luet_pkg.DefaultPackage, 0)
	pList := []luet_pkg.Package{}
	for idx, p := range list {

		r, err := reciperBuild.GetDatabase().FindPackage(list[idx])
		if err != nil {
			return ans, errors.New(
				fmt.Sprintf("Error on resolve package %s", p.HumanReadableString()),
			)
		}

		// TODO: R is with broken requires!!!
		//mMissings[p.HumanReadableString()] = r.(*luet_pkg.DefaultPackage)
		mMissings[p.HumanReadableString()] = list[idx]
		pList = append(pList, r)
	}

	// Create stage4 worker
	worker := converter.Stage4Worker{
		Levels:  converter.NewStage4LevelsWithSize(1),
		Map:     make(map[string]*luet_pkg.DefaultPackage, 0),
		Changed: make(map[string]*luet_pkg.DefaultPackage, 0),
	}

	// Quiet output
	worker.Levels.Quiet = true

	for _, p := range pList {
		err := pc.Stage4AddDeps2Levels(p.(*luet_pkg.DefaultPackage),
			nil, &worker, 1, []string{},
		)
		if err != nil {
			return ans, errors.New("Error on initialize stage4 tree: " + err.Error())
		}
	}

	if LuetCfg.GetGeneral().Debug {
		pc.Stage4LevelsDumpWrapper(worker.Levels, "")
	}

	if withResolve {
		pc.Stage4AlignLevel1(&worker)
		worker.Levels.Resolve()
	}

	ans = c.retrieveMissingOrdered(&worker, mMissings)

	return ans, nil
}

func (c *RepoList) retrieveMissingOrdered(w *converter.Stage4Worker, missings map[string]*luet_pkg.DefaultPackage) []*luet_pkg.DefaultPackage {

	ans := []*luet_pkg.DefaultPackage{}
	processedDeps := make(map[string]bool, 0)

	// TODO: Check why we have leaf with version >=0
	// Temporary workaround that handle missing package with cat/pkg
	for _, v := range missings {
		missings[fmt.Sprintf("%s/%s", v.GetCategory(), v.GetName())] = v
	}

	for i := len(w.Levels.Levels) - 1; i >= 0; i-- {

		for _, dep := range w.Levels.Levels[i].Deps {

			key := fmt.Sprintf("%s/%s", dep.GetCategory(), dep.GetName())
			if _, ok := processedDeps[key]; !ok {

				// Check if the package is one of the package missing
				if v, ok := missings[key]; ok {

					// Package to build
					ans = append(ans, v)

				}

				// Add the package to the processed map
				processedDeps[key] = true

			}

		}

	}

	return ans
}
