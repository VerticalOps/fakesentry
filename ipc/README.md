

# ipc [![GoDoc](https://godoc.org/github.com/VerticalOps/fakesentry/ipc?status.svg)](https://godoc.org/github.com/VerticalOps/fakesentry/ipc)
`import "github.com/VerticalOps/fakesentry/ipc"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Examples](#pkg-examples)

## <a name="pkg-overview">Overview</a>
Package ipc implements a simple set of intra-process-communication net.Conn types, specifically a net.Listener
and an associated Dialer that can be used in place of http Servers/Clients or anything that needs a net.Conn to work
or test within the same process. This package is built around net.Pipe and so the network is not actually involved.
If a buffered implementation is needed there are alternatives.




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [type Dialer](#Dialer)
* [type Listener](#Listener)
  * [func NewListener() Listener](#NewListener)
  * [func (l Listener) Accept() (net.Conn, error)](#Listener.Accept)
  * [func (addr Listener) Addr() net.Addr](#Listener.Addr)
  * [func (l Listener) Close() error](#Listener.Close)
  * [func (Listener) Network() string](#Listener.Network)
  * [func (l Listener) NewDialer() Dialer](#Listener.NewDialer)
  * [func (Listener) String() string](#Listener.String)

#### <a name="pkg-examples">Examples</a>
* [Package](#example_)

#### <a name="pkg-files">Package files</a>
[ipc.go](/src/github.com/VerticalOps/fakesentry/ipc/ipc.go) 



## <a name="pkg-variables">Variables</a>
``` go
var ErrListenerClosed = errors.New("fakesentry/ipc: Listener was closed")
```
ErrListenerClosed is returned from Dial attempts after its associated Listener has been closed.
It is also returned from Listener.Accept and Listener.Close if the Listener has been closed once already.




## <a name="Dialer">type</a> [Dialer](/src/target/ipc.go?s=836:987#L19)
``` go
type Dialer interface {
    Dial(network, address string) (net.Conn, error)
    DialContext(ctx context.Context, network, address string) (net.Conn, error)
}
```
Dialer can be used in place of many Dial-related functions of the net package and others.
The returned net.Conn is always related to the Listener that returned the Dialer via Listener.NewDialer.
Therefore the passed network and address are ignored. The given context is not however, if the context
expires before the connection is made then ctx.Err is returned.










## <a name="Listener">type</a> [Listener](/src/target/ipc.go?s=1902:1976#L59)
``` go
type Listener struct {
    // contains filtered or unexported fields
}
```
Listener is a simple intra-process-communication type that returns in-memory net.Conn's
for testing and other uses. It implements net.Listener though it does not listen on anything
network related. NewDialer may be used to return a Dialer that communicates with the Listener.







### <a name="NewListener">func</a> [NewListener](/src/target/ipc.go?s=2038:2065#L67)
``` go
func NewListener() Listener
```
NewListener returns a new intra-process Listener for use.





### <a name="Listener.Accept">func</a> (Listener) [Accept](/src/target/ipc.go?s=2491:2535#L79)
``` go
func (l Listener) Accept() (net.Conn, error)
```
Accept returns the next intra-process connection for use.




### <a name="Listener.Addr">func</a> (Listener) [Addr](/src/target/ipc.go?s=1566:1601#L52)
``` go
func (addr Listener) Addr() net.Addr
```



### <a name="Listener.Close">func</a> (Listener) [Close](/src/target/ipc.go?s=2898:2929#L91)
``` go
func (l Listener) Close() error
```
Close closes the Listener, it stops all Accept calls as well as all Dials to the Listener.
Note that it does not close any connections created by the listener, nor is this method safe
to call concurrently. It is however fine to call more than once.




### <a name="Listener.Network">func</a> (Listener) [Network](/src/target/ipc.go?s=1467:1498#L48)
``` go
func (Listener) Network() string
```



### <a name="Listener.NewDialer">func</a> (Listener) [NewDialer](/src/target/ipc.go?s=3254:3290#L105)
``` go
func (l Listener) NewDialer() Dialer
```
NewDialer returns a new Dialer that can be used to dail the given Listener.




### <a name="Listener.String">func</a> (Listener) [String](/src/target/ipc.go?s=1517:1547#L50)
``` go
func (Listener) String() string
```







- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
