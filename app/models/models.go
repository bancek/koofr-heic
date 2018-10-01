package models

import "github.com/revel/revel"

func init() {
	revel.OnAppStart(InitKoofrOAuthConfig)
}
