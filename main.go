package main

import "github.com/prynnekey/gin-im/router"

func main() {
	r := router.Init()

	r.Run(":8080")
}
