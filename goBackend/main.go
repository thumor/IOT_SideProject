package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"go.bug.st/serial"
	"iotDashboard/goBackend/controllers"
	"iotDashboard/goBackend/services"
)

func main() {
	// 序列埠配置 (模擬工業 Gateway)
	// 請確保 COM3 是您 Arduino 實際連接的編號
	const serialPortName = "COM3" 
	
	// 啟動背景序列埠監聽器 (非同步處理邊緣端遙測數據)
	go func() {
		mode := &serial.Mode{BaudRate: 9600}
		port, err := serial.Open(serialPortName, mode)
		if err != nil {
			log.Printf("[系統警告] 無法開啟序列埠 %s: %v", serialPortName, err)
			return
		}
		defer port.Close()

		log.Printf("[系統資訊] 工業數據網關已在 %s 啟動", serialPortName)

		scanner := bufio.NewScanner(port)
		for scanner.Scan() {
			line := scanner.Text()
			
			// 識別產線設備的遙測封包 (Telemetry Packet)
			if strings.HasPrefix(line, "TELEMETRY:") {
				jsonStr := strings.TrimPrefix(line, "TELEMETRY: ")
				
				var data struct {
					Temp float32 `json:"motorTemperature"`
				}

				if err := json.Unmarshal([]byte(jsonStr), &data); err == nil {
					log.Printf("[即時數據] 馬達溫度: %.1f C", data.Temp)
					
					// 核心邏輯：呼叫 Gemini 進行異常推理 (Anomaly Reasoning)
					advice := services.GetAIInstruction(data.Temp)
					log.Printf("[AI 建議] 推理結果: %s", advice)
					
					// 閉環回傳：將 AI 指令發送回 Arduino HMI 顯示器
					port.Write([]byte(advice + "\n"))
				}
			}
		}
	}()

	// 註冊工業級遙測 API 路由
	http.HandleFunc("/api/telemetry", controllers.TelemetryHandler)

	fmt.Println("==================================================")
	log.Println(" Industrial-IoT-LLM-Agent Orchestrator Started")
	log.Println(" 後端中控核心已啟動，監聽 8080 埠")
	fmt.Println("==================================================")

	// 啟動 Web 服務監聽 8080 埠
	log.Fatal(http.ListenAndServe(":8080", nil))
}