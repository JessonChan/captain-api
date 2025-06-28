<template>
  <div class="header-settings">
    <button class="settings-btn" @click="$emit('open')" title="Manage Header Collections">
      <span class="btn-icon">⚙️</span>
      <span class="btn-text">Headers</span>
    </button>
    
    <!-- Header Settings Modal -->
    <div v-if="showManageModal" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Header Collections</h3>
          <button @click="closeModal" class="close-btn">×</button>
        </div>
        
        <div class="modal-body">
          <!-- Header Collections List -->
          <div class="header-collections-list">
            <div class="list-header">
              <h4>Your Collections</h4>
              <div class="list-actions">
                <div class="collection-count">{{ headerCollections.length }} collections</div>
                <div class="sort-control">
                  <select v-model="sortOption" class="sort-select">
                    <option value="name-asc">Name (A-Z)</option>
                    <option value="name-desc">Name (Z-A)</option>
                    <option value="date-desc">Newest first</option>
                    <option value="date-asc">Oldest first</option>
                  </select>
                  <span class="sort-icon">⌄</span>
                </div>
                <button @click="triggerImportFile" class="import-btn" title="Import collection">
                  <span class="btn-icon">📥</span>
                  <span class="btn-text">Import</span>
                </button>
                <input 
                  type="file" 
                  ref="fileInput" 
                  @change="importHeaderCollection" 
                  accept=".json" 
                  style="display: none;"
                />
              </div>
            </div>
            
            <div class="search-container">
              <input 
                type="text" 
                v-model="searchQuery" 
                placeholder="Search collections..." 
                class="search-input"
              />
              <span class="search-icon">🔍</span>
            </div>
            
            <div v-if="headerCollections.length === 0" class="empty-collections">
              <div class="empty-icon">📋</div>
              <div class="empty-message">No header collections yet</div>
              <div class="empty-description">Create your first header collection to reuse headers across requests</div>
            </div>
            
            <div v-else v-for="collection in filteredCollections" :key="collection.id" class="header-collection-item">
              <div class="header-collection-info">
                <div class="header-collection-name">{{ collection.name }}</div>
                <div class="header-collection-description">{{ collection.description }}</div>
                <div class="header-collection-stats">
                  <span class="header-count-badge">{{ Object.keys(collection.headers || {}).length }} headers</span>
                </div>
              </div>
              <div class="header-collection-actions">
                <button @click="editHeaderCollection(collection)" class="edit-btn" title="Edit collection">
                  <span class="btn-icon">✏️</span>
                  <span class="btn-text">Edit</span>
                </button>
                <button @click="exportHeaderCollection(collection)" class="export-btn" title="Export collection">
                  <span class="btn-icon">📤</span>
                  <span class="btn-text">Export</span>
                </button>
                <button @click="deleteHeaderCollection(collection.id)" class="delete-btn" title="Delete collection">
                  <span class="btn-icon">🗑️</span>
                  <span class="btn-text">Delete</span>
                </button>
              </div>
            </div>
          </div>

          <!-- Add/Edit Header Collection Form -->
          <div class="header-collection-form">
            <div class="form-header">
              <h4>{{ editingCollection ? 'Edit Header Collection' : 'Add New Header Collection' }}</h4>
              <div v-if="editingCollection" class="editing-badge">Editing</div>
            </div>
            <form @submit.prevent="saveHeaderCollection">
              <div class="form-group">
                <div class="form-label-row">
                  <label for="collection-name">Collection Name:</label>
                  <span class="char-counter" :class="{ 'near-limit': formData.name.length > 40, 'at-limit': formData.name.length >= 50 }">
                    {{ formData.name.length }}/50
                  </span>
                </div>
                <input 
                  id="collection-name"
                  v-model="formData.name" 
                  type="text" 
                  required 
                  maxlength="50"
                  class="form-input" 
                  placeholder="Enter collection name"
                  autofocus
                />
              </div>
              <div class="form-group">
                <div class="form-label-row">
                  <label for="collection-description">Description:</label>
                  <span class="char-counter" :class="{ 'near-limit': formData.description.length > 80, 'at-limit': formData.description.length >= 100 }">
                    {{ formData.description.length }}/100
                  </span>
                </div>
                <input 
                  id="collection-description"
                  v-model="formData.description" 
                  type="text" 
                  maxlength="100"
                  class="form-input" 
                  placeholder="Optional description" 
                />                
              </div>
              
              <!-- Headers Editor -->
              <div class="headers-editor">
                <div class="headers-editor-header">
                  <h5>Headers</h5>
                  <div class="header-count">{{ formData.headers.filter(h => h.key && h.value).length }} valid headers</div>
                </div>
                
                <div class="header-table">
                  <div class="header-table-head">
                    <div class="header-name-col">Header Name</div>
                    <div class="header-value-col">Value</div>
                    <div class="header-action-col"></div>
                  </div>
                  
                  <div class="header-rows">
                    <div v-if="formData.headers.length === 0" class="empty-headers">
                      <div class="empty-headers-message">No headers added yet</div>
                      <div class="empty-headers-description">Click "Add Header" below to add your first header</div>
                    </div>
                    
                    <div 
                      class="header-row" 
                      v-for="(header, index) in formData.headers" 
                      :key="index"
                      :class="{'is-valid': header.key && header.value, 'is-invalid': (!header.key && header.value) || (header.key && !header.value)}"
                    >
                      <div class="input-wrapper">
                        <input 
                          v-model="header.key" 
                          type="text" 
                          placeholder="Header name"
                          class="header-input"
                          :class="{'invalid-input': !header.key && header.value}"
                        />
                        <div v-if="!header.key && header.value" class="validation-message">Name required</div>
                      </div>
                      <div class="input-wrapper">
                        <input 
                          v-model="header.value" 
                          type="text" 
                          placeholder="Header value"
                          class="header-input"
                          :class="{'invalid-input': header.key && !header.value}"
                        />
                        <div v-if="header.key && !header.value" class="validation-message">Value required</div>
                      </div>
                      <button 
                        type="button" 
                        @click="removeHeader(index)" 
                        class="remove-header-btn"
                        title="Remove header"
                      >
                        <span class="remove-icon">×</span>
                      </button>
                    </div>
                  </div>
                </div>
                
                <button type="button" @click="addHeader" class="add-header-btn">
                  <span class="add-icon">+</span> Add Header
                </button>
              </div>
              
              <div class="form-actions">
                <button type="button" @click="cancelEdit" class="cancel-btn">
                  <span class="btn-icon">✕</span>
                  <span class="btn-text">Cancel</span>
                </button>
                <button type="submit" class="save-btn">
                  <span class="btn-icon">{{ editingCollection ? '✓' : '+' }}</span>
                  <span class="btn-text">{{ editingCollection ? 'Update Collection' : 'Add Collection' }}</span>
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { GetHeaderCollections, CreateHeaderCollection, UpdateHeaderCollection, DeleteHeaderCollection, AddHeaderTemplate, UpdateHeaderTemplate, DeleteHeaderTemplate } from '../../bindings/captain-api/headerservice'

