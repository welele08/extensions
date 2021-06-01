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
	"path/filepath"
	"regexp"

	"github.com/Luet-lab/extensions/extensions/repo-devkit/pkg/backends"
	specs "github.com/Luet-lab/extensions/extensions/repo-devkit/pkg/specs"

	tmtools "github.com/geaaru/time-master/pkg/tools"
	artifact "github.com/mudler/luet/pkg/compiler/types/artifact"
	. "github.com/mudler/luet/pkg/logger"
	luet_pkg "github.com/mudler/luet/pkg/package"
	luet_tree "github.com/mudler/luet/pkg/tree"
)

type RepoKnife struct {
	Specs          *specs.LuetRDConfig
	BackendHandler specs.RepoBackendHandler
	ReciperRuntime luet_tree.Builder

	PkgsMap      map[string]string
	MetaMap      map[string]*artifact.PackageArtifact
	Files2Remove []string
	Verbose      bool
}

func NewRepoKnife(s *specs.LuetRDConfig,
	backend, path string, opts map[string]string) (*RepoKnife, error) {
	var err error
	var handler specs.RepoBackendHandler

	ans := &RepoKnife{
		Specs:          s,
		ReciperRuntime: luet_tree.NewInstallerRecipe(luet_pkg.NewInMemoryDatabase(false)),
		PkgsMap:        make(map[string]string, 0),
		MetaMap:        make(map[string]*artifact.PackageArtifact, 0),
	}

	switch backend {
	case "local":
		handler, err = backends.NewBackendLocal(s, path)
	case "mottainai":
		handler, err = backends.NewBackendMottainai(s, path, opts)
	default:
		return nil, errors.New("Invalid backend")
	}

	if err != nil {
		return nil, err
	}
	ans.BackendHandler = handler
	return ans, nil
}

func (c *RepoKnife) LoadTrees(treePath []string) error {

	// Load trees
	for _, t := range treePath {
		if c.Verbose {
			InfoC(fmt.Sprintf(":evergreen_tree: Loading tree %s...", t))
		} else {
			DebugC(fmt.Sprintf(":evergreen_tree: Loading tree %s...", t))
		}
		err := c.ReciperRuntime.Load(t)
		if err != nil {
			return errors.New("Error on load tree" + err.Error())
		}
	}

	return nil
}

func (c *RepoKnife) Analyze() error {

	// Reset previous values
	c.PkgsMap = make(map[string]string, 0)
	c.MetaMap = make(map[string]*artifact.PackageArtifact, 0)
	c.Files2Remove = []string{}

	// Retrieve the list of the files
	files, err := c.BackendHandler.GetFilesList()
	if err != nil {
		return err
	}

	if c.Specs.GetCleaner().HasExcludes() {
		files, err = c.GetFilteredList(files)
		if err != nil {
			return err
		}
	}

	// Exclude repository files
	repoRegex := []string{
		"repository.meta.yaml.tar.*|repository.meta.yaml",
		"repository.yaml",
		"tree.tar.*|tree.tar",
		"compilertree.tar.*|compilertree.tar",
	}

	metaFilesRegex := []string{
		".*metadata.yaml",
	}

	pkgFilesRegex := []string{
		".*package.tar|.*package.tar.*",
	}

	for _, f := range files {
		if tmtools.RegexEntry(f, repoRegex) {
			DebugC(fmt.Sprintf("Ignoring repository file %s", f))
			continue
		}

		if c.Verbose {
			InfoC(fmt.Sprintf("[%s] Analyzing...", f))
		} else {
			DebugC(fmt.Sprintf("[%s] Analyzing...", f))
		}

		if tmtools.RegexEntry(f, metaFilesRegex) {

			art, err := c.BackendHandler.GetMetadata(f)
			if err != nil {
				return err
			}

			c.MetaMap[f] = art
		} else if tmtools.RegexEntry(f, pkgFilesRegex) {

			replaceRegex := regexp.MustCompile(
				`.package.tar$|.package.tar.gz$|.package.tar.zst$`,
			)

			metaFile := replaceRegex.ReplaceAllString(f, ".metadata.yaml")
			c.PkgsMap[f] = metaFile

		} else {
			// POST: file to remove
			c.Files2Remove = append(c.Files2Remove, f)
		}
	}

	// Check if there are all package for every metafile
	meta2Remove := []string{}
	for f, art := range c.MetaMap {
		pkg := filepath.Base(art.Path)

		if _, ok := c.PkgsMap[pkg]; !ok {
			if c.Verbose {
				InfoC(fmt.Sprintf(
					"No tarball found for metafile %s. I delete metafile.",
					f))
			} else {
				DebugC(fmt.Sprintf(
					"No tarball found for metafile %s. I delete metafile.",
					f))
			}
			c.Files2Remove = append(c.Files2Remove, f)
			meta2Remove = append(meta2Remove, f)
		}
	}

	for _, m := range meta2Remove {
		delete(c.PkgsMap, m)
	}

	// Check if there are all metadata for every package tarball
	for f, meta := range c.PkgsMap {
		if _, ok := c.MetaMap[meta]; !ok {
			if c.Verbose {
				InfoC(fmt.Sprintf(
					"No tarball file available for meta %s. I delete the tarball.",
					f))
			} else {
				DebugC(fmt.Sprintf(
					"No tarball file available for meta %s. I delete the tarball.",
					f))
			}
			c.Files2Remove = append(c.Files2Remove, f)
		}
	}

	//
	err = c.CheckFilesWithTrees()
	if err != nil {
		return err
	}

	return nil
}

func (c *RepoKnife) CheckFilesWithTrees() error {

	for m, art := range c.MetaMap {

		pkg := luet_pkg.NewPackage(
			art.CompileSpec.Package.Name,
			art.CompileSpec.Package.Version,
			[]*luet_pkg.DefaultPackage{},
			[]*luet_pkg.DefaultPackage{},
		)
		pkg.Category = art.CompileSpec.Package.Category

		p, _ := c.ReciperRuntime.GetDatabase().FindPackage(pkg)
		if p == nil {

			pkgFile := filepath.Base(art.Path)

			if c.Verbose {
				InfoC(fmt.Sprintf(
					"[%s] No more available in the repo. I will delete it.",
					pkg.HumanReadableString(),
				))
			} else {
				DebugC(fmt.Sprintf(
					"[%s] No more available in the repo. I will delete it.",
					pkg.HumanReadableString(),
				))
			}

			c.Files2Remove = append(c.Files2Remove, m)
			c.Files2Remove = append(c.Files2Remove, pkgFile)
		}

	}

	return nil
}

func (c *RepoKnife) GetFilteredList(files []string) ([]string, error) {
	ans := []string{}

	for _, f := range files {
		if !tmtools.RegexEntry(f, c.Specs.GetCleaner().Excludes) {
			ans = append(ans, f)
		}
	}

	return ans, nil
}
