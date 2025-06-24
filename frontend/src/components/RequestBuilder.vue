<template>
  <div class="request-builder">
    <div class="request-header">
      <div class="method-url-row">
        <select v-model="request.method" class="method-select">
          <option v-for="method in methods" :key="method" :value="method">
            {{ method }}
          </option>
        </select>
        <input 
          v-model="request.url" 
          type="text" 
          placeholder="Enter request URL" 
          class="url-input"
          @keyup.enter="sendRequest"
        />
        <button @click="sendRequest" :disabled="loading" class="send-btn">
          {{ loading ? 'Sending...' : 'Send' }}
        </button>
      </div>
    </div>

    <div class="request-tabs">
      <button 
        v-for="tab in tabs" 
        :key="tab" 
        @click="activeTab = tab"
        :class="['tab-btn', { active: activeTab === tab }]"
      >
        {{ tab }}
      </button>
    </div>

    <div class="tab-content">
      <!-- Headers Tab -->
      <div v-if="activeTab === 'Headers'" class="headers-section">
        <div class="header-row" v-for="(header, index) in headersList" :key="index">
          <input 
            v-model="header.key" 
            type="text" 
            placeholder="Header name"
            class="header-input"
          />
          <input 
            v-model="header.value" 
            type="text" 
            placeholder="Header value"
            class="header-input"
          />
          <button @click="removeHeader(index)" class="remove-btn">×</button>
        </div>
        <button @click="addHeader" class="add-btn">+ Add Header</button>
      </div>

      <!-- Body Tab -->
      <div v-if="activeTab === 'Body'" class="body-section">
        <div class="body-type-selector">
          <label>
            <input type="radio" v-model="bodyType" value="none" /> None
          </label>
          <label>
            <input type="radio" v-model="bodyType" value="json" /> JSON
          </label>
          <label>
            <input type="radio" v-model="bodyType" value="text" /> Text
          </label>
        </div>
        <textarea 
          v-if="bodyType !== 'none'"
          v-model="request.body"
          placeholder="Enter request body"
          class="body-textarea"
          rows="10"
        ></textarea>
      </div>
    </div>

    <!-- Save Request Section -->
    <div class="save-section">
      <input 
        v-model="requestName" 
        type="text" 
        placeholder="Request name (optional)"
        class="name-input"
      />
      <button @click="saveRequest" class="save-btn">Save Request</button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { HTTPService, CollectionService } from '../../bindings/changeme'

const emit = defineEmits(['response-received'])

const request = ref({
  method: 'GET',
  url: '',
  headers: {},
  body: ''
})

const methods = ref(['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS'])
const activeTab = ref('Headers')
const tabs = ref(['Headers', 'Body'])
const loading = ref(false)
const bodyType = ref('none')
const requestName = ref('')

// Headers management
const headersList = ref([{ key: '', value: '' }])

const addHeader = () => {
  headersList.value.push({ key: '', value: '' })
}

const removeHeader = (index) => {
  if (headersList.value.length > 1) {
    headersList.value.splice(index, 1)
  }
}

// Convert headers list to object
const headersObject = computed(() => {
  const headers = {}
  headersList.value.forEach(header => {
    if (header.key && header.value) {
      headers[header.key] = header.value
    }
  })
  return headers
})

// Watch body type changes
watch(bodyType, (newType) => {
  if (newType === 'json' && !request.value.body) {
    request.value.body = '{\n  \n}'
  } else if (newType === 'none') {
    request.value.body = ''
  }
})

// Send HTTP request
const sendRequest = async () => {
  if (!request.value.url) {
    alert('Please enter a URL')
    return
  }

  loading.value = true
  
  try {
    const requestData = {
      method: request.value.method,
      url: request.value.url,
      headers: headersObject.value,
      body: bodyType.value === 'none' ? '' : request.value.body
    }

    const response = await HTTPService.SendRequest(requestData)
    emit('response-received', response)
  } catch (error) {
    console.error('Request failed:', error)
    emit('response-received', {
      statusCode: 0,
      status: 'Error',
      headers: {},
      body: error.message || 'Request failed',
      duration: 0,
      size: 0
    })
  } finally {
    loading.value = false
  }
}

// Save request to collection
const saveRequest = async () => {
  try {
    const requestItem = {
      id: '',
      name: requestName.value || `${request.value.method} ${request.value.url}`,
      method: request.value.method,
      url: request.value.url,
      headers: headersObject.value,
      body: bodyType.value === 'none' ? '' : request.value.body,
      description: ''
    }

    await CollectionService.SaveRequest('default', requestItem)
    alert('Request saved successfully!')
    requestName.value = ''
  } catch (error) {
    console.error('Failed to save request:', error)
    alert('Failed to save request')
  }
}

// Load request from props
const loadRequest = (requestData) => {
  request.value = {
    method: requestData.method || 'GET',
    url: requestData.url || '',
    headers: requestData.headers || {},
    body: requestData.body || ''
  }

  // Convert headers object to list
  headersList.value = Object.entries(request.value.headers).map(([key, value]) => ({ key, value }))
  if (headersList.value.length === 0) {
    headersList.value = [{ key: '', value: '' }]
  }

  // Set body type
  if (request.value.body) {
    try {
      JSON.parse(request.value.body)
      bodyType.value = 'json'
    } catch {
      bodyType.value = 'text'
    }
  } else {
    bodyType.value = 'none'
  }
}

defineExpose({ loadRequest })
</script>

<style scoped>
.request-builder {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.method-url-row {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.method-select {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: white;
  min-width: 100px;
}

.url-input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.send-btn {
  padding: 8px 20px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
}

.send-btn:hover:not(:disabled) {
  background: #0056b3;
}

.send-btn:disabled {
  background: #6c757d;
  cursor: not-allowed;
}

.request-tabs {
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
}

.tab-btn.active {
  border-bottom-color: #007bff;
  color: #007bff;
  font-weight: 500;
}

.header-row {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
  align-items: center;
}

.header-input {
  flex: 1;
  padding: 6px 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.remove-btn {
  padding: 6px 10px;
  background: #dc3545;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.add-btn {
  padding: 8px 16px;
  background: #28a745;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.body-type-selector {
  display: flex;
  gap: 20px;
  margin-bottom: 15px;
}

.body-type-selector label {
  display: flex;
  align-items: center;
  gap: 5px;
  cursor: pointer;
}

.body-textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  resize: vertical;
}

.save-section {
  display: flex;
  gap: 10px;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.name-input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.save-btn {
  padding: 8px 16px;
  background: #17a2b8;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.save-btn:hover {
  background: #138496;
}
</style>