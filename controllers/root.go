package controllers

import (
	"log"
	"net/http"

	"github.com/unrolled/render"
)

type RootController struct {
	render        *render.Render
	auth0Domain   string
	auth0ClientID string
}

func NewRootController(r *render.Render, auth0Domain, auth0ClientID string) *RootController {
	return &RootController{
		render:        r,
		auth0Domain:   auth0Domain,
		auth0ClientID: auth0ClientID,
	}
}

func (rc *RootController) Root(r *http.Request) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Auth0Domain   string
			Auth0ClientID string
		}{
			Auth0Domain:   rc.auth0Domain,
			Auth0ClientID: rc.auth0ClientID,
		}

		if err := rc.render.HTML(w, http.StatusOK, "root", data); err != nil {
			log.Printf("[ERROR] Failed to render 'root' template: %v", err)
		}
	})
}
