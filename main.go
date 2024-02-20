package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jpcairesf/rinha-2024-q1-go/internal"
)

func main() {
	router := gin.Default()
	fmt.Println("Rinha is ready to Go")

	router.GET("/health", func(ctx *gin.Context) {
		ctx.Writer.Write([]byte("Rinha is running..."))
	})

	clientes := router.Group("/clientes/:id")
	{
		clientes.GET("/extrato", internal.GetExtrato)
		clientes.POST("/transacoes", internal.PostTransacao)
	}

	router.Run("8080")
}
