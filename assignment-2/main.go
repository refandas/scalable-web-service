package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/refandas/scalable-web-service/assignment-2/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

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
	var order models.CreateOrderInput
	if err := c.ShouldBindJSON(&order); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid order payload",
		})
		return
	}

	if err := validate.Struct(order); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid order payload",
		})
		return
	}

	err := order.CreateOrder()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func getOrder(c *gin.Context) {
	orders, err := models.GetOrders()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func updateOrder(c *gin.Context) {
	orderIdPayload := c.Param("orderId")
	var order models.UpdateOrderInput
	if err := c.ShouldBindJSON(&order); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid order payload",
		})
		return
	}

	if err := validate.Struct(order); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid order payload",
		})
		return
	}

	orderId, _ := strconv.ParseInt(orderIdPayload, 10, 64)

	updatedOrder, err := order.UpdateOrder(orderId)
	if err != nil {
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

	c.JSON(http.StatusOK, updatedOrder)
}

func deleteOrder(c *gin.Context) {
	orderId, err := strconv.ParseInt(c.Param("orderId"), 10, 64)
	if err != nil || orderId <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid order ID",
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
