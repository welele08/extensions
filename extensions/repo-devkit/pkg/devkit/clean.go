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
)

type RepoCleaner struct {
	*RepoKnife
	DryRun bool
}

func NewRepoCleaner(s *specs.LuetRDConfig,
	backend, path string, opts map[string]string,
	dryRun bool) (*RepoCleaner, error) {

	knife, err := NewRepoKnife(s, backend, path, opts)
	if err != nil {
		return nil, err
	}

	ans := &RepoCleaner{
		RepoKnife: knife,
		DryRun:    dryRun,
	}

	return ans, nil
}

func (c *RepoCleaner) Run() error {

	err := c.RepoKnife.Analyze()

	if len(c.Files2Remove) > 0 {
		for _, f := range c.Files2Remove {
			if c.DryRun {
				InfoC(fmt.Sprintf("[%s] Could be removed.", f))
			} else {
				err = c.BackendHandler.CleanFile(f)
				if err != nil {
					Error(fmt.Sprintf("[%s] Error on removing file: %s", f, err.Error()))
				} else {
					InfoC(fmt.Sprintf("[%s] Removed.", f))
				}
			}
		}
	} else {
		InfoC("No files to remove.")
	}

	return nil
}
