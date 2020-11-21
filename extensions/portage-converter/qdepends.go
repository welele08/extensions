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
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	helpers "github.com/MottainaiCI/lxd-compose/pkg/helpers"

	gentoo "github.com/Sabayon/pkgs-checker/pkg/gentoo"
)

type QDependsResolver struct{}

func NewQDependsResolver() *QDependsResolver {
	return &QDependsResolver{}
}

func retrieveVersion(solution *PortageSolution) error {
	var outBuffer, errBuffer bytes.Buffer

	cmd := []string{"qdepends", "-qC"}

	pkg := solution.Package.GetPackageName()
	if solution.Package.Slot != "0" {
		pkg = fmt.Sprintf("%s:%s", pkg, solution.Package.Slot)
	}
	cmd = append(cmd, pkg)

	qdepends := exec.Command(cmd[0], cmd[1:]...)
	qdepends.Stdout = helpers.NewNopCloseWriter(&outBuffer)
	qdepends.Stderr = helpers.NewNopCloseWriter(&errBuffer)

	err := qdepends.Start()
	if err != nil {
		return err
	}

	err = qdepends.Wait()
	if err != nil {
		return err
	}

	ans := qdepends.ProcessState.ExitCode()
	if ans != 0 {
		return errors.New("Error on running rdepends for package " + pkg + ": " + errBuffer.String())
	}

	out := outBuffer.String()
	if len(out) > 0 {

		words := strings.Split(out, " ")
		p := strings.TrimSuffix(words[0], "\n")

		gp, err := gentoo.ParsePackageStr(p[:len(p)-1])
		if err != nil {
			return errors.New("On convert pkg " + p + ": " + err.Error())
		}

		solution.Package.Version = gp.Version
		solution.Package.VersionSuffix = gp.VersionSuffix

	} else {
		return errors.New("No version found for package " + solution.Package.GetPackageName())
	}

	return nil
}

func SanitizeSlot(pkg *gentoo.GentooPackage) {
	if strings.Index(pkg.Slot, "/") > 0 {
		pkg.Slot = pkg.Slot[0:strings.Index(pkg.Slot, "/")]
	}

	if pkg.Slot == "*" {
		pkg.Slot = "0"
	}
}

func runQdepends(solution *PortageSolution, runtime bool) error {
	var outBuffer, errBuffer bytes.Buffer

	cmd := []string{"qdepends", "-qC", "-F", "deps"}

	if runtime {
		cmd = append(cmd, "-r")
	} else {
		cmd = append(cmd, "-bd")
	}

	pkg := solution.Package.GetPackageName()
	if solution.Package.Slot != "0" {
		pkg = fmt.Sprintf("%s:%s", pkg, solution.Package.Slot)
	}
	cmd = append(cmd, pkg)

	qdepends := exec.Command(cmd[0], cmd[1:]...)
	qdepends.Stdout = helpers.NewNopCloseWriter(&outBuffer)
	qdepends.Stderr = helpers.NewNopCloseWriter(&errBuffer)

	err := qdepends.Start()
	if err != nil {
		return err
	}

	err = qdepends.Wait()
	if err != nil {
		return err
	}

	ans := qdepends.ProcessState.ExitCode()
	if ans != 0 {
		return errors.New("Error on running rdepends for package " + pkg + ": " + errBuffer.String())
	}

	out := outBuffer.String()
	if len(out) > 0 {
		// Drop prefix
		out = out[6:]

		// Multiple match returns multiple rows. I get the first.
		rows := strings.Split(out, "\n")
		if len(rows) > 1 {
			out = rows[0]
		}

		deps := strings.Split(out, " ")

		for _, dep := range deps {

			originalDep := dep

			// Drop garbage string
			if len(dep) == 0 {
				continue
			}

			dep = strings.Trim(dep, "\n")
			dep = strings.Trim(dep, "\r")
			// Ignore ! / conflict for now ... not well supported by pkgs-checker now.
			if strings.Index(dep, "!") >= 0 {
				continue
			}

			if strings.Index(dep, ":") > 0 {

				depWithoutSlot := dep[0:strings.Index(dep, ":")]
				slot := dep[strings.Index(dep, ":")+1:]
				// i found slot but i want drop all subslot
				if strings.Index(slot, "/") > 0 {
					slot = slot[0:strings.Index(slot, "/")]
				}
				dep = depWithoutSlot + ":" + slot
			}

			// Ignoring use flags for now
			// >=dev-python/setuptools-42.0.2[python_targets_python3_7(+),-python_single_target_python3_6(+),-python_single_target_python3_7(+),-python_single_target_python3_8(+)]
			// it's not supported by pkgs-checker
			if strings.Index(dep, "[") > 0 {
				dep = dep[0:strings.Index(dep, "[")]
			}

			gp, err := gentoo.ParsePackageStr(dep)
			if err != nil {
				return errors.New("On convert dep " + dep + ": " + err.Error())
			}

			fmt.Println(fmt.Sprintf("[%s] Resolving dep '%s' -> %s...",
				solution.Package.GetPackageName(), originalDep,
				gp.GetPackageName()))
			SanitizeSlot(gp)
			if runtime {
				solution.RuntimeDeps = append(solution.RuntimeDeps, *gp)
			} else {
				solution.BuildDeps = append(solution.BuildDeps, *gp)
			}
		}

	} else {
		typeDeps := "build-time"
		if runtime {
			typeDeps = "runtime"
		}
		fmt.Println(fmt.Sprintf("No %s dependencies found for package %s.",
			typeDeps, solution.Package.GetPackageName()))
	}

	return nil
}

func (r *QDependsResolver) Resolve(pkg string) (*PortageSolution, error) {
	ans := &PortageSolution{
		BuildDeps:   make([]gentoo.GentooPackage, 0),
		RuntimeDeps: make([]gentoo.GentooPackage, 0),
	}

	gp, err := gentoo.ParsePackageStr(pkg)
	if err != nil {
		return nil, err
	}

	ans.Package = *gp

	// Retrive last version
	err = retrieveVersion(ans)
	if err != nil {
		// If with slot trying to use a package without slot
		if strings.Index(pkg, ":") > 0 {
			pkg = pkg[0:strings.Index(pkg, ":")]
			gp, err = gentoo.ParsePackageStr(pkg)
			if err != nil {
				return nil, err
			}

			ans.Package = *gp
			err = retrieveVersion(ans)
			if err != nil {
				return nil, err
			}

		} else {
			return nil, err
		}
	}

	// Retrieve runtime deps
	err = runQdepends(ans, true)
	if err != nil {
		return nil, err
	}

	// Retrieve build-time deps
	err = runQdepends(ans, false)
	if err != nil {
		return nil, err
	}

	// Sanitize slot
	SanitizeSlot(&ans.Package)

	return ans, nil
}