// Define props and emits
const props = defineProps({
  show: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close', 'open'])

const headerCollections = ref([])
const showManageModal = computed(() => props.show)
const editingCollection = ref(null)
const searchQuery = ref('')
const sortOption = ref('name-asc') // Default sort option
const formData = ref({
  name: '',
  description: '',
  headers: [{ key: '', value: '' }]
})

// Filter and sort collections based on search query and sort option
const filteredCollections = computed(() => {
  // First filter by search query
  let filtered = headerCollections.value
  
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(collection => 
      collection.name.toLowerCase().includes(query) || 
      (collection.description && collection.description.toLowerCase().includes(query))
    )
  }
  
  // Then sort based on selected option
  return [...filtered].sort((a, b) => {
    switch (sortOption.value) {
      case 'name-asc':
        return a.name.localeCompare(b.name)
      case 'name-desc':
        return b.name.localeCompare(a.name)
      case 'date-desc':
        // Assuming newer collections have higher IDs or using timestamp if available
        return b.id.localeCompare(a.id)
      case 'date-asc':
        return a.id.localeCompare(b.id)
      default:
        return 0
    }
  })
})

// Handler function for the open-header-settings event
const handleOpenHeaderSettings = () => {
  showManageModal.value = true
}

// Load header collections on component mount and set up event listener
onMounted(async () => {
  await loadHeaderCollections()
  
  // Add event listener for opening header settings from other components
  document.addEventListener('open-header-settings', handleOpenHeaderSettings)
})

// Clean up event listener when component is unmounted
onUnmounted(() => {
  document.removeEventListener('open-header-settings', handleOpenHeaderSettings)
})

// Load header collections from backend
const loadHeaderCollections = async () => {
  try {
    const collections = await GetHeaderCollections()
    
    // Convert backend format to frontend format
    headerCollections.value = collections.map(collection => {
      // Extract headers from headerTemplates if available
      const headers = []
      if (collection.headerTemplates && collection.headerTemplates.length > 0) {
        // Use the first template's headers
        const template = collection.headerTemplates[0]
        if (template && template.headers) {
          Object.entries(template.headers).forEach(([key, value]) => {
            headers.push({ key, value })
          })
        }
      }
      
      return {
        id: collection.id,
        name: collection.name,
        description: collection.description,
        headers
      }
    })
  } catch (error) {
    console.error('Failed to load header collections:', error)
  }
}

// Add a new header field
const addHeader = () => {
  formData.value.headers.push({ key: '', value: '' })
}

// Remove a header field
const removeHeader = (index) => {
  if (formData.value.headers.length > 1) {
    formData.value.headers.splice(index, 1)
  }
}

// Edit a header collection
const editHeaderCollection = (collection) => {
  editingCollection.value = collection
  formData.value = {
    name: collection.name,
    description: collection.description || '',
    headers: Object.entries(collection.headers || {}).map(([key, value]) => ({ key, value }))
  }
  
  // Ensure there's at least one header row
  if (formData.value.headers.length === 0) {
    formData.value.headers = [{ key: '', value: '' }]
  }
}

// Cancel editing
const cancelEdit = () => {
  editingCollection.value = null
  formData.value = {
    name: '',
    description: '',
    headers: [{ key: '', value: '' }]
  }
}

// Validate form data
const validateForm = () => {
  // Check if name is provided
  if (!formData.value.name.trim()) {
    alert('Please provide a collection name')
    return false
  }
  
  // Check for duplicate header keys
  const validHeaders = formData.value.headers.filter(h => h.key && h.value)
  const headerKeys = validHeaders.map(h => h.key.toLowerCase())
  const uniqueKeys = new Set(headerKeys)
  
  if (headerKeys.length !== uniqueKeys.size) {
    alert('Duplicate header names found. Each header name must be unique.')
    return false
  }
  
  // Check if any header is incomplete (has key but no value or vice versa)
  const incompleteHeaders = formData.value.headers.filter(
    h => (h.key && !h.value) || (!h.key && h.value)
  )
  
  if (incompleteHeaders.length > 0) {
    alert('Some headers are incomplete. Please provide both name and value for each header.')
    return false
  }
  
  return true
}

// Save header collection
const saveHeaderCollection = async () => {
  // Validate form before saving
  if (!validateForm()) return
  
  try {
    // Convert headers array to object
    const headersObject = {}
    formData.value.headers.forEach(header => {
      if (header.key && header.value) {
        headersObject[header.key] = header.value
      }
    })
    
    // Create a header template from the headers
    const template = {
      id: editingCollection.value?.id ? `${editingCollection.value.id}_template` : '',
      name: `${formData.value.name.trim()} Template`,
      description: formData.value.description.trim(),
      headers: headersObject
    }
    
    // Create the collection with the template
    const collection = {
      id: editingCollection.value?.id || '',
      name: formData.value.name.trim(),
      description: formData.value.description.trim(),
      headerTemplates: [template]
    }
    
    if (editingCollection.value) {
      // Update existing collection
      await UpdateHeaderCollection(collection)
      
      // Update in local array
      const index = headerCollections.value.findIndex(c => c.id === collection.id)
      if (index !== -1) {
        headerCollections.value[index] = {
          ...collection,
          headers: formData.value.headers
        }
      }
    } else {
      // Create new collection
      await CreateHeaderCollection(collection)
      
      // Reload collections to get the server-generated ID
      await loadHeaderCollections()
    }
    
    // Reset form
    cancelEdit()
  } catch (error) {
    console.error('Failed to save header collection:', error)
    alert('Failed to save header collection')
  }
}

// Delete header collection with confirmation
const deleteHeaderCollection = async (collectionId) => {
  const collection = headerCollections.value.find(c => c.id === collectionId)
  if (!collection) return
  
  // Show confirmation dialog
  const confirmDelete = confirm(`Are you sure you want to delete the "${collection.name}" header collection? This action cannot be undone.`)
  
  if (confirmDelete) {
    try {
      await DeleteHeaderCollection(collectionId)
      
      // Remove from local array
      headerCollections.value = headerCollections.value.filter(c => c.id !== collectionId)
    } catch (error) {
      console.error('Failed to delete header collection:', error)
      alert('Failed to delete header collection')
    }
  }
}

// Export header collection as JSON file
const exportHeaderCollection = (collection) => {
  try {
    // Create a formatted JSON string with indentation
    const jsonData = JSON.stringify(collection, null, 2)
    
    // Create a blob with the JSON data
    const blob = new Blob([jsonData], { type: 'application/json' })
    
    // Create a URL for the blob
    const url = URL.createObjectURL(blob)
    
    // Create a temporary anchor element
    const a = document.createElement('a')
    a.href = url
    a.download = `${collection.name.replace(/\s+/g, '_')}_headers.json`
    
    // Trigger the download
    document.body.appendChild(a)
    a.click()
    
    // Clean up
    setTimeout(() => {
      document.body.removeChild(a)
      URL.revokeObjectURL(url)
    }, 100)
  } catch (error) {
    console.error('Failed to export header collection:', error)
    alert('Failed to export header collection')
  }
}

// Trigger file input click to open file dialog
const fileInput = ref(null)
const triggerImportFile = () => {
  fileInput.value.click()
}

// Import header collection from JSON file
const importHeaderCollection = (event) => {
  const file = event.target.files[0]
  if (!file) return
  
  const reader = new FileReader()
  reader.onload = (e) => {
    try {
      const importedData = JSON.parse(e.target.result)
      
      // Validate the imported data
      if (!importedData.name || !importedData.headers) {
        throw new Error('Invalid header collection format')
      }
      
      // Create a new collection with a unique ID
      const newCollection = {
        id: 'header_col_' + Date.now(),
        name: importedData.name,
        description: importedData.description || '',
        headers: importedData.headers
      }
      
      // Check for duplicate name
      const existingNames = headerCollections.value.map(c => c.name.toLowerCase())
      if (existingNames.includes(newCollection.name.toLowerCase())) {
        newCollection.name += ' (Imported)'
      }
      
      // Add to collections
      headerCollections.value.push(newCollection)
      
      // Show success message
      alert(`Successfully imported "${newCollection.name}" collection`)
    } catch (error) {
      console.error('Failed to import header collection:', error)
      alert('Failed to import header collection: Invalid format')
    }
    
    // Reset file input
    event.target.value = ''
  }
  
  reader.onerror = () => {
    alert('Error reading file')
    event.target.value = ''
  }
  
  reader.readAsText(file)
}

// Close modal
const closeModal = (e) => {
  emit('close')
}
</script>

<style>
.header-settings {
  display: flex;
  align-items: center;
}

.settings-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1rem;
  padding: 6px 10px;
  border-radius: 4px;
  transition: background-color 0.2s;
  display: flex;
  align-items: center;
  gap: 5px;
  color: #4a5568;
}

