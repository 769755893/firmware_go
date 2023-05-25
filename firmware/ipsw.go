package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)


const INFO_URL = "https://api.ipsw.me/v4/device/%s?type=ipsw"


func getResponseBody(resp *http.Response) io.ReadCloser {
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}
	return reader
}


func getIpswInfo(product string) (ipsw *IPSWFirm, err error) {
	url := fmt.Sprintf(INFO_URL, product)
	fmt.Printf("url: %v\n", url)

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("error get ")
		return nil, err
	}

	defer res.Body.Close()

	buffer, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}
	
	ips, err := decodeIpsw(buffer)

	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, err
	}

	return ips, nil
}