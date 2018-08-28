package controllers

import (
	"github.com/zpatrick/fireball"
)

type RootController struct{}

func NewRootController() *RootController {
	return &RootController{}
}

func (r *RootController) Routes() []*fireball.Route {
	routes := []*fireball.Route{
		{
			Path: "/",
			Handlers: fireball.Handlers{
				"GET": r.index,
			},
		},
	}

	return routes
}

func (r *RootController) index(c *fireball.Context) (fireball.Response, error) {
	return c.HTML(200, "index.html", nil)
}
