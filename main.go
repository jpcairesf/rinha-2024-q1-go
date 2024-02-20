package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	controllers2 "github.com/jpcairesf/rinha-2024-q1-go/controllers"
)

func main() {
	router := gin.Default()
	fmt.Println("Rinha is now running...")

	router.GET("/health", func(ctx *gin.Context) {
		ctx.Writer.Write([]byte("Rinha is running..."))
	})

	clientes := router.Group("/clientes/:id")
	{
		clientes.GET("/extrato", controllers2.GetExtrato)
		clientes.POST("/transacoes", controllers2.PostTransacao)
	}

	router.Run("8080")
}
