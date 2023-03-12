package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pborman/getopt/v2"
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
	// Cli options.
	help := getopt.BoolLong("help", 'h', "Show this help text")
	version := getopt.BoolLong("version", 'v', "Show program version")
	source := getopt.StringLong("source", 's', "", "Source IP address to send packages from.")
	count := getopt.IntLong("count", 'c', 0, "Stop after this many packages has been sent (and received). If this option is not specified, pinger will operate until interrupted.")
	timeout := getopt.DurationLong("timeout", 't', 0, "Exit the program after this time is reached, regardless of how many packets have been received.")
	interval := getopt.DurationLong("interval", 'i', time.Second*2, "Wait time between each packet send")
	proto_udp := getopt.BoolLong("udp", 0, "Send UDP ping instead of a raw IMCP ping, IMCP required super-user privileges")
	record_rtt := getopt.BoolLong("rtt", 0, "Keep a record of rtts of all received packets.")

	getopt.SetParameters("<ip>")
	getopt.Parse()

	if *version {
		fmt.Println("Version 0.0.1")
		os.Exit(0)
	}

	if *help || len(getopt.Args()) < 1 {
		getopt.Usage()
		os.Exit(1)
	}

	// Setup pinger.
	pinger, err := probing.NewPinger(getopt.Arg(0))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	pinger.SetPrivileged(!*proto_udp)
	pinger.Interval = *interval
	pinger.Source = *source
	pinger.Count = *count
	pinger.RecordRtts = *record_rtt
	if timeout != nil && *timeout > 0 {
		pinger.Timeout = *timeout
	}

	// Enter event loop in another goroutine.
	go eventLoop(pinger)

	// Run pinger in main thread.
	fmt.Println("PING", pinger.Addr())
	if err = pinger.Run(); err != nil {
		fmt.Println("Error:", err)
		return
	}
}
