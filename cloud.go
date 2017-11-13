package cloud

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	Url      *url.URL
	Username string
	Password string
}

type Error struct {
	Exception string `xml:"exception"`
	Message   string `xml:"message"`
}

func NewClient(host, username, password string) (*Client, error) {
	url, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	return &Client{
		Url:      url,
		Username: username,
		Password: password,
	}, nil
}

func (c *Client) Mkdir(folder string) error {

	// Create the https request

	folderUrl, err := url.Parse(folder)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("MKCOL", c.Url.ResolveReference(folderUrl).String(), nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.Username, c.Password)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(body) > 0 {
		error := Error{}
		err = xml.Unmarshal(body, &error)
		if err != nil {
			return fmt.Errorf("Error during XML Unmarshal for response %s. The error was %s", body, err)
		}
		if error.Exception != "" {
			return fmt.Errorf("Exception: %s, Message: %s", error.Exception, error.Message)
		}

	}

	return nil
}

func (c *Client) Delete(folder string) error {

	// Create the https request

	folderUrl, err := url.Parse(folder)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", c.Url.ResolveReference(folderUrl).String(), nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.Username, c.Password)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(body) > 0 {
		error := Error{}
		err = xml.Unmarshal(body, &error)
		if err != nil {
			return err
		}
		if error.Exception != "" {
			return fmt.Errorf("Exception: %s, Message: %s", error.Exception, error.Message)
		}

	}

	return nil
}

func (c *Client) Upload(src []byte, dest string) error {

	destUrl, err := url.Parse(dest)
	if err != nil {
		return err
	}

	// Create the https request

	client := &http.Client{}
	req, err := http.NewRequest("PUT", c.Url.ResolveReference(destUrl).String(), bytes.NewReader(src))
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.Username, c.Password)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(body) > 0 {
		error := Error{}
		err = xml.Unmarshal(body, &error)
		if err != nil {
			return fmt.Errorf("Error during XML Unmarshal for response %s. The error is %s", body, err)
		}
		if error.Exception != "" {
			return fmt.Errorf("Exception: %s, Message: %s", error.Exception, error.Message)
		}

	}

	return nil
}

func (c *Client) Download(path string) ([]byte, error) {

	pathUrl, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	// Create the https request

	client := &http.Client{}
	req, err := http.NewRequest("GET", c.Url.ResolveReference(pathUrl).String(), nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Username, c.Password)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	error := Error{}
	err = xml.Unmarshal(body, &error)
	if err == nil {
		if error.Exception != "" {
			return nil, fmt.Errorf("Exception: %s, Message: %s", error.Exception, error.Message)
		}
	}

	return body, nil
}
