package cmd

import (
	"fmt"
	"github.com/go-ping/ping"
	"github.com/spf13/cobra"
)

var PingUrl string
var PingCount int

func init() {
	pingCommand.Flags().StringVarP(&PingUrl, "url", "u", "", "set the URL")
	pingCommand.Flags().IntVarP(&PingCount, "count", "c", 5, "set the ping count")
	rootCmd.AddCommand(pingCommand)
}

var pingCommand = &cobra.Command{
	Use:   "ping",
	Short: "Perform PING operations",
	RunE: func(cmd *cobra.Command, args []string) error {
		pinger, err := ping.NewPinger(PingUrl)
		if err != nil {
			return err
		}

		pinger.Count = PingCount

		pinger.OnRecv = func(pkt *ping.Packet) {
			fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
		}

		pinger.OnDuplicateRecv = func(pkt *ping.Packet) {
			fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
		}

		pinger.OnFinish = func(stats *ping.Statistics) {
			fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
			fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
				stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
			fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
				stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
		}
		fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())

		err = pinger.Run() // Blocks until finished.
		if err != nil {
			return err
		}

		return nil
	},
}
