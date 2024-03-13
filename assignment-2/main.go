package main

import (
	"github.com/gin-gonic/gin"
	"github.com/refandas/scalable-web-service/assignment-2/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	r := gin.Default()

	r.POST("/orders", createOrder)
	r.GET("/orders", getOrder)
	r.PUT("/orders/:orderId", updateOrder)
	r.DELETE("/orders/:orderId", deleteOrder)

	if err := r.Run(); err != nil {
		log.Fatalf("error running server %s", err.Error())
	}
}

func createOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid order payload",
		})
		return
	}

	err := order.CreateOrder()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func getOrder(c *gin.Context) {
	orders, err := models.GetOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func updateOrder(c *gin.Context) {
	orderId := c.Param("orderId")
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid order payload",
		})
		return
	}

	order.ID, _ = strconv.ParseInt(orderId, 10, 64)

	updatedOrder, err := order.UpdateOrder()
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Order is not found",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, updatedOrder)
}

func deleteOrder(c *gin.Context) {
	orderId, err := strconv.ParseInt(c.Param("orderId"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	if err := models.DeleteOrder(orderId); err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "order is not found",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, "Success delete")
}
