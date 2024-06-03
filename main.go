package main

import (
	"errors"
	"fmt"
	"context"
	"net"
	"net/http"
	"io"

)

const keyServerAddr = "serverAddr"

func dirRoot(write http.ResponseWriter, req *http.Request) {

	cntx := req.Context()

	fmt.Printf("%s: /root directory request received\n", cntx.Value(keyServerAddr))
	io.WriteString(write, "Content served to the browser from root directory\n" )
}

func dirTest(write http.ResponseWriter, req *http.Request) {

	cntx := req.Context()

	fmt.Printf("%s: /test directory request received\n", cntx.Value(keyServerAddr))
	io.WriteString(write, "Content served at route /test\n" )
}

func main() {


	mux := http.NewServeMux()

	mux.HandleFunc("/", dirRoot)
	mux.HandleFunc("/test", dirTest)

	cntx, cancelCntx := context.WithCancel(context.Background())
	serverOne := &http.Server {
		Addr: ":3000",
		Handler: mux, 
		BaseContext: func(listener net.Listener) context.Context {
			cntx = context.WithValue(cntx, keyServerAddr, listener.Addr().String())
			return cntx
		},
	}

	serverTwo := &http.Server {
		Addr: ":4000", 
		Handler: mux, 
		BaseContext: func(listener net.Listener) context.Context {
			cntx = context.WithValue(cntx, keyServerAddr, listener.Addr().String())
			return cntx
		},
	}

	go func() {
		err := serverOne.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Server 1 was closed\n")
		} else if err != nil {
			fmt.Printf("Error listening on Server 1 &s\n", err)
		}
		cancelCntx()
	}()

	go func() {
		err := serverTwo.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Server 2 closed\n")
		} else if err != nil {
			fmt.Printf("Error listening on Server 2 %s\n", err)
		}
		cancelCntx()
	}()

	<-cntx.Done()

	// blocking call (like async await in javacript)
	// err := http.ListenAndServe(":3000", mux)

	// if errors.Is(err, http.ErrServerClosed) {
	// 	fmt.Printf("Server was closed\n")
	// } else if err != nil {
	// 	fmt.Printf("Error occured when starting server %s\n", err)
	// 	os.Exit(1)
	// }
}