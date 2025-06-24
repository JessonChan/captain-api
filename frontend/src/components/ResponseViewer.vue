<template>
  <div class="response-viewer">
    <div v-if="!response" class="no-response">
      <p>Send a request to see the response here</p>
    </div>
    
    <div v-else class="response-content">
      <!-- Response Status -->
      <div class="response-status">
        <div class="status-info">
          <span :class="['status-code', getStatusClass(response.statusCode)]">
            {{ response.statusCode }}
          </span>
          <span class="status-text">{{ response.status }}</span>
        </div>
        <div class="response-meta">
          <span class="duration">{{ response.duration }}ms</span>
          <span class="size">{{ formatSize(response.size) }}</span>
        </div>
      </div>

      <!-- Response Tabs -->
      <div class="response-tabs">
        <button 
          v-for="tab in responseTabs" 
          :key="tab" 
          @click="activeTab = tab"
          :class="['tab-btn', { active: activeTab === tab }]"
        >
          {{ tab }}
          <span v-if="tab === 'Headers'" class="tab-count">
            ({{ Object.keys(response.headers || {}).length }})
          </span>
        </button>
      </div>

      <!-- Tab Content -->
      <div class="tab-content">
        <!-- Body Tab -->
        <div v-if="activeTab === 'Body'" class="body-section">
          <div class="body-controls">
            <button 
              @click="formatResponse"
              :disabled="!canFormat"
              class="format-btn"
            >
              Format JSON
            </button>
            <button @click="copyResponse" class="copy-btn">
              Copy Response
            </button>
          </div>
          <pre class="response-body">{{ displayBody }}</pre>
        </div>

        <!-- Headers Tab -->
        <div v-if="activeTab === 'Headers'" class="headers-section">
          <div class="headers-list">
            <div 
              v-for="(value, key) in response.headers" 
              :key="key"
              class="header-item"
            >
              <span class="header-key">{{ key }}:</span>
              <span class="header-value">{{ value }}</span>
            </div>
          </div>
        </div>

        <!-- Raw Tab -->
        <div v-if="activeTab === 'Raw'" class="raw-section">
          <pre class="raw-response">{{ rawResponse }}</pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { HTTPService } from '../../bindings/captain-api'

const props = defineProps({
  response: {
    type: Object,
    default: null
  }
})

const activeTab = ref('Body')
const responseTabs = ref(['Body', 'Headers', 'Raw'])
const formattedBody = ref('')
const isFormatted = ref(false)

// Watch for response changes
watch(() => props.response, (newResponse) => {
  if (newResponse) {
    activeTab.value = 'Body'
    formattedBody.value = ''
    isFormatted.value = false
  }
}, { immediate: true })

// Computed properties
const displayBody = computed(() => {
  if (!props.response) return ''
  return isFormatted.value ? formattedBody.value : props.response.body
})

const canFormat = computed(() => {
  if (!props.response?.body) return false
  try {
    JSON.parse(props.response.body)
    return true
  } catch {
    return false
  }
})

const rawResponse = computed(() => {
  if (!props.response) return ''
  
  let raw = `HTTP/1.1 ${props.response.statusCode} ${props.response.status}\n`
  
  // Add headers
  Object.entries(props.response.headers || {}).forEach(([key, value]) => {
    raw += `${key}: ${value}\n`
  })
  
  raw += '\n'
  raw += props.response.body
  
  return raw
})

// Methods
const getStatusClass = (statusCode) => {
  if (statusCode >= 200 && statusCode < 300) return 'success'
  if (statusCode >= 300 && statusCode < 400) return 'redirect'
  if (statusCode >= 400 && statusCode < 500) return 'client-error'
  if (statusCode >= 500) return 'server-error'
  return 'unknown'
}

const formatSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatResponse = async () => {
  if (!canFormat.value) return
  
  try {
    const formatted = await HTTPService.FormatJSON(props.response.body)
    formattedBody.value = formatted
    isFormatted.value = true
  } catch (error) {
    console.error('Failed to format JSON:', error)
  }
}

const copyResponse = async () => {
  try {
    await navigator.clipboard.writeText(displayBody.value)
    // You could add a toast notification here
    console.log('Response copied to clipboard')
  } catch (error) {
    console.error('Failed to copy response:', error)
  }
}
</script>

<style scoped>
.response-viewer {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  height: 100%;
  display: flex;
  flex-direction: column;
}

.no-response {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 200px;
  color: #6c757d;
  font-style: italic;
}

.response-content {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.response-status {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  background: #f8f9fa;
  border-radius: 6px;
  margin-bottom: 20px;
}

.status-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.status-code {
  padding: 4px 8px;
  border-radius: 4px;
  font-weight: bold;
  font-size: 14px;
}

.status-code.success {
  background: #d4edda;
  color: #155724;
}

.status-code.redirect {
  background: #fff3cd;
  color: #856404;
}

.status-code.client-error {
  background: #f8d7da;
  color: #721c24;
}

.status-code.server-error {
  background: #f5c6cb;
  color: #721c24;
}

.status-code.unknown {
  background: #e2e3e5;
  color: #383d41;
}

.status-text {
  font-weight: 500;
  color: #495057;
}

.response-meta {
  display: flex;
  gap: 15px;
  font-size: 14px;
  color: #6c757d;
}

.duration::before {
  content: '⏱ ';
}

.size::before {
  content: '📦 ';
}

.response-tabs {
  display: flex;
  border-bottom: 1px solid #ddd;
  margin-bottom: 20px;
}

.tab-btn {
  padding: 10px 20px;
  border: none;
  background: none;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  display: flex;
  align-items: center;
  gap: 5px;
}

.tab-btn.active {
  border-bottom-color: #007bff;
  color: #007bff;
  font-weight: 500;
}

.tab-count {
  font-size: 12px;
  color: #6c757d;
}

.tab-content {
  flex: 1;
  overflow: auto;
}

.body-controls {
  display: flex;
  gap: 10px;
  margin-bottom: 15px;
}

.format-btn, .copy-btn {
  padding: 6px 12px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.format-btn:hover:not(:disabled), .copy-btn:hover {
  background: #f8f9fa;
}

.format-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.response-body, .raw-response {
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 4px;
  padding: 15px;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.4;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 400px;
  overflow: auto;
  margin: 0;
}

.headers-list {
  max-height: 400px;
  overflow: auto;
}

.header-item {
  display: flex;
  padding: 8px 0;
  border-bottom: 1px solid #f1f3f4;
}

.header-item:last-child {
  border-bottom: none;
}

.header-key {
  font-weight: 600;
  color: #495057;
  min-width: 200px;
  margin-right: 15px;
}

.header-value {
  color: #6c757d;
  word-break: break-all;
}
</style>