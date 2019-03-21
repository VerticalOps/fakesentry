/*
Package ipc implements a simple set of intra-process-communication net.Conn types, specifically a net.Listener
and an associated Dialer that can be used in place of http Servers/Clients or anything that needs a net.Conn to work
or test within the same process. This package is built around net.Pipe and so the network is not actually involved.
If a buffered implementation is needed there are alternatives.
*/
package ipc

import (
	"context"
	"errors"
	"net"
)

//A Dailer can be used in place of many Dial-related functions of the net package and others.
//The returned net.Conn is always related to the Listener that returned the Dialer via Listener.NewDialer.
//Therefore the passed network and address are ignored. The given context is not however, if the context
//expires before the connection is made then ctx.Err is returned.
type Dailer interface {
	Dial(network, address string) (net.Conn, error)
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

type dialer struct {
	nc   chan<- net.Conn
	done <-chan struct{}
}

func (d dialer) Dial(network, address string) (net.Conn, error) {
	return d.DialContext(context.Background(), network, address)
}

func (d dialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	c, s := net.Pipe()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-d.done:
		return nil, ErrListenerClosed
	case d.nc <- s:
		return c, nil
	}
}

type ipcAddr struct{}

func (ipcAddr) Network() string { return "ipc" }

func (ipcAddr) String() string { return "ipc" }

func (addr ipcAddr) Addr() net.Addr {
	return addr
}

//Listener is a simple intra-process-communication type that returns in-memory net.Conn's
//for testing and other uses. It implements net.Listener though it does not listen on anything
//network related. NewDialer may be used to return a Dialer that communicates with the Listener.
type Listener struct {
	ipcAddr

	nc   chan net.Conn
	done chan struct{}
}

//NewListener returns a new intra-process Listener for use.
func NewListener() Listener {
	return Listener{
		nc:   make(chan net.Conn),
		done: make(chan struct{}),
	}
}

//ErrListenerClosed is returned from Dial attempts after its associated Listener has been closed.
//It is also returned from Listener.Accept and Listener.Close if the Listener has been closed once already.
var ErrListenerClosed = errors.New("fakesentry/ipc: Listener was closed")

//Accept returns the next intra-process connection for use.
func (l Listener) Accept() (net.Conn, error) {
	select {
	case conn := <-l.nc:
		return conn, nil
	case <-l.done:
		return nil, ErrListenerClosed
	}
}

//Close closes the Listener, it stops all Accept calls as well as all Dials to the Listener.
//Note that it does not close any connections created by the listener, nor is this method safe
//to call concurrently. It is however fine to call more than once.
func (l Listener) Close() error {
	//Not safe being called concurrently, as they'll race on closing the channel.
	//It's fine to call several times from the same goroutine though.
	select {
	case <-l.done:
		return ErrListenerClosed
	default:
		close(l.done)
	}

	return nil
}

//NewDialer returns a new Dialer that can be used to dail the given Listener.
func (l Listener) NewDialer() Dailer {
	return dialer{
		nc:   l.nc,
		done: l.done,
	}
}
