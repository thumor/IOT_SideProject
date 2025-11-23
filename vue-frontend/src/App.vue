<!-- 檔案路徑: vue-frontend/src/App.vue -->
<script setup>
import { ref, onMounted, reactive } from 'vue'
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
} from 'chart.js'
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
)

// --- 即時數據狀態 ---
const temperature = ref(0)
const humidity = ref(0)
const lastUpdated = ref(null)
const error = ref(null)

// --- 圖表數據狀態 ---
const tempChartData = reactive({
  labels: [],
  datasets: [
    {
      label: '溫度 (°C)',
      backgroundColor: '#f87979',
      borderColor: '#f87979',
      data: [],
      tension: 0.3,
    },
  ],
})

const humiChartData = reactive({
  labels: [],
  datasets: [
    {
      label: '濕度 (%)',
      backgroundColor: '#3498db',
      borderColor: '#3498db',
      data: [],
      tension: 0.3,
    },
  ],
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
}

// --- 非同步函式 ---

const fetchRealtimeData = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/sensor-data')
    if (!response.ok) throw new Error('後端伺服器回應錯誤')
    const data = await response.json()
    temperature.value = data.temperature.toFixed(2)
    humidity.value = data.humidity.toFixed(2)
    lastUpdated.value = new Date().toLocaleTimeString()
    error.value = null
  } catch (e) {
    console.error(e)
    error.value = '無法獲取即時資料。請確認後端伺服器是否已啟動。'
  }
}

const fetchHistoryData = async () => {
  try {
    const response = await fetch('/mock-data.csv');
    const csvText = await response.text();
    
    const rows = csvText.trim().split('\n');
    // 確保 CSV 檔案不是空的
    if (rows.length <= 1) return;

    const headers = rows.shift().split(',');
    
    const labels = [];
    const temps = [];
    const humis = [];

    rows.forEach(row => {
      const values = row.split(',');
      if (values.length === headers.length) {
        const rowData = headers.reduce((obj, header, index) => {
          obj[header.trim()] = values[index].trim();
          return obj;
        }, {});
        
        labels.push(rowData.timestamp);
        temps.push(parseFloat(rowData.temperature));
        humis.push(parseFloat(rowData.humidity));
      }
    });

    tempChartData.labels = labels;
    tempChartData.datasets[0].data = temps;
    
    humiChartData.labels = labels;
    humiChartData.datasets[0].data = humis;

  } catch (e) {
    console.error("讀取歷史數據失敗:", e);
  }
}

// --- 生命週期鉤子 ---
onMounted(() => {
  fetchRealtimeData()
  fetchHistoryData()
  setInterval(fetchRealtimeData, 3000)
})
</script>

<template>
  <div id="dashboard">
    <header>
      <h1>房間溫濕度儀表板</h1>
      <p v-if="lastUpdated" class="last-updated">即時數據更新於: {{ lastUpdated }}</p>
    </header>
    
    <main>
      <div v-if="error" class="error-box">
        <p>{{ error }}</p>
      </div>
      
      <div v-else class="cards-container">
        <div class="card temperature-card">
          <h2>即時溫度</h2>
          <div class="value">
            <span>{{ temperature }}</span>
            <span class="unit">°C</span>
          </div>
        </div>
        
        <div class="card humidity-card">
          <h2>即時濕度</h2>
          <div class="value">
            <span>{{ humidity }}</span>
            <span class="unit">%</span>
          </div>
        </div>
      </div>

      <!-- 使用 v-if 確保在 labels 陣列被填充數據後，才渲染圖表容器 -->
      <div v-if="tempChartData.labels.length > 0" class="charts-container">
        <div class="chart-wrapper">
          <h3>歷史溫度走勢</h3>
          <Line :data="tempChartData" :options="chartOptions" />
        </div>
        <div class="chart-wrapper">
          <h3>歷史濕度走勢</h3>
          <Line :data="humiChartData" :options="chartOptions" />
        </div>
      </div>
    </main>
  </div>
</template>

<style scoped>
#dashboard {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
  max-width: 1000px;
  margin: 40px auto;
  padding: 20px;
  background-color: #f4f7f6;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  color: #333;
}

header {
  text-align: center;
  border-bottom: 1px solid #e0e0e0;
  padding-bottom: 20px;
  margin-bottom: 20px;
}

h1 {
  font-size: 2.5em;
  color: #2c3e50;
  font-weight: 600;
}

.last-updated {
  color: #7f8c8d;
  font-size: 0.9em;
}

.cards-container {
  display: flex;
  justify-content: space-around;
  gap: 20px;
  margin-bottom: 40px;
}

.card {
  background-color: #ffffff;
  padding: 30px;
  border-radius: 10px;
  text-align: center;
  flex: 1;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 15px rgba(0, 0, 0, 0.1);
}

h2 {
  margin-top: 0;
  font-size: 1.2em;
  color: #34495e;
  font-weight: 500;
}

.value {
  font-size: 4em;
  font-weight: bold;
}

.temperature-card .value {
  color: #e74c3c;
}

.humidity-card .value {
  color: #3498db;
}

.unit {
  font-size: 0.4em;
  vertical-align: super;
  margin-left: 5px;
}

.error-box {
  background-color: #ffebee;
  color: #c62828;
  border: 1px solid #ef9a9a;
  padding: 15px;
  border-radius: 8px;
  text-align: center;
}

.charts-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 30px;
}

.chart-wrapper {
  background: #fff;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  height: 400px;
}

.chart-wrapper h3 {
  text-align: center;
  margin-top: 0;
  margin-bottom: 15px;
  font-weight: 500;
  color: #34495e;
}
</style>
