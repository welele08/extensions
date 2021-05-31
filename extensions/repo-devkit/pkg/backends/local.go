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
package backends

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Luet-lab/extensions/extensions/repo-devkit/pkg/specs"
	artifact "github.com/mudler/luet/pkg/compiler/types/artifact"
	. "github.com/mudler/luet/pkg/logger"
)

type BackendLocal struct {
	Specs *specs.LuetRDConfig
	Path  string
}

func NewBackendLocal(specs *specs.LuetRDConfig, path string) (*BackendLocal, error) {
	if path == "" {
		return nil, errors.New("Invalid path")
	}

	_, err := os.Stat(path)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf(
				"Error on retrieve stat of the path %s: %s",
				path, err.Error(),
			))
	}

	if os.IsNotExist(err) {
		return nil, errors.New("The path doesn't exist!")
	}

	ans := &BackendLocal{
		Path: path,
	}

	return ans, nil
}

func (b *BackendLocal) GetFilesList() ([]string, error) {
	ans := []string{}

	files, err := ioutil.ReadDir(b.Path)
	if err != nil {
		return ans, err
	}

	for _, f := range files {
		DebugC("Cheking file ", f.Name())
		if f.IsDir() {
			// Ignoring directories at the moment.
			DebugC(fmt.Sprintf("Ignoring directory %s", f.Name()))
			continue
		}

		ans = append(ans, f.Name())
	}

	return ans, nil
}

func (b *BackendLocal) GetMetadata(file string) (*artifact.PackageArtifact, error) {
	metafile := filepath.Join(b.Path, file)
	content, err := ioutil.ReadFile(metafile)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("Error on open file %s", metafile))
	}

	return artifact.NewPackageArtifactFromYaml(content)
}

func (b *BackendLocal) CleanFile(file string) error {
	absFile := filepath.Join(b.Path, file)
	return os.Remove(absFile)
}
