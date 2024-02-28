package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/chenliu1993/proxy/pkg/dispatcher"
	"github.com/spf13/pflag"
)

var (
	mode string
)

func init() {
	pflag.StringVarP(&mode, "mode", "m", "tcp", "type of packets to transfer")
}

func main() {
	pflag.Parse()
	stopCh := make(chan struct{}, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT)

	dispatcher := dispatcher.NewDispatcher(10, 10)

	print := func(data interface{}) {
		fmt.Println(data)
	}
	dispatcher.Register(print)

	// TODO: currently order is not ensured
	dispatcher.Run()

	data := []string{"hello", "world", "my", "name", "is", "lc"}
	for i := 0; i < len(data); i++ {
		dispatcher.Put(data[i])
	}

	fmt.Println("Starting...")
	<-sigCh
	fmt.Println("Stopping...")
	stopCh <- struct{}{}
	dispatcher.Stop()
}
