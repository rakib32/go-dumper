// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-dumper/dump"
	"github.com/go-dumper/load"
)

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "It will create backup file",
	Long:  `It will create backup file and then it will convert it to tar file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dump called")
		dumpAndLoad()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Inside rootCmd PostRun with args: %v\n", args)
		//TODO::
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dumpCmd.PersistentFlags().String("foo", "", "A help for foo")
	dumpCmd.PersistentFlags().Bool("skip-bucket-store", false, "Skip cloud bucket store")
	viper.BindPFlag("skip-bucket-store", dumpCmd.PersistentFlags().Lookup("skip-bucket-store"))

	dumpCmd.PersistentFlags().Bool("skip-restore", false, "Skip DB load to destination")
	viper.BindPFlag("skip-restore", dumpCmd.PersistentFlags().Lookup("skip-restore"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dumpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func dumpAndLoad() {
	result, err := dump.Run()

	if err != nil {
		logger.Error(err)
	} else {

		logger.Info("------------------Dump success!!!-------------------")

		isSkipRestore := viper.GetBool("skip-restore")

		if !isSkipRestore {
			err := load.Run(result.Path)

			if err != nil {
				logger.Error(err)
			} else {
				logger.Info("------------------Data Loaded Successfully!!!-------------------")
			}
		}

		deleteDumpFiles := viper.GetBool("delete-dump-files")

		if deleteDumpFiles {
			logger.Info("----------------Removing Local backup sql file---------------")
			os.Remove(result.Path)
		}

		logger.Info(result.Path)
	}
}
