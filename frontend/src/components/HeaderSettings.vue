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
                <div class="collection-count">{{ collection.headerCollections.length }} collections</div>
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

            
            <div v-if="collection.headerCollections.length === 0" class="empty-collections">
              <div class="empty-icon">📋</div>
              <div class="empty-message">No header collections yet</div>
              <div class="empty-description">Create your first header collection to reuse headers across requests</div>
            </div>
            
            <div v-else v-for="headerCollection in collection.headerCollections" :key="headerCollection.id" class="header-collection-item">
              <div class="header-collection-info">
                <div class="header-collection-name">{{ headerCollection.name }}</div>
                <div class="header-collection-description">{{ headerCollection.description }}</div>
                <div class="header-collection-stats">
                  <span class="header-count-badge">{{ Object.keys(headerCollection.headers || {}).length }} headers</span>
                </div>
              </div>
              <div class="header-collection-actions">
                <button @click="editHeaderCollection(headerCollection)" class="edit-btn" title="Edit collection">
                  <span class="btn-icon">✏️</span>
                  <span class="btn-text">Edit</span>
                </button>
                <button @click="exportHeaderCollection(headerCollection)" class="export-btn" title="Export collection">
                  <span class="btn-icon">📤</span>
                  <span class="btn-text">Export</span>
                </button>
                <button @click="deleteHeaderCollection(headerCollection)" class="delete-btn" title="Delete collection">
                  <span class="btn-icon">🗑️</span>
                  <span class="btn-text">Delete</span>
                </button>
              </div>
            </div>
          </div>

          <!-- Add/Edit Header Collection Form -->
          <HeaderCollectionForm
            v-model:collection="formData"
            :is-editing="!!editingCollection"
            @submit="saveHeaderCollection"
            @cancel="cancelEdit"
          />
        </div>
      </div>
    </div>
    <ConfirmDialog 
      ref="confirmDialog"
      title="Delete Header Collection"
      :message="confirmMessage"
      confirm-text="Delete"
      cancel-text="Cancel"
      @confirm="handleConfirmDelete"
      @cancel="handleCancelDelete"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import { CollectionService, EventBusService } from '../../bindings/captain-api'
import { HeaderCollection as ModelHeaderCollection } from '../../bindings/captain-api/models'
import HeaderCollectionForm from './HeaderCollectionForm.vue'
import ConfirmDialog from './ConfirmDialog.vue'
import popupService from './popupService'
// Using local lightweight types to avoid relying on generated model types at build time
type HeaderCollection = {
  id: string
  name: string
  description?: string
  headers: Record<string, string>
  createdAt?: string | null
  updatedAt?: string | null
}
type Collection = {
  id: string
  headerCollections: any[]
}

// Types
type HeaderKV = { key: string; value: string }

const confirmDialog = ref<InstanceType<typeof ConfirmDialog> | null>(null)
const collectionToDelete = ref<HeaderCollection | null>(null)

const confirmMessage = computed(() => {
  if (collectionToDelete.value) {
    return `Are you sure you want to delete the "${collectionToDelete.value.name}" header collection? This action cannot be undone.`
  }
  return 'Are you sure? This action cannot be undone.'
})

// Define props and emits
const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  collection: {
    type: Object as () => Collection,
    required: true
  }
})

const emit = defineEmits(['close', 'open'])

const showManageModal = computed(() => props.show)
const editingCollection = ref<HeaderCollection | null>(null)
const searchQuery = ref('')
const sortOption = ref('name-asc') // Default sort option
// Form data for new collections
const newCollectionTemplate: { id: string; name: string; description: string; headers: HeaderKV[] } = {
  id: '',
  name: '',
  description: '',
  headers: [{ key: '', value: '' }]
}

// Current form data (reactive)
const formData = ref<{ id: string; name: string; description: string; headers: HeaderKV[] }>({ ...newCollectionTemplate })

