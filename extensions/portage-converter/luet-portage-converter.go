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
	"encoding/json"
	"fmt"
	"os"
	"strings"

	gentoo "github.com/Sabayon/pkgs-checker/pkg/gentoo"
	luet_pkg "github.com/mudler/luet/pkg/package"
	"github.com/spf13/cobra"
)

const (
	cliName = `Copyright (c) 2020 - Daniele Rondina

Portage/Overlay converter for Luet specs.`

	version = "0.1.0"
)

// Build time and commit information. This code is get from: https://github.com/mudler/luet/
//
// ⚠️ WARNING: should only be set by "-ldflags".
var (
	BuildTime   string
	BuildCommit string
)

type PortageResolver interface {
	Resolve(pkg string) (*PortageSolution, error)
}

type PortageSolution struct {
	Package          gentoo.GentooPackage   `json:"package"`
	PackageDir       string                 `json:"package_dir"`
	BuildDeps        []gentoo.GentooPackage `json:"build-deps,omitempty"`
	RuntimeDeps      []gentoo.GentooPackage `json:"runtime-deps,omitempty"`
	RuntimeConflicts []gentoo.GentooPackage `json:"runtime-conflicts,omitempty"`
	BuildConflicts   []gentoo.GentooPackage `json:"build-conflicts,omitempty"`
}

func SanitizeCategory(cat string, slot string) string {
	ans := cat
	if slot != "0" {
		// Ignore sub-slot
		if strings.Contains(slot, "/") {
			slot = slot[0:strings.Index(slot, "/")]
		}

		if slot != "0" {
			ans = fmt.Sprintf("%s-%s", cat, slot)
		}
	}
	return ans
}

func (s *PortageSolution) ToPack(runtime bool) *luet_pkg.DefaultPackage {

	// TODO: handle particular use cases
	version := fmt.Sprintf("%s%s", s.Package.Version, s.Package.VersionSuffix)

	labels := make(map[string]string, 0)
	labels["original.package.name"] = s.Package.GetPackageName()
	labels["original.package.version"] = s.Package.GetPVR()

	ans := &luet_pkg.DefaultPackage{
		Name:     s.Package.Name,
		Category: SanitizeCategory(s.Package.Category, s.Package.Slot),
		Version:  version,
		UseFlags: s.Package.UseFlags,
		Labels:   labels,
	}

	deps := s.BuildDeps
	if runtime {
		deps = s.RuntimeDeps
	}

	for _, req := range deps {

		dep := &luet_pkg.DefaultPackage{
			Name:     req.Name,
			Category: SanitizeCategory(req.Category, req.Slot),
			UseFlags: req.UseFlags,
		}
		if req.Version != "" && req.Condition != gentoo.PkgCondNot &&
			req.Condition != gentoo.PkgCondAnyRevision &&
			req.Condition != gentoo.PkgCondMatchVersion {

			// TODO: to complete
			dep.Version = fmt.Sprintf("%s%s%s",
				req.Condition.String(), req.Version, req.VersionSuffix)

		} else {
			dep.Version = ">=0"
		}

		ans.PackageRequires = append(ans.PackageRequires, dep)
	}

	if runtime && len(s.RuntimeConflicts) > 0 {

		for _, req := range s.RuntimeConflicts {

			dep := &luet_pkg.DefaultPackage{
				Name:     req.Name,
				Category: SanitizeCategory(req.Category, req.Slot),
				UseFlags: req.UseFlags,
			}
			if req.Version != "" && req.Condition == gentoo.PkgCondNot {
				// TODO: to complete
				dep.Version = fmt.Sprintf("%s%s%s",
					req.Condition.String(), req.Version, req.VersionSuffix)

			} else {
				dep.Version = ">=0"
			}

			ans.PackageConflicts = append(ans.PackageConflicts, dep)
		}

	} else if !runtime && len(s.BuildConflicts) > 0 {

		for _, req := range s.BuildConflicts {

			dep := &luet_pkg.DefaultPackage{
				Name:     req.Name,
				Category: SanitizeCategory(req.Category, req.Slot),
				UseFlags: req.UseFlags,
			}
			if req.Version != "" && req.Condition == gentoo.PkgCondNot {
				// TODO: to complete
				dep.Version = fmt.Sprintf("%s%s%s",
					req.Condition.String(), req.Version, req.VersionSuffix)

			} else {
				dep.Version = ">=0"
			}

			ans.PackageConflicts = append(ans.PackageConflicts, dep)
		}

	}

	return ans
}

func (s *PortageSolution) String() string {
	data, _ := json.Marshal(*s)
	return string(data)
}

func Execute() {
	var rootCmd = &cobra.Command{
		Use:     "luet-portage-converter --",
		Short:   cliName,
		Version: fmt.Sprintf("%s-g%s %s", version, BuildCommit, BuildTime),
		PreRun: func(cmd *cobra.Command, args []string) {
			to, _ := cmd.Flags().GetString("to")
			if to == "" {
				fmt.Println("Missing --to argument")
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			treePath, _ := cmd.Flags().GetStringArray("tree")
			to, _ := cmd.Flags().GetString("to")
			rulesFile, _ := cmd.Flags().GetString("rules")
			override, _ := cmd.Flags().GetBool("override")

			converter := NewPortageConverter(to)
			converter.Override = override
			err := converter.LoadRules(rulesFile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = converter.LoadTrees(treePath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = converter.Generate()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		},
	}

	rootCmd.Flags().StringArrayP("tree", "t", []string{}, "Path of the tree to use.")
	rootCmd.Flags().String("to", "", "Targer tree where bump new specs.")
	rootCmd.Flags().String("rules", "", "Rules file.")
	rootCmd.Flags().Bool("override", false, "Override existing specs if already present.")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
