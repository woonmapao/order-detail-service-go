package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello Mon")

	r := gin.Default()
	r.Run()
}
