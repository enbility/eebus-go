package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/grandcat/zeroconf"
)

var myService *service.EEBUSService

func browseMDNS() {
	fmt.Println("Browsing for services...")

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatal(err)
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			fmt.Println("Service discovered: ", entry.ServiceInstanceName())

			fmt.Printf("Connecting to %s:%d\n", entry.HostName, entry.Port)
			myService.ConnectToService(entry.HostName, strconv.Itoa(int(entry.Port)))
			fmt.Printf("\n\n")
		}
		log.Println("No more entries.")
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err = resolver.Browse(ctx, "_ship._tcp", "local.", entries); err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
}

func usage() {
	fmt.Println("Usage: {} <command>", os.Args[0])
	fmt.Println("Commands:")
	fmt.Println("  server <serverport>")
	fmt.Println("  browse")
	fmt.Println("  connect <serverport> <host> <port>")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	certificate, err := service.CreateCertificate()
	if err != nil {
		log.Fatal(err)
	}

	myService = &service.EEBUSService{
		Certificate: certificate,
	}

	command := os.Args[1]
	switch command {
	case "server":
		if len(os.Args) < 3 {
			fmt.Println("Usage: {} server <serverport>", os.Args[0])
			return
		}
		myService.Port, _ = strconv.Atoi(os.Args[2])
		myService.Start()
	case "connect":
		if len(os.Args) < 4 {
			fmt.Println("Usage: {} connect <host> <port>", os.Args[0])
			return
		}
		myService.Start()
		host := os.Args[2]
		port := os.Args[3]
		_ = myService.ConnectToService(host, port)
	case "browse":
		browseMDNS()
	default:
		usage()
		return
	}

	select {}

}
