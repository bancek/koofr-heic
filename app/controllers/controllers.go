package controllers

import "github.com/revel/revel"

func init() {
	revel.InterceptFunc(setuser, revel.BEFORE, &App{})
}
