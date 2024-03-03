package udp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/chenliu1993/proxy/pkg/dispatcher"
)

type UDPProxy struct {
	Src        *net.UDPConn
	Dst        *net.UDPConn
	Dispatcher *dispatcher.Dispatcher
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewUDPProxy(numOfProcs, bufferSize int) *UDPProxy {
	ctx, cancel := context.WithCancel(context.Background())
	return &UDPProxy{
		Dispatcher: dispatcher.NewDispatcher(numOfProcs, bufferSize),
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (udp *UDPProxy) RegisterSrcUDP(addr string) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}

	udp.Src = udpConn
	return nil
}

func (udp *UDPProxy) RegisterDstUDP(addr string) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}

	udp.Dst = udpConn
	return nil
}

func (udp *UDPProxy) RegisterHandler(handler func(interface{})) {
	udp.Dispatcher.Register(handler)
}

func (udp *UDPProxy) DefaultUDPHandler(data interface{}) {
	_, err := udp.Dst.Write([]byte(data.(string)))
	if err != nil {
		fmt.Println("writing udp errored:", err)
	}
}

func (udp *UDPProxy) Run(errCh chan error) {
	udp.Dispatcher.Run()

	for {
		select {
		case <-udp.ctx.Done():
			return
		default:
			buf := make([]byte, 1024)
			var (
				n    int
				addr *net.UDPAddr
				err  error
			)

			n, addr, err = udp.Src.ReadFromUDP(buf)

			if err != nil {
				errCh <- errors.New(fmt.Sprintf("reading udp %v errored on %v:", addr, err))
				time.Sleep(5 * time.Second)
				continue
			}

			if n == 0 {
				continue
			}

			udp.Dispatcher.Put(string(buf[:n]))
		}
	}
}

func (udp *UDPProxy) Stop() {
	udp.cancel()
	udp.Src.Close()
	udp.Dst.Close()
	udp.Dispatcher.Stop()
}
