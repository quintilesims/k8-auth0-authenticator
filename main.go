package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/quintilesims/auth0"
	"github.com/quintilesims/k8-auth0-authenticator/controllers"
	"github.com/urfave/cli"
	"github.com/zpatrick/fireball"
	"github.com/zpatrick/helpers"
)

/*
TODO:
  - Make auth0 lock configurable?
  - Docs for setting up w/ minicube (use ngrok?)
*/

const (
	FlagPort          = "port"
	FlagDebug         = "debug"
	FlagAuth0Domain   = "auth0-domain"
	FlagAuth0ClientID = "auth0-client-id"
)

const (
	EVPort          = "KAA_PORT"
	EVDebug         = "KAA_DEBUG"
	EVAuth0Domain   = "KAA_AUTH0_DOMAIN"
	EVAuth0ClientID = "KAA_AUTH0_CLIENT_ID"
)

func main() {
	app := cli.NewApp()
	app.Name = "k8-auth0-authenticator"
	app.Usage = "Auth0 Authenticator for Kubernetes"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:   FlagPort,
			EnvVar: EVPort,
			Value:  80,
		},
		cli.BoolFlag{
			Name:   FlagDebug,
			EnvVar: EVDebug,
		},
		cli.StringFlag{
			Name:   FlagAuth0Domain,
			EnvVar: EVAuth0Domain,
		},
		cli.StringFlag{
			Name:   FlagAuth0ClientID,
			EnvVar: EVAuth0ClientID,
		},
	}

	app.Before = func(c *cli.Context) error {
		requiredFlags := []string{
			FlagAuth0Domain,
			FlagAuth0ClientID,
		}

		for _, flag := range requiredFlags {
			if !c.IsSet(flag) {
				return fmt.Errorf("Required flag '%s' is not set!", flag)
			}
		}

		w := helpers.NewLogWriter(c.Bool(FlagDebug))
		log.SetOutput(w)

		return nil
	}

	app.Action = func(c *cli.Context) error {
		routes := []*fireball.Route{}
		client := auth0.NewClient(c.String(FlagAuth0Domain))
		getProfile := func(token string) (*auth0.Profile, error) {
			log.Printf("[DEBUG] Attempting to validate token '%s'", token)
			profile, err := client.GetProfile(token)
			if err != nil {
				log.Printf("[ERROR] Failed to validate token '%s': %v", token, err)
				return nil, err
			}

			log.Printf("[DEBUG] Successfully validated token '%s' (owner: '%s')", token, profile.Email)
			return profile, nil
		}

		authenticateController := controllers.NewAuthenticateController(getProfile)
		routes = append(routes, authenticateController.Routes()...)

		rootController := controllers.NewRootController(c.String(FlagAuth0Domain), c.String(FlagAuth0ClientID))
		routes = append(routes, rootController.Routes()...)

		tokenController := controllers.NewTokenController(getProfile)
		routes = append(routes, tokenController.Routes()...)
		routes = fireball.Decorate(routes, fireball.LogDecorator())

		app := fireball.NewApp(routes)
		http.Handle("/", app)

		http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "static/favicon.png")
		})

		addr := fmt.Sprintf("0.0.0.0:%d", c.Int(FlagPort))
		log.Printf("[INFO] Listening on %s\n", addr)
		return http.ListenAndServe(addr, nil)
	}

	if err := app.Run(os.Args); err != nil {
		helpers.Exit(1, err)
	}
}
