<template>
  <div class="collection-sidebar">
    <div class="sidebar-header">
      <h3>Collections</h3>
      <button @click="showNewCollectionModal = true" class="new-collection-btn">
        +
      </button>
    </div>

    <div class="collections-list">
      <div 
        v-for="collection in collections" 
        :key="collection.id"
        class="collection-item"
      >
        <div 
          class="collection-header"
          @click="toggleCollection(collection.id)"
        >
          <span class="collection-icon">
            {{ expandedCollections.has(collection.id) ? '📂' : '📁' }}
          </span>
          <span class="collection-name">{{ collection.name }}</span>
          <span class="request-count">({{ collection.requests.length }})</span>
        </div>
        
        <div 
          v-if="expandedCollections.has(collection.id)"
          class="requests-list"
        >
          <div 
            v-for="request in collection.requests" 
            :key="request.id"
            class="request-item"
            @click="loadRequest(request)"
          >
            <span :class="['method-badge', request.method.toLowerCase()]">
              {{ request.method }}
            </span>
            <span class="request-name">{{ request.name }}</span>
            <button 
              @click.stop="deleteRequest(collection.id, request.id)"
              class="delete-btn"
              title="Delete request"
            >
              ×
            </button>
          </div>
          
          <div v-if="collection.requests.length === 0" class="empty-collection">
            No requests saved
          </div>
        </div>
      </div>
      
      <div v-if="collections.length === 0" class="empty-collections">
        <p>No collections yet</p>
        <p class="empty-hint">Create a collection to organize your requests</p>
      </div>
    </div>

    <!-- New Collection Modal -->
    <div v-if="showNewCollectionModal" class="modal-overlay" @click="closeModal">
      <div class="modal" @click.stop>
        <div class="modal-header">
          <h4>Create New Collection</h4>
          <button @click="closeModal" class="close-btn">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>Collection Name</label>
            <input 
              v-model="newCollection.name" 
              type="text" 
              placeholder="Enter collection name"
              class="form-input"
              @keyup.enter="createCollection"
            />
          </div>
          <div class="form-group">
            <label>Description (optional)</label>
            <textarea 
              v-model="newCollection.description" 
              placeholder="Enter description"
              class="form-textarea"
              rows="3"
            ></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="closeModal" class="cancel-btn">Cancel</button>
          <button @click="createCollection" class="create-btn">Create</button>
        </div>
      </div>
    </div>
    
    <!-- Confirm Dialog -->
    <ConfirmDialog
      ref="confirmDialog"
      title="Delete Request"
      message="Are you sure you want to delete this request? This action cannot be undone."
      confirm-text="Delete"
      cancel-text="Cancel"
      @confirm="confirmDelete"
      @cancel="cancelDelete"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { CollectionService } from '../../bindings/changeme'
import ConfirmDialog from './ConfirmDialog.vue'

const emit = defineEmits(['load-request'])

const collections = ref([])
const expandedCollections = ref(new Set(['default']))
const showNewCollectionModal = ref(false)
const newCollection = ref({
  name: '',
  description: ''
})
const confirmDialog = ref(null)
let pendingDeleteRequest = null

// Load collections on mount
onMounted(async () => {
  await loadCollections()
})

// Methods
const loadCollections = async () => {
  try {
    const allCollections = await CollectionService.GetAllCollections()
    collections.value = allCollections
    
    // Create default collection if none exist
    if (collections.value.length === 0) {
      await CollectionService.CreateCollection('Default Collection', 'Default collection for requests')
      await loadCollections()
    }
  } catch (error) {
    console.error('Failed to load collections:', error)
  }
}

const toggleCollection = (collectionId) => {
  if (expandedCollections.value.has(collectionId)) {
    expandedCollections.value.delete(collectionId)
  } else {
    expandedCollections.value.add(collectionId)
  }
}

const loadRequest = (request) => {
  emit('load-request', request)
}

const deleteRequest = (collectionId, requestId) => {
  // Store the pending delete request
  pendingDeleteRequest = { collectionId, requestId }
  // Show the confirmation dialog
  confirmDialog.value.show()
}

const confirmDelete = async () => {
  if (pendingDeleteRequest) {
    try {
      await CollectionService.DeleteRequest(pendingDeleteRequest.collectionId, pendingDeleteRequest.requestId)
      await loadCollections()
      console.log('Request deleted successfully')
    } catch (error) {
      console.error('Failed to delete request:', error)
      alert('Failed to delete request')
    } finally {
      pendingDeleteRequest = null
    }
  }
}

const cancelDelete = () => {
  pendingDeleteRequest = null
}

const createCollection = async () => {
  if (!newCollection.value.name.trim()) {
    alert('Please enter a collection name')
    return
  }
  
  try {
    await CollectionService.CreateCollection(
      newCollection.value.name.trim(),
      newCollection.value.description.trim()
    )
    await loadCollections()
    closeModal()
  } catch (error) {
    console.error('Failed to create collection:', error)
    alert('Failed to create collection')
  }
}

const closeModal = () => {
  showNewCollectionModal.value = false
  newCollection.value = {
    name: '',
    description: ''
  }
}

// Expose refresh method
const refresh = () => {
  loadCollections()
}

defineExpose({ refresh })
</script>

<style scoped>
.collection-sidebar {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #eee;
}

.sidebar-header h3 {
  margin: 0;
  color: #333;
  font-size: 18px;
}

.new-collection-btn {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  border: 1px solid #007bff;
  background: #007bff;
  color: white;
  cursor: pointer;
  font-size: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.new-collection-btn:hover {
  background: #0056b3;
}

.collections-list {
  flex: 1;
  overflow-y: auto;
}

.collection-item {
  margin-bottom: 15px;
}

.collection-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
  cursor: pointer;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.collection-header:hover {
  background: #f8f9fa;
}

.collection-icon {
  font-size: 16px;
}

.collection-name {
  font-weight: 500;
  color: #333;
  flex: 1;
}

.request-count {
  font-size: 12px;
  color: #6c757d;
}

.requests-list {
  margin-left: 24px;
  border-left: 2px solid #f1f3f4;
  padding-left: 15px;
}

.request-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 8px;
  margin: 4px 0;
  cursor: pointer;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.request-item:hover {
  background: #f8f9fa;
}

.method-badge {
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 10px;
  font-weight: bold;
  text-transform: uppercase;
  min-width: 40px;
  text-align: center;
}

.method-badge.get {
  background: #d4edda;
  color: #155724;
}

.method-badge.post {
  background: #fff3cd;
  color: #856404;
}

.method-badge.put {
  background: #cce5ff;
  color: #004085;
}

.method-badge.delete {
  background: #f8d7da;
  color: #721c24;
}

.method-badge.patch {
  background: #e2e3e5;
  color: #383d41;
}

.method-badge.head,
.method-badge.options {
  background: #f0f0f0;
  color: #666;
}

.request-name {
  flex: 1;
  font-size: 13px;
  color: #495057;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.delete-btn {
  width: 20px;
  height: 20px;
  border: none;
  background: #dc3545;
  color: white;
  border-radius: 50%;
  cursor: pointer;
  font-size: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s;
}

.request-item:hover .delete-btn {
  opacity: 1;
}

.delete-btn:hover {
  background: #c82333;
}

.empty-collection,
.empty-collections {
  text-align: center;
  color: #6c757d;
  font-style: italic;
  padding: 20px;
}

.empty-hint {
  font-size: 12px;
  margin-top: 5px;
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

.modal {
  background: white;
  border-radius: 8px;
  width: 90%;
  max-width: 400px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #eee;
}

.modal-header h4 {
  margin: 0;
  color: #333;
}

.close-btn {
  width: 30px;
  height: 30px;
  border: none;
  background: none;
  font-size: 20px;
  cursor: pointer;
  color: #6c757d;
}

.modal-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: 500;
  color: #333;
}

.form-input,
.form-textarea {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
}

.form-textarea {
  resize: vertical;
  font-family: inherit;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 20px;
  border-top: 1px solid #eee;
}

.cancel-btn,
.create-btn {
  padding: 8px 16px;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.cancel-btn {
  background: white;
  color: #6c757d;
}

.cancel-btn:hover {
  background: #f8f9fa;
}

.create-btn {
  background: #007bff;
  color: white;
  border-color: #007bff;
}

.create-btn:hover {
  background: #0056b3;
}
</style>