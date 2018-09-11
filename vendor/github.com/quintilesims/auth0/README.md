# Auth0
[![Go Doc](https://godoc.org/github.com/quintilesims/auth0?status.svg)](https://godoc.org/github.com/quintilesims/auth0)

Minimalist client for the [Auth0 API](https://auth0.com/docs/api/info)

## Usage

```go
package main

import (
        "log"

        "github.com/quintilesims/auth0"
)

func main() {
        client := auth0.NewClient("domain.auth0.com")

        req := auth0.GetOAuthTokenRequest{
                GrantType:    "password",
                ClientID:     "client_id",
                ClientSecret: "client_secret",
                Username:     "username",
                Password:     "password",
        }

        resp, err := client.GetOAuthToken(req)
        if err != nil {
                log.Fatalf("Failed to get ID token: %v", err)
        }

        log.Printf("Access Token: %s", resp.AccessToken)
}
```
