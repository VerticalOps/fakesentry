

# fakesentry [![GoDoc](https://godoc.org/github.com/VerticalOps/fakesentry?status.svg)](https://godoc.org/github.com/VerticalOps/fakesentry)
`import "github.com/VerticalOps/fakesentry"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Examples](#pkg-examples)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>
Package fakesentry provides various functionality to run a 'fake sentry' server internal to a Go process
off or on the network. The primary motivation of this package was to provide an easy way of viewing data
your Go program may send to Sentry during run, but in a development mode without sending data to a production
DSN or otherwise setting up an entire Sentry server instance with its dependencies. As such, this should be
guarded under a development mode or build flag unless care is taken and the user has read this packages documentation.

By default, the Server type in this package will start an internal HTTP server without using the network stack
and provides a Dialer and *http.Transport that can be configured with your Sentry/Raven client to make requests to.
The Server will simply pretty-print the JSON sent by the client to STDERR. Much of this behavior can be added to
or changed. It is recommended to look at the package example and read each types documentation for more configuration.




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [func FromContext(ctx context.Context) (jb json.RawMessage, ok bool)](#FromContext)
* [func FromRequest(r *http.Request) (json.RawMessage, bool)](#FromRequest)
* [type ErrorHandler](#ErrorHandler)
* [type Handler](#Handler)
  * [func NewHandler(opts ...Option) Handler](#NewHandler)
  * [func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request)](#Handler.ServeHTTP)
* [type Logger](#Logger)
* [type Option](#Option)
  * [func AsMiddleware(next http.Handler) Option](#AsMiddleware)
  * [func WithErrorHandler(eh ErrorHandler) Option](#WithErrorHandler)
  * [func WithLogger(logger Logger) Option](#WithLogger)
* [type Server](#Server)
  * [func NewServer() Server](#NewServer)
  * [func NewUnstartedServer() Server](#NewUnstartedServer)
  * [func (s Server) Dialer() ipc.Dialer](#Server.Dialer)
  * [func (s Server) Listener() ipc.Listener](#Server.Listener)
  * [func (s Server) Transport() *http.Transport](#Server.Transport)

#### <a name="pkg-examples">Examples</a>
* [Package](#example_)

#### <a name="pkg-files">Package files</a>
[fakesentry.go](/src/github.com/VerticalOps/fakesentry/fakesentry.go) [options.go](/src/github.com/VerticalOps/fakesentry/options.go) [server.go](/src/github.com/VerticalOps/fakesentry/server.go) 



## <a name="pkg-variables">Variables</a>
``` go
var (
    //ErrBadMethod is passed into an ErrorHandler if the client uses an unexpected HTTP method.
    ErrBadMethod = errors.New("fakesentry: Bad HTTP method from client")

    //ErrBadContentLength is passed into an ErrorHandler if the client sends an invalid content length.
    ErrBadContentLength = errors.New("fakesentry: Bad ContentLength from client")

    //ErrBadContentType is passed into an ErrorHandler if the clients sends an unexpected content type.
    ErrBadContentType = errors.New("fakesentry: Bad ContentType from client")
)
```


## <a name="FromContext">func</a> [FromContext](/src/target/options.go?s=1116:1183#L48)
``` go
func FromContext(ctx context.Context) (jb json.RawMessage, ok bool)
```
FromContext returns a json.RawMessage that was saved if AsMiddleware
is used with the Handler type. Look to FromRequest for more information.



## <a name="FromRequest">func</a> [FromRequest](/src/target/options.go?s=1459:1516#L56)
``` go
func FromRequest(r *http.Request) (json.RawMessage, bool)
```
FromRequest retrieves the saved json.RawMessage Handler.ServeHTTP saves
into the request context if AsMiddleware is used. Other http.Handlers can
use this to retrieve the value and do with it as they will.




## <a name="ErrorHandler">type</a> [ErrorHandler](/src/target/options.go?s=677:742#L29)
``` go
type ErrorHandler func(http.ResponseWriter, *http.Request, error)
```
ErrorHandler is called when an error occurs within Handler.ServeHTTP.
The error is guaranteed to be non-nil, the handler may use the passed http types
to communicate to the client what happened.










## <a name="Handler">type</a> [Handler](/src/target/fakesentry.go?s=2072:2149#L45)
``` go
type Handler struct {
    // contains filtered or unexported fields
}
```
Handler implements a basic ServeHTTP method that can be used to accept and parse HTTP requests from
Sentry clients. The Handler type is meant for testing, usually guarded by a development or build flag.
It is generally unsafe to use in production unless guarded by other HTTP Handlers, such as something
that checks the requests auth state.







### <a name="NewHandler">func</a> [NewHandler](/src/target/fakesentry.go?s=2565:2604#L57)
``` go
func NewHandler(opts ...Option) Handler
```
NewHandler returns a new Handler with defaults unless Options are set.
Unless otherwise changed, the default logger is a stdlib 'log' that prints to STDERR.
The default ErrorHandler prints the error to the log and writes a 400 or 500 code to the client.
If the Handler is not used in middleware then the default action is to dump the HTTP request
as well as pretty print the Sentry JSON into the logger.





### <a name="Handler.ServeHTTP">func</a> (Handler) [ServeHTTP](/src/target/fakesentry.go?s=3336:3402#L90)
``` go
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request)
```
ServeHTTP implements http.Handler. Please refer to the Handler type definiton
NewHandler function for more information.




## <a name="Logger">type</a> [Logger](/src/target/options.go?s=223:291#L13)
``` go
type Logger interface {
    Printf(format string, arg ...interface{})
}
```
Logger is the basic logging interface used by this package.










## <a name="Option">type</a> [Option](/src/target/options.go?s=133:159#L10)
``` go
type Option func(*Handler)
```
Option is a functional option to help configure a Handler.







### <a name="AsMiddleware">func</a> [AsMiddleware](/src/target/options.go?s=1791:1834#L63)
``` go
func AsMiddleware(next http.Handler) Option
```
AsMiddleware takes the next http.Handler and changes the behavior of
Handler.ServeHTTP, in that if a http.Handler is set, the parsed Sentry JSON
value will be saved into the http.Requests context for use by the next http.Handler.


### <a name="WithErrorHandler">func</a> [WithErrorHandler](/src/target/options.go?s=801:846#L32)
``` go
func WithErrorHandler(eh ErrorHandler) Option
```
WithErrorHandler returns an Option for an ErrorHandler


### <a name="WithLogger">func</a> [WithLogger](/src/target/options.go?s=357:394#L18)
``` go
func WithLogger(logger Logger) Option
```
WithLogger returns an Option that sets a Logger on a Handler.





## <a name="Server">type</a> [Server](/src/target/server.go?s=562:669#L14)
``` go
type Server struct {
    *http.Server
    // contains filtered or unexported fields
}
```
Server is a convenience type that uses the default configurations of a *http.Server,
this packages Handler type, this repos ipc package Listener/Dialer types, and *http.Transport,
to start an HTTP server that only listens for intra-process connections and does not use the network.
For more configuration use this Server's methods, or use each of the aforementioned types
individually. Refer to the package example and each types documentation for more information.







### <a name="NewServer">func</a> [NewServer](/src/target/server.go?s=1341:1364#L39)
``` go
func NewServer() Server
```
NewServer is like NewUnstartedServer but it does start a goroutine to service connections
from the available Dialer and Transport methods. Server.Close should be called when done.


### <a name="NewUnstartedServer">func</a> [NewUnstartedServer](/src/target/server.go?s=863:895#L24)
``` go
func NewUnstartedServer() Server
```
NewUnstartedServer returns a new Server with defaults mentioned on the Server types definition.
It does not start a goroutine serving connections with the internal intra-process listener.





### <a name="Server.Dialer">func</a> (Server) [Dialer](/src/target/server.go?s=1783:1818#L54)
``` go
func (s Server) Dialer() ipc.Dialer
```
Dialer returns the intra-process Dialer from this repos ipc package. Look to
that packages documentation for more information.




### <a name="Server.Listener">func</a> (Server) [Listener](/src/target/server.go?s=1588:1627#L48)
``` go
func (s Server) Listener() ipc.Listener
```
Listener returns the intra-process connection Listener from this repos ipc package. Look to
that packages documentation for more information.




### <a name="Server.Transport">func</a> (Server) [Transport](/src/target/server.go?s=1940:1983#L59)
``` go
func (s Server) Transport() *http.Transport
```
Transport returns an *http.Transport that can be used with *http.Clients to speak to the Server.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
