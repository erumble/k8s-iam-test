package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/jessevdk/go-flags"
)

// Opts struct to allow input from user
var Opts struct {
	Port string `short:"p" long:"port" description:"Port to run server on" required:"false" default:"8080"`
}

type stsAPI interface {
	GetCallerIdentity(*sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error)
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

	svc := sts.New(session.New())

	log.Println("starting server, listening on port " + Opts.Port)

	http.HandleFunc("/", echoHandler())
	http.HandleFunc("/getCallerIdentity", getCallerIdentityHandler(svc))
	http.ListenAndServe(":"+Opts.Port, nil)
}

func echo(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Echoing back request made to %s to client (%s)", request.URL.Path, request.RemoteAddr)

	writer.Header().Set("Access-Control-Allow-Origin", "*")

	// allow pre-flight headers
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Range, Content-Disposition, Content-Type, ETag")

	request.Write(writer)
}

func echoHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		echo(writer, request)
	}
}

func getCallerIdentityHandler(stsClient stsAPI) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		input := &sts.GetCallerIdentityInput{}

		// Only log the output from the GetCallerIdentity API Call, no need to share the info with the world.
		res, err := stsClient.GetCallerIdentity(input)
		if err != nil {
			log.Printf("Error calling sts.GetCallerIdentity: %v", err)
		}

		log.Printf("sts.GetCallerIdentityResults: %v", res)

		echo(writer, request)
	}
}
