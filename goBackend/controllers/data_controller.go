package controllers

import (
	"encoding/json"
	"net/http"

    "iotDashboard/goBackend/models"
    "iotDashboard/goBackend/services"
)

// SensorDataHandler 保持不變
func SensorDataHandler(w http.ResponseWriter, r *http.Request) {
	services.GlobalDataStore.Mu.Lock()
	latestData := services.GlobalDataStore.Data
	services.GlobalDataStore.Mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
	// 允許 POST 方法和 Content-Type header，為前端做準備
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	json.NewEncoder(w).Encode(latestData)
}

// CommandHandler 是新的處理器，用來接收控制指令
func CommandHandler(w http.ResponseWriter, r *http.Request) {
	// 設定 CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// 瀏覽器在發送 POST 請求前，會先發送一個 OPTIONS "預檢"請求，我們直接回 200 OK
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 限制只能用 POST 方法發送指令
	if r.Method != http.MethodPost {
		http.Error(w, "僅允許 POST 方法", http.StatusMethodNotAllowed)
		return
	}

	var cmd models.Command
	// 解碼請求 body 中的 JSON 到 cmd 變數
	err := json.NewDecoder(r.Body).Decode(&cmd)
	if err != nil {
		http.Error(w, "無效的請求 body", http.StatusBadRequest)
		return
	}

	// 呼叫 service 層的函式來發送指令
	err = services.SendCommandToArduino(cmd.Command)
	if err != nil {
		http.Error(w, "發送指令失敗", http.StatusInternalServerError)
		return
	}

	// 回應成功訊息
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("指令已成功發送"))
}
