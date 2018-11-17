package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joeshaw/envdecode"
	"github.com/matryer/go-oauth/oauth"
)

var (
	authSetupOnce sync.Once
	httpClient    *http.Client
	authClient    *oauth.Client
	creds         *oauth.Credentials
	client        *http.Client
	conn          net.Conn
	reader        io.ReadCloser
	shutdown      bool
)

// Method closes and opens new connections (to the address on the named network - twitter) continuosuly so if a connection dies
// we can re-dial without worrying about zombie connections.
func dial(netw, addr string) (net.Conn, error) {

	fmt.Println(netw, addr)

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

// read environment variables
func SetupTwitterAuth() {
	var ts struct {
		ConsumerKey    string `env:"SP_TWITTER_KEY,required"`
		ConsumerSecret string `env:"SP_TWITTER_SECRET,required"`
		AccessToken    string `env:"SP_TWITTER_ACCESSTOKEN,required"`
		AccessSecret   string `env:"SP_TWITTER_ACCESSSECRET,required"`
	}
	if err := envdecode.Decode(&ts); err != nil {
		log.Fatalln(err)
	}
	creds = &oauth.Credentials{
		Token:  ts.AccessToken,
		Secret: ts.AccessSecret,
	}
	authClient = &oauth.Client{
		Credentials: oauth.Credentials{
			Token:  ts.ConsumerKey,
			Secret: ts.ConsumerSecret,
		}}

}

// Build the authorized request and return the response
func makeRequest(url string, params url.Values) (*http.Response,
	error) {

	// //use sync.Once to ensure our initialization code gets run only once despite the number of times we call makeRequest
	// authSetupOnce.Do(func() {
	// 	setupTwitterAuth()
	// 	httpClient = &http.Client{
	// 		Transport: &http.Transport{
	// 			Dial: dial,
	// 		}}
	// })

	formEnc := params.Encode()
	req, err := http.NewRequest("POST", url, strings.NewReader(formEnc))

	if err != nil {
		log.Println("creating filter request failed:", err)
		return nil, err
	}

	req.Header.Set("Authorization", authClient.AuthorizationHeader(creds, "POST", req.URL, params))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(formEnc)))

	return httpClient.Do(req)

}

type tweet struct {
	Text string
}

// read from twitter
func ReadTwitter(votes chan string, options []string) {

	SetupTwitterAuth()

	// Talk to services over http. include transport struct that used by clients to manage
	// the underlying TCP connection and it’s Dialer is a struct that manages the
	// establishment of the connection.
	httpClient = &http.Client{
		Transport: &http.Transport{
			Dial: dial,
		}}

	shutdown = false

	// continuosly loop forever
	for {

		// make the url
		u, _ := url.Parse("https://stream.twitter.com/1.1/statuses/filter.json")

		hashtags := make([]string, len(options))
		for i := range options {
			hashtags[i] = "#" + strings.ToLower(options[i])
		}
		form := url.Values{"track": {strings.Join(hashtags, ",")}}

		resp, err := makeRequest(u.String(), form)
		if err != nil {
			log.Println("making request failed:", err)
		}

		// this is a nice way to see what the error actually is:
		if resp.StatusCode != http.StatusOK {
			s := bufio.NewScanner(resp.Body)
			s.Scan()
			log.Println(s.Text())
			log.Println(hashtags)
			log.Println("StatusCode =", resp.StatusCode)
			continue
		}
		reader = resp.Body
		decoder := json.NewDecoder(reader)

		for {
			if shutdown {
				log.Println("Twitter shutdown")
				return
			}
			fmt.Println("SHTDOWN: ", shutdown)
			var t tweet
			if err := decoder.Decode(&t); err == nil {
				for _, option := range options {
					if strings.Contains(
						strings.ToLower(t.Text),
						strings.ToLower(option),
					) {
						votes <- option
					}
				}
			} else {
				break
			}
		}
	}
}

func closeConn() {
	if conn != nil {
		conn.Close()
	}
	if reader != nil {
		reader.Close()
	}
}

func ShutDownTwitter() {
	shutdown = true
}
