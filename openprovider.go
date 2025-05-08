package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// Define constants
const (
	baseurl = "http://api.sandbox.openprovider.nl:8480/v1beta/"
)

// Client struct for grouping key information meant for a client
type Client struct {
	// HTTP client for sending and receiving http requests
	client *http.Client

	// Bearer token used during communication against endpoints
	Token string

	// Services used for talking to OpenProvider API endpoints
	// Dns DNSServiceInterface
	// Domains DomainServiceInterface
}

// The NewClient function creates a new client with the provided username and password
func NewClient(username string, password string) (*Client, error) {

	client, err := CreateClient(username, password)
	if err != nil {
		return nil, err
	}

	return client, err
}

// The function CreateClient creates a new client with the provided username and password,
// authenticates the client, and returns the client object along with any potential errors.
func CreateClient(username string, password string) (*Client, error) {
	// Instance c of Client struct
	c := &Client{}

	// Defining http client so that individual functions can send and receive http requests
	c.client = &http.Client{}

	// Acquiring BearerToken
	c.Token = auth(username, password)

	// All services are connected to Client
	// c.Domains = &DomainService{client: c}

	return c, nil
}

// The function auth gathers and returns a token by sending a http post request to OpenProvider's authentication endpoint
func auth(username string, password string) string {
	loginData := map[string]string{
		"username": username,
		"password": password,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		panic(err)
	}

	url := baseurl + "auth/login"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var res map[string]any
	{
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		panic(err)
	}

	token := ""
	if data, ok := res["data"].(map[string]any); ok {
		if t, ok := data["token"].(string); ok {
			token = t
		}
	} else {
		panic("No token returned. Check your credentials.")
	}

	return token
}

// NewRequest handles all HTTP requests
func (c *Client) NewRequest(method string, path string, payload any) (*http.Request, error) {

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// OpenProvider always expects a Bearer token and a json payload
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// Do handles all HTTP requests, sends them and parses the response into a dictionary
func (c *Client) Do(req *http.Request) (map[string]any, error) {

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res map[string]any

	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
