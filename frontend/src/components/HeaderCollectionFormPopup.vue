<template>
  <div v-if="isVisible" class="form-popup-overlay" @click="handleOverlayClick">
    <div class="form-popup header-collection-popup" @click.stop>
      <!-- Header -->
      <div class="popup-header">
        <h3>{{ isEditing ? 'Edit Header Collection' : 'Add Header Collection' }}</h3>
        <button @click="$emit('close')" class="close-btn" title="Close">
          ×
        </button>
      </div>

      <!-- Body -->
      <div class="popup-body">
        <div class="form-group">
          <label>Collection Name *</label>
          <input 
            v-model="formData.name"
            type="text" 
            placeholder="Enter collection name"
            class="form-input"
          />
        </div>
        
        <div class="form-group">
          <label>Description</label>
          <input 
            v-model="formData.description"
            type="text" 
            placeholder="Optional description"
            class="form-input"
          />
        </div>
        
        <div class="form-group">
          <label>Headers</label>
          <div class="headers-list">
            <div 
              v-for="(header, index) in formData.headers" 
              :key="index"
              class="header-item"
            >
              <div class="header-inputs">
                <input 
                  v-model="header.key"
                  type="text" 
                  placeholder="Header name"
                  class="header-input key-input"
                />
                <input 
                  v-model="header.value"
                  type="text" 
                  placeholder="Header value"
                  class="header-input value-input"
                />
                <button 
                  @click="removeHeader(index)"
                  class="remove-header-btn"
                  title="Remove header"
                  :disabled="formData.headers.length === 1"
                >
                  ×
                </button>
              </div>
            </div>
          </div>
          <button @click="addHeader" class="add-header-btn">
            <span class="btn-icon">+</span>
            Add Header
          </button>
        </div>
      </div>

      <!-- Footer -->
      <div class="popup-footer">
        <button @click="$emit('close')" class="btn btn-secondary">
          Cancel
        </button>
        <button 
          @click="handleSubmit" 
          class="btn btn-primary"
          :disabled="!isFormValid"
        >
          {{ isEditing ? 'Update' : 'Create' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'

// Props
const props = defineProps({
  isVisible: {
    type: Boolean,
    default: false
  },
  headerCollection: {
    type: Object,
    default: null
  }
})

// Emits
const emit = defineEmits(['close', 'submit'])

// Form data
const formData = ref({
  id: '',
  name: '',
  description: '',
  headers: [{ key: '', value: '' }]
})

// Computed
const isEditing = computed(() => props.headerCollection && props.headerCollection.id)

const isFormValid = computed(() => {
  if (!formData.value.name.trim()) return false
  
  // Check for duplicate header keys
  const headers = formData.value.headers.filter(h => h.key && h.value)
  const keys = headers.map(h => h.key.toLowerCase())
  const uniqueKeys = new Set(keys)
  
  return keys.length === uniqueKeys.size
})

// Methods
const handleSubmit = () => {
  if (!isFormValid.value) return
  
  const validHeaders = formData.value.headers.filter(h => h.key && h.value)
  
  const submitData = {
    ...formData.value,
    name: formData.value.name.trim(),
    description: formData.value.description.trim(),
    headers: validHeaders
  }
  
  emit('submit', submitData)
}

const handleOverlayClick = () => {
  emit('close')
}

const addHeader = () => {
  formData.value.headers.push({ key: '', value: '' })
}

const removeHeader = (index) => {
  if (formData.value.headers.length > 1) {
    formData.value.headers.splice(index, 1)
  }
}

const resetForm = () => {
  formData.value = {
    id: '',
    name: '',
    description: '',
    headers: [{ key: '', value: '' }]
  }
}

// Watch for header collection changes
watch(() => props.headerCollection, (newCollection) => {
  if (newCollection) {
    // Convert headers object to array for editing
    let headers = []
    if (newCollection.headers) {
      if (Array.isArray(newCollection.headers)) {
        headers = newCollection.headers.map(h => ({ ...h }))
      } else {
        // Convert object to array
        headers = Object.entries(newCollection.headers).map(([key, value]) => ({
          key,
          value: typeof value === 'object' ? JSON.stringify(value) : String(value)
        }))
      }
    }
    
    if (headers.length === 0) {
      headers = [{ key: '', value: '' }]
    }
    
    formData.value = {
      id: newCollection.id || '',
      name: newCollection.name || '',
      description: newCollection.description || '',
      headers
    }
  } else {
    resetForm()
  }
}, { immediate: true })

// Handle escape key
watch(() => props.isVisible, (newVal) => {
  if (newVal) {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        emit('close')
      }
    }
    document.addEventListener('keydown', handleEscape)
    
    return () => {
      document.removeEventListener('keydown', handleEscape)
    }
  }
})
</script>

