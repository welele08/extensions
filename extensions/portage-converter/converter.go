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
	"fmt"
	"path/filepath"
	"strings"

	gentoo "github.com/Sabayon/pkgs-checker/pkg/gentoo"
	luet_config "github.com/mudler/luet/pkg/config"
	luet_helpers "github.com/mudler/luet/pkg/helpers"
	luet_pkg "github.com/mudler/luet/pkg/package"
	luet_tree "github.com/mudler/luet/pkg/tree"
)

type PortageConverter struct {
	Config         *luet_config.LuetConfig
	Cache          map[string]*gentoo.GentooPackage
	ReciperBuild   luet_tree.Builder
	ReciperRuntime luet_tree.Builder
	Specs          *PortageConverterSpecs
	TargetDir      string
}

func NewPortageConverter(targetDir string) *PortageConverter {
	return &PortageConverter{
		// TODO: we use it as singleton
		Config:         luet_config.LuetCfg,
		Cache:          make(map[string]*gentoo.GentooPackage, 0),
		ReciperBuild:   luet_tree.NewCompilerRecipe(luet_pkg.NewInMemoryDatabase(false)),
		ReciperRuntime: luet_tree.NewInstallerRecipe(luet_pkg.NewInMemoryDatabase(false)),
		TargetDir:      targetDir,
	}
}

func (pc *PortageConverter) LoadTrees(treePath []string) error {

	// Load trees
	for _, t := range treePath {
		err := pc.ReciperBuild.Load(t)
		if err != nil {
			return errors.New("Error on load tree" + err.Error())
		}
		err = pc.ReciperRuntime.Load(t)
		if err != nil {
			return errors.New("Error on load tree" + err.Error())
		}
	}

	return nil
}

func (pc *PortageConverter) LoadRules(file string) error {
	spec, err := LoadSpecsFile(file)
	if err != nil {
		return err
	}

	pc.Specs = spec

	return nil
}

func (pc *PortageConverter) GetSpecs() *PortageConverterSpecs {
	return pc.Specs
}

func (pc *PortageConverter) Generate() error {
	resolver := NewQDependsResolver()

	listSolutions := []*PortageSolution{}

	// Resolve all packages
	for _, artefact := range pc.Specs.GetArtefacts() {
		for _, pkg := range artefact.GetPackages() {

			solution, err := resolver.Resolve(pkg)
			if err != nil {
				return err
			}

			pkgDir := fmt.Sprintf("%s/%s/%s/",
				filepath.Join(pc.TargetDir, artefact.GetTree()),
				solution.Package.Category, solution.Package.Name)

			if solution.Package.Slot != "0" {
				slot := solution.Package.Slot
				// Ignore sub-slot
				if strings.Contains(solution.Package.Slot, "/") {
					slot = solution.Package.Slot[0:strings.Index(slot, "/")]
				}

				pkgDir = fmt.Sprintf("%s/%s-%s/%s",
					filepath.Join(pc.TargetDir, artefact.GetTree()),
					solution.Package.Category, slot, solution.Package.Name)
			}

			// Check if specs is already present
			if luet_helpers.Exists(filepath.Join(pkgDir, "definition.yaml")) {
				// Nothing to do
				continue
			}

			solution.PackageDir = pkgDir

			listSolutions = append(listSolutions, solution)

		}
	}

	for _, pkg := range listSolutions {

		defFile := filepath.Join(pkg.PackageDir, "definition.yaml")

		// Convert solution to luet package
		pack := pkg.ToPack()

		// Write definition.yaml
		luet_tree.WriteDefinitionFile(pack, defFile)

		// TODO: create build.yaml

	}

	return nil
}
