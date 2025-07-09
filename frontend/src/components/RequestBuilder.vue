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
        <button @click="updateRequest" :disabled="!isDirty" class="btn btn--save" :class="{ 'is-disabled': !isDirty }">
          Save
        </button>
        <button @click="cloneRequest" class="btn btn--clone">
          Clone
        </button>
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
          <div class="header-row" v-for="(header, index) in allHeaders" :key="index">
            <input 
              :value="header.key"
              @input="updateHeader(index, 'key', $event.target.value)"
              type="text" 
              placeholder="Header name"
              class="header-input"
              :readonly="header.readonly"
            />
            <input 
              :value="header.value"
              @input="updateHeader(index, 'value', $event.target.value)"
              type="text" 
              placeholder="Header value"
              class="header-input"
              :readonly="header.readonly"
            />
            <button @click="removeHeader(index)" class="remove-btn" :disabled="header.readonly">×</button>
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

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { HTTPService, CollectionService, EventBusService } from '../../bindings/captain-api'
import { Events } from '@wailsio/runtime'

const props = defineProps({
  collectionId: {
    type: String,
    default: 'default'
  },
  virtualHeaders: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['response-received', 'request-updated', 'load-request'])

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
const currentRequestId = ref(null)
const isValidJSON = ref(true)
const jsonValidationMessage = ref('')
const jsonValidationClass = ref('')

// Headers management
const headersList = ref([])

const virtualHeadersCount = computed(() => Object.keys(props.virtualHeaders || {}).length)

const allHeaders = computed(() => {
  const combined = []
  if (props.virtualHeaders) {
    for (const [key, value] of Object.entries(props.virtualHeaders)) {
      combined.push({ key, value, readonly: true })
    }
  }
  headersList.value.forEach(h => {
      combined.push({ ...h, readonly: false })
  })
  return combined
})

const addHeader = () => {
  headersList.value.push({ key: '', value: '' })
}

const removeHeader = (index) => {
  const headerToRemove = allHeaders.value[index]
  if (headerToRemove.readonly) return

  const headersListIndex = index - virtualHeadersCount.value
  
  if (headersList.value.length > 0) {
    headersList.value.splice(headersListIndex, 1)
  }
}

const updateHeader = (index, field, value) => {
    const headerToUpdate = allHeaders.value[index]
    if (headerToUpdate.readonly) return

    const headersListIndex = index - virtualHeadersCount.value
    if (headersList.value[headersListIndex]) {
        headersList.value[headersListIndex][field] = value
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

// Watch for changes in the headers list and update the headers object
watch(headersList, (newHeaders) => {
  emitRequestUpdate()
}, { deep: true })

watch(requestName, () => {
  emitRequestUpdate()
})

watch(request, () => {
  emitRequestUpdate()
}, { deep: true })

// Watch body type changes
watch(bodyType, (newType) => {
  if (newType === 'json') {
    // Add Content-Type header for JSON
    const hasContentType = headersList.value.some(h => h.key === 'Content-Type')
    if (!hasContentType) {
      headersList.value.push({ key: 'Content-Type', value: 'application/json' })
    } else {
      // Update existing Content-Type header
      const index = headersList.value.findIndex(h => h.key === 'Content-Type')
      if (index !== -1) {
        headersList.value[index].value = 'application/json'
      }
    }
  } else if (newType === 'text') {
    // Add Content-Type header for text
    const hasContentType = headersList.value.some(h => h.key === 'Content-Type')
    if (!hasContentType) {
      headersList.value.push({ key: 'Content-Type', value: 'text/plain' })
    } else {
      // Update existing Content-Type header
      const index = headersList.value.findIndex(h => h.key === 'Content-Type')
      if (index !== -1) {
        headersList.value[index].value = 'text/plain'
      }
    }
  } else {
    // Remove Content-Type header for none
    const index = headersList.value.findIndex(h => 
      h.key === 'Content-Type' && h.value === 'application/json'
    )
    if (index !== -1) {
      headersList.value.splice(index, 1)
    }
  }
  
  // Emit request update when body type changes
  emitRequestUpdate()
})
// Emit request updates to parent for tab management
const emitRequestUpdate = () => {
  emit('request-updated', {
    request: {
      method: request.value.method,
      url: request.value.url,
      headers: headersObject.value,
      body: request.value.body
    },
    name: requestName.value,
    bodyType: bodyType.value
  })
}

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
    // Prepare body based on body type
    const bodyContent = bodyType.value === 'none' ? '' : request.value.body
    
    const combinedHeaders = { ...props.virtualHeaders, ...headersObject.value };
    console.log('Combined headers:', combinedHeaders)

    // Create a request object with synchronized headers
    const requestToSend = {
      method: request.value.method,
      url: request.value.url,
      headers: combinedHeaders,
      body: bodyContent
    }

    const response = await HTTPService.SendRequest(requestToSend)

    // Emit the response to the parent component
    emit('response-received', response)

    // Publish event for logs to refresh
    // Events.Emit('request-sent')
    EventBusService.EmitEvent('request-sent', response)
  } catch (error) {
    console.error('Request failed:', error)
    
    // Create an error response object
    const errorResponse = {
      statusCode: 0,
      statusText: 'Error',
      headers: {},
      body: JSON.stringify({ error: error.message || 'Unknown error occurred' }),
      error: error.message || 'Unknown error occurred'
    }
    
    // Emit the error response
    emit('response-received', errorResponse)
  } finally {
    loading.value = false
  }
}

async function cloneRequest() {
  const newName = `${requestName.value || 'Request'} (Copy)`;

  const requestToSave = {
    name: newName,
    method: request.value.method,
    url: request.value.url,
    headers: headersObject.value,
    body: request.value.body,
    description: ''
  };

  try {
    loading.value = true;
    const clonedRequest = await CollectionService.SaveRequest(props.collectionId, requestToSave);
    
    emit('load-request', clonedRequest);
    
    EventBusService.EmitEvent('request-saved');

  } catch (error) {
    console.error('Failed to clone request:', error);
  } finally {
    loading.value = false;
  }
}

// Save request to collection
async function saveRequest() {
  // Generate a default name if none provided
  if (!requestName.value) {
    const url = new URL(request.value.url, 'http://dummy.com') // Create URL to parse path
    const pathParts = url.pathname.split('/').filter(Boolean)
    const lastPart = pathParts[pathParts.length - 1] || 'request'
    requestName.value = `${request.value.method} ${lastPart}`.toUpperCase()
  }

  const requestToSave = {
    name: requestName.value,
    method: request.value.method,
    url: request.value.url,
    headers: headersObject.value,
    body: request.value.body,
    description: ''
  }

  try {
    loading.value = true
    console.log('Saving request:', requestToSave)
    
    // Call the backend service to save the request
    const savedRequest = await CollectionService.SaveRequest(props.collectionId, requestToSave)

    console.log('Request saved successfully:', savedRequest)
    
    // Update the current request ID if we get one back
    if (savedRequest && savedRequest.id) {
      currentRequestId.value = savedRequest.id
    }
    
    // Update the original request snapshot to reflect the saved state
    originalRequest.value = { ...request.value }
    originalRequestSnapshot.value = createSnapshot(request.value, requestName.value, headersList.value)
    
    // Notify the parent component that the request was saved
    // emit('request-saved', savedRequest)
    EventBusService.EmitEvent('request-saved', savedRequest)
    
    return savedRequest
  } catch (error) {
    console.error('Failed to save request:', error)
    // You might want to show an error message to the user here
    throw error
  } finally {
    loading.value = false
  }
}
  
async function updateRequest() {
  if (!currentRequestId.value) {
    console.log('No current request ID, redirecting to save')
    return saveRequest()
  }
  
  if (!requestName.value) {
    alert('Request name is required.')
    return
  }

  // Get the current timestamp for updatedAt
  const now = new Date().toISOString()
  
  // Create the complete request object with all required fields
  const requestToUpdate = {
    id: currentRequestId.value,
    name: requestName.value,
    method: request.value.method,
    url: request.value.url,
    headers: { ...headersObject.value }, // Create a new object to ensure reactivity
    body: request.value.body,
    description: '',
    updatedAt: now,
    // Preserve the original createdAt if it exists
    ...(originalRequest.value?.createdAt && { createdAt: originalRequest.value.createdAt })
  }

  console.log('Updating request with data:', requestToUpdate)
  
  try {
    // First try to update the request
    await CollectionService.UpdateRequest(props.collectionId, currentRequestId.value, requestToUpdate)
    console.log('Update successful, emitting request-saved event')
    // Notify parent components
    EventBusService.EmitEvent('request-saved')
    // Update the original request and snapshot
    originalRequest.value = { ...requestToUpdate }
    originalRequestSnapshot.value = createSnapshot(request.value, requestName.value, headersList.value)
  } catch (error) {
    console.error('Failed to update request:', error)
  }
}

// Reference to the original request data
const originalRequest = ref(null)
const originalRequestSnapshot = ref('')

const createSnapshot = (req, name, headersList) => {
  console.log('Creating snapshot for request:', req,name,headersList)
  const headers = {}
  if (headersList) {
    headersList.forEach(h => {
      if (h.key) headers[h.key] = h.value
    })
  }
  const sortedHeaders = Object.keys(headers).sort().reduce((obj, key) => {
    obj[key] = headers[key]
    return obj
  }, {})

  const json = JSON.stringify({
    name,
    method: req.method,
    url: req.url,
    headers: sortedHeaders,
    body: req.body
  })
  console.log('Snapshot created:', json)
  return json
}

const isDirty = computed(() => {
  if (!originalRequest.value) return false
  const currentSnapshot = createSnapshot(request.value, requestName.value, headersList.value)
  return currentSnapshot !== originalRequestSnapshot.value
})

const loadRequest = (requestData) => {
  console.log('Loading request data:', requestData)
  
  // Store the original request data and create a snapshot
  originalRequest.value = { ...requestData }
  
  // Set basic request data first (without body)
  request.value = {
    method: requestData.method || 'GET',
    url: requestData.url || '',
    headers: { ...(requestData.headers || {}) },
    body: requestData.body || ''
  }
  
  // Set the current request ID and name
  currentRequestId.value = requestData.id || ''
  if (requestData.name) {
    requestName.value = requestData.name
  }

  // Convert headers object to list
  const headers = requestData.headers || {}
  headersList.value = Object.entries(headers).map(([key, value]) => ({ 
    key, 
    value: String(value) 
  }))

  // Set body type based on content
  const bodyContent = requestData.body || ''
  if (bodyContent) {
    try {
      // Try to parse as JSON to determine if it's valid JSON
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
  
  // Create snapshot after all values are set
  originalRequestSnapshot.value = createSnapshot(request.value, requestName.value, headersList.value)
  
  console.log('Finished loading request. Current state:', {
    request: request.value,
    requestName: requestName.value,
    currentRequestId: currentRequestId.value,
    headers: headersList.value,
    bodyType: bodyType.value
  })
  
  // Emit request-updated event to sync with parent component's tab state
  emitRequestUpdate()
}

function newRequest() {
  // Reset request data
  request.value = {
    method: 'GET',
    url: '',
    headers: {},
    body: ''
  }
  
  // Reset UI state
  headersList.value = []
  bodyType.value = 'none'
  requestName.value = ''
  currentRequestId.value = null
  activeTab.value = 'Headers'
  
  // Reset validation state
  isValidJSON.value = true
  jsonValidationMessage.value = ''
  jsonValidationClass.value = ''
  
  // Reset dirty state tracking
  originalRequest.value = { ...request.value }
  originalRequestSnapshot.value = createSnapshot(request.value, requestName.value, headersList.value)
  
  // Notify parent component of the reset
  emitRequestUpdate()
}

defineExpose({ loadRequest, newRequest })
</script>


<style scoped lang="scss">
@use "sass:map";
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

.header-input[readonly] {
  background-color: #e9ecef;
  cursor: not-allowed;
}

.remove-btn:disabled {
  background-color: #6c757d;
  cursor: not-allowed;
  opacity: 0.5;
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

// Button styles
$btn-colors: (
  save: (
    default: #17a2b8,
    hover: #138496,
    disabled: #ccc
  ),
  clone: (
    default: #6c757d,
    hover: #5a6268,
    disabled: #d1d5db
  )
);

// Base button styles
.btn {
  padding: 8px 16px;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s ease-in-out;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 80px;
  
  // Disabled state for all buttons
  &:disabled,
  &.is-disabled,
  &.disabled {
    opacity: 0.7;
  }
  
  // Generate button variants using @each loop
  $variants: (save, clone);
  
  @each $variant in $variants {
    &--#{$variant} {
      background: map.get(map.get($btn-colors, $variant), default);
      
      &:hover:not(:disabled) {
        background: map.get(map.get($btn-colors, $variant), hover);
      }
      
      &:disabled,
      &.is-disabled,
      &.disabled {
        background: map.get(map.get($btn-colors, $variant), disabled);
        cursor: not-allowed;
      }
    }
  }
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