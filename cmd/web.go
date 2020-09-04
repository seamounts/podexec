package main

import "github.com/seamounts/pod-exec/cmd/app"

func main() {
	cmd := app.NewAppCMD()
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
