<template>
  <div class="request-logs">
    <div class="logs-header">
      <h3>Request Logs</h3>
      <button @click="clearLogs" class="clear-logs-btn" :disabled="logs.length === 0">
        Clear All
      </button>
    </div>

    <div class="logs-list">
      <div v-if="logs.length === 0" class="empty-logs">
        <p>No requests logged yet</p>
        <p class="empty-hint">API requests will appear here as you make them</p>
      </div>
      
      <div 
        v-for="log in logs" 
        :key="log.id"
        class="log-item"
        @click="loadLogRequest(log)"
      >
        <div class="log-header">
          <span :class="['method-badge', log.method.toLowerCase()]">
            {{ log.method }}
          </span>
          <span class="log-url">{{ log.url }}</span>
          <span :class="['status-badge', getStatusClass(log.status)]">
            {{ log.status }}
          </span>
        </div>
        
        <div class="log-details">
          <span class="log-time">{{ formatTime(log.timestamp) }}</span>
          <span class="log-duration">{{ log.duration }}ms</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { HTTPService } from '../../bindings/captain-api'

const emit = defineEmits(['load-request'])

const logs = ref([])
const loading = ref(false)

// Load logs from backend
const loadLogs = async () => {
  try {
    loading.value = true
    const requestLogs = await HTTPService.GetRequestLogs()
    // Convert timestamp strings to Date objects and map status from response
    logs.value = requestLogs.map(log => ({
      ...log,
      timestamp: new Date(log.timestamp),
      status: log.response.statusCode
    }))
  } catch (error) {
    console.error('Failed to load logs:', error)
    logs.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadLogs()
})

// Expose loadLogs method for parent components to refresh logs
defineExpose({
  loadLogs
})

const clearLogs = async () => {
  try {
    await HTTPService.ClearRequestLogs()
    logs.value = []
  } catch (error) {
    console.error('Failed to clear logs:', error)
  }
}

const loadLogRequest = (log) => {
  // Convert log to request format and emit complete log data
  const request = {
    method: log.request.method,
    url: log.request.url,
    headers: log.request.headers,
    body: log.request.body
  }
  // Emit both request and response data
  emit('load-request', { request, response: log.response })
}

const getStatusClass = (status) => {
  if (status >= 200 && status < 300) return 'success'
  if (status >= 300 && status < 400) return 'redirect'
  if (status >= 400 && status < 500) return 'client-error'
  if (status >= 500) return 'server-error'
  return 'unknown'
}

const formatTime = (timestamp) => {
  return new Intl.DateTimeFormat('en-US', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  }).format(timestamp)
}
</script>

<style scoped>
.request-logs {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.logs-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  border-bottom: 1px solid #e5e7eb;
  background: #f9fafb;
}

.logs-header h3 {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: #374151;
}

.clear-logs-btn {
  padding: 0.5rem 1rem;
  background: #ef4444;
  color: white;
  border: none;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.clear-logs-btn:hover:not(:disabled) {
  background: #dc2626;
}

.clear-logs-btn:disabled {
  background: #d1d5db;
  cursor: not-allowed;
}

.logs-list {
  flex: 1;
  overflow-y: auto;
  padding: 0.5rem;
}

.empty-logs {
  text-align: center;
  padding: 2rem;
  color: #6b7280;
}

.empty-logs p {
  margin: 0.5rem 0;
}

.empty-hint {
  font-size: 0.875rem;
  color: #9ca3af;
}

.log-item {
  padding: 0.75rem;
  margin-bottom: 0.5rem;
  border: 1px solid #e5e7eb;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: all 0.2s;
  background: white;
}

.log-item:hover {
  border-color: #3b82f6;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.log-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
}

.method-badge {
  padding: 0.25rem 0.5rem;
  border-radius: 0.25rem;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  min-width: 3.5rem;
  text-align: center;
}

.method-badge.get {
  background: #dcfce7;
  color: #166534;
}

.method-badge.post {
  background: #dbeafe;
  color: #1e40af;
}

.method-badge.put {
  background: #fef3c7;
  color: #92400e;
}

.method-badge.patch {
  background: #f3e8ff;
  color: #7c3aed;
}

.method-badge.delete {
  background: #fee2e2;
  color: #dc2626;
}

.log-url {
  flex: 1;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 0.875rem;
  color: #374151;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.status-badge {
  padding: 0.25rem 0.5rem;
  border-radius: 0.25rem;
  font-size: 0.75rem;
  font-weight: 600;
  min-width: 2.5rem;
  text-align: center;
}

.status-badge.success {
  background: #dcfce7;
  color: #166534;
}

.status-badge.redirect {
  background: #fef3c7;
  color: #92400e;
}

.status-badge.client-error {
  background: #fee2e2;
  color: #dc2626;
}

.status-badge.server-error {
  background: #fecaca;
  color: #991b1b;
}

.status-badge.unknown {
  background: #f3f4f6;
  color: #6b7280;
}

.log-details {
  display: flex;
  justify-content: space-between;
  font-size: 0.75rem;
  color: #6b7280;
}

.log-time {
  font-family: 'Monaco', 'Menlo', monospace;
}

.log-duration {
  font-weight: 500;
}
</style>