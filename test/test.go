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

	switch strings.ToUpper(*method) {
	case "CREATE":
		*method = "POST"
	case "LIST":
		*method = "GET"
	case "REMOVE":
		*method = "DELETE"
	case "UPDATE":
		*method = "UPDATE"
	default:
		log.Fatal("Incorrect Method.")
	}

	if strings.ToUpper(*method) == "CREATE" {
		*method = "POST"
	}

	client := &http.Client{}

	url, err := urllib.Parse(*urlFlag)
	if err != nil {
		log.Fatal(err)
	}

	buff := []byte(*data)
	req := &http.Request{
		Method:        strings.ToUpper(*method),
		URL:           url,
		Body:          ioutil.NopCloser(bytes.NewBuffer(buff)),
		ContentLength: int64(len(buff)),
	}

	if *year != "" {
		params := req.URL.Query()
		params.Add("year", *year)
		req.URL.RawQuery = params.Encode()
	}

	log.Infof("%+v", req)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	buff, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buff))
}
