/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"github.com/seanlan/goutils/xlconfig"
	"github.com/seanlan/xlvein/app/exchange"
	"github.com/seanlan/xlvein/app/router"
	"github.com/seanlan/xlvein/app/transport"
	"github.com/seanlan/xlvein/conf"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func startFunc(cmd *cobra.Command, args []string) {
	var (
		err error
		ctx = context.Background()
		_exchange exchange.Exchange
	)

	exchangeType := xlconfig.GetString("exchange", "type")
	switch exchangeType {
	case "local":
		_exchange, err = exchange.NewLocalExchange(zap.S())
	case "rabbitmq":
		_exchange, err = exchange.NewRabbitMQExchange(xlconfig.GetString("exchange", "rabbitmq"), zap.S())

	}
	if err != nil {
		zap.S().Fatal("new exchange error", zap.Error(err))
	}
	transport.InitHub(ctx, _exchange, zap.S())
	router.Setup(conf.DebugMode)
	router.Run(xlconfig.GetString("host"))
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: startFunc,
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
