package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/quintilesims/auth0"
	"github.com/zpatrick/fireball"
	authentication "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AuthenticateController struct {
	getProfile func(string) (*auth0.Profile, error)
}

func NewAuthenticateController(getProfile func(string) (*auth0.Profile, error)) *AuthenticateController {
	return &AuthenticateController{
		getProfile: getProfile,
	}
}

func (a *AuthenticateController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/authenticate",
			Handlers: fireball.Handlers{
				"POST": a.authenticate,
			},
		},
	}
}

func (a *AuthenticateController) authenticate(c *fireball.Context) (fireball.Response, error) {
	defer c.Request.Body.Close()

	var tr authentication.TokenReview
	if err := json.NewDecoder(c.Request.Body).Decode(&tr); err != nil {
		log.Printf("[ERROR] Failed to decode TokenReview: %v", err)
		return tokenReview(authentication.TokenReviewStatus{Error: err.Error()})
	}

	profile, err := a.getProfile(tr.Spec.Token)
	if err != nil {
		log.Printf("[ERROR] Failed to get Auth0 profile: %v", err)
		return tokenReview(authentication.TokenReviewStatus{Error: err.Error()})
	}

	status := authentication.TokenReviewStatus{
		Authenticated: true,
		User: authentication.UserInfo{
			UID:      profile.Email,
			Username: profile.Email,
		},
	}

	return tokenReview(status)
}

func tokenReview(status authentication.TokenReviewStatus) (fireball.Response, error) {
	statusCode := http.StatusUnauthorized
	if status.Authenticated {
		statusCode = http.StatusOK
	}

	resp := authentication.TokenReview{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "authentication.k8s.io/v1beta1",
			Kind:       "TokenReview",
		},
		Status: status,
	}

	return fireball.NewJSONResponse(statusCode, resp)
}
