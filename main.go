package main

import (
	"errors"
	"fmt"
	"context"
	"net"
	"net/http"
	"io"
	"io/ioutil"

)

const keyServerAddr = "serverAddr"

func dirRoot(write http.ResponseWriter, req *http.Request) {

	cntx := req.Context()

	hasQueryOne := req.URL.Query().Has("QueryOne")
	QueryOne := req.URL.Query().Get("QueryOne")
	
	hasQueryTwo := req.URL.Query().Has("QueryTwo")
	QueryTwo := req.URL.Query().Get("QueryTwo")

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("Could not read request body. Error given: %s\n", err)
	}


	fmt.Printf("%s: /root directory request received\n QueryOne(%t)=%s, QueryTwo(%t)=%s, Req Body: \n %s \n", 
		cntx.Value(keyServerAddr),
		hasQueryOne, QueryOne,
		hasQueryTwo, QueryTwo, 
		body)
	io.WriteString(write, "Content served to the browser from root directory\n" )
}

func dirTest(write http.ResponseWriter, req *http.Request) {

	cntx := req.Context()

	fmt.Printf("%s: /test directory request received\n", cntx.Value(keyServerAddr))

	thisName := req.PostFormValue("thisName")
	if thisName == "" {
		write.Header().Set("x-missing-filed", "thisName")
		write.WriteHeader(http.StatusBadRequest)
		return 		
	}
	io.WriteString(write, fmt.Sprintf("Hello, %s \n", thisName))
}

func main() {


	mux := http.NewServeMux()

	mux.HandleFunc("/", dirRoot)
	mux.HandleFunc("/test", dirTest)

	cntx := context.Background()

	server := &http.Server {
		Addr: ":3000",
		Handler: mux, 
		BaseContext: func(listener net.Listener) context.Context {
			cntx = context.WithValue(cntx, keyServerAddr, listener.Addr().String())
			return cntx
		},
	}

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server was closed\n")
	} else if err != nil {
		fmt.Printf("Error listening to server, Error given: %s\n", err)
	}

}