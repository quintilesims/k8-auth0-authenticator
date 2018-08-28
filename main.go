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
Routes:
  / 			-> Has link for getting a token
  /authenticate		-> Does K8 Authentication
  /token		-> Does Auth0 handshake

Features:
  - Caching (w/ configurable timeout)
  - Logging
  - Dependency Inversion
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
		client := auth0.NewClient(c.String(FlagAuth0Domain))

		rootController := controllers.NewRootController(c.String(FlagAuth0Domain), c.String(FlagAuth0ClientID))
		tokenController := controllers.NewTokenController(client)

		routes := append(rootController.Routes(), tokenController.Routes()...)
		app := fireball.NewApp(routes)
		http.Handle("/", app)

		addr := fmt.Sprintf("0.0.0.0:%d", c.Int(FlagPort))
		log.Printf("[INFO] Listening on %s\n", addr)
		return http.ListenAndServe(addr, nil)
	}

	if err := app.Run(os.Args); err != nil {
		helpers.Exit(1, err)
	}
}
