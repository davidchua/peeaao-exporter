package pinger

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func MakePing(target string, locations []string, authToken string) (io.ReadCloser, int, error) {

	locationsDelimited := strings.Join(locations, ",")
	formData := url.Values{
		"target": {target},
	}

	if locationsDelimited != "" {
		formData["location"] = []string{locationsDelimited}
	}

	posturl := "https://api.peeaao.com/api/ping"

	r, err := http.NewRequest("POST", posturl, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, 0, err
	}

	r.Header.Add("Authorization", "Bearer abcdefghij")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return nil, 0, err
	}

	fmt.Printf("status code is %#v", res.StatusCode)

	return res.Body, res.StatusCode, nil

}