// Filter and sort collections based on search query and sort option
const filteredCollections = computed(() => {
  // First filter by search query
  let filtered = props.collection.headerCollections
  
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
  emit('open')
}

// Load header collections on component mount and set up event listener
onMounted(async () => {
  // Add event listener for opening header settings from other components
  document.addEventListener('open-header-settings', handleOpenHeaderSettings)
})

// Clean up event listener when component is unmounted
onUnmounted(() => {
  document.removeEventListener('open-header-settings', handleOpenHeaderSettings)
})

// Edit a header collection
const editHeaderCollection = (collection: HeaderCollection) => {
  console.log('Starting edit of collection:', collection)
  
  // Get the latest version of the collection from the list
  const latestCollection = props.collection.headerCollections.find(c => c.id === collection.id) || collection
  
  // Store the collection being edited
  editingCollection.value = { ...latestCollection }
  
  // Create a clean copy of the headers
  let headers: HeaderKV[] = []
  if (latestCollection.headers) {
    if (Array.isArray(latestCollection.headers)) {
      headers = (latestCollection.headers as any[]).map((h: any) => ({ key: String(h.key || ''), value: String(h.value || '') }))
    } else if (typeof latestCollection.headers === 'object') {
      headers = Object.entries(latestCollection.headers as Record<string, string>).map(([key, value]) => ({ key, value }))
    }
  }
  
  // Ensure there's at least one header row
  if (headers.length === 0) {
    headers = [{ key: '', value: '' }]
  }
  
  // Update form data in a single assignment
  formData.value = {
    id: latestCollection.id || '',
    name: latestCollection.name || '',
    description: latestCollection.description || '',
    headers
  }
  
  console.log('Form data after edit:', { ...formData.value })
}

// Cancel editing
const cancelEdit = () => {
  editingCollection.value = null
  formData.value = { ...newCollectionTemplate }
}

// Validate form data
const validateForm = async (): Promise<boolean> => {
  // Check if name is provided
  if (!formData.value.name.trim()) {
    await popupService.alert('Please provide a collection name', {
      severity: 'warning'
    })
    return false
  }
  
  // Check for duplicate header keys
  const validHeaders = formData.value.headers.filter(h => h.key && h.value)
  const headerKeys = validHeaders.map(h => h.key.toLowerCase())
  const uniqueKeys = new Set(headerKeys)
  
  if (headerKeys.length !== uniqueKeys.size) {
    await popupService.alert('Duplicate header names found. Each header name must be unique.', {
      severity: 'warning'
    })
    return false
  }
  
  // Check if any header is incomplete (has key but no value or vice versa)
  const incompleteHeaders = formData.value.headers.filter(
    h => (h.key && !h.value) || (!h.key && h.value)
  )
  
  if (incompleteHeaders.length > 0) {
    await popupService.alert('Some headers are incomplete. Please provide both name and value for each header.', {
      severity: 'warning'
    })
    return false
  }
  
  return true
}

// Save header collection
const saveHeaderCollection = async (collectionData: { id: string; name: string; description: string; headers: HeaderKV[] }) => {
  try {
    console.log('Saving collection with data:', collectionData)
    // Validate before saving
    if (!(await validateForm())) {
      return
    }
    
    // Convert headers array to object for the template
    const headersObject = {}
    const validHeaders = collectionData.headers.filter(h => h.key && h.value)
    
    validHeaders.forEach(header => {
      headersObject[header.key] = header.value
    })
    
    const now = new Date().toISOString()
    
    const collectionPayload = ModelHeaderCollection.createFrom({
      id: collectionData.id || `collection_${Date.now()}`,
      collectionId: props.collection.id,
      name: collectionData.name.trim(),
      description: (collectionData.description || '').trim(),
      headers: headersObject,
      createdAt: editingCollection.value?.createdAt || null,
      updatedAt: now
    })
    
    if (editingCollection.value) {
      // Update existing collection
      console.log('Updating collection:', collectionPayload)
      const updatedResult = await CollectionService.UpdateHeaderCollection(props.collection.id, collectionPayload)
      
      // Update in local array
      const index = props.collection.headerCollections.findIndex((c: any) => c.id === collectionPayload.id)
      if (index !== -1 && updatedResult) {
        props.collection.headerCollections[index] = updatedResult
        EventBusService.EmitEvent('header-collection-updated', updatedResult)
      }
    } else {
      // Create new collection
      const newCollectionResult = await CollectionService.CreateHeaderCollection(props.collection.id, collectionPayload)
      if (newCollectionResult) {
        props.collection.headerCollections.push(newCollectionResult)
        EventBusService.EmitEvent('header-collection-updated', newCollectionResult)
      } else {
        throw new Error('Failed to create collection: No data returned from server')
      }
    }
    
    // Reset form and close modal
    cancelEdit()
    
    // Show success message
    console.log('Header collection saved successfully')
  } catch (error) {
    console.error('Failed to save header collection:', error)
  }
}

// Delete header collection with confirmation
const deleteHeaderCollection = (collection: HeaderCollection) => {
  collectionToDelete.value = collection
  confirmDialog.value?.show()
}

const handleConfirmDelete = async () => {
  if (collectionToDelete.value) {
    try {
      await CollectionService.DeleteHeaderCollection(props.collection.id, collectionToDelete.value.id)
      EventBusService.EmitEvent('header-collection-updated')
    } catch (error) {
      console.error('Failed to delete header collection:', error)
    } finally {
      collectionToDelete.value = null
    }
  }
}

const handleCancelDelete = () => {
  collectionToDelete.value = null
}

// Export header collection as JSON file
const exportHeaderCollection = async (collection: HeaderCollection) => {
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
    await popupService.alert('Failed to export header collection', {
      severity: 'error'
    })
  }
}

