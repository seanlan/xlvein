/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/seanlan/xlvein/internal/common/transport"
	"github.com/seanlan/xlvein/pkg/veinsdk"
	"time"

	"github.com/spf13/cobra"
)

func clientFunc(cmd *cobra.Command, args []string) {
	fmt.Println("client called")
	sdk := veinsdk.New("http://127.0.0.1:8090", "test", "j8jasd98efan9sdfj89asjdf")
	token, _ := sdk.ProduceIMToken("user_1")
	uri := fmt.Sprintf("ws://127.0.0.1:8090/ws/connect?app_id=test&token=%s", token)
	client := veinsdk.NewClient(uri, func(msg transport.Message) {
		fmt.Println(msg)
	})
	go func() {
		time.Sleep(5 * time.Second)
		client.Send(transport.Message{
			To: "user_1",
			Data: map[string]interface{}{
				"text": "123123123123123",
			},
		})
	}()
	client.Run()
}

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: clientFunc,
}

func init() {
	rootCmd.AddCommand(clientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
