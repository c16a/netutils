package cmd

import (
	"fmt"
	"github.com/c16a/netutils/lib"
	"github.com/spf13/cobra"
)

var TracerouteUrl string

func init() {
	traceRouteCommand.Flags().StringVarP(&TracerouteUrl, "url", "u", "", "set the URL")
	rootCmd.AddCommand(traceRouteCommand)
}

var traceRouteCommand = &cobra.Command{
	Use:   "traceroute",
	Short: "Perform traceroute operations",
	RunE: func(cmd *cobra.Command, args []string) error {
		options := lib.TracerouteOptions{}
		c := make(chan lib.TracerouteHop, 0)
		go func() {
			for {
				hop, ok := <-c
				if !ok {
					fmt.Println()
					return
				}
				printHop(hop)
			}
		}()

		_, err := lib.Traceroute(TracerouteUrl, &options, c)
		if err != nil {
			return err
		}

		return nil
	},
}

func printHop(hop lib.TracerouteHop) {
	addr := fmt.Sprintf("%v.%v.%v.%v", hop.Address[0], hop.Address[1], hop.Address[2], hop.Address[3])
	hostOrAddr := addr
	if hop.Host != "" {
		hostOrAddr = hop.Host
	}
	if hop.Success {
		fmt.Printf("%-3d %v (%v)  %v\n", hop.TTL, hostOrAddr, addr, hop.ElapsedTime)
	} else {
		fmt.Printf("%-3d *\n", hop.TTL)
	}
}
