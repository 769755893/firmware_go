package api

import (
	"compress/gzip"
	"fmt"
	"io"
	"main/model"
	"net/http"
)

const InfoUrl = "https://api.ipsw.me/v4/device/%s?type=ipsw"

func GetResponseBody(resp *http.Response) io.ReadCloser {
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(resp.Body)
		defer func(reader io.ReadCloser) {
			err := reader.Close()
			if err != nil {
				panic(err)
			}
		}(reader)
	default:
		reader = resp.Body
	}
	return reader
}

func GetIpswInfo(product string) (ipsw *model.IPSWFirm, err error) {
	url := fmt.Sprintf(InfoUrl, product)
	fmt.Printf("url: %v\n", url)

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("error get ")
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)

	buffer, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	ips, err := model.DecodeIpsw(buffer)

	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, err
	}

	return ips, nil
}
