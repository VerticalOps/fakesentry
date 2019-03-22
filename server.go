package fakesentry

import (
	"net/http"

	"github.com/VerticalOps/fakesentry/ipc"
)

//Server is a convenience type that uses the default configurations of a *http.Server,
//this packages Handler type, this repos ipc package Listener/Dialer types, and *http.Transport,
//to start an HTTP server that only listens for intra-process connections and does not use the network.
//For more configuration use this Server's methods, or use each of the aforementioned types
//individually. Refer to the package example and each types documentation for more information.
type Server struct {
	*http.Server

	tr       *http.Transport
	listener ipc.Listener
	dialer   ipc.Dialer
}

//NewUnstartedServer returns a new Server with defaults mentioned on the Server types definition.
//It does not start a goroutine serving connections with the internal intra-process listener.
func NewUnstartedServer() Server {
	listener := ipc.NewListener()
	dialer := listener.NewDialer()
	handler := NewHandler()

	return Server{
		Server:   &http.Server{Handler: handler},
		tr:       &http.Transport{DialContext: dialer.DialContext},
		dialer:   dialer,
		listener: listener,
	}
}

//NewServer is like NewUnstartedServer but it does start a goroutine to service connections
//from the available Dialer and Transport methods. Server.Close should be called when done.
func NewServer() Server {
	srv := NewUnstartedServer()
	go srv.Serve(srv.Listener())

	return srv
}

//Listener returns the intra-process connection Listener from this repos ipc package. Look to
//that packages documentation for more information.
func (s Server) Listener() ipc.Listener {
	return s.listener
}

//Dialer returns the intra-process Dialer from this repos ipc package. Look to
//that packages documentation for more information.
func (s Server) Dialer() ipc.Dialer {
	return s.dialer
}

//Transport returns an *http.Transport that can be used with *http.Clients to speak to the Server.
func (s Server) Transport() *http.Transport {
	return s.tr
}
