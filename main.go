package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	controllers2 "github.com/jpcairesf/rinha-2024-q1-go/controllers"
)

var r *gin.Engine

func main() {
	r := gin.Default()
	fmt.Println("Rinha is now running...")

	r.GET("/health", func(c *gin.Context) {
		c.Writer.Write([]byte("Rinha is running..."))
	})

	clientes := r.Group("/clientes/:id")
	{
		clientes.GET("/extrato", controllers2.GetExtrato)
		clientes.POST("/transacoes", controllers2.PostTransacao)
	}

	r.Run("8080")
}
