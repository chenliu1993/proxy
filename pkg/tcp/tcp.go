package tcp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/chenliu1993/proxy/pkg/dispatcher"
)

type TCPProxy struct {
	Src        *net.TCPConn
	Dst        *net.TCPConn
	Dispatcher *dispatcher.Dispatcher
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewTCPProxy(numOfProcs, bufferSize int) *TCPProxy {
	ctx, cancel := context.WithCancel(context.Background())
	return &TCPProxy{
		Dispatcher: dispatcher.NewDispatcher(numOfProcs, bufferSize),
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (tcp *TCPProxy) RegisterSrcTCP(addr string) error {
	TCPAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return err
	}

	TCPListener, err := net.ListenTCP("tcp", TCPAddr)
	if err != nil {
		return err
	}

	TCPConn, err := TCPListener.AcceptTCP()
	if err != nil {
		return err
	}

	TCPConn.SetKeepAlive(true)
	tcp.Src = TCPConn
	return nil
}

func (tcp *TCPProxy) RegisterDstTCP(addr string) error {
	TCPAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return err
	}

	TCPConn, err := net.DialTCP("tcp", nil, TCPAddr)
	if err != nil {
		return err
	}

	TCPConn.SetKeepAlive(true)
	tcp.Dst = TCPConn
	return nil
}

func (tcp *TCPProxy) RegisterHandler(handler func(interface{})) {
	tcp.Dispatcher.Register(handler)
}

func (tcp *TCPProxy) DefaultTCPHandler(data interface{}) {
	_, err := tcp.Dst.Write([]byte(data.(string)))
	if err != nil {
		fmt.Println("writing TCP errored:", err)
	}
}

func (tcp *TCPProxy) Run(errCh chan error) {
	tcp.Dispatcher.Run()

	for {
		select {
		case <-tcp.ctx.Done():
			return
		default:
			buf := make([]byte, 1024)
			var (
				n   int
				err error
			)
			n, err = tcp.Src.Read(buf)
			if err != nil {
				errCh <- errors.New(fmt.Sprintf("reading tcp  errored on %v:", err))
				time.Sleep(5 * time.Second)
				continue
			}

			if n == 0 {
				continue
			}

			tcp.Dispatcher.Put(string(buf[:n]))

		}
	}
}

func (tcp *TCPProxy) Stop() {
	tcp.cancel()
	tcp.Src.Close()
	tcp.Dst.Close()
	tcp.Dispatcher.Stop()
}
