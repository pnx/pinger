package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/prometheus-community/pro-bing"
)

func printStats(pinger *probing.Pinger) {
	stats := pinger.Statistics()
	fmt.Printf("%d packets transmitted, %d packets received, %.2f%% packet loss\n",
		stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
	fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
		stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
}

func eventLoop(pinger *probing.Pinger) {
	ticker := time.NewTicker(time.Second * 4)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sig:
			pinger.Stop()
			return
		case <-ticker.C:
			printStats(pinger)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <ip>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	pinger, err := probing.NewPinger(os.Args[1])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	pinger.SetPrivileged(true)
	pinger.Interval = time.Second * 2

	go eventLoop(pinger)

	fmt.Println("PING", pinger.Addr())
	if err = pinger.Run(); err != nil {
		fmt.Println("Error:", err)
		return
	}
}
