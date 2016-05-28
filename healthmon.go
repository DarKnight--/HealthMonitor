package main

import (
	"flag"
	"sync"

	_ "health_monitor/api"
	_ "health_monitor/config"
	"health_monitor/disk"
	"health_monitor/live"
	"health_monitor/utils"
	"health_monitor/webui"
)

// Flags holds the health_monitor command line arguments
type Flags struct {
	NoWebUI *bool
	NoCLI   *bool
	Quite   *bool
}

func main() {
	var (
		wg    sync.WaitGroup
		flags Flags
	)

	flags.NoWebUI = flag.Bool("nowebui", false, "Disables the web ui")
	flags.NoCLI = flag.Bool("nocli", false, "Disables cli")
	flags.Quite = flag.Bool("quite", false, "Disables all notifications except email")

	flag.Parse()

	go webui.RunServer("8009")

	signal := make(chan utils.Status)
	wg.Add(1)
	go live.Live(signal, &wg)
	wg.Add(1)
	go disk.Disk(signal, &wg)
	wg.Wait()
}
