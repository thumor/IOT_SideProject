# 即時溫濕度監控儀表板 (Real-time IoT Dashboard)

這是一個結合了韌體開發、Go 後端與 Vue3 前端技術的全端物聯網 (IoT) Side Project。系統透過 Arduino 上的 DHT11 感測器採集我房間內的即時溫濕度數據，經由 Go 後端進行處理與決策，最終在一個動態的網頁儀表板上進行視覺化呈現與自動化控制。

這個專案的主要目的是深入學習 Go 語言在後端系統中的應用，並實踐韌體與後端服務之間的雙向通訊與閉環自動化控制。
![Arduino 照片](https://github.com/user-attachments/assets/1d16d430-528f-4dfb-8832-93da1705b1fe)

## 儀表板預覽
![儀錶板](https://github.com/user-attachments/assets/7328d9f9-62e8-44bf-8fd8-c5e87f8c0bdd)

## 核心功能 (Key Features)

即時數據監控： 透過網頁儀表板即時顯示感測器回傳的溫度與濕度數據。

歷史數據可視化： 以折線圖呈現歷史溫濕度走勢，便於觀察數據變化趨勢。

雙向通訊機制： 後端伺服器不僅能接收數據，更能主動向韌體下達控制指令 (如：手動觸發 LED 亮滅)。

閉環自動化控制： 後端具備決策能力，當溫度超過預設閾值時，會自動發送指令給 Arduino，觸發硬體警報 (LCD 顯示警告 & 板載 LED 閃爍)。

非阻塞式韌體設計： Arduino 程式碼採用 millis() 計時器，實現了多任務處理（數據採集、序列埠通訊、LED 閃爍），避免使用 delay() 造成的程式阻塞。

## 系統架構 (System Architecture)
本專案採用標準的 IoT 三層式架構，各層職責分明，透過序列埠與 RESTful API 進行通訊。

### 感知層 (Perception Layer):

硬體： Arduino Uno, DHT11 溫濕度感測器, I2C LCD 顯示器。

職責： 負責採集環境數據，並執行來自後端的指令。韌體採用 JSON 格式透過序列埠 (Serial Port) 將數據上報。

網路層 / 應用層 (Network / Application Layer):

後端： 使用 Go 語言開發。

職責： 作為系統的中樞。它同時執行兩個核心任務：

透過 Goroutine 持續監聽序列埠，即時接收並解析 Arduino 的數據。

建立一個 HTTP Web Server，對外提供 RESTful API，供前端獲取數據與發送控制指令。

### 表現層 (Presentation Layer):

前端： 使用 Vue3 框架開發。

職責： 向使用者呈現一個互動式的儀表板。它會定期向 Go 後端 API 請求最新數據來更新圖表，並提供介面讓使用者能下達控制指令。

技術棧 (Technical Stack)
韌體 (Firmware):

語言： C/C++ (Arduino)

硬體： Arduino Uno, DHT11, I2C LCD 1602

函式庫： DHT.h, LiquidCrystal_PCF8574.h

後端 (Backend):

語言： Go (Golang)

## 核心技術：

net/http: 建立 RESTful API 伺服器

encoding/json: 處理 JSON 數據

go.bug.st/serial: 進行序列埠通訊

Goroutines & Mutex: 實現高效率的並行處理與執行緒安全的數據存取

前端 (Frontend):

框架： Vue3 (Composition API)

工具： Vite, npm

圖表庫： Chart.js, vue-chartjs

版本控制： Git & GitHub

## 技術亮點與學習總結
這個專案不僅是一個功能展示，更是一次深入的技術實踐。以下是我在此專案中特別關注並實現的技術重點：

### 1. Go 後端：高並行與執行緒安全
為了同時處理「對硬體的持續監聽」和「對前端的 HTTP 請求」，我運用了 Go 語言最核心的並行 (Concurrency) 特性：

Goroutine: 我將序列埠的讀取任務放在一個獨立的 Goroutine 中執行，使其在背景運行而不阻塞主執行緒的 Web 服務。這完美地模擬了真實世界中需要同時處理多個 I/O 任務的後端架構。

Mutex (互斥鎖): 為了確保背景的讀取 Goroutine 和處理 API 請求的 Goroutine 在存取共享的感測器數據時不會發生衝突 (Race Condition)，我使用了 sync.Mutex 來保護共享變數，確保了數據的讀寫一致性與執行緒安全。

### 2. 韌體：非阻塞式設計與指令解析
為了讓 Arduino 能同時處理多個任務（定時上報數據、監聽指令、閃爍 LED），我避免了會「卡死」整個程式的 delay() 函式，轉而採用基於 millis() 的非阻塞式狀態機設計：

狀態機模型: 透過 isWarningActive 等狀態變數，讓 loop() 函式可以根據不同狀態執行不同任務（正常顯示 vs 警告模式），實現了高效的多工處理。

指令解析: 設計並實作了一個簡單的字元指令協定 (R, 1, 0, W1, W0)，讓韌體能解析來自後端的指令並執行對應的硬體操作，完成了從監控到控制的閉環。

### 3. 全端整合：從硬體到網頁的完整數據流
本專案打通了從硬體感測、序列埠通訊、後端處理、API 封裝到前端視覺化的完整數據鏈路。這讓我深刻理解到一個 IoT 產品從底層到上層是如何協同工作的，並學習到如何在不同技術棧之間定義清晰的通訊協定 (JSON) 來解耦系統。

## 未來可擴展方向 (Future Work)
資料庫整合： 將 Go 後端接收到的數據存入 InfluxDB 或 PostgreSQL 等時序資料庫，使歷史圖表能呈現真實的長期數據。
無線化 (Wireless): 將 Arduino Uno 替換為 ESP32 或 ESP8266，改用 Wi-Fi (MQTT 或 HTTP) 進行數據傳輸，使之成為一個真正的無線 IoT 節點。
二進位通訊協定： 設計並實作一個更高效的二進位通訊協定取代 JSON，以降低單晶片的運算負擔並提升傳輸效率。
前端功能增強： 在前端加入設定溫度閾值的介面，讓使用者可以透過網頁動態調整自動化控制的參數。

## 本機安裝與執行 (Setup & Run)

### 硬體需求
Arduino Uno x1
DHT11 溫濕度感測器 x1
I2C LCD 1602 顯示器 x1
杜邦線若干

### 軟體需求
Arduino IDE
Go (1.22 或以上版本)
Node.js (LTS 版本)
Git

### 執行步驟
Clone 專案：
git clone https://github.com/thumor/IOT_SideProject.git
cd iot-dashboard

韌體設定：
使用 Arduino IDE 打開 arduino-firmware/ 目錄下的 .ino 檔案。
安裝所需的函式庫 (DHT sensor library by Adafruit, LiquidCrystal_PCF8574)。
將程式碼上傳到您的 Arduino Uno。

後端啟動：
打開一個終端機，進入後端專案目錄。
cd go-backend

(首次執行) 安裝依賴套件：
go get go.bug.st/serial
修改 main.go 中的 serialPortName 為您 Arduino 的實際 COM Port 名稱。
啟動後端伺服器：
go run .
伺服器將運行在 http://localhost:8080。
前端啟動：
打開另一個新的終端機，進入前端專案目錄。
cd vue-frontend
(首次執行) 安裝依賴套件：
npm install
啟動前端開發伺服器：
npm run dev
前端頁面將運行在 http://localhost:5173 (或終端機提示的其他埠號)。
查看結果：
打開瀏覽器，訪問前端開發伺服器的網址，即可看到即時儀表板。
