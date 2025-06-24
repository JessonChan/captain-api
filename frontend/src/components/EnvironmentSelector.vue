<template>
  <div class="environment-selector">
    <div class="env-header">
      <label class="env-label">Environment:</label>
      <select v-model="selectedEnvironment" @change="switchEnvironment" class="env-select">
        <option v-for="env in environments" :key="env.id" :value="env.id">
          {{ env.name }} ({{ env.baseUrl }})
        </option>
      </select>
      <button @click="showManageModal = true" class="manage-btn" title="Manage Environments">
        ⚙️
      </button>
    </div>

    <!-- Environment Management Modal -->
    <div v-if="showManageModal" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Manage Environments</h3>
          <button @click="closeModal" class="close-btn">×</button>
        </div>
        
        <div class="modal-body">
          <!-- Environment List -->
          <div class="env-list">
            <div v-for="env in environments" :key="env.id" class="env-item">
              <div class="env-info">
                <div class="env-name">{{ env.name }}</div>
                <div class="env-url">{{ env.baseUrl }}</div>
                <div class="env-description">{{ env.description }}</div>
              </div>
              <div class="env-actions">
                <button @click="editEnvironment(env)" class="edit-btn">Edit</button>
                <button @click="deleteEnvironment(env.id)" class="delete-btn" :disabled="environments.length <= 1">
                  Delete
                </button>
              </div>
            </div>
          </div>

          <!-- Add/Edit Environment Form -->
          <div class="env-form">
            <h4>{{ editingEnv ? 'Edit Environment' : 'Add New Environment' }}</h4>
            <form @submit.prevent="saveEnvironment">
              <div class="form-group">
                <label>Name:</label>
                <input v-model="formData.name" type="text" required class="form-input" />
              </div>
              <div class="form-group">
                <label>Base URL:</label>
                <input v-model="formData.baseUrl" type="url" required class="form-input" placeholder="https://api.example.com" />
              </div>
              <div class="form-group">
                <label>Description:</label>
                <input v-model="formData.description" type="text" class="form-input" placeholder="Optional description" />
              </div>
              <div class="form-actions">
                <button type="submit" class="save-btn">{{ editingEnv ? 'Update' : 'Add' }}</button>
                <button type="button" @click="cancelEdit" class="cancel-btn">Cancel</button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { GetEnvironments, SetActiveEnvironment, AddEnvironment, UpdateEnvironment, DeleteEnvironment } from '../../bindings/captain-api/environmentservice'

const environments = ref([])
const selectedEnvironment = ref('')
const showManageModal = ref(false)
const editingEnv = ref(null)
const formData = ref({
  name: '',
  baseUrl: '',
  description: ''
})

// Load environments on component mount
onMounted(async () => {
  await loadEnvironments()
})

// Load environments from backend
const loadEnvironments = async () => {
  try {
    const envs = await GetEnvironments()
    environments.value = envs
    
    // Set selected environment to the active one
    const activeEnv = envs.find(env => env.isActive)
    if (activeEnv) {
      selectedEnvironment.value = activeEnv.id
    }
  } catch (error) {
    console.error('Failed to load environments:', error)
    alert('Failed to load environments')
  }
}

// Switch to selected environment
const switchEnvironment = async () => {
  try {
    await SetActiveEnvironment(selectedEnvironment.value)
    await loadEnvironments() // Reload to update active status
  } catch (error) {
    console.error('Failed to switch environment:', error)
    alert('Failed to switch environment')
  }
}

// Edit environment
const editEnvironment = (env) => {
  editingEnv.value = env
  formData.value = {
    name: env.name,
    baseUrl: env.baseUrl,
    description: env.description
  }
}

// Save environment (add or update)
const saveEnvironment = async () => {
  try {
    if (editingEnv.value) {
      // Update existing environment
      await UpdateEnvironment({
        id: editingEnv.value.id,
        name: formData.value.name,
        baseUrl: formData.value.baseUrl,
        description: formData.value.description,
        isActive: editingEnv.value.isActive
      })
    } else {
      // Add new environment
      await AddEnvironment({
        name: formData.value.name,
        baseUrl: formData.value.baseUrl,
        description: formData.value.description,
        isActive: false
      })
    }
    
    await loadEnvironments()
    cancelEdit()
  } catch (error) {
    console.error('Failed to save environment:', error)
    alert('Failed to save environment')
  }
}

// Delete environment
const deleteEnvironment = async (envId) => {
  if (confirm('Are you sure you want to delete this environment?')) {
    try {
      await DeleteEnvironment(envId)
      await loadEnvironments()
    } catch (error) {
      console.error('Failed to delete environment:', error)
      alert('Failed to delete environment')
    }
  }
}

// Cancel edit
const cancelEdit = () => {
  editingEnv.value = null
  formData.value = {
    name: '',
    baseUrl: '',
    description: ''
  }
}

// Close modal
const closeModal = () => {
  showManageModal.value = false
  cancelEdit()
}
</script>

<style scoped>
.environment-selector {
  margin-bottom: 20px;
}

.env-header {
  display: flex;
  align-items: center;
  gap: 10px;
}

.env-label {
  font-weight: 500;
  color: #374151;
}

.env-select {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  background: white;
  font-size: 14px;
}

.env-select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.manage-btn {
  padding: 8px 12px;
  background: #f3f4f6;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  cursor: pointer;
  font-size: 16px;
}

.manage-btn:hover {
  background: #e5e7eb;
}

/* Modal Styles */
.modal-overlay {
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

.modal-content {
  background: white;
  border-radius: 8px;
  width: 90%;
  max-width: 600px;
  max-height: 80vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e5e7eb;
}

.modal-header h3 {
  margin: 0;
  color: #111827;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #6b7280;
}

.close-btn:hover {
  color: #374151;
}

.modal-body {
  padding: 20px;
}

/* Environment List */
.env-list {
  margin-bottom: 30px;
}

.env-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  margin-bottom: 10px;
}

.env-info {
  flex: 1;
}

.env-name {
  font-weight: 500;
  color: #111827;
  margin-bottom: 4px;
}

.env-url {
  color: #3b82f6;
  font-size: 14px;
  margin-bottom: 4px;
}

.env-description {
  color: #6b7280;
  font-size: 12px;
}

.env-actions {
  display: flex;
  gap: 8px;
}

.edit-btn, .delete-btn {
  padding: 6px 12px;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.edit-btn {
  background: #f3f4f6;
  color: #374151;
}

.edit-btn:hover {
  background: #e5e7eb;
}

.delete-btn {
  background: #fef2f2;
  color: #dc2626;
  border-color: #fecaca;
}

.delete-btn:hover:not(:disabled) {
  background: #fee2e2;
}

.delete-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Form Styles */
.env-form {
  border-top: 1px solid #e5e7eb;
  padding-top: 20px;
}

.env-form h4 {
  margin: 0 0 15px 0;
  color: #111827;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: 500;
  color: #374151;
}

.form-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
}

.form-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-actions {
  display: flex;
  gap: 10px;
  margin-top: 20px;
}

.save-btn, .cancel-btn {
  padding: 8px 16px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
}

.save-btn {
  background: #3b82f6;
  color: white;
  border-color: #3b82f6;
}

.save-btn:hover {
  background: #2563eb;
}

.cancel-btn {
  background: #f3f4f6;
  color: #374151;
}

.cancel-btn:hover {
  background: #e5e7eb;
}
</style>