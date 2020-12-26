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
)

type RepoScanSpec struct {
	CacheDataVersion string                  `json:"cache_data_version"`
	Atoms            map[string]RepoScanAtom `json:"atoms"`
}

type RepoScanAtom struct {
	Atom string `json:"atom"`

	Category string   `json:"category"`
	Package  string   `json:"package"`
	Revision string   `json:"revision"`
	CatPkg   string   `json:"catpkg"`
	Eclasses []string `json:"eclasses"`

	Kit    string `json:"kit"`
	Branch string `json:"branch"`

	// Relations contains the list of the keys defined on
	// relations_by_kind. The values could be RDEPEND, DEPEND, BDEPEND
	Relations       []string            `json:"relations"`
	RelationsByKind map[string][]string `json:"relations_by_kind"`

	// Metadata contains ebuild variables.
	// Ex: SLOT, SRC_URI, HOMEPAGE, etc.
	Metadata    map[string]string `json:"metadata"`
	MetadataOut string            `json:"metadata_out"`

	ManifestMd5 string `json:"manifest_md5"`
	Md5         string `json:"md5"`

	Files []RepoScanFile `json:"files"`
}

type RepoScanFile struct {
	SrcUri []string          `json:"src_uri"`
	Size   string            `json:"size"`
	Hashes map[string]string `json:"hashes"`
	Name   string            `json:"name"`
}

type RepoScanResolver struct {
	JsonSources []string
}

func NewRepoScanResolver() *RepoScanResolver {
	return &RepoScanResolver{
		JsonSources: make([]string, 0),
	}
}

func (r *RepoScanResolver) Resolve(pkg string) (*PortageSolution, error) {
	return nil, errors.New("Not implemented")
}
