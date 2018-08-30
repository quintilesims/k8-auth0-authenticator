package controllers

import (
	"fmt"

	"github.com/quintilesims/auth0"
	"github.com/zpatrick/fireball"
)

type TokenController struct {
	getProfile func(string) (*auth0.Profile, error)
}

func NewTokenController(getProfile func(string) (*auth0.Profile, error)) *TokenController {
	return &TokenController{
		getProfile: getProfile,
	}
}

func (t *TokenController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/token",
			Handlers: fireball.Handlers{
				"POST": t.postToken,
			},
		},
	}
}

func (t *TokenController) postToken(c *fireball.Context) (fireball.Response, error) {
	accessToken := c.Request.FormValue("access_token")
	if accessToken == "" {
		return nil, fmt.Errorf("Required value 'access_token' not included in form")
	}

	profile, err := t.getProfile(accessToken)
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
