package services

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"strings"

	"go.bug.st/serial"
	"iotDashboard/goBackend/models"
)

var (
	port serial.Port
	GlobalDataStore struct {
		Mu   sync.Mutex
		Data models.SensorData
	}
	// isFanOn 變數現在用來追蹤 LCD 警告的狀態
	isFanOn bool
	fanControlMutex sync.Mutex
)

// SendCommandToArduino 函式負責向序列埠寫入指令
func SendCommandToArduino(command string) error {
	if port == nil {
		return errors.New("序列埠未初始化或未連接")
	}
	_, err := port.Write([]byte(command))
	if err != nil {
		log.Printf("向 Arduino 發送指令失敗: %v", err)
		return err
	}
	log.Printf("成功發送指令: %s", command)
	return nil
}

// 檢查溫度並決定是否要觸發 LCD 警告
func checkTemperatureAndManageWarning(temp float32) {
	fanControlMutex.Lock()
	defer fanControlMutex.Unlock()

	// 定義溫度的觸發閾值
	const upperThreshold float32 = 28.0
	const lowerThreshold float32 = 27.0

	// 如果溫度高於上限，且目前警告是關閉的
	if temp > upperThreshold && !isFanOn {
		log.Println("溫度過高！正在發送「LCD 顯示警告」指令...")
		err := SendCommandToArduino("W1")
		if err == nil {
			isFanOn = true
		}
	} else if temp < lowerThreshold && isFanOn {
		log.Println("溫度已降低。正在發送「LCD 清除警告」指令...")
		err := SendCommandToArduino("W0")
		if err == nil {
			isFanOn = false
		}
	}
}

// StartSerialReader 函式負責在背景持續讀取序列埠
func StartSerialReader(portName string) {
	mode := &serial.Mode{
		BaudRate: 9600,
	}

	var err error
	port, err = serial.Open(portName, mode)
	if err != nil {
		log.Fatalf("無法打開序列埠 %s: %v", portName, err)
	}
	log.Printf("成功打開序列埠 %s，開始監聽...", portName)

	scanner := bufio.NewScanner(port)
	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "{") || !strings.HasSuffix(line, "}") {
			log.Printf("收到格式不符的資料，已跳過: %s", line)
			continue
		}

		var data models.SensorData
		err := json.Unmarshal([]byte(line), &data)
		if err != nil {
			log.Printf("解析 JSON 失敗: %v, 原始資料: %s", err, line)
			continue
		}

		GlobalDataStore.Mu.Lock()
		GlobalDataStore.Data = data
		GlobalDataStore.Mu.Unlock()

		// 每次收到新數據，就呼叫決策函式
		go checkTemperatureAndManageWarning(data.Temperature)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("從序列埠讀取時發生錯誤: %v", err)
	}
}
