package main

import (
	"fmt"
	"proxy-server/handler"
)

func main() {
	handler := handler.NewHandler()

	fmt.Printf("proxy server started -> http://localhost:8080")
	handler.InitRoutes()
}
