package mdnsservice

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grandcat/zeroconf"
)

var (
	name    = flag.String("name", "Clipboard Sync", "The name for the mDNS service")
	service = flag.String("service", "_clipboardsync._tcp", "Sets the service category to look for device")
	port    = flag.Int("port", 40222, "Set the port the service is listening to")
)

func SearchForDevices() ([]*zeroconf.ServiceEntry, error) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}

	devicesDiscovered := []*zeroconf.ServiceEntry{}
	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			devicesDiscovered = append(devicesDiscovered, entry)
		}
		log.Println("No more discoverable devices available.")
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(10))
	defer cancel()
	err = resolver.Browse(ctx, "_clipboardsync._tcp", "local", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
		return nil, err
	}

	<-ctx.Done()

	return devicesDiscovered, nil
}

func RegisterDeviceForDiscovery() {
	hostname, _ := os.Hostname()
	server, err := zeroconf.Register(
		*name,
		*service,
		"local.",
		*port,
		[]string{
			"version=1",
			"device=desktop",
			"hostname=" + hostname,
		},
		nil,
	)
	if err != nil {
		panic(err)
	}
	defer server.Shutdown()

	// Clean exit.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	select {
	case <-sig:
		// Exit by user
	}

	log.Println("Shutting down.")
}

func main() {
	flag.Parse()
}
