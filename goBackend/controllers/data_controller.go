package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	// 確保路徑與您的 go.mod 一致
	"iotDashboard/goBackend/services"
)

// TelemetryHandler 處理來自前端或 HMI 的數據請求
func TelemetryHandler(w http.ResponseWriter, r *http.Request) {
	// 處理跨域請求
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		return
	}

	var input struct {
		Temp float32 `json:"motorTemperature"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid Payload", http.StatusBadRequest)
		return
	}

	// 呼叫 Gemini AI 推理邏輯
	instruction := services.GetAIInstruction(input.Temp)
	
	log.Printf("[SYSTEM] 遙測數據接收: %.1f C -> AI 決策指令: %s", input.Temp, instruction)

	// 回傳 AI 的操作指令
	json.NewEncoder(w).Encode(map[string]string{
		"operatorInstruction": instruction,
	})
}