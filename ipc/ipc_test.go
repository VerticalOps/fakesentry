package ipc

import (
	"io"
	"testing"

	"encoding/json"
)

func TestBasicUsage(t *testing.T) {
	listener := NewListener()
	defer listener.Close()

	dialer := listener.NewDialer()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}

			go io.Copy(conn, conn)
		}
	}()

	conn, err := dialer.Dial("thisgets", "ignored")
	if err != nil {
		t.Fatalf("Unable to dial: %v", err)
	}

	if _, err = conn.Write([]byte(`{"Msg":"Hello World"}`)); err != nil {
		t.Fatalf("Unable to write to conn: %v", err)
	}

	type m struct {
		Msg string
	}

	msg := new(m)
	if err = json.NewDecoder(conn).Decode(msg); err != nil {
		t.Fatalf("Unable to decode json: %v", err)
	}
	t.Logf("Got Message: %+v", msg)

	//close listener early
	if err = listener.Close(); err != nil {
		t.Fatalf("Unable to close listener: %v", err)
	}

	if _, err = dialer.Dial("this", "shouldfail"); err == nil {
		t.Fatal("Dialer did not return an error after Listener.Close, it should have")
	}
}
