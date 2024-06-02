package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

)

func root(write http.ResponseWriter, req *http.Request) {
	fmt.Printf("Root request recieved")
	io.WriteString(write, "Write given to server" )
}