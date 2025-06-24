<template>
  <div v-if="isVisible" class="confirm-overlay" @click="onCancel">
    <div class="confirm-dialog" @click.stop>
      <div class="confirm-header">
        <h3>{{ title }}</h3>
      </div>
      <div class="confirm-body">
        <p>{{ message }}</p>
      </div>
      <div class="confirm-actions">
        <button class="btn btn-secondary" @click="onCancel">
          {{ cancelText }}
        </button>
        <button class="btn btn-danger" @click="onConfirm">
          {{ confirmText }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  title: {
    type: String,
    default: 'Confirm Action'
  },
  message: {
    type: String,
    required: true
  },
  confirmText: {
    type: String,
    default: 'Confirm'
  },
  cancelText: {
    type: String,
    default: 'Cancel'
  }
})

const emit = defineEmits(['confirm', 'cancel'])

const isVisible = ref(false)

const show = () => {
  isVisible.value = true
}

const hide = () => {
  isVisible.value = false
}

const onConfirm = () => {
  emit('confirm')
  hide()
}

const onCancel = () => {
  emit('cancel')
  hide()
}

// Expose methods for parent component
defineExpose({
  show,
  hide
})
</script>

<style scoped>
.confirm-overlay {
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

.confirm-dialog {
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  min-width: 320px;
  max-width: 500px;
  margin: 20px;
}

.confirm-header {
  padding: 20px 20px 0 20px;
  border-bottom: 1px solid #e9ecef;
}

.confirm-header h3 {
  margin: 0 0 15px 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.confirm-body {
  padding: 20px;
}

.confirm-body p {
  margin: 0;
  color: #666;
  line-height: 1.5;
}

.confirm-actions {
  padding: 0 20px 20px 20px;
  display: flex;
  gap: 12px;
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

.btn-secondary {
  background: #6c757d;
  color: white;
}

.btn-secondary:hover {
  background: #5a6268;
}

.btn-danger {
  background: #dc3545;
  color: white;
}

.btn-danger:hover {
  background: #c82333;
}
</style>