<style scoped>
.form-popup-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.6);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 3000;
  backdrop-filter: blur(3px);
  animation: fadeIn 0.2s ease-out;
}

.form-popup {
  background-color: #fff;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  animation: slideIn 0.3s ease-out;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
}

.header-collection-popup {
  max-width: 700px;
}

.popup-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 25px;
  border-bottom: 1px solid #e9ecef;
  background-color: #f8f9fa;
  border-top-left-radius: 12px;
  border-top-right-radius: 12px;
}

.popup-header h3 {
  margin: 0;
  font-size: 1.25rem;
  color: #212529;
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #6c757d;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: all 0.2s ease;
}

.close-btn:hover {
  background-color: #e9ecef;
  color: #343a40;
}

.popup-body {
  flex: 1;
  overflow-y: auto;
  padding: 25px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group:last-child {
  margin-bottom: 0;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #495057;
  font-size: 0.9rem;
}

.form-input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #ced4da;
  border-radius: 6px;
  font-size: 0.9rem;
  transition: all 0.2s ease;
  font-family: inherit;
}

.form-input:focus {
  outline: none;
  border-color: #0d6efd;
  box-shadow: 0 0 0 3px rgba(13, 110, 253, 0.1);
}

.form-input::placeholder {
  color: #6c757d;
}

.headers-list {
  margin-bottom: 12px;
}

.header-item {
  margin-bottom: 12px;
}

.header-item:last-child {
  margin-bottom: 0;
}

.header-inputs {
  display: flex;
  gap: 8px;
  align-items: center;
}

.header-input {
  flex: 1;
  padding: 10px 12px;
  border: 1px solid #ced4da;
  border-radius: 6px;
  font-size: 0.85rem;
  transition: all 0.2s ease;
  font-family: inherit;
}

.header-input:focus {
  outline: none;
  border-color: #0d6efd;
  box-shadow: 0 0 0 3px rgba(13, 110, 253, 0.1);
}

.header-input::placeholder {
  color: #6c757d;
}

.remove-header-btn {
  background: none;
  border: 1px solid #dc3545;
  color: #dc3545;
  border-radius: 4px;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 16px;
  flex-shrink: 0;
}

.remove-header-btn:hover:not(:disabled) {
  background-color: #dc3545;
  color: white;
}

.remove-header-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.add-header-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background-color: #f8f9fa;
  border: 1px solid #dee2e6;
  border-radius: 6px;
  color: #495057;
  font-size: 0.85rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.add-header-btn:hover {
  background-color: #e9ecef;
  border-color: #adb5bd;
}

.btn-icon {
  font-size: 14px;
}

.popup-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px 25px;
  border-top: 1px solid #e9ecef;
  background-color: #f8f9fa;
  border-bottom-left-radius: 12px;
  border-bottom-right-radius: 12px;
}

.btn {
  padding: 10px 16px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 6px;
}

.btn-secondary {
  background-color: #e9ecef;
  color: #495057;
}

.btn-secondary:hover {
  background-color: #dee2e6;
  transform: translateY(-1px);
}

.btn-primary {
  background-color: #0d6efd;
  color: white;
}

.btn-primary:hover {
  background-color: #0b5ed7;
  transform: translateY(-1px);
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideIn {
  from { transform: translateY(20px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}
</style>