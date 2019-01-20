package sleepiq

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

// httpGet conducts a GET request with the provided url. The function returns the response from the
// service as a byte array.
func httpGet(url string, cookies []*http.Cookie, headers map[string]string) ([]byte, error) {
	var response []byte

	// Create a new http client
	client := http.Client{
		Timeout: 20 * time.Second,
	}

	// Create the request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return response, err
	}

	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Add cookies
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	// Make the request
	res, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer res.Body.Close()

	// Read the response
	response, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return response, err
	}

	return response, nil
}

// httpPut conducts a PUT request with the provided url. The function returns the response from the
// service as a byte array.
func httpPut(url string, payload []byte, cookies []*http.Cookie) ([]byte, []*http.Cookie, error) {
	var response []byte

	// Create a new http client
	client := http.Client{
		Timeout: 20 * time.Second,
	}

	// Create the request
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(payload))
	if err != nil {
		return response, cookies, err
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")

	// Add cookies
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	// Make the request
	res, err := client.Do(req)
	if err != nil {
		return response, cookies, err
	}
	defer res.Body.Close()

	// Read the response
	response, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return response, cookies, err
	}

	cookies = res.Cookies()

	return response, cookies, nil
}

// httpPost conducts a POST request with the provided url. The function returns the response from the
// service as a byte array.
func httpPost(url string, payload []byte) ([]byte, error) {
	var response []byte

	// Create a new http client
	client := http.Client{
		Timeout: 20 * time.Second,
	}

	// Create the request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return response, err
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", "3c924e14923642baa1c4ad1d5096a1c5")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")

	// Make the request
	res, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer res.Body.Close()

	// Read the response
	response, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return response, err
	}

	return response, nil
}
