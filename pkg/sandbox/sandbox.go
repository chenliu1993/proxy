package sandbox

import (
	"fmt"

	"github.com/chenliu1993/proxy/pkg/config"
	"github.com/chenliu1993/proxy/pkg/tcp"
	"github.com/chenliu1993/proxy/pkg/udp"
)

type Sandbox struct {
	UDPConns map[string]*udp.UDPProxy
	TCPConns map[string]*tcp.TCPProxy
	Config   *config.Config
	ErrCh    chan error
}

func NewSandbox() *Sandbox {
	return &Sandbox{
		UDPConns: map[string]*udp.UDPProxy{},
		TCPConns: map[string]*tcp.TCPProxy{},
		Config:   config.NewConfig(),
		ErrCh:    make(chan error, 1),
	}
}

func (s *Sandbox) StartUDPConns() error {
	if len(s.Config.UDPConns) == 0 {
		fmt.Println("no udp specified")
		return nil
	}
	for _, udpConfig := range s.Config.UDPConns {
		udpProxy := udp.NewUDPProxy(s.Config.NumOfProcs, s.Config.BufSizePerProc)
		err := udpProxy.RegisterSrcUDP(udpConfig.SrcAddr)
		if err != nil {
			return err
		}
		err = udpProxy.RegisterDstUDP(udpConfig.DstAddr)
		if err != nil {
			return err
		}

		udpProxy.RegisterHandler(udpProxy.DefaultUDPHandler)

		go udpProxy.Run(s.ErrCh)
		s.UDPConns[fmt.Sprintf("%sTo%s", udpConfig.SrcAddr, udpConfig.DstAddr)] = udpProxy
		fmt.Println("Starting UDPConn ", fmt.Sprintf("%sTo%s", udpConfig.SrcAddr, udpConfig.DstAddr))

	}
	return nil
}

func (s *Sandbox) StopUDPConns() {
	for k, v := range s.UDPConns {
		fmt.Printf("Stopping UDPConn %s\n", k)
		v.Stop()
	}
}

func (s *Sandbox) StartTCPConns() error {
	if len(s.Config.TCPConns) == 0 {
		fmt.Println("no tcp specified")
		return nil
	}

	for _, tcpConfig := range s.Config.TCPConns {
		tcpProxy := tcp.NewTCPProxy(s.Config.NumOfProcs, s.Config.BufSizePerProc)
		err := tcpProxy.RegisterSrcTCP(tcpConfig.SrcAddr)
		if err != nil {
			return err
		}
		err = tcpProxy.RegisterDstTCP(tcpConfig.DstAddr)
		if err != nil {
			return err
		}

		tcpProxy.RegisterHandler(tcpProxy.DefaultTCPHandler)

		go tcpProxy.Run(s.ErrCh)
		s.TCPConns[fmt.Sprintf("%sTo%s", tcpConfig.SrcAddr, tcpConfig.DstAddr)] = tcpProxy
	}
	return nil
}

func (s *Sandbox) StopTCPConns() {
	for k, v := range s.TCPConns {
		fmt.Printf("Stopping TCPConn %s\n", k)
		v.Stop()
	}
}
