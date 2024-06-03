package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

)

func dirRoot(write http.ResponseWriter, req *http.Request) {
	fmt.Printf("/root directory request received\n")
	io.WriteString(write, "Content served to the browser from root directory\n" )
}

func dirTest(write http.ResponseWriter, req *http.Request) {
	fmt.Printf("/test directory request received\n")
	io.WriteString(write, "Content served at route /hello\n" )
}

func main() {
	http.HandleFunc("/", dirRoot)
	http.HandleFunc("/test", dirTest)
	
	// blocking call (like async await in javacript)
	err := http.ListenAndServe(":3000", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server was closed\n")
	} else if err != nil {
		fmt.Printf("Error occured when starting server %s\n", err)
		os.Exit(1)
	}
}