package cmd

import (
	"fmt"
	"github.com/c16a/netutils/utils"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"log"
)

var WsUrl string
var WsTimeout int
var WsHeaders map[string]string

func init() {
	websocketCommand.Flags().StringVarP(&WsUrl, "url", "u", "", "set the URL for this operation")
	websocketCommand.Flags().IntVarP(&WsTimeout, "timeout", "t", 5000, "set the timeout in milliseconds for this operation")
	websocketCommand.Flags().StringToStringVarP(&WsHeaders, "headers", "", nil, "set the headers for this operation")

	rootCmd.AddCommand(websocketCommand)
}

var websocketCommand = &cobra.Command{
	Use:   "websocket",
	Short: "Perform Websocket operations",
	RunE: func(cmd *cobra.Command, args []string) error {
		headers := utils.ToMultiValuedMap(WsHeaders)
		c, _, err := websocket.DefaultDialer.Dial(WsUrl, headers)
		if err != nil {
			return err
		}
		log.Println("Bi-directional WS connection is now open")
		log.Println("Enter messages to send to server and press ENTER key.")
		fmt.Println("----------------------")
		defer c.Close()

		go func() {
			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					return
				}
				log.Printf("SERVER> %s", message)
			}
		}()

		for {
			var message string
			fmt.Scanln(&message)

			c.WriteMessage(websocket.TextMessage, []byte(message))
		}
	},
}
