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
	"github.com/go-dumper/load"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "It will load the dump file to destination DB",
	Long:  `It will load the dump file to destination DB.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("load called")
		loadBackup()
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loadCmd.PersistentFlags().String("foo", "", "A help for foo")
	loadCmd.PersistentFlags().String("dump-path", "", "Dump path")
	loadCmd.MarkPersistentFlagRequired("dump-path")
	viper.BindPFlag("dump-path", loadCmd.PersistentFlags().Lookup("dump-path"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func loadBackup() {
	logger.Info("------Loading backup-----")
	logger.Info(viper.GetString("dump-path"))
	err := load.Run(viper.GetString("dump-path"))

	if err != nil {
		logger.Error(err)
	} else {
		logger.Info("Load success!!!")
	}
}
