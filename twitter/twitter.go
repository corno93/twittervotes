package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

func twitter() {
	fmt.Println("hey")
}

// Read authorisation fields from environment variables
type authClient struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

var (
	authSetupOnce sync.Once
	httpClient    *http.Client
	auth          authClient // authorisation fields
)

// read authorisation keys and secrets from env variables
func readAuth() error {

	auth = authClient{ConsumerKey: os.Getenv("SP_TWITTER_KEY"),
		ConsumerSecret: os.Getenv("SP_TWITTER_SECRET"),
		AccessToken:    os.Getenv("SP_TWITTER_ACCESSTOKEN"),
		AccessSecret:   os.Getenv("SP_TWITTER_ACCESSSECRET")}

	if auth.AccessSecret == "" || auth.ConsumerSecret == "" || auth.AccessToken == "" ||
		auth.AccessSecret == "" {
		return errFailedAuth
	}

	return nil
}

var conn net.Conn

func dial(netw, addr string) (net.Conn, error) {
	if conn != nil {
		conn.Close()
		conn = nil
	}
	netc, err := net.DialTimeout(netw, addr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	conn = netc
	return netc, nil
}

// Build the authorized request and return the response
func makeRequest(req *http.Request, params url.Values) (*http.Response,
	error) {
	//use sync.Once to ensure our initialization code gets run only once despite the number of times we call makeRequest
	authSetupOnce.Do(func() {
		readAuth()
		httpClient = &http.Client{
			Transport: &http.Transport{
				Dial: dial,
			}}
	})

	formEnc := params.Encode()
	req.Header.Set("Content-Type", "application/x-www-form- urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(formEnc)))
	req.Header.Set("Authorization", authClient.AuthorizationHeader(creds, "POST", req.URL, params))
	return httpClient.Do(req)
}
