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

package cmd

import (
	"fmt"
	"os"

	devkit "github.com/Luet-lab/extensions/extensions/repo-devkit/pkg/devkit"
	specs "github.com/Luet-lab/extensions/extensions/repo-devkit/pkg/specs"

	cobra "github.com/spf13/cobra"
)

func NewCleanCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "clean [OPTIONS]",
		Short: "Clean repository files.",
		PreRun: func(cmd *cobra.Command, args []string) {
			treePath, _ := cmd.Flags().GetStringArray("tree")

			if len(treePath) == 0 {
				fmt.Println("At least one tree path is needed.")
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			var s *specs.LuetRDConfig
			var err error

			specsFile, _ := cmd.Flags().GetString("specs-file")
			backend, _ := cmd.Flags().GetString("backend")
			path, _ := cmd.Flags().GetString("path")
			treePath, _ := cmd.Flags().GetStringArray("tree")
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			mottainaiProfile, _ := cmd.Flags().GetString("mottainai-profile")
			mottainaiMaster, _ := cmd.Flags().GetString("mottainai-master")
			mottainaiApiKey, _ := cmd.Flags().GetString("mottainai-apikey")
			mottainaiNamespace, _ := cmd.Flags().GetString("mottainai-namespace")

			if specsFile == "" {
				s = specs.NewLuetRDConfig()
			} else {
				s, err = specs.LoadSpecsFile(specsFile)
				if err != nil {
					fmt.Println("Error on load specs: " + err.Error())
					os.Exit(1)
				}
			}

			opts := make(map[string]string, 0)
			if backend == "mottainai" {
				if mottainaiProfile != "" {
					opts["mottainai-profile"] = mottainaiProfile
				}
				if mottainaiMaster != "" {
					opts["mottainai-master"] = mottainaiMaster
				}
				if mottainaiApiKey != "" {
					opts["mottainai-apikey"] = mottainaiApiKey
				}
				if mottainaiNamespace != "" {
					opts["mottainai-namespace"] = mottainaiNamespace
				}
			}

			repoCleaner, err := devkit.NewRepoCleaner(s, backend, path, opts, dryRun)
			if err != nil {
				fmt.Println("Error on initialize repo cleaner: " + err.Error())
				os.Exit(1)
			}

			// Loading tree in memory
			err = repoCleaner.LoadTrees(treePath)
			if err != nil {
				fmt.Println("Erro on loading trees: " + err.Error())
				os.Exit(1)
			}

			err = repoCleaner.Run()
			if err != nil {
				fmt.Println("Error on clean repository: " + err.Error())
				os.Exit(1)
			}

			fmt.Println("All done.")
		},
	}

	var flags = cmd.Flags()
	flags.StringP("backend", "b", "local", "Select backend repository: local|mottainai.")
	flags.StringP("path", "p", "", "Path of the repository artefacts.")
	flags.Bool("dry-run", false, "Only check files to remove.")
	flags.String("mottainai-profile", "", "Set mottainai profile to use.")
	flags.String("mottainai-master", "", "Set mottainai Server to use.")
	flags.String("mottainai-apikey", "", "Set mottainai API Key to use.")
	flags.String("mottainai-namespace", "", "Set mottainai namespace to use.")

	return cmd
}
