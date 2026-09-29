package main

//демо сервис ресторана --- отдает ошибки ~20%

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateOrderRequest struct {
	OrderID string      `json:"order_id" binding:"required"`
	Items   []OrderItem `json:"items" binding:"required,min=1"`
	Total   int64       `json:"total_price" binding:"required,min=0"`
}

type OrderItem struct {
	Name       string `json:"name" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
	UnitPrice  int64  `json:"unit_price" binding:"required,min=0"`
	TotalPrice int64  `json:"total_price" binding:"required,min=0"`
}

type CreateOrderResponse struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

func main() {
	router := gin.Default()

	group := router.Group("/api/v1")
	group.POST("/orders", createOrder)

	if err := router.Run(":8081"); err != nil {
		panic(err)
	}
}

func createOrder(c *gin.Context) {
	var req CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		log.Println("error json bind")
		return
	}

	// Имитируем временную ошибку
	if rand.Intn(100) < 20 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "restaurant is temporarily unavailable",
		})
		log.Println("error rand 20")
		return
	}

	c.JSON(http.StatusOK, CreateOrderResponse{
		OrderID: req.OrderID,
		Status:  "accepted",
	})
}
