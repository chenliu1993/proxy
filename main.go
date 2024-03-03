package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof"

	"github.com/chenliu1993/proxy/pkg/config"
	"github.com/chenliu1993/proxy/pkg/sandbox"
	"github.com/spf13/pflag"
)

var (
	configFile string
)

func init() {
	pflag.StringVarP(&configFile, "config", "c", "config.yaml", "path to config file<yaml>")
}

func main() {
	pflag.Parse()
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	stopCh := make(chan struct{}, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT)

	sandbox := sandbox.NewSandbox()

	if configFile == "" {
		fmt.Println("config file not specified")
		os.Exit(1)
	}

	configContent, err := config.ParseConfigFile(configFile)
	if err != nil {
		fmt.Printf("parsing config file errored: %v\n", err)
		os.Exit(1)
	}
	sandbox.Config = configContent

	err = sandbox.StartUDPConns()
	if err != nil {
		fmt.Printf("starting udp conns errored: %v\n", err)
		os.Exit(1)
	}

	err = sandbox.StartTCPConns()
	if err != nil {
		fmt.Printf("starting tcp conns errored: %v\n", err)
		os.Exit(1)
	}

	go func(chan struct{}) {
		for {
			select {
			case err := <-sandbox.ErrCh:
				fmt.Fprintf(os.Stdout, "error: %v\n", err)
			case <-stopCh:
				return
			}
		}
	}(stopCh)

	fmt.Println("Starting...")
	<-sigCh
	fmt.Println("Stopping all conns...")
	stopCh <- struct{}{}
	sandbox.StopTCPConns()
	sandbox.StopUDPConns()
}
