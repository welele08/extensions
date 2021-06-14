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

package main

import (
	"fmt"
	"os"
	"strings"

	devkitcmd "github.com/Luet-lab/extensions/extensions/repo-devkit/pkg/cmd"
	"github.com/Luet-lab/extensions/extensions/repo-devkit/pkg/devkit"

	. "github.com/mudler/luet/pkg/config"
	. "github.com/mudler/luet/pkg/logger"
	"github.com/spf13/cobra"
)

const (
	cliName = `Copyright (c) 2020-2021 - Daniele Rondina

Luet Repository Devkit`
)

func initConfig() error {
	LuetCfg.Viper.SetEnvPrefix("LUET")
	LuetCfg.Viper.AutomaticEnv() // read in environment variables that match

	// Create EnvKey Replacer for handle complex structure
	replacer := strings.NewReplacer(".", "__")
	LuetCfg.Viper.SetEnvKeyReplacer(replacer)
	LuetCfg.Viper.SetTypeByDefaultValue(true)

	err := LuetCfg.Viper.Unmarshal(&LuetCfg)
	if err != nil {
		return err
	}

	InitAurora()
	NewSpinner()

	return nil
}

func Execute() {
	var rootCmd = &cobra.Command{
		Use:   "luet-repo-devkit --",
		Short: cliName,
		Version: fmt.Sprintf("%s-g%s %s",
			devkit.Version, devkit.BuildCommit, devkit.BuildTime,
		),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			err := initConfig()
			if err != nil {
				fmt.Println("Error on setup luet config/logger: " + err.Error())
				os.Exit(1)
			}

			debug, _ := cmd.Flags().GetBool("debug")
			if debug {
				LuetCfg.GetGeneral().Debug = true
			}

		},
	}

	rootCmd.PersistentFlags().StringArrayP("tree", "t", []string{}, "Path of the tree to use.")
	rootCmd.PersistentFlags().StringP("specs-file", "s", "", "Path of the devkit specification file.")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug logging.")

	rootCmd.AddCommand(
		devkitcmd.NewCleanCommand(),
		devkitcmd.NewPkgsCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