// Trigger file input click to open file dialog
const fileInput = ref<HTMLInputElement | null>(null)
const triggerImportFile = () => {
  fileInput.value?.click()
}

// Import header collection from JSON file
const importHeaderCollection = (event: Event) => {
  const input = event.target as HTMLInputElement | null
  const file = input?.files?.[0]
  if (!file) return
  
  const reader = new FileReader()
  reader.onload = async (e) => {
    try {
      const result = e.target ? (e.target as FileReader).result : null
      const importedData = JSON.parse((result as string) || '{}')
      
      // Validate the imported data
      if (!importedData.name || !importedData.headers) {
        throw new Error('Invalid header collection format')
      }
      
      // Create a new collection with a unique ID
      const newCollectionPayload = ModelHeaderCollection.createFrom({
        id: 'header_col_' + Date.now(),
        collectionId: props.collection.id,
        name: importedData.name,
        description: (importedData.description || '').toString(),
        headers: importedData.headers,
        createdAt: null,
        updatedAt: new Date().toISOString()
      })
      
      // Check for duplicate name
      const existingNames = props.collection.headerCollections.map((c: any) => (c.name || '').toLowerCase())
      if (existingNames.includes((newCollectionPayload.name || '').toLowerCase())) {
        newCollectionPayload.name = `${newCollectionPayload.name} (Imported)`
      }
      
      // Persist via backend and update local state
      const created = await CollectionService.CreateHeaderCollection(props.collection.id, newCollectionPayload)
      if (created) {
        props.collection.headerCollections.push(created)
        EventBusService.EmitEvent('header-collection-updated', created)
      }
      
      // Show success message
      await popupService.alert(`Successfully imported "${newCollectionPayload.name}" collection`, {
        severity: 'success'
      })
    } catch (error) {
      console.error('Failed to import header collection:', error)
      await popupService.alert('Failed to import header collection: Invalid format', {
        severity: 'error'
      })
    }
    
    // Reset file input
    if (input) input.value = ''
  }
  
  reader.onerror = () => {
    popupService.alert('Error reading file', {
      severity: 'error'
    })
    if (input) input.value = ''
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