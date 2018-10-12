package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/TheSp1der/goerror"
)

func webRoot(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "text/plain")
	if req.Method == "GET" {
		resp.Write([]byte("Welcome to the site\n"))
		resp.Write([]byte("Request Method: " + req.Method + "\n"))
		resp.Write([]byte("Request Host: " + req.Host + "\n"))
		resp.Write([]byte("Request Header: " + fmt.Sprintf("%+v", req.Header) + "\n"))
		resp.Write([]byte("Request Protocol: " + req.Proto + "\n"))
		resp.Write([]byte("Request Remote Address: " + req.RemoteAddr + "\n"))
		resp.Write([]byte("Request Request URI: " + req.RequestURI + "\n"))
	}
}

func webUUID(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "text/plain")
	if req.Method == "GET" {
		resp.Write([]byte("Please enjoy this free UUID: " + uuidV5Gen(1000, time.Now().Local().Format("2006-01-02 15:04:05")) + "\n"))
	}
}

func webListener(port int) {
	ws := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: ws,
		// time to read request headers
		ReadTimeout: time.Duration(15 * time.Second),
		// time from accept to end of response
		WriteTimeout: time.Duration(10 * time.Second),
		// time a Keep-Alive connection will be kept idle
		IdleTimeout: time.Duration(120 * time.Second),
	}

	ws.HandleFunc("/", webRoot)
	ws.HandleFunc("/uuid", webUUID)

	if err := srv.ListenAndServe(); err != nil {
		goerror.Fatal(err)
	}
}
