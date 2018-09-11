package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/quintilesims/auth0"
	"github.com/quintilesims/k8-auth0-authenticator/controllers"
	"github.com/unrolled/render"
	"github.com/urfave/cli"
	"github.com/zpatrick/helpers"
	"github.com/zpatrick/router"
	"github.com/zpatrick/rye"
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
		client := auth0.NewClient(c.String(FlagAuth0Domain))
		render := render.New(render.Options{
			Directory:  "templates",
			Extensions: []string{".tmpl"},
			Layout:     "layout",
		})

		authenticateController := controllers.NewAuthenticateController(client.GetProfile)
		rootController := controllers.NewRootController(render, c.String(FlagAuth0Domain), c.String(FlagAuth0ClientID))
		tokenController := controllers.NewTokenController(render, client.GetProfile)

		rm := router.RouteMap{
			"/": router.MethodHandlers{
				http.MethodGet: rye.ToHandler(rootController.Root),
			},
			"/authenticate": router.MethodHandlers{
				http.MethodPost: rye.ToHandler(authenticateController.Authenticate),
			},
			"/token": router.MethodHandlers{
				http.MethodPost: rye.ToHandler(tokenController.HandleCallback),
			},
		}

		// todo: log middleware
		r := router.NewRouter(rm.StringMatch())
		http.Handle("/", r)

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
