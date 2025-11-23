package main

import (
	"log"
	"net/http"

    "iotDashboard/goBackend/controllers"
    "iotDashboard/goBackend/services"
)

func main() {
	const serialPortName = "COM3"

	go services.StartSerialReader(serialPortName)

	http.HandleFunc("/api/sensor-data", controllers.SensorDataHandler)
	http.HandleFunc("/api/command", controllers.CommandHandler)

	log.Println("伺服器已啟動於 http://localhost:8080")
	log.Println("GET /api/sensor-data - 獲取感測器資料")
	log.Println("POST /api/command - 發送控制指令")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
