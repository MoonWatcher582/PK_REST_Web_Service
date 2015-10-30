package main

import (
	"bytes"
	"flag"
	"fmt"
	log "github.com/golang/glog"
	"io/ioutil"
	"net/http"
	urllib "net/url"
	"strings"
)

//	List of possible input flags
var (
	urlFlag = flag.String("url", "", "URL for server")
	method  = flag.String("method", "", "HTTP method to use - create, get, remove, update")
	data    = flag.String("data", "", "JSON data to pass in")
	year    = flag.String("year", "", "year to pass in")
)

func main() {
	flag.Parse()

	if *urlFlag == "" {
		log.Fatal("URL not defined.")
	}
	if *method == "" {
		log.Fatal("Method not defined.")
	}

	//	Switch on method
	methodName := strings.ToUpper(*method)
	opMap := map[string]string{
		"CREATE": "POST",
		"LIST":   "GET",
		"REMOVE": "DELETE",
		"UPDATE": "UPDATE",
	}
	httpOperation, ok := opMap[methodName]
	if !ok {
		log.Fatal("No proper method specified")
		return
	}

	//	Create a basic client object
	client := &http.Client{}

	//	Parse the url flag
	url, err := urllib.Parse(*urlFlag)
	if err != nil {
		log.Fatal(err)
	}

	// Cast content of data flag into a byte array
	buff := []byte{}
	if httpOperation == "POST" {
		buff = []byte(*data)
	}

	//	Build HTTP request
	req := &http.Request{
		Method:        httpOperation,
		URL:           url,
		Body:          ioutil.NopCloser(bytes.NewBuffer(buff)),
		ContentLength: int64(len(buff)),
	}

	//	Parse year flag for DELETE
	if httpOperation == "DELETE" && *year != "" {
		params := req.URL.Query()
		params.Add("year", *year)
		req.URL.RawQuery = params.Encode()
	}

	//	log the request before we make it
	log.Infof("%+v", req)

	//	run the request with the client
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	//	store output into a buffer and print it
	buff, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buff))
}
