package main

import "sumwhere/app"

func main() {
	application := app.NewApp()
	if err := application.Run(); err != nil {
		panic(err)
	}
}
