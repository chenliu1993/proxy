package main

import (
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

}
