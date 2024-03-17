package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type Monitor struct {
	Status `json:"status"`
}

func main() {
	ticker := time.NewTicker(15 * time.Second)

	refreshData()

	go func() {
		for range ticker.C {
			refreshData()
		}
	}()

	// Set up the services
	router := gin.Default()

	router.StaticFile("/", "index.html")
	router.GET("/data", getData)

	port := 8080
	err := router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalln("Error starting server", err)
	}
	log.Printf("Starting server on port %d\n", port)
}

func refreshData() {
	water := rand.Intn(100) + 1
	wind := rand.Intn(100) + 1

	data := Monitor{
		Status{
			Water: water,
			Wind:  wind,
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return
	}

	err = os.WriteFile("monitor.json", jsonData, 0644)
	if err != nil {
		log.Println("Error writing JSON file:", err)
		return
	}

	statusWater, statusWind := checkStatus(water, wind)
	log.Printf(
		"\nStatus water: %s | Ketinggian: %d meter \n"+
			"Status wind: %s | Kecepatan: %d meter/detik",
		statusWater, water, statusWind, wind,
	)
}

func checkStatus(water int, wind int) (string, string) {
	var waterStatus string
	if water < 5 {
		waterStatus = "aman"
	} else if water >= 6 && water <= 8 {
		waterStatus = "siaga"
	} else {
		waterStatus = "bahaya"
	}

	var windStatus string
	if wind < 6 {
		windStatus = "aman"
	} else if wind >= 7 && wind <= 15 {
		windStatus = "siaga"
	} else {
		windStatus = "bahaya"
	}

	return waterStatus, windStatus
}

func getData(c *gin.Context) {
	jsonData, err := os.ReadFile("monitor.json")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read JSON data",
		})
		return
	}

	//c.JSON(http.StatusOK, jsonData)
	c.Data(http.StatusOK, "application/json", jsonData)
}
