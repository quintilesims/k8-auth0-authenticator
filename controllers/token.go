package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/quintilesims/auth0"
	"github.com/unrolled/render"
	"github.com/zpatrick/rye"
)

type TokenController struct {
	render     *render.Render
	getProfile func(string) (*auth0.Profile, error)
}

func NewTokenController(r *render.Render, getProfile func(string) (*auth0.Profile, error)) *TokenController {
	return &TokenController{
		render:     r,
		getProfile: getProfile,
	}
}

func (tc *TokenController) HandleCallback(r *http.Request) http.Handler {
	accessToken := r.FormValue("access_token")
	if accessToken == "" {
		return rye.Error(400, fmt.Errorf("Required value 'access_token' not included in form"))
	}

	profile, err := tc.getProfile(accessToken)
	if err != nil {
		// todo: we should return 400 if the token is invalid, 500 otherwise
		log.Printf("[ERROR] Failed to get profile for token '%s': %v", accessToken, err)
		return rye.Error(400, err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Token   string
			Profile auth0.Profile
		}{
			Token:   accessToken,
			Profile: *profile,
		}

		if err := tc.render.HTML(w, http.StatusOK, "token", data); err != nil {
			log.Printf("[ERROR] Failed to render 'token' template: %v", err)
			http.Error(w, err.Error(), 500)
		}
	})
}
