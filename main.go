package main

import (
	"sumwhere/app"
)

var ServiceVersion string

func main() {
	application := app.NewApp()
	if err := application.Run(); err != nil {
		panic(err)
	}
}
