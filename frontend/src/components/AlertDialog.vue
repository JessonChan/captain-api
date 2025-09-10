<template>
  <div v-if="isVisible" class="alert-overlay" @click="onClose">
    <div class="alert-dialog" @click.stop>
      <div class="alert-header">
        <h3>{{ currentTitle }}</h3>
      </div>
      <div class="alert-body">
        <p>{{ currentMessage }}</p>
        <div v-if="currentDetails" class="alert-details">
          {{ currentDetails }}
        </div>
      </div>
      <div class="alert-actions">
        <button class="btn btn-primary" @click="onClose">
          {{ currentCloseText }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  title: {
    type: String,
    default: 'Alert'
  },
  message: {
    type: String,
    required: true
  },
  details: {
    type: String,
    default: ''
  },
  closeText: {
    type: String,
    default: 'OK'
  },
  severity: {
    type: String,
    default: 'info',
    validator: (value) => ['info', 'success', 'warning', 'error'].includes(value)
  }
})

const emit = defineEmits(['close'])

const isVisible = ref(false)
const currentTitle = ref('')
const currentMessage = ref('')
const currentDetails = ref('')
const currentCloseText = ref('')
const currentSeverity = ref('info')
let closeAction = null

const severityIcon = computed(() => {
  switch (currentSeverity.value) {
    case 'success':
      return '✅'
    case 'warning':
      return '⚠️'
    case 'error':
      return '❌'
    default:
      return 'ℹ️'
  }
})

const severityClass = computed(() => {
  return `severity-${currentSeverity.value}`
})

const show = (title, message, details, closeText, severity, onCloseCallback) => {
  currentTitle.value = title || props.title
  currentMessage.value = message || props.message
  currentDetails.value = details || props.details
  currentCloseText.value = closeText || props.closeText
  currentSeverity.value = severity || props.severity
  closeAction = onCloseCallback
  isVisible.value = true
}

const hide = () => {
  isVisible.value = false
}

const onClose = () => {
  if (closeAction) {
    closeAction()
  }
  emit('close')
  hide()
}

// Expose methods for parent component
defineExpose({
  show,
  hide
})
</script>

<style scoped>
.alert-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.alert-dialog {
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  min-width: 320px;
  max-width: 500px;
  margin: 20px;
}

.alert-header {
  padding: 20px 20px 0 20px;
  border-bottom: 1px solid #e9ecef;
  display: flex;
  align-items: center;
  gap: 10px;
}

.alert-header h3 {
  margin: 0 0 15px 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.alert-icon {
  font-size: 24px;
}

.alert-body {
  padding: 20px;
}

.alert-body p {
  margin: 0 0 10px 0;
  color: #666;
  line-height: 1.5;
}

.alert-details {
  background: #f8f9fa;
  padding: 12px;
  border-radius: 4px;
  font-size: 14px;
  color: #495057;
  border-left: 4px solid #6c757d;
}

.severity-error .alert-details {
  border-left-color: #dc3545;
  background: #f8d7da;
}

.severity-warning .alert-details {
  border-left-color: #ffc107;
  background: #fff3cd;
}

.severity-success .alert-details {
  border-left-color: #28a745;
  background: #d4edda;
}

.severity-info .alert-details {
  border-left-color: #17a2b8;
  background: #d1ecf1;
}

.alert-actions {
  padding: 0 20px 20px 20px;
  display: flex;
  justify-content: flex-end;
}

.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: background-color 0.2s;
}

.btn-primary {
  background: #007bff;
  color: white;
}

.btn-primary:hover {
  background: #0056b3;
}
</style>