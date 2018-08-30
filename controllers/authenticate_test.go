package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/quintilesims/auth0"
	"github.com/stretchr/testify/assert"
	"github.com/zpatrick/fireball"
	authentication "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAuthenticateControllerAuthenticate(t *testing.T) {
	c := &fireball.Context{
		Request: newTokenReviewReq(t, "some_token"),
	}

	profile := auth0.Profile{Email: "email@domain.com"}
	controller := NewAuthenticateController(func(token string) (*auth0.Profile, error) {
		assert.Equal(t, "some_token", token)
		return &profile, nil
	})

	resp, err := controller.authenticate(c)
	if err != nil {
		t.Fatal(err)
	}

	var result authentication.TokenReview
	recorder := fireball.RecordJSONResponse(t, resp, &result)

	expected := authentication.TokenReview{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "authentication.k8s.io/v1beta1",
			Kind:       "TokenReview",
		},
		Status: authentication.TokenReviewStatus{
			Authenticated: true,
			User: authentication.UserInfo{
				UID:      profile.Email,
				Username: profile.Email,
			},
		},
	}

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, expected, result)
}

func TestAuthenticateControllerAuthenticateError_badRequestBody(t *testing.T) {
	c := &fireball.Context{
		Request: &http.Request{
			Body: ioutil.NopCloser(bytes.NewBuffer(nil)),
		},
	}

	controller := NewAuthenticateController(nil)
	resp, err := controller.authenticate(c)
	if err != nil {
		t.Fatal(err)
	}

	var result authentication.TokenReview
	recorder := fireball.RecordJSONResponse(t, resp, &result)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.False(t, result.Status.Authenticated)
}

func TestAuthenticateControllerAuthenticateError_getProfileError(t *testing.T) {
	c := &fireball.Context{
		Request: newTokenReviewReq(t, "some_token"),
	}

	controller := NewAuthenticateController(func(token string) (*auth0.Profile, error) {
		return nil, errors.New("some error")
	})

	resp, err := controller.authenticate(c)
	if err != nil {
		t.Fatal(err)
	}

	var result authentication.TokenReview
	recorder := fireball.RecordJSONResponse(t, resp, &result)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.False(t, result.Status.Authenticated)
}

func newTokenReviewReq(t *testing.T, token string) *http.Request {
	tr := authentication.TokenReview{
		Spec: authentication.TokenReviewSpec{Token: token},
	}

	b, err := json.Marshal(tr)
	if err != nil {
		t.Fatal(err)
	}

	return &http.Request{
		Body: ioutil.NopCloser(bytes.NewBuffer(b)),
	}
}
