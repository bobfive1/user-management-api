package main

import "github.com/bobfive1/user-management-api/internal/app"

func main() {
	apiServer, wg := app.Bootstrap()

	go apiServer.Start()
	wg.Wait()

}
