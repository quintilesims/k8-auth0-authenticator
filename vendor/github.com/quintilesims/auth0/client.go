package auth0

import (
	"encoding/json"
	"net/http"

	"github.com/zpatrick/rclient"
)

type Client struct {
	*rclient.RestClient
}

func NewClient(domain string) *Client {
	return &Client{
		rclient.NewRestClient("https://"+domain, rclient.Reader(auth0ResponseReader)),
	}
}

func (c *Client) GetProfile(accessToken string) (*Profile, error) {
	var resp Profile
	header := rclient.Header("Authorization", "Bearer "+accessToken)
	if err := c.Get("/userinfo", &resp, header); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) GetOAuthToken(req GetOAuthTokenRequest) (*GetOAuthTokenResponse, error) {
	var resp GetOAuthTokenResponse
	if err := c.Post("/oauth/token", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func auth0ResponseReader(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	switch {
	case resp.StatusCode < 200, resp.StatusCode > 299:
		var auth0Error Error
		if err := json.NewDecoder(resp.Body).Decode(&auth0Error); err != nil {
			return err
		}

		return auth0Error
	case v == nil:
		return nil
	default:
		return json.NewDecoder(resp.Body).Decode(v)
	}
}
