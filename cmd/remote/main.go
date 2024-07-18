package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/usecases/cem/evcc"
	"github.com/enbility/eebus-go/usecases/cem/evsecc"
	eglpc "github.com/enbility/eebus-go/usecases/eg/lpc"
	"github.com/enbility/eebus-go/usecases/ma/mpc"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/cert"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type eebusConfiguration struct {
	vendorCode   string
	deviceBrand  string
	deviceModel  string
	serialNumber string
}

func loadCertificate(config eebusConfiguration, crtPath, keyPath string) tls.Certificate {
	certificate, err := tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		certificate, err = cert.CreateCertificate(config.vendorCode, config.deviceModel, "DE", config.serialNumber)
		if err != nil {
			log.Fatal(err)
		}

		if err = WriteKey(certificate, keyPath); err != nil {
			log.Fatal(err)
		}
		if err = WriteCertificate(certificate, crtPath); err != nil {
			log.Fatal(err)
		}
	}

	return certificate
}

func main() {
	config := eebusConfiguration{}

	iface := flag.String("iface", "",
		"Optional network interface the EEBUS connection should be limited to")
	flag.StringVar(&config.vendorCode, "vendor", "", "EEBus vendor code")
	flag.StringVar(&config.deviceBrand, "brand", "", "EEBus device brand")
	flag.StringVar(&config.deviceModel, "model", "", "EEBus device model")
	flag.StringVar(&config.serialNumber, "serial", "", "EEBus device serial")

	flag.Parse()

	if config.serialNumber == "" {
		serialNumber, err := os.Hostname()
		if err != nil {
			log.Fatal(err)
		}
		config.serialNumber = serialNumber
	}

	if config.vendorCode == "" || config.deviceBrand == "" || config.deviceModel == "" {
		flag.Usage()
		return
	}

	certificate := loadCertificate(config, "cert.pem", "key.pem")

	configuration, err := api.NewConfiguration(
		config.vendorCode, config.deviceBrand, config.deviceModel, config.serialNumber,
		[]shipapi.DeviceCategoryType{
			shipapi.DeviceCategoryTypeEnergyManagementSystem,
		},
		model.DeviceTypeTypeEnergyManagementSystem,
		[]model.EntityTypeType{
			model.EntityTypeTypeGridGuard,
			model.EntityTypeTypeCEM,
		},
		23292, certificate, time.Second*4)
	if *iface != "" {
		configuration.SetInterfaces([]string{*iface})
		log.Printf("waiting until %v is up", iface)
		for {
			ifi, err := net.InterfaceByName(*iface)
			if err != nil {
				log.Fatal(err)
			}

			// wait until interface is up and available for multicast
			flags := net.FlagUp | net.FlagMulticast
			if (ifi.Flags & flags) == flags {
				break
			}
			time.Sleep(1 * time.Second)
		}
		log.Printf("interface online, continuing")
	}

	r, err := NewRemote(configuration)
	if err != nil {
		log.Fatal(err)
	}

	err = r.RegisterUseCase(model.EntityTypeTypeCEM, "EG-LPC", func(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) api.UseCaseInterface {
		return eglpc.NewLPC(localEntity, eventCB)
	})
	if err != nil {
		log.Fatal(err)
	}

	err = r.RegisterUseCase(model.EntityTypeTypeCEM, "MA-MPC", func(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) api.UseCaseInterface {
		return mpc.NewMPC(localEntity, eventCB)
	})
	if err != nil {
		log.Fatal(err)
	}

	err = r.RegisterUseCase(model.EntityTypeTypeCEM, "CEM-EVCC", func(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) api.UseCaseInterface {
		return evcc.NewEVCC(r.service, localEntity, eventCB)
	})
	if err != nil {
		log.Fatal(err)
	}

	err = r.RegisterUseCase(model.EntityTypeTypeCEM, "CEM-EVSECC", func(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) api.UseCaseInterface {
		return evsecc.NewEVSECC(localEntity, eventCB)
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancelCtx := context.WithCancel(context.Background())
	if err = r.Listen(ctx, "tcp", net.JoinHostPort("::", strconv.Itoa(3393))); err != nil {
		log.Fatal(err)
	}
	log.Print("Started")

	// Clean exit to make sure mdns shutdown is invoked
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	// User exit

	cancelCtx()
}
