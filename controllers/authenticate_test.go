package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/quintilesims/auth0"
	"github.com/stretchr/testify/assert"
	auth "k8s.io/api/authentication/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newTokenReviewRequest(t *testing.T, token string) *http.Request {
	tr := auth.TokenReview{
		Spec: auth.TokenReviewSpec{Token: token},
	}

	b, err := json.Marshal(tr)
	if err != nil {
		t.Fatal(err)
	}

	return &http.Request{
		Body: ioutil.NopCloser(bytes.NewBuffer(b)),
	}
}

// Bad Request error returns correct token review
// GetProfile error returns token review
func TestAuthenticate(t *testing.T) {
	profile := &auth0.Profile{Email: "email@domain.com"}
	getProfile := func(token string) (*auth0.Profile, error) {
		assert.Equal(t, "some_token", token)
		return profile, nil
	}

	recorder := httptest.NewRecorder()
	r := newTokenReviewRequest(t, "some_token")
	NewAuthenticateController(getProfile).
		Authenticate(r).
		ServeHTTP(recorder, r)

	var result auth.TokenReview
	if err := json.NewDecoder(recorder.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	expected := auth.TokenReview{
		TypeMeta: meta.TypeMeta{
			APIVersion: "authentication.k8s.io/v1beta1",
			Kind:       "TokenReview",
		},
		Status: auth.TokenReviewStatus{
			Authenticated: true,
			User: auth.UserInfo{
				UID:      profile.Email,
				Username: profile.Email,
			},
		},
	}

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, expected, result)
}

func TestAuthenticateError_badRequestBody(t *testing.T) {
	recorder := httptest.NewRecorder()
	r := &http.Request{
		Body: ioutil.NopCloser(bytes.NewBuffer(nil)),
	}

	NewAuthenticateController(nil).
		Authenticate(r).
		ServeHTTP(recorder, r)

	var result auth.TokenReview
	if err := json.NewDecoder(recorder.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.False(t, result.Status.Authenticated)
	assert.NotNil(t, result.Status.Error)
}

func TestAuthenticateError_badToken(t *testing.T) {
	getProfile := func(token string) (*auth0.Profile, error) {
		return nil, errors.New("bad token")
	}

	recorder := httptest.NewRecorder()
	r := newTokenReviewRequest(t, "")
	NewAuthenticateController(getProfile).
		Authenticate(r).
		ServeHTTP(recorder, r)

	var result auth.TokenReview
	if err := json.NewDecoder(recorder.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.False(t, result.Status.Authenticated)
	assert.NotNil(t, result.Status.Error)
}
