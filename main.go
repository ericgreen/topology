package main

import (
	"flag"
	"fmt"
	"github.com/SpirentOrion/luddite"
	"os"
	"os/signal"
	"syscall"
)

var cfg = Config{}

var service luddite.Service

var shutdown = false

func Cleanup() {
	shutdown = true
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [-c topology.yaml]\n", os.Args[0])
}

func main() {
	var cfgFile string
	var err error

	fs := flag.NewFlagSet("topology", flag.ExitOnError)
	fs.StringVar(&cfgFile, "c", "topology.yaml", "Path to config file")
	fs.Usage = usage
	if err = fs.Parse(os.Args[1:]); err != nil {
		os.Exit(1)
	}

	if err = luddite.ReadConfig(cfgFile, &cfg); err != nil {
		panic(err)
	}

	service, err = luddite.NewService(&cfg.Service)
	if err != nil {
		panic(err)
	}

	InitApp(service.Router())

	go discover()

	go func() {
		service.Logger().Info("Starting to listen on " + cfg.Service.Addr)
		if err = service.Run(); err != nil {
			service.Logger().Error(err.Error())
		}

	}()

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	service.Logger().Info("System Shutting Down")
	Cleanup()

	os.Exit(0)
}
