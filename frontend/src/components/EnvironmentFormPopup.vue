<template>
  <div v-if="isVisible" class="form-popup-overlay" @click="handleOverlayClick">
    <div class="form-popup" @click.stop>
      <!-- Header -->
      <div class="popup-header">
        <h3>{{ isEditing ? 'Edit Environment' : 'Add Environment' }}</h3>
        <button @click="$emit('close')" class="close-btn" title="Close">
          ×
        </button>
      </div>

      <!-- Body -->
      <div class="popup-body">
        <div class="form-group">
          <label>Environment Name *</label>
          <input 
            v-model="formData.name"
            type="text" 
            placeholder="e.g., Production, Staging, Development"
            class="form-input"
            @keyup.enter="handleSubmit"
          />
        </div>
        
        <div class="form-group">
          <label>Base URL *</label>
          <input 
            v-model="formData.baseUrl"
            type="text" 
            placeholder="https://api.example.com"
            class="form-input"
            @keyup.enter="handleSubmit"
          />
        </div>
        
        <div class="form-group">
          <label>Description</label>
          <textarea 
            v-model="formData.description"
            placeholder="Optional description for this environment"
            class="form-textarea"
            rows="3"
          ></textarea>
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
  environment: {
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
  baseUrl: '',
  description: ''
})

// Computed
const isEditing = computed(() => props.environment && props.environment.id)

const isFormValid = computed(() => {
  return formData.value.name.trim() && formData.value.baseUrl.trim()
})

// Methods
const handleSubmit = () => {
  if (!isFormValid.value) return
  
  const submitData = {
    ...formData.value,
    name: formData.value.name.trim(),
    baseUrl: formData.value.baseUrl.trim(),
    description: formData.value.description.trim()
  }
  
  emit('submit', submitData)
}

const handleOverlayClick = () => {
  emit('close')
}

const resetForm = () => {
  formData.value = {
    id: '',
    name: '',
    baseUrl: '',
    description: ''
  }
}

// Watch for environment changes
watch(() => props.environment, (newEnv) => {
  if (newEnv) {
    formData.value = { ...newEnv }
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
  max-width: 500px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
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

.form-input,
.form-textarea {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #ced4da;
  border-radius: 6px;
  font-size: 0.9rem;
  transition: all 0.2s ease;
  font-family: inherit;
}

.form-input:focus,
.form-textarea:focus {
  outline: none;
  border-color: #0d6efd;
  box-shadow: 0 0 0 3px rgba(13, 110, 253, 0.1);
}

.form-input::placeholder,
.form-textarea::placeholder {
  color: #6c757d;
}

.form-textarea {
  resize: vertical;
  min-height: 80px;
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