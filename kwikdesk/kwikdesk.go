// Copyright 2014 KwikDesk. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kwikdesk

import (
    "bytes"
    "fmt"
    "strings"
    "strconv"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "os"
)

// A set of variables we use to represent endpoints and the
// header content-type we need to use thourough the library.
const (
    kdHTTPContentType = "application/json"

    _KD_TOKEN      = "/token"
    _KD_SEARCH     = "/search"
    _KD_CHANNEL    = "/channel"
    _KD_MESSAGES   = "/messages"
    _KD_SERVERTIME = "/server-time"
)

// The HTTP Client structure that contains the
// request handler.
type Client struct {
    Host       string
    FullHost   string
    XToken   string
    HttpClient *http.Client
    ContentType string
}

// The Message structure that contains the json
// translation. This is used to store a message before
// passing it to the Requester.
type Message struct {
    Content string `json:"content"`
    Delete  int `json:"delete"`
    Private bool `json:"private"`
}

// Not used.
type Error struct {
    Message interface{}
    error int
}

// NewClient creates a new Client object.
// The token can be empty if you plan on creating
// the token using the `CreateToken` func. If you want
// to use your existing token, you need to pass it as the paramter.
// It returns a Client object that is pre-configured for usage.
func NewClient(token string) *Client {

    // Perhaps we should allow people to configure this through some
    // variable. It would make the package almost 100% backwards compatible
    // with the public api (https://developer.kwikdesk.com)
    var (
        host = "platform.kwikdesk.com"
        base = "https://" + host
        req, _ = http.NewRequest("GET", base, nil)
        proxy, _ = http.ProxyFromEnvironment(req)
        transport *http.Transport
    )

    transport = &http.Transport{}
    if proxy != nil {
        transport = &http.Transport{
            Proxy: http.ProxyURL(proxy),
        }
    }

    return &Client{
        Host: host,
        FullHost: base,
        HttpClient: &http.Client{
            Transport: transport,
        },
        ContentType: kdHTTPContentType,
        XToken: token,
    }
}

// This handles parsing the body from the API whenever there's an ERROR
// that happens and is returned by the API. It returns the parsed error and the 
// error code that's been retrieve. 
func ResponseErrorHandler(res *http.Response, res_err error) (ret interface{}, err error) {
    var (
        message = map[string]interface{}{}
        responseBody []byte
    )

    responseBody, _ = ioutil.ReadAll(res.Body)
    json.Unmarshal(responseBody, &message)

    return message, res_err
}

// This is the most used func of the package. It is used to execute every request on the platform endpoints.
// You pass an endpoint (i.e. /channel), a request type (i.e. "GET"), the body content to post (JSON-formatted string),
// and the extra headers you want to add (i.e. X-API-Token, X-Appname, etc.) and it'll add them.
// It will then return an interface containing the results, and the error code if there's an error.
func (c *Client) Requester(endpoint string, reqType string, bodyContent string, headers map[string]interface{}) (ret interface{}, err error) {
    var (
        request *http.Request
        response *http.Response
        responseBody []byte
        responseJson = map[string]interface{}{}
        url = fmt.Sprintf("%s/%s", c.FullHost, endpoint)
    )

    // I'm too dumb to find out how to transform a string "" into nil so fuck it.
    if len(bodyContent) > 0 {
        request, err = http.NewRequest(strings.ToUpper(reqType), url, bytes.NewBufferString(bodyContent))
    } else {
        request, err = http.NewRequest(strings.ToUpper(reqType), url, nil)
    }

    request.Header.Set("Content-Type", c.ContentType)
    for key, value := range headers {
        request.Header.Set(fmt.Sprintf("%v", key), fmt.Sprintf("%v", value))
    }

    if response, err = c.HttpClient.Do(request); err != nil {
        return ResponseErrorHandler(response, err)
    }

    // Transpose to string because I don't know how to use "HasPrefix" with an int.
    code := strconv.Itoa(response.StatusCode)

    // Here we check if the prefix starts with "20". That's a 200, 201, etc.
    if !strings.HasPrefix(code, "20") {
        err = fmt.Errorf("HTTP Response Code: %v", response.StatusCode)
        return ResponseErrorHandler(response, err)
    }

    // We then identify whether or not there's an error whilst reading the
    // response body from the response object.
    if responseBody, err = ioutil.ReadAll(response.Body); err != nil {
        return ResponseErrorHandler(response, err)
    }

    // Furthermore, we look for errors when unmarshalling the json.
    if err = json.Unmarshal(responseBody, &responseJson); err != nil {
        return ResponseErrorHandler(response, err)
    }

    // Everything went well so we just send the JSON-interface object back to the user with no error.
    return responseJson, nil
}

// This function is used to create Messages as described in:
// https://partners.kwikdesk.com/#create
// You pass the content string, the delete time integer value and whether or not it's private (true false)
// This will return the response from the API and the error if something happened.
func (c *Client) Messages(messageContent string, deleteTime int, privateFlag bool) (ret interface{}, err error) {
    var (
        endpoint = _KD_MESSAGES
        requestType = "POST"
        bodyContent = &Message{Content: messageContent, Delete: deleteTime, Private: privateFlag}
        headers = make(map[string]interface{})
    )

    headers["X-API-Token"] = c.XToken

    body, err := json.Marshal(bodyContent)
    response, err := c.Requester(endpoint, requestType, string(body), headers)


    results := response.(map[string]interface{})
    return results, err
}

// This function is used to retrieve channel-messages as described in:
// https://partners.kwikdesk.com/#channel
// All you require to pass is the token which is stored in the Client object
// that was either saved at instance creation, or after you've invoked the CreateToken func.
// This will return the response from the API and the error if something happened.
func (c *Client) Channel() (ret interface{}, err error) {
    var (
        endpoint = _KD_CHANNEL
        requestType = "GET"
        bodyContent = ""
        headers = make(map[string]interface{})
    )

    headers["X-API-Token"] = c.XToken

    response, err := c.Requester(endpoint, requestType, bodyContent, headers)
    results := response.(map[string]interface{})["results"]

    return results, err
}

// This function is used to place searches on hashtags associated to your token
// and marked as private being false, as described in:
// https://partners.kwikdesk.com/#search
// You only need to pass the search term you are interested in searching for
// and the function will return the response from the API and the error 
// if something happened.
func (c* Client) Search(term string) (ret interface{}, err error) {
    var (
        endpoint = fmt.Sprintf("%s?q=%s", _KD_SEARCH, term)
        requestType = "GET"
        bodyContent = ""
        headers = make(map[string]interface{})
    )

    headers["X-API-Token"] = c.XToken

    response, err := c.Requester(endpoint, requestType, bodyContent, headers)
    results := response.(map[string]interface{})["results"]

    stopOnError(err)
    return results, err
}

// Set the Xtoken Client object variable.
func (c *Client) SetToken(token string) {
    c.XToken = token
}

// Retrieve the token value from the Client object variable XToken
func (c *Client) GetToken() string {
    return c.XToken
}

// This function is used to create Tokens as described in:
// https://partners.kwikdesk.com/#token
// You are required to pass the application name (or an email address). Upon success, 
// the value of the token will be set in the Client object which you can then reuse 
// easily with the other funcs like Message, Search, etc.
// This will return the response from the API and the error if something happened.
func (c *Client) CreateToken(appName string) (ret interface{}, err error) {
    var (
        endpoint = _KD_TOKEN
        requestType = "POST"
        bodyContent = ""
        headers = make(map[string]interface{})
    )

    headers["X-Appname"] = appName

    response, err := c.Requester(endpoint, requestType, bodyContent, headers)
    token := response.(map[string]interface{})["token"].(string)

    c.SetToken(token)
    results := response.(map[string]interface{})

    return results, err
}

// This queries the KwikDesk Platform to retrieve the current
// server time. Even though HTTP has a native mechanism for accomplishing
// precisely that, we want to insure consistency between and this also gives
// you different formats to work with.
// This will return the response from the API and the error if something happened.
func (c* Client) ServerTime() (ret interface{}, err error) {
    var (
        endpoint = _KD_SERVERTIME
        requestType = "GET"
        bodyContent = ""
        headers = make(map[string]interface{})
    )

    headers["X-API-Token"] = c.XToken

    response, err := c.Requester(endpoint, requestType, bodyContent, headers)
    stopOnError(err)

    results := response.(map[string]interface{})
    return results, err
}

// Stop on error and print the message if something happened. This is just
// like a fake "catch".
func stopOnError(err error) {
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error caught during request: ", err)
        os.Exit(1)
    }
}
