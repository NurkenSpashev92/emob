// @title EMob API
// @version 1.0
// @description API documentation for EMob
// @host localhost:8080
// @BasePath /api/v1
package main

import "github.com/nurkenspashev92/emob/cmd/app"

func main() {
	app := new(app.App)
	app.Run()
}
