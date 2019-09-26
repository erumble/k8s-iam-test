package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jessevdk/go-flags"
)

// Opts struct to allow input from user
var Opts struct {
	Port string `short:"p" long:"port" description:"Port to run server on" required:"false" default:"8080"`
}

func main() {
	// Parse the command line args
	if _, err := flags.NewParser(&Opts, flags.HelpFlag).Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			fmt.Println(err)
			return
		}

		log.Fatal(err)
	}

	log.Println("starting server, listening on port " + Opts.Port)

	http.HandleFunc("/", EchoHandler)
	http.ListenAndServe(":"+Opts.Port, nil)
}

// EchoHandler echos back the request as a response
func EchoHandler(writer http.ResponseWriter, request *http.Request) {

	log.Printf("Echoing back request made to %s to client (%s)", request.URL.Path, request.RemoteAddr)

	writer.Header().Set("Access-Control-Allow-Origin", "*")

	// allow pre-flight headers
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Range, Content-Disposition, Content-Type, ETag")

	request.Write(writer)
}