.settings-btn:hover {
  background-color: rgba(0, 0, 0, 0.1);
}

.settings-btn .btn-icon {
  font-size: 1.2rem;
}

.settings-btn .btn-text {
  font-size: 0.9rem;
  font-weight: 500;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.6);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  backdrop-filter: blur(3px);
  animation: fadeIn 0.2s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modal-content {
  background-color: #fff;
  border-radius: 12px;
  width: 90%;
  max-width: 900px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
  from { transform: translateY(20px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 18px 25px;
  border-bottom: 1px solid #e9ecef;
  background-color: #f8f9fa;
  border-top-left-radius: 12px;
  border-top-right-radius: 12px;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.3rem;
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

.modal-body {
  padding: 25px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 25px;
}

.header-collections-list {
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 0;
  max-height: 500px;
  overflow-y: auto;
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.list-header {
  padding: 15px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.list-header h4 {
  margin: 0;
  font-size: 1rem;
  color: #343a40;
}

.list-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.collection-count {
  font-size: 0.85rem;
  color: #6c757d;
  background-color: #e9ecef;
  padding: 3px 8px;
  border-radius: 12px;
}

.sort-control {
  position: relative;
}

.sort-select {
  appearance: none;
  background-color: #fff;
  border: 1px solid #ced4da;
  border-radius: 4px;
  padding: 4px 24px 4px 8px;
  font-size: 0.85rem;
  color: #495057;
  cursor: pointer;
  transition: all 0.2s ease;
}

.sort-select:focus {
  outline: none;
  border-color: #86b7fe;
  box-shadow: 0 0 0 2px rgba(13, 110, 253, 0.1);
}

.sort-icon {
  position: absolute;
  right: 8px;
  top: 50%;
  transform: translateY(-50%);
  pointer-events: none;
  color: #6c757d;
  font-size: 0.9rem;
}

.search-container {
  position: relative;
  padding: 10px 15px;
  border-bottom: 1px solid #e9ecef;
}

.search-input {
  width: 100%;
  padding: 8px 12px 8px 32px;
  border: 1px solid #ced4da;
  border-radius: 20px;
  font-size: 0.9rem;
  transition: all 0.2s ease;
  color: #495057;
}

.search-input:focus {
  outline: none;
  border-color: #86b7fe;
  box-shadow: 0 0 0 3px rgba(13, 110, 253, 0.1);
}

.search-input::placeholder {
  color: #adb5bd;
}

.search-icon {
  position: absolute;
  left: 25px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 0.9rem;
  color: #6c757d;
  pointer-events: none;
}

.header-collection-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  border-bottom: 1px solid #f1f3f5;
  transition: background-color 0.2s ease;
}

.header-collection-item:hover {
  background-color: #f8f9fa;
}

.header-collection-item:last-child {
  border-bottom: none;
}

.header-collection-info {
  flex: 1;
}

.header-collection-name {
  font-weight: 600;
  margin-bottom: 5px;
  color: #212529;
  font-size: 1rem;
}

.header-collection-description {
  font-size: 0.9rem;
  color: #6c757d;
  margin-bottom: 8px;
}

.header-collection-stats {
  display: flex;
  gap: 8px;
  margin-top: 5px;
}

.header-count-badge {
  display: inline-block;
  background-color: #e9ecef;
  color: #495057;
  font-size: 0.8rem;
  padding: 2px 8px;
  border-radius: 12px;
  font-weight: 500;
}

.empty-collections {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 15px;
  color: #adb5bd;
}

.empty-message {
  font-size: 1.1rem;
  font-weight: 600;
  color: #495057;
  margin-bottom: 8px;
}

.empty-description {
  font-size: 0.9rem;
  color: #6c757d;
  max-width: 300px;
  line-height: 1.5;
}

.header-collection-actions {
  display: flex;
  gap: 8px;
}

.edit-btn, .delete-btn, .export-btn, .import-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.85rem;
  font-weight: 500;
  transition: all 0.2s ease;
}

.btn-icon {
  font-size: 0.9rem;
}

.edit-btn {
  background-color: #e9ecef;
  color: #495057;
}

.edit-btn:hover {
  background-color: #dee2e6;
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.export-btn {
  background-color: #d1e7dd;
  color: #146c43;
}

.export-btn:hover {
  background-color: #badbcc;
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(20, 108, 67, 0.1);
}

.import-btn {
  background-color: #cff4fc;
  color: #055160;
}

.import-btn:hover {
  background-color: #b6effb;
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(5, 81, 96, 0.1);
}

.delete-btn {
  background-color: #f8d7da;
  color: #dc3545;
}

.delete-btn:hover {
  background-color: #f5c2c7;
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(220, 53, 69, 0.1);
}

.header-collection-form {
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 20px;
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.form-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.form-header h4 {
  margin: 0;
  font-size: 1.1rem;
  color: #343a40;
}

.editing-badge {
  background-color: #cff4fc;
  color: #055160;
  font-size: 0.8rem;
  font-weight: 600;
  padding: 3px 10px;
  border-radius: 12px;
}

.form-group {
  margin-bottom: 20px;
}

.form-label-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.form-group label {
  display: block;
  font-weight: 500;
  color: #495057;
  font-size: 0.9rem;
}

.char-counter {
  font-size: 0.75rem;
  color: #6c757d;
  transition: all 0.2s ease;
}

.char-counter.near-limit {
  color: #fd7e14;
}

.char-counter.at-limit {
  color: #dc3545;
  font-weight: 600;
}

.form-input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ced4da;
  border-radius: 6px;
  font-size: 0.95rem;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
  color: #212529;
}

.form-input:focus {
  outline: none;
  border-color: #86b7fe;
  box-shadow: 0 0 0 3px rgba(13, 110, 253, 0.1);
}

.form-input::placeholder {
  color: #adb5bd;
}

.headers-editor {
  margin-top: 25px;
  margin-bottom: 25px;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  overflow: hidden;
}

.headers-editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 15px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
}

.headers-editor-header h5 {
  margin: 0;
  font-size: 1rem;
  color: #343a40;
}

.header-count {
  font-size: 0.8rem;
  color: #6c757d;
  background-color: #e9ecef;
  padding: 3px 8px;
  border-radius: 12px;
}

.header-table {
  width: 100%;
}

.header-table-head {
  display: grid;
  grid-template-columns: 1fr 1fr 40px;
  gap: 10px;
  padding: 10px 15px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
  font-size: 0.85rem;
  font-weight: 600;
  color: #495057;
}

.header-rows {
  max-height: 250px;
  overflow-y: auto;
  padding: 5px 15px;
}

.empty-headers {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
  text-align: center;
  background-color: #f8f9fa;
  border-radius: 6px;
  margin: 10px 0;
}

.empty-headers-message {
  font-size: 0.95rem;
  font-weight: 500;
  color: #495057;
  margin-bottom: 5px;
}

.empty-headers-description {
  font-size: 0.85rem;
  color: #6c757d;
}

.header-row {
  display: grid;
  grid-template-columns: 1fr 1fr 40px;
  gap: 10px;
  margin-bottom: 10px;
  align-items: start;
  padding: 5px;
  border-radius: 6px;
  transition: background-color 0.2s ease;
}

.header-row.is-valid {
  background-color: rgba(25, 135, 84, 0.05);
}

.header-row.is-invalid {
  background-color: rgba(220, 53, 69, 0.05);
}

.input-wrapper {
  position: relative;
}

.header-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ced4da;
  border-radius: 6px;
  font-size: 0.9rem;
  transition: all 0.2s ease;
}

.header-input:focus {
  outline: none;
  border-color: #86b7fe;
  box-shadow: 0 0 0 3px rgba(13, 110, 253, 0.1);
}

.header-input.invalid-input {
  border-color: #dc3545;
}

.header-input.invalid-input:focus {
  box-shadow: 0 0 0 3px rgba(220, 53, 69, 0.1);
}

.validation-message {
  position: absolute;
  font-size: 0.75rem;
  color: #dc3545;
  margin-top: 2px;
}

.remove-header-btn {
  background-color: #f8d7da;
  color: #dc3545;
  border: none;
  border-radius: 6px;
  width: 30px;
  height: 30px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.remove-header-btn:hover {
  background-color: #f5c2c7;
  transform: translateY(-1px);
}

.remove-icon {
  font-size: 1.2rem;
  line-height: 1;
}

.add-header-btn {
  background-color: #e9ecef;
  border: 1px dashed #ced4da;
  border-radius: 6px;
  padding: 10px;
  width: calc(100% - 30px);
  margin: 15px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 0.9rem;
  color: #495057;
  font-weight: 500;
  transition: all 0.2s ease;
}

.add-header-btn:hover {
  background-color: #dee2e6;
  border-color: #adb5bd;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 25px;
  padding-top: 20px;
  border-top: 1px solid #e9ecef;
}

.save-btn, .cancel-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all 0.2s ease;
}

.save-btn {
  background-color: #0d6efd;
  color: white;
}

.save-btn:hover {
  background-color: #0b5ed7;
  transform: translateY(-1px);
  box-shadow: 0 2px 5px rgba(13, 110, 253, 0.2);
}

.cancel-btn {
  background-color: #e9ecef;
  color: #495057;
}

.cancel-btn:hover {
  background-color: #dee2e6;
  transform: translateY(-1px);
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.05);
}

@media (max-width: 768px) {
  .modal-body {
    grid-template-columns: 1fr;
  }
}
</style>