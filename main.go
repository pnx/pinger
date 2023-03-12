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

// printStats prints statistics from a ping.
func printStats(pinger *probing.Pinger) {
	stats := pinger.Statistics()

	fmt.Printf("%d packets transmitted, %d packets received, %.2f%% packet loss\n",
		stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)

	fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
		stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
}

func eventLoop(pinger *probing.Pinger) {
	ticker := time.NewTicker(time.Second * 4)

	// Setup signals on term and interrupt.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		// Got signal. stop pinger and exit goroutine
		case <-sig:
			ticker.Stop()
			pinger.Stop()
			return
		// Ticker ticks. print stats.
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

	// Setup pinger.
	pinger, err := probing.NewPinger(os.Args[1])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	pinger.SetPrivileged(true)
	pinger.Interval = time.Second * 2

	// Enter event loop in another goroutine.
	go eventLoop(pinger)

	// Run pinger in main thread.
	fmt.Println("PING", pinger.Addr())
	if err = pinger.Run(); err != nil {
		fmt.Println("Error:", err)
		return
	}
}
