package controllers

import (
	"fmt"

	"github.com/quintilesims/auth0"
	"github.com/zpatrick/fireball"
)

// todo: use dependency inversion
type TokenController struct {
	client *auth0.Client
}

func NewTokenController(client *auth0.Client) *TokenController {
	return &TokenController{
		client: client,
	}
}

func (t *TokenController) Routes() []*fireball.Route {
	routes := []*fireball.Route{
		{
			Path: "/token",
			Handlers: fireball.Handlers{
				"POST": t.getToken,
			},
		},
	}

	return routes
}

func (t *TokenController) getToken(c *fireball.Context) (fireball.Response, error) {
	accessToken := c.Request.FormValue("access_token")
	if accessToken == "" {
		return nil, fmt.Errorf("Required value 'access_token' not included in form")
	}

	profile, err := t.client.GetProfile(accessToken)
	if err != nil {
		return nil, err
	}

	data := struct {
		Token   string
		Profile auth0.Profile
	}{
		Token:   accessToken,
		Profile: *profile,
	}

	return c.HTML(200, "token.html", data)
}
