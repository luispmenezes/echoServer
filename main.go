package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-xmlfmt/xmlfmt"
	"github.com/pkg/errors"
	"github.com/yosssi/gohtml"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const defaultPort = 8080

var responseStatusCode *int
var rawBody *bool
var responseBody string

func main() {
	port := flag.Int("port", defaultPort, "server port")
	responseStatusCode = flag.Int("respCode", http.StatusOK, "responseBody status code")
	rawBody = flag.Bool("raw", false, "show unformatted request body")
	responseFile := flag.String("respFile", "", "path to file containing responseBody")
	help := flag.Bool("h", false, "print help")
	flag.Parse()

	if len(*responseFile) > 0 {
		fileBytes, err := ioutil.ReadFile(*responseFile)
		if err != nil {
			log.Fatal(errors.Wrap(err, "failed to open file containing response body"))
		}
		responseBody = string(fileBytes)
	}

	if *help {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	log.Printf("Starting server on port %d...", *port)

	http.HandleFunc("/", handleRequest)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)

	if err != nil {
		log.Fatal(err)
	}
}

func handleRequest(w http.ResponseWriter, req *http.Request) {
	log.Printf(">>> %s %s", req.Method, req.RequestURI)
	fmt.Printf("Protocol: %s ContentLength: %d RemoteAddr:%s\n", req.Proto, req.ContentLength, req.RemoteAddr)
	fmt.Println("Headers:")
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Printf("\t%v: %v\n", name, h)
		}
	}
	if req.ContentLength > 0 {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Body:")

		contentType, contentTypeExists := req.Header["Content-Type"]

		if contentTypeExists && !*rawBody {
			if strings.Contains(contentType[0], "application/json") {
				formattedBody := &bytes.Buffer{}
				if err := json.Indent(formattedBody, body, "", "  "); err != nil {
					log.Fatal(formattedBody)
				}
				fmt.Println(formattedBody.String())
			} else if strings.Contains(contentType[0], "application/xml") {
				fmt.Println(xmlfmt.FormatXML(string(body), "", "  "))
			} else if strings.Contains(contentType[0], "application/html") {
				fmt.Println(gohtml.Format(string(body)))
			} else {
				fmt.Println(string(body))
			}
		} else {
			fmt.Println(string(body))
		}
	}
	http.Error(w, responseBody, *responseStatusCode)
}
