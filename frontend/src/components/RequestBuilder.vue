<template>
  <div class="request-builder">
    <div class="request-header">
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
      
      <div class="method-url-row">
        <select v-model="request.method" class="method-select">
          <option v-for="method in methods" :key="method" :value="method">
            {{ method }}
          </option>
        </select>
        <input 
          v-model="request.url" 
          type="text" 
          placeholder="Enter request URL (use relative paths like /api/users with environments)" 
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

        
        <div class="headers-list">
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
        <div v-if="bodyType !== 'none'" class="body-input-container">
          <div v-if="bodyType === 'json'" class="json-controls">
            <button @click="formatJSON" class="format-btn" type="button">Format JSON</button>
            <span v-if="jsonValidationMessage" :class="['validation-message', jsonValidationClass]">
              {{ jsonValidationMessage }}
            </span>
          </div>
          <textarea 
            spellcheck="false"
            v-model="request.body"
            :placeholder="bodyType === 'json' ? 'Enter JSON body' : 'Enter request body'"
            :class="['body-textarea', { 'json-invalid': bodyType === 'json' && !isValidJSON }]"
            rows="10"
            @input="onBodyInput"
          ></textarea>
        </div>
      </div>
    </div>


  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { HTTPService, CollectionService } from '../../bindings/captain-api'

const props = defineProps({
  collectionId: {
    type: String,
    default: 'default'
  }
})

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
const isValidJSON = ref(true)
const jsonValidationMessage = ref('')
const jsonValidationClass = ref('')

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
  if (newType === 'json') {
    if (!request.value.body) {
      request.value.body = '{\n  \n}'
    }
    // Auto-set Content-Type header for JSON
    const hasContentType = headersList.value.some(h => 
      h.key.toLowerCase() === 'content-type'
    )
    if (!hasContentType) {
      headersList.value.push({ key: 'Content-Type', value: 'application/json' })
    }
  } else if (newType === 'none') {
    request.value.body = ''
    // Remove Content-Type header if it was auto-added
    const index = headersList.value.findIndex(h => 
      h.key === 'Content-Type' && h.value === 'application/json'
    )
    if (index !== -1) {
      headersList.value.splice(index, 1)
    }
  }
})

// Validate JSON format
const validateJSON = (jsonString) => {
  try {
    JSON.parse(jsonString)
    return { valid: true, error: null }
  } catch (error) {
    return { valid: false, error: error.message }
  }
}

// Format JSON in textarea
const formatJSON = () => {
  if (request.value.body.trim()) {
    const validation = validateJSON(request.value.body)
    if (validation.valid) {
      try {
        const parsed = JSON.parse(request.value.body)
        request.value.body = JSON.stringify(parsed, null, 2)
        updateJSONValidation(true, 'JSON formatted successfully')
      } catch (error) {
        updateJSONValidation(false, `Format error: ${error.message}`)
      }
    } else {
      updateJSONValidation(false, `Cannot format: ${validation.error}`)
    }
  }
}

// Handle body input changes for real-time JSON validation
const onBodyInput = () => {
  if (bodyType.value === 'json' && request.value.body.trim()) {
    const validation = validateJSON(request.value.body)
    updateJSONValidation(validation.valid, validation.error)
  } else {
    updateJSONValidation(true, '')
  }
}

// Update JSON validation state
const updateJSONValidation = (valid, message) => {
  isValidJSON.value = valid
  jsonValidationMessage.value = message || ''
  jsonValidationClass.value = valid ? 'valid' : 'invalid'
  
  // Clear success message after 2 seconds
  if (valid && message) {
    setTimeout(() => {
      if (jsonValidationMessage.value === message) {
        jsonValidationMessage.value = ''
      }
    }, 2000)
  }
}

// Send HTTP request
const sendRequest = async () => {
  if (!request.value.url) {
    alert('Please enter a URL')
    return
  }

  // Validate JSON if body type is JSON
  if (bodyType.value === 'json' && request.value.body.trim()) {
    const validation = validateJSON(request.value.body)
    if (!validation.valid) {
      alert(`Invalid JSON format: ${validation.error}`)
      return
    }
  }

  loading.value = true
  
  try {
    const requestData = {
      method: request.value.method,
      url: request.value.url,
      headers: headersObject.value,
      body: bodyType.value === 'none' ? '' : request.value.body,
      collectionId: props.collectionId
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
      name: requestName.value || `${request.value.url}`,
      method: request.value.method,
      url: request.value.url,
      headers: headersObject.value,
      body: bodyType.value === 'none' ? '' : request.value.body,
      description: ''
    }

    await CollectionService.SaveRequest(props.collectionId, requestItem)
    alert('Request saved successfully!')
    requestName.value = ''
  } catch (error) {
    console.error('Failed to save request:', error)
    alert('Failed to save request')
  }
}

// Load request from props
const loadRequest = (requestData) => {
  // Set basic request data first (without body)
  request.value = {
    method: requestData.method || 'GET',
    url: requestData.url || '',
    headers: requestData.headers || {},
    body: ''
  }

  // Convert headers object to list
  headersList.value = Object.entries(request.value.headers).map(([key, value]) => ({ key, value }))
  if (headersList.value.length === 0) {
    headersList.value = [{ key: '', value: '' }]
  }

  // Set body type first to avoid watcher conflicts
  const bodyContent = requestData.body || ''
  if (bodyContent) {
    try {
      JSON.parse(bodyContent)
      bodyType.value = 'json'
    } catch {
      bodyType.value = 'text'
    }
  } else {
    bodyType.value = 'none'
  }

  // Set body content after body type is set
  request.value.body = bodyContent
}

defineExpose({ loadRequest })
</script>

<style scoped>
.request-builder {
  background: white;
  border-radius: 8px;
  padding-left: 20px;
  padding-right: 20px;
  padding-bottom: 10px;
  padding-top: 10px;
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
  width: 100px;
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



.headers-list {
  margin-top: 10px;
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
  /* Prevent smart quote conversion */
  font-variant-ligatures: none;
  text-rendering: optimizeSpeed;
  -webkit-font-feature-settings: "liga" off;
  font-feature-settings: "liga" off;
}

.save-section {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
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

/* JSON-specific styles */
.body-input-container {
  width: 100%;
}

.json-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  padding: 8px;
  background: #f8f9fa;
  border-radius: 4px;
  border: 1px solid #e9ecef;
}

.format-btn {
  padding: 6px 12px;
  background: #6f42c1;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  font-weight: 500;
}

.format-btn:hover {
  background: #5a32a3;
}

.validation-message {
  font-size: 12px;
  font-weight: 500;
  padding: 4px 8px;
  border-radius: 3px;
}

.validation-message.valid {
  color: #155724;
  background: #d4edda;
  border: 1px solid #c3e6cb;
}

.validation-message.invalid {
  color: #721c24;
  background: #f8d7da;
  border: 1px solid #f5c6cb;
}

.body-textarea.json-invalid {
  border-color: #dc3545;
  box-shadow: 0 0 0 0.2rem rgba(220, 53, 69, 0.25);
}

.body-textarea:focus {
  outline: none;
  border-color: #007bff;
  box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
}

.body-textarea.json-invalid:focus {
  border-color: #dc3545;
  box-shadow: 0 0 0 0.2rem rgba(220, 53, 69, 0.25);
}
</style>