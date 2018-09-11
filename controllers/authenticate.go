package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/quintilesims/auth0"
	"github.com/zpatrick/rye"
	auth "k8s.io/api/authentication/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AuthenticateController struct {
	getProfile func(string) (*auth0.Profile, error)
}

func NewAuthenticateController(getProfile func(string) (*auth0.Profile, error)) *AuthenticateController {
	return &AuthenticateController{
		getProfile: getProfile,
	}
}

func (a *AuthenticateController) Authenticate(r *http.Request) http.Handler {
	defer r.Body.Close()
	var tr auth.TokenReview
	if err := json.NewDecoder(r.Body).Decode(&tr); err != nil {
		log.Printf("[ERROR] Failed to decode TokenReview: %v", err)
		return tokenReview(auth.TokenReviewStatus{Error: err.Error()})
	}

	profile, err := a.getProfile(tr.Spec.Token)
	if err != nil {
		log.Printf("[ERROR] Failed to get profile for token '%s': %v", tr.Spec.Token, err)
		return tokenReview(auth.TokenReviewStatus{Error: err.Error()})
	}

	trs := auth.TokenReviewStatus{
		Authenticated: true,
		User: auth.UserInfo{
			UID:      profile.Email,
			Username: profile.Email,
		},
	}

	return tokenReview(trs)
}

func tokenReview(trs auth.TokenReviewStatus) http.Handler {
	status := http.StatusUnauthorized
	if trs.Authenticated {
		status = http.StatusOK
	}

	tr := auth.TokenReview{
		TypeMeta: meta.TypeMeta{
			APIVersion: "auth.k8s.io/v1beta1",
			Kind:       "TokenReview",
		},
		Status: trs,
	}

	return rye.JSON(status, tr)
}
