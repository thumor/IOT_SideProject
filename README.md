# Industrial-IoT-LLM-Agent: 工業級 GenAI 邊緣運算原型 (PoC)

這是一個結合了 IIoT (工業物聯網) 與 LLM Agent (大語言模型代理) 技術的邊緣運算原型專案。系統模擬產線馬達的運行狀態，透過 Go 後端驅動 Google Gemini API 進行專業的異常推理 (Anomaly Reasoning)，並將決策指令即時回饋至現場 HMI (人機介面)。

本專案旨在展示如何將 Generative AI 引入智慧製造場景，從傳統的「數據監控」轉向「自主決策輔助」，非常適合用於展示工業 4.0 轉型至 5.0 的技術整合能力。


## 成果展示

### 1. 現場 HMI 與硬體配置

這是我實際搭建的工業邊緣端原型，包含 Arduino Uno 與 I2C LCD 顯示器，模擬現場馬達狀態。
![Arduino 照片](https://github.com/user-attachments/assets/c09d27a2-7355-4a51-8338-4c5af31d7896)
(註：此處展示了當溫度異常時，LCD 第一行顯示即時數據，第二行即時更新來自 AI 的指令建議)

### 2. LLM 推理過程 (終端機輸出)

這是後端 Orchestrator 呼叫 Gemini API 並獲得推理建議的完整日誌紀錄。
![後端日誌](https://github.com/user-attachments/assets/178e8df8-f52a-48c0-8c59-9a5a86841518)    
(註：展示了 Telemetry 接收、AI 推理延遲以及最終生成的簡短指令)

## 核心功能 (Key Features)

LLM 異常推理 (Anomaly Reasoning)： 捨棄傳統的固定閥值 (Threshold) 判斷，利用 LLM 的語境感知能力，針對不同溫度波動給出具備「預測性」的操作指令。

工業級閉環控制 (Closed-loop Control)： 實現了從「設備端數據上報 -> 中控端 AI 決策 -> 設備端指令執行」的完整自動化閉環。

HMI 即時反饋： 專為工業現場設計，將複雜的 AI 推理結果壓縮為 12 字元以內的精確指令，直接顯示於現場顯示器。

邊緣網關設計 (Edge Gateway)： Go 後端扮演了工業網關的角色，同時處理非同步的序列通訊 (Serial Port) 與外部 AI API 串接。

穩定性機制： 在 AI 呼叫過程中實作了指數退避 (Exponential Backoff) 機制，確保在工業網路環境下的高可用性。

## 系統架構 (System Architecture)

專案採用三層式工業物聯網架構，強調邊緣與雲端決策的協同：

### 1. 邊緣感知層 (Edge/Perception Layer):

硬體： Arduino Uno (模擬馬達), I2C LCD 1602 (HMI)。

職責： 模擬馬達溫度數據上報，接收並呈現 AI 指令。

### 2. AI 中控層 (Orchestrator Layer):

後端： Go (Golang) 高性能後端。

職責： 透過 Goroutine 監聽序列埠數據，驅動 LLM 邏輯，並封裝 RESTful API 供前端監控。

### 3. 認知推理層 (Cognitive/AI Core):

引擎： Google Gemini 2.5 Flash。

職責： 根據 Prompt Engineering 進行工業專家級的決策判斷。

## 技術棧 (Technical Stack)

韌體 (Firmware): C/C++ (Arduino), LiquidCrystal_PCF8574 (HMI Control).

後端 (Backend): Go (Golang), Google Generative AI SDK, go.bug.st/serial (Industrial Serial Comm).

AI 技術： Prompt Engineering (System Instruction), Zero-shot Inference.

版本控制： Git & GitHub.

## 技術亮點與學習總結

### 1. 工業級 Prompt Engineering

為了讓 AI 建議能符合工業現場 16x2 LCD 的物理限制，我設計了特定的 System Prompt，強制模型在進行複雜推理後，僅輸出 1-3 個單詞的操作指令（如：CHECK COOLANT, EMERGENCY STOP），這體現了對硬體限制與使用者體驗的權衡。

### 2. Go 語言的高併發優勢

利用 Go 的 Goroutine 與 Mutex，我成功地讓系統能在背景持續掃描序列埠數據流，同時不影響對前端 Dashboard 的 API 回應速度。這種非同步設計模擬了真實工業環境中高吞吐量的資料傳輸需求。

### 3. 通訊協議的穩健性 (Robustness)

在 Go 後端中，我實作了錯誤處理與 API 重試機制。針對不穩定的網路環境，透過退避算法避免了單次失敗導致的生產線控制中斷，這是從「實驗室專案」邁向「工業應用」的關鍵細節。

## 本機安裝與執行 (Setup & Run)

### 軟硬體需求

Arduino Uno & I2C LCD 面板

Go (1.22+)

Google Gemini API Key

### 執行步驟

#### 韌體部署：

使用 Arduino IDE 上傳 Industrial-IoT-LLM-Agent.ino 至開發板。

#### 後端配置：

進入 goBackend 資料夾。

執行 go get 安裝 Gemini SDK 與序列埠套件。

執行 go run . 啟動 AI 中控核心。

#### 系統聯調：

觀察 Arduino LCD 上的溫度變動。

確認後端日誌是否有收到 TELEMETRY 封包，並產出 AI 建議。

## 未來擴展 (Roadmap)

多設備聯網 (Multi-node)： 導入 MQTT 協定，支持多個馬達節點同時由一個 AI Agent 控管。

邊緣模型部署： 嘗試將小型模型 (如 TinyML) 部署至邊緣端，減少對雲端 API 的依賴。