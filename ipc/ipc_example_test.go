package ipc_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/VerticalOps/fakesentry/ipc"
)

func Example() {
	listener := ipc.NewListener()
	defer listener.Close()

	dialer := listener.NewDialer()

	srv := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello World!")
		}),
	}

	//Client setup to use Dialer from Listener
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: dialer.DialContext,
		},
	}

	//Server setup to use Listener
	go srv.Serve(listener)
	defer srv.Close()

	//Does not actually connect to localhost:80
	resp, err := client.Get("http://localhost/myurl")
	if err != nil {
		//Actually handle error
		log.Fatalf("Client HTTP GET: %v", err)
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
	//Output: Hello World!
}
