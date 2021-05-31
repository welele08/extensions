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
	artifact "github.com/mudler/luet/pkg/compiler/types/artifact"
)

type LuetRDConfig struct {
	Cleaner LuetRDCCleaner `json:"cleaner,omitempty" yaml:"cleaner,omitempty"`
}

type LuetRDCCleaner struct {
	Excludes []string `json:"excludes,omitempty" yaml:"excludes,omitempty"`
}

type RepoBackendHandler interface {
	GetFilesList() ([]string, error)
	GetMetadata(string) (*artifact.PackageArtifact, error)
	CleanFile(string) error
}
