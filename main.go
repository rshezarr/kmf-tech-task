package main

import "proxy-server/handler"

func main() {
	handler := handler.NewHandler()

	handler.InitRoutes()
}
