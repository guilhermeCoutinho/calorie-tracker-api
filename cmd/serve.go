// Copyright © 2021 NAME HERE <EMAIL ADDRESS>
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
	"github.com/guilhermeCoutinho/api-studies/dal"
	"github.com/guilhermeCoutinho/api-studies/server/http"
	"github.com/guilhermeCoutinho/api-studies/usecase"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		startHTTPServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func startHTTPServer() {
	logger := getLogger(config)

	db, err := connectToPG("db")
	if err != nil {
		logger.Fatalln(err.Error())
	}

	dal := dal.NewUser(config, db)
	usecase := usecase.NewUsecase(config, dal, logger)

	app, err := http.NewApp(config, logger, usecase)
	if err != nil {
		logger.Fatalln(err.Error())
	}
	app.ListenAndServe()
}
