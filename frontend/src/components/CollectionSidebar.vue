<template>
  <div class="collection-sidebar">
    <div class="sidebar-header">
      <div class="view-toggle">
        <button 
          @click="switchView('collections')" 
          :class="['toggle-btn', { active: currentView === 'collections' }]"
        >
          <span class="toggle-icon">📁</span>
          Collections
        </button>
        <button 
          @click="switchView('logs')" 
          :class="['toggle-btn', { active: currentView === 'logs' }]"
        >
          <span class="toggle-icon">📋</span>
          Logs
        </button>
      </div>
      <div class="header-actions">
        <button 
          @click="$emit('new-request')" 
          class="new-collection-btn"
          title="Create new request"
          :disabled="currentView === 'logs'"
        >
          <span class="btn-icon">+</span>
          <span class="btn-text">Request</span>
        </button>
        <button 
          @click="showNewCollectionModal = true" 
          class="new-collection-btn"
          title="Create new collection"
          :disabled="currentView === 'logs'"
        >
          <span class="btn-icon">+</span>
          <span class="btn-text">Collection</span>
        </button>
      </div>
    </div>

    <!-- Collections View -->
    <div v-if="currentView === 'collections'" class="collections-list">
      <div 
        v-for="collection in collections" 
        :key="collection.id"
        class="collection-item"
        :class="{ 'expanded': expandedCollections.has(collection.id) }"
      >
        <div class="collection-header">
          <div class="collection-main" @click="toggleCollection(collection.id)">
            <span class="collection-icon">
              {{ expandedCollections.has(collection.id) ? '📂' : '📁' }}
            </span>
            <div class="collection-info">
              <span class="collection-name">{{ collection.name }}</span>
              <span class="collection-meta">
                <span class="request-count">{{ collection.requests.length }} requests</span>
                <span v-if="getActiveEnvironment(collection.id)" class="active-env">
                  • {{ getActiveEnvironment(collection.id).name }}
                </span>
              </span>
            </div>
          </div>
          <div class="collection-actions">
            <button 
              @click.stop="openCollectionMenu(collection.id, $event)"
              class="action-btn menu-btn"
              title="Collection options"
            >
              ⋯
            </button>
          </div>
        </div>
        
        <!-- Collection Settings (Compact) -->
        <div v-if="expandedCollections.has(collection.id)" class="collection-settings">
          <div class="settings-grid">
            <!-- Environment Selector -->
            <div class="setting-item">
              <label class="setting-label">Environment</label>
              <select 
                :value="getActiveEnvironment(collection.id)?.id || ''"
                @change="setActiveEnvironment(collection.id, $event.target.value)"
                @click.stop
                class="setting-select"
              >
                <option value="">None</option>
                <option 
                  v-for="env in collection.environments || []"
                  :key="env.id"
                  :value="env.id"
                >
                  {{ env.name }}
                </option>
              </select>
            </div>
            
            <!-- Header Collection Selector -->
            <div class="setting-item">
              <label class="setting-label">Headers</label>
              <select 
                :value="activeHeaderCollections.get(collection.id) || ''"
                @change="setActiveHeaderCollection(collection.id, $event.target.value)"
                @click.stop
                class="setting-select"
              >
                <option value="">None</option>
                <option 
                  v-for="headerCollection in headerCollections"
                  :key="headerCollection.id"
                  :value="headerCollection.id"
                >
                  {{ headerCollection.name }}
                </option>
              </select>
            </div>
          </div>
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
        <div class="empty-state">
          <div class="empty-icon">📁</div>
          <h3>No collections yet</h3>
          <p>Create your first collection to organize your API requests</p>
          <button @click="showNewCollectionModal = true" class="empty-action-btn">
            Create Collection
          </button>
        </div>
      </div>
    </div>

    <!-- Collection Context Menu -->
    <div 
      v-if="showCollectionMenu" 
      class="collection-menu"
      :style="{ top: menuPosition.y + 'px', left: menuPosition.x + 'px' }"
      @click.stop
    >
      <div class="menu-item" @click="openEnvironmentModal(selectedCollectionId)">
        <span class="menu-icon">⚙️</span>
        Manage Environments
      </div>
      <div class="menu-item" @click="openHeaderSettingsModal(selectedCollectionId)">
        <span class="menu-icon">📋</span>
        Manage Headers
      </div>
      <div class="menu-divider"></div>
      <div class="menu-item" @click="renameCollection(selectedCollectionId)">
        <span class="menu-icon">✏️</span>
        Rename
      </div>
      <div class="menu-item" @click="duplicateCollection(selectedCollectionId)">
        <span class="menu-icon">📋</span>
        Duplicate
      </div>
      <div class="menu-item" @click="exportCollection(selectedCollectionId)">
        <span class="menu-icon">📤</span>
        Export
      </div>
      <div class="menu-divider"></div>
      <div class="menu-item danger" @click="confirmDeleteCollection(selectedCollectionId)">
        <span class="menu-icon">🗑️</span>
        Delete
      </div>
    

    <!-- Logs View -->
    <RequestLogs 
      v-if="currentView === 'logs'"
      ref="requestLogsRef"
      @load-request="loadRequest"
    />

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
    
    <!-- Environment Management Modal -->
    <div v-if="showEnvironmentModal" class="modal-overlay" @click="closeEnvironmentModal">
      <div class="modal environment-modal" @click.stop>
        <div class="modal-header">
          <h4>Manage Environments - {{ getCollectionName(selectedCollectionId) }}</h4>
          <button @click="closeEnvironmentModal" class="close-btn">×</button>
        </div>
        <div class="modal-body">
          <div class="environments-list">
            <div 
              v-for="env in currentEnvironments"
              :key="env.id"
              class="environment-item"
              :class="{ active: env.isActive }"
            >
              <div class="env-info">
                <div class="env-name">{{ env.name }}</div>
                <div class="env-url">{{ env.baseURL }}</div>
                <div class="env-description">{{ env.description }}</div>
              </div>
              <div class="env-actions">
                <button 
                  v-if="!env.isActive"
                  @click="activateEnvironment(env.id)"
                  class="activate-btn"
                >
                  Activate
                </button>
                <button 
                  @click="editEnvironment(env)"
                  class="edit-btn"
                >
                  Edit
                </button>
                <button 
                  v-if="currentEnvironments.length > 1"
                  @click="deleteEnvironment(env.id)"
                  class="delete-env-btn"
                >
                  Delete
                </button>
              </div>
            </div>
          </div>
          
          <div class="add-environment">
            <button @click="showAddEnvironmentForm" class="add-env-btn">
              + Add Environment
            </button>
          </div>
          
          <!-- Add/Edit Environment Form -->
          <div v-if="showEnvForm" class="env-form">
            <h5>{{ editingEnvironment ? 'Edit' : 'Add' }} Environment</h5>
            <div class="form-group">
              <label>Name</label>
              <input 
                v-model="envForm.name" 
                type="text" 
                placeholder="Environment name"
                class="form-input"
              />
            </div>
            <div class="form-group">
              <label>Base URL</label>
              <input 
                v-model="envForm.baseURL" 
                type="text" 
                placeholder="https://api.example.com"
                class="form-input"
              />
            </div>
            <div class="form-group">
              <label>Description</label>
              <input 
                v-model="envForm.description" 
                type="text" 
                placeholder="Environment description"
                class="form-input"
              />
            </div>
            
            <div class="form-actions">
              <button @click="cancelEnvForm" class="cancel-btn">Cancel</button>
              <button @click="saveEnvironment" class="save-btn">Save</button>
            </div>
          </div>
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
    
    <!-- Header Settings Modal -->
    <div v-if="showHeaderSettingsModal" class="modal-overlay" @click="closeHeaderSettingsModal">
      <div class="modal-content header-settings-modal" @click.stop>
        <HeaderSettings :show="showHeaderSettingsModal" @close="closeHeaderSettingsModal" />
      </div>
    </div>

<!-- Logs View -->
<RequestLogs 
  v-if="currentView === 'logs'"
  ref="requestLogsRef"
  @load-request="loadRequest"
/>

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
  
<!-- Environment Management Modal -->
<div v-if="showEnvironmentModal" class="modal-overlay" @click="closeEnvironmentModal">
  <div class="modal environment-modal" @click.stop>
    <div class="modal-header">
      <h4>Manage Environments - {{ getCollectionName(selectedCollectionId) }}</h4>
      <button @click="closeEnvironmentModal" class="close-btn">×</button>
    </div>
    <div class="modal-body">
      <div class="environments-list">
        <div 
          v-for="env in currentEnvironments"
          :key="env.id"
          class="environment-item"
          :class="{ active: env.isActive }"
        >
          <div class="env-info">
            <div class="env-name">{{ env.name }}</div>
            <div class="env-url">{{ env.baseURL }}</div>
            <div class="env-description">{{ env.description }}</div>
          </div>
          <div class="env-actions">
            <button 
              v-if="!env.isActive"
              @click="activateEnvironment(env.id)"
              class="activate-btn"
            >
              Activate
            </button>
            <button 
              @click="editEnvironment(env)"
              class="edit-btn"
            >
              Edit
            </button>
            <button 
              v-if="currentEnvironments.length > 1"
              @click="deleteEnvironment(env.id)"
              class="delete-env-btn"
            >
              Delete
            </button>
          </div>
        </div>
      </div>
      
      <div class="add-environment">
        <button @click="showAddEnvironmentForm" class="add-env-btn">
          + Add Environment
        </button>
      </div>
      
      <!-- Add/Edit Environment Form -->
      <div v-if="showEnvForm" class="env-form">
        <h5>{{ editingEnvironment ? 'Edit' : 'Add' }} Environment</h5>
        <div class="form-group">
          <label>Name</label>
          <input 
            v-model="envForm.name" 
            type="text" 
            placeholder="Environment name"
            class="form-input"
          />
        </div>
        <div class="form-group">
          <label>Base URL</label>
          <input 
            v-model="envForm.baseURL" 
            type="text" 
            placeholder="https://api.example.com"
            class="form-input"
          />
        </div>
        <div class="form-group">
          <label>Description</label>
          <input 
            v-model="envForm.description" 
            type="text" 
            placeholder="Environment description"
            class="form-input"
          />
        </div>
        
        <div class="form-actions">
          <button @click="cancelEnvForm" class="cancel-btn">Cancel</button>
          <button @click="saveEnvironment" class="save-btn">Save</button>
        </div>
      </div>
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
  
<!-- Header Settings Modal -->
<div v-if="showHeaderSettingsModal" class="modal-overlay" @click="closeHeaderSettingsModal">
  <div class="modal-content header-settings-modal" @click.stop>
    <HeaderSettings :show="showHeaderSettingsModal" @close="closeHeaderSettingsModal" />
  </div>
</div>
</div>
</div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { CollectionService } from '../../bindings/captain-api'
import { GetHeaderCollections } from '../../bindings/captain-api/headerservice'
import ConfirmDialog from './ConfirmDialog.vue'
import RequestLogs from './RequestLogs.vue'
import { Events } from '@wailsio/runtime'
import HeaderSettings from './HeaderSettings.vue'

const emit = defineEmits(['load-request', 'new-request'])

const currentView = ref('collections')
const collections = ref([])
const expandedCollections = ref(new Set(['default']))
const showNewCollectionModal = ref(false)
const newCollection = ref({
  name: '',
  description: ''
})

// Header collection management
const headerCollections = ref([])
const activeHeaderCollections = ref(new Map()) // collectionId -> headerCollectionId
const showHeaderSettingsModal = ref(false)

// Environment management
const showEnvironmentModal = ref(false)
const selectedCollectionId = ref('')
const currentEnvironments = ref([])
const showEnvForm = ref(false)
const editingEnvironment = ref(null)
const envForm = ref({
  id: '',
  name: '',
  baseURL: '',
  description: ''
})

// Collection menu
const showCollectionMenu = ref(false)
const menuPosition = ref({ x: 0, y: 0 })

const confirmDialog = ref(null)
const requestLogsRef = ref(null)
let pendingDeleteRequest = null
// Track event subscriptions
const eventSubscriptions = []

onMounted(() => {
  console.log('CollectionSidebar mounted, loading collections...')
  loadCollections()
  loadHeaderCollections()
  document.addEventListener('click', closeMenu)

  // Setup event listeners
  console.log('Setting up event listeners')
  
  // Handle request-saved event
  const requestSavedUnsubscribe = Events.On('request-saved', async () => {
    console.log('request-saved event received, reloading collections...')
    await loadCollections()
    
    // Small delay to ensure UI has updated
    setTimeout(() => {
      // Emit an event to notify other components that collections were updated
      Events.Emit('collections-updated')
    }, 100)
  })
  
  if (requestSavedUnsubscribe) {
    eventSubscriptions.push(requestSavedUnsubscribe)
  }
  
  // Handle collections-updated event (in case another component updates collections)
  const collectionsUpdatedUnsubscribe = Events.On('collections-updated', async () => {
    console.log('collections-updated event received, reloading collections...')
    await loadCollections()
  })
  
  if (collectionsUpdatedUnsubscribe) {
    eventSubscriptions.push(collectionsUpdatedUnsubscribe)
  }
})

// Clean up event listeners
onUnmounted(() => {
  console.log('CollectionSidebar unmounted, cleaning up event listeners')
  document.removeEventListener('click', closeMenu)
  
  // Unsubscribe from all events
  eventSubscriptions.forEach(unsubscribe => {
    if (typeof unsubscribe === 'function') {
      unsubscribe()
    }
  })
  eventSubscriptions.length = 0
})

onUnmounted(() => {
  document.removeEventListener('click', closeMenu)
  if (unsubscribeRequestSaved) {
    unsubscribeRequestSaved()
  }
})

// Load collections and update the UI
async function loadCollections() {
  console.log('Loading collections...')
  try {
    const allCollections = await CollectionService.GetAllCollections()
    console.log('Fetched collections:', allCollections)
    
    // Update the collections array in a way that triggers Vue's reactivity
    collections.value = [...allCollections]
    
    // Create default collection if none exist
    if (collections.value.length === 0) {
      console.log('No collections found, creating default collection...')
      await CollectionService.CreateCollection('Default Collection', 'Default collection for requests')
      await loadCollections() // Recursively load after creating default collection
      return
    }
    
    console.log('Collections loaded successfully:', collections.value)
    
    // If there's a selected collection, ensure it's still valid
    if (selectedCollectionId.value) {
      const collectionExists = collections.value.some(c => c.id === selectedCollectionId.value)
      if (!collectionExists && collections.value.length > 0) {
        selectedCollectionId.value = collections.value[0].id
      }
    } else if (collections.value.length > 0) {
      selectedCollectionId.value = collections.value[0].id
    }
    
    // Force update the selected request if needed
    if (selectedRequestId.value) {
      const requestExists = collections.value.some(c => 
        c.requests && c.requests.some(r => r.id === selectedRequestId.value)
      )
      if (!requestExists) {
        selectedRequestId.value = ''
      }
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

const switchView = async (view) => {
  currentView.value = view
  // Refresh logs when switching to logs view
  if (view === 'logs' && requestLogsRef.value) {
    await requestLogsRef.value.loadLogs()
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

// Environment management methods
const getActiveEnvironment = (collectionId) => {
  const collection = collections.value.find(c => c.id === collectionId)
  if (!collection || !collection.environments) return null
  return collection.environments.find(env => env.isActive)
}

const getCollectionName = (collectionId) => {
  const collection = collections.value.find(c => c.id === collectionId)
  return collection ? collection.name : ''
}

const setActiveEnvironment = async (collectionId, envId) => {
  try {
    await CollectionService.SetActiveCollectionEnvironment(collectionId, envId)
    await loadCollections()
  } catch (error) {
    console.error('Failed to set active environment:', error)
    alert('Failed to set active environment')
  }
}

const openEnvironmentModal = async (collectionId) => {
  selectedCollectionId.value = collectionId
  try {
    const environments = await CollectionService.GetCollectionEnvironments(collectionId)
    currentEnvironments.value = environments
    showEnvironmentModal.value = true
  } catch (error) {
    console.error('Failed to load environments:', error)
    alert('Failed to load environments')
  }
}

const closeEnvironmentModal = () => {
  showEnvironmentModal.value = false
  selectedCollectionId.value = ''
  currentEnvironments.value = []
  showEnvForm.value = false
  editingEnvironment.value = null
  resetEnvForm()
}

const showAddEnvironmentForm = () => {
  editingEnvironment.value = null
  resetEnvForm()
  showEnvForm.value = true
}

const editEnvironment = (env) => {
  editingEnvironment.value = env
  envForm.value = {
    id: env.id,
    name: env.name,
    baseURL: env.baseURL,
    description: env.description
  }
  
  showEnvForm.value = true
}

const resetEnvForm = () => {
  envForm.value = {
    id: '',
    name: '',
    baseURL: '',
    description: ''
  }
}

const cancelEnvForm = () => {
  showEnvForm.value = false
  editingEnvironment.value = null
  resetEnvForm()
}

const saveEnvironment = async () => {
  if (!envForm.value.name.trim() || !envForm.value.baseURL.trim()) {
    alert('Please fill in name and base URL')
    return
  }
  
  try {
    const envData = {
      id: envForm.value.id,
      name: envForm.value.name.trim(),
      baseURL: envForm.value.baseURL.trim(),
      description: envForm.value.description.trim(),
      isActive: false
    }
    
    if (editingEnvironment.value) {
      await CollectionService.UpdateCollectionEnvironment(selectedCollectionId.value, envData)
    } else {
      await CollectionService.AddCollectionEnvironment(selectedCollectionId.value, envData)
    }
    
    // Refresh environments
    const environments = await CollectionService.GetCollectionEnvironments(selectedCollectionId.value)
    currentEnvironments.value = environments
    await loadCollections()
    
    cancelEnvForm()
  } catch (error) {
    console.error('Failed to save environment:', error)
    alert('Failed to save environment')
  }
}

const activateEnvironment = async (envId) => {
  await setActiveEnvironment(selectedCollectionId.value, envId)
  // Refresh environments
  const environments = await CollectionService.GetCollectionEnvironments(selectedCollectionId.value)
  currentEnvironments.value = environments
}

const deleteEnvironment = async (envId) => {
  if (!confirm('Are you sure you want to delete this environment?')) {
    return
  }
  
  try {
    await CollectionService.DeleteCollectionEnvironment(selectedCollectionId.value, envId)
    // Refresh environments
    const environments = await CollectionService.GetCollectionEnvironments(selectedCollectionId.value)
    currentEnvironments.value = environments
    await loadCollections()
  } catch (error) {
    console.error('Failed to delete environment:', error)
    alert('Failed to delete environment')
  }
}

// Header collection methods
const loadHeaderCollections = async () => {
  try {
    const collections = await GetHeaderCollections()
    
    // Convert backend format to frontend format
    headerCollections.value = collections.map(collection => {
      // Extract headers from headerTemplates if available
      const headers = {}
      if (collection.headerTemplates && collection.headerTemplates.length > 0) {
        // Use the first template's headers
        const template = collection.headerTemplates[0]
        if (template && template.headers) {
          console.log('Template headers:', template.headers,typeof template.headers)
          Object.entries(template.headers).forEach(([key, value]) => {
            console.log('list Key:', key, 'Value:', value)
          })
          for(let key in template.headers){
            console.log('in Key:', key, 'Value:', template.headers[key])
            headers[key]=template.headers[key]
          }
          for(let key of template.headers){
            console.log('of Key:', key.key, 'Value:', key.value)
          }
          // Convert each header value to string to prevent [object Object] display
          Object.entries(template.headers).forEach(([key, value]) => {
           // headers[key] = typeof value === 'object' ? JSON.stringify(value) : String(value)
          })
        }
      }
      
      return {
        id: collection.id,
        name: collection.name,
        headers
      }
    })
  } catch (error) {
    console.error('Failed to load header collections:', error)
  }
}

const getActiveHeaderCollection = (collectionId) => {
  return activeHeaderCollections.value.get(collectionId) || ''
}

const getActiveHeaderCollectionData = (collectionId) => {
  const activeId = activeHeaderCollections.value.get(collectionId)
  return headerCollections.value.find(hc => hc.id === activeId)
}

const setActiveHeaderCollection = (collectionId, headerCollectionId) => {
  if (headerCollectionId) {
    activeHeaderCollections.value.set(collectionId, headerCollectionId)
  } else {
    activeHeaderCollections.value.delete(collectionId)
  }
}

// Open header manager modal
const openHeaderManager = () => {
  showHeaderSettingsModal.value = true
}

const openHeaderSettingsModal = (collectionId) => {
  selectedCollectionId.value = collectionId
  showHeaderSettingsModal.value = true
}

const closeHeaderSettingsModal = () => {
  showHeaderSettingsModal.value = false
  selectedCollectionId.value = ''
}

// Header manager function removed to simplify the application

// Collection menu methods
const openCollectionMenu = (collectionId, event) => {
  selectedCollectionId.value = collectionId
  menuPosition.value = {
    x: event.clientX,
    y: event.clientY
  }
  showCollectionMenu.value = true
  
  // Close menu when clicking outside
  const closeMenu = () => {
    showCollectionMenu.value = false
    document.removeEventListener('click', closeMenu)
  }
  setTimeout(() => {
    document.addEventListener('click', closeMenu)
  }, 0)
}

const renameCollection = async (collectionId) => {
  const collection = collections.value.find(c => c.id === collectionId)
  if (!collection) return
  
  const newName = prompt('Enter new collection name:', collection.name)
  if (newName && newName.trim() && newName.trim() !== collection.name) {
    try {
      await CollectionService.UpdateCollection(collectionId, {
        ...collection,
        name: newName.trim()
      })
      await loadCollections()
    } catch (error) {
      console.error('Failed to rename collection:', error)
      alert('Failed to rename collection')
    }
  }
  showCollectionMenu.value = false
}

const duplicateCollection = async (collectionId) => {
  const collection = collections.value.find(c => c.id === collectionId)
  if (!collection) return
  
  try {
    await CollectionService.CreateCollection(
      `${collection.name} (Copy)`,
      collection.description
    )
    await loadCollections()
  } catch (error) {
    console.error('Failed to duplicate collection:', error)
    alert('Failed to duplicate collection')
  }
  showCollectionMenu.value = false
}

const exportCollection = async (collectionId) => {
  const collection = collections.value.find(c => c.id === collectionId)
  if (!collection) return
  
  try {
    const exportData = {
      collection: collection,
      exportedAt: new Date().toISOString(),
      version: '1.0'
    }
    
    const blob = new Blob([JSON.stringify(exportData, null, 2)], {
      type: 'application/json'
    })
    
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${collection.name.replace(/[^a-z0-9]/gi, '_').toLowerCase()}_collection.json`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  } catch (error) {
    console.error('Failed to export collection:', error)
    alert('Failed to export collection')
  }
  showCollectionMenu.value = false
}

const confirmDeleteCollection = (collectionId) => {
  const collection = collections.value.find(c => c.id === collectionId)
  if (!collection) return
  
  if (confirm(`Are you sure you want to delete "${collection.name}"? This will also delete all requests in this collection. This action cannot be undone.`)) {
    deleteCollection(collectionId)
  }
  showCollectionMenu.value = false
}

const deleteCollection = async (collectionId) => {
  try {
    await CollectionService.DeleteCollection(collectionId)
    await loadCollections()
  } catch (error) {
    console.error('Failed to delete collection:', error)
    alert('Failed to delete collection')
  }
}

// Expose refresh method
const refresh = () => {
  loadCollections()
}

defineExpose({ 
  refresh,
  requestLogsRef
})
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
  gap: 12px;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.sidebar-header h3 {
  margin: 0;
  color: #333;
  font-size: 18px;
}

.new-collection-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  border: none;
  background: #007bff;
  color: white;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  transition: all 0.2s;
  box-shadow: 0 1px 3px rgba(0, 123, 255, 0.3);
}

.new-collection-btn:hover {
  background: #0056b3;
  transform: translateY(-1px);
  box-shadow: 0 2px 6px rgba(0, 123, 255, 0.4);
}

.new-collection-btn:disabled {
  background: #6c757d;
  color: #adb5bd;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
  opacity: 0;
}

.new-collection-btn:disabled:hover {
  background: #6c757d;
  transform: none;
  box-shadow: none;
}

.btn-icon {
  font-size: 14px;
  font-weight: bold;
}

.btn-text {
  font-size: 12px;
}

.toggle-icon {
  margin-right: 6px;
  font-size: 14px;
}

.collections-list {
  flex: 1;
  overflow-y: auto;
}

.collection-item {
  margin-bottom: 12px;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  background: white;
  overflow: hidden;
  transition: all 0.2s;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.collection-item:hover {
  border-color: #007bff;
  box-shadow: 0 2px 8px rgba(0, 123, 255, 0.1);
}

.collection-item.expanded {
  border-color: #007bff;
  box-shadow: 0 2px 12px rgba(0, 123, 255, 0.15);
}

.collection-header {
  display: flex;
  align-items: center;
  padding: 14px 16px;
  background: linear-gradient(135deg, #f8f9fa 0%, #ffffff 100%);
}

.collection-main {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  cursor: pointer;
  transition: all 0.2s;
}

.collection-main:hover {
  color: #007bff;
}

.collection-icon {
  font-size: 18px;
  color: #6c757d;
  transition: transform 0.2s;
}

.collection-item.expanded .collection-icon {
  transform: scale(1.1);
  color: #007bff;
}

.collection-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.collection-name {
  font-weight: 600;
  color: #333;
  font-size: 14px;
}

.collection-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 11px;
  color: #6c757d;
}

.request-count {
  font-weight: 500;
}

.active-env {
  color: #007bff;
  font-weight: 500;
}

.collection-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s;
}

.collection-item:hover .collection-actions {
  opacity: 1;
}

.action-btn {
  width: 28px;
  height: 28px;
  border: none;
  background: transparent;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  transition: all 0.2s;
}

.action-btn:hover {
  background: rgba(0, 123, 255, 0.1);
  transform: scale(1.1);
}

.env-btn {
  color: #28a745;
}

.menu-btn {
  color: #6c757d;
  font-weight: bold;
}

/* Environment Section Styles */
.environment-section {
  padding: 16px;
  background: #f8f9fa;
  border-top: 1px solid #e9ecef;
}

.environment-card {
  background: white;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.env-selector-wrapper {
  margin-bottom: 12px;
}

.environment-label {
  display: block;
  font-size: 12px;
  font-weight: 600;
  color: #495057;
  margin-bottom: 6px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.custom-select {
  position: relative;
  display: flex;
  align-items: center;
  flex: 1;
}

.environment-select {
  flex: 1;
  padding: 8px 32px 8px 12px;
  border: 1px solid #ced4da;
  border-radius: 6px;
  font-size: 13px;
  background: white;
  cursor: pointer;
  transition: all 0.2s;
  appearance: none;
}

.environment-select:focus {
  outline: none;
  border-color: #007bff;
  box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
}

.select-arrow {
  position: absolute;
  right: 8px;
  top: 50%;
  transform: translateY(-50%);
  pointer-events: none;
  font-size: 8px;
  color: #6c757d;
}

.environment-details {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #e9ecef;
}

.env-detail-item {
  display: flex;
  margin-bottom: 6px;
  font-size: 12px;
}

.detail-label {
  font-weight: 600;
  color: #6c757d;
  min-width: 70px;
  margin-right: 8px;
}

.detail-value {
  color: #495057;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  word-break: break-all;
}

.no-environment {
  text-align: center;
  padding: 20px;
  color: #6c757d;
}

.no-env-text {
  display: block;
  margin-bottom: 12px;
  font-size: 13px;
}

.create-env-btn {
  padding: 8px 16px;
  border: 2px dashed #007bff;
  background: transparent;
  color: #007bff;
  border-radius: 6px;
  cursor: pointer;
  font-size: 12px;
  font-weight: 500;
  transition: all 0.2s;
}

.create-env-btn:hover {
  background: rgba(0, 123, 255, 0.1);
  transform: translateY(-1px);
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

.empty-collection {
  text-align: center;
  color: #6c757d;
  font-style: italic;
  padding: 20px;
}

/* Empty State Styles */
.empty-collections {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 300px;
  padding: 40px 20px;
}

.empty-state {
  text-align: center;
  max-width: 280px;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.6;
}

.empty-state h3 {
  margin: 0 0 8px 0;
  color: #495057;
  font-size: 18px;
  font-weight: 600;
}

.empty-state p {
  margin: 0 0 20px 0;
  color: #6c757d;
  font-size: 14px;
  line-height: 1.5;
}

.empty-action-btn {
  padding: 12px 24px;
  border: none;
  background: linear-gradient(135deg, #007bff 0%, #0056b3 100%);
  color: white;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s;
  box-shadow: 0 2px 8px rgba(0, 123, 255, 0.3);
}

.empty-action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 123, 255, 0.4);
}

/* Collection Context Menu Styles */
.collection-menu {
  position: fixed;
  background: white;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  z-index: 1000;
  min-width: 160px;
  overflow: hidden;
  animation: menuFadeIn 0.15s ease-out;
}

@keyframes menuFadeIn {
  from {
    opacity: 0;
    transform: scale(0.95) translateY(-5px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  cursor: pointer;
  font-size: 13px;
  color: #495057;
  transition: all 0.2s;
  border: none;
  background: none;
  width: 100%;
  text-align: left;
}

.menu-item:hover {
  background: #f8f9fa;
  color: #007bff;
}

.menu-item.danger {
  color: #dc3545;
}

.menu-item.danger:hover {
  background: #f8d7da;
  color: #721c24;
}

.menu-icon {
  font-size: 14px;
  width: 16px;
  text-align: center;
}

.menu-divider {
  height: 1px;
  background: #e9ecef;
  margin: 4px 0;
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

/* Environment Management Styles */
.environment-section {
  margin: 10px 0;
  padding: 10px;
  background: #f8f9fa;
  border-radius: 6px;
  border-left: 3px solid #007bff;
}

.environment-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 5px;
}

.environment-label {
  font-size: 12px;
  font-weight: 500;
  color: #6c757d;
  min-width: 80px;
}

.environment-select {
  flex: 1;
  padding: 4px 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 12px;
  background: white;
}

.env-manage-btn {
  padding: 4px 8px;
  border: none;
  background: #007bff;
  color: white;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.env-manage-btn:hover {
  background: #0056b3;
}

.environment-info {
  margin-top: 5px;
}

.env-url {
  font-size: 11px;
  color: #6c757d;
  font-family: monospace;
}

/* Environment Modal Styles */
.environment-modal {
  width: 700px;
  max-height: 85vh;
  overflow-y: auto;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
}

.environments-list {
  margin-bottom: 20px;
}

.environment-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 18px;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  margin-bottom: 12px;
  background: white;
  transition: all 0.2s;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.environment-item:hover {
  border-color: #007bff;
  box-shadow: 0 2px 8px rgba(0, 123, 255, 0.1);
  transform: translateY(-1px);
}

.environment-item.active {
  border-color: #007bff;
  background: linear-gradient(135deg, #f0f8ff 0%, #ffffff 100%);
  box-shadow: 0 2px 12px rgba(0, 123, 255, 0.15);
}

.env-info {
  flex: 1;
}

.env-name {
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
}

.env-url {
  font-family: monospace;
  color: #007bff;
  font-size: 12px;
  margin-bottom: 4px;
}

.env-description {
  color: #6c757d;
  font-size: 12px;
}

.env-actions {
  display: flex;
  gap: 8px;
}

.activate-btn, .edit-btn, .delete-env-btn {
  padding: 8px 14px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 12px;
  font-weight: 500;
  transition: all 0.2s;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.activate-btn {
  background: linear-gradient(135deg, #28a745 0%, #20c997 100%);
  color: white;
}

.activate-btn:hover {
  background: linear-gradient(135deg, #218838 0%, #1ea085 100%);
  transform: translateY(-1px);
  box-shadow: 0 2px 6px rgba(40, 167, 69, 0.3);
}

.edit-btn {
  background: linear-gradient(135deg, #ffc107 0%, #fd7e14 100%);
  color: #212529;
}

.edit-btn:hover {
  background: linear-gradient(135deg, #e0a800 0%, #e8590c 100%);
  transform: translateY(-1px);
  box-shadow: 0 2px 6px rgba(255, 193, 7, 0.3);
}

.delete-env-btn {
  background: linear-gradient(135deg, #dc3545 0%, #e74c3c 100%);
  color: white;
}

.delete-env-btn:hover {
  background: linear-gradient(135deg, #c82333 0%, #c0392b 100%);
  transform: translateY(-1px);
  box-shadow: 0 2px 6px rgba(220, 53, 69, 0.3);
}

.add-environment {
  text-align: center;
  margin-bottom: 20px;
}

.add-env-btn {
  padding: 12px 24px;
  border: 2px dashed #007bff;
  background: transparent;
  color: #007bff;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s;
  width: 100%;
}

.add-env-btn:hover {
  background: rgba(0, 123, 255, 0.1);
  border-style: solid;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 123, 255, 0.2);
}

.env-form {
  border-top: 1px solid #ddd;
  padding-top: 20px;
}

.env-form h5 {
  margin: 0 0 15px 0;
  color: #333;
}

.form-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  margin-top: 15px;
}

.save-btn {
  padding: 10px 20px;
  background: linear-gradient(135deg, #007bff 0%, #0056b3 100%);
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  transition: all 0.2s;
  box-shadow: 0 2px 6px rgba(0, 123, 255, 0.3);
}

.save-btn:hover {
  background: linear-gradient(135deg, #0056b3 0%, #004085 100%);
  transform: translateY(-1px);
  box-shadow: 0 3px 10px rgba(0, 123, 255, 0.4);
}

.save-btn:disabled {
  background: #6c757d;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

/* View Toggle Styles */
.view-toggle {
  display: flex;
  background: #f8f9fa;
  border-radius: 6px;
  padding: 2px;
  gap: 2px;
}

.toggle-btn {
  padding: 8px 16px;
  border: none;
  background: transparent;
  color: #6c757d;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s;
  flex: 1;
  text-align: center;
}

.toggle-btn:hover {
  color: #495057;
  background: rgba(0, 123, 255, 0.1);
}

/* Header Management Styles */
.headers-section {
  border: 1px solid #e9ecef;
  border-radius: 6px;
  padding: 15px;
  background: #f8f9fa;
}

.header-row {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
  align-items: center;
}

.header-input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 13px;
  background: white;
  transition: border-color 0.2s;
}

.header-select {
  width: 100%;
  padding: 6px 24px 6px 10px;
  border: 1px solid #ced4da;
  border-radius: 4px;
  background-color: white;
  font-size: 13px;
  color: #495057;
  appearance: none;
  cursor: pointer;
}

.header-input:focus {
  outline: none;
  border-color: #007bff;
  box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.1);
}

.remove-header-btn {
  padding: 6px 10px;
  background: #dc3545;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  font-weight: bold;
  transition: all 0.2s;
  min-width: 32px;
}

.remove-header-btn:hover:not(:disabled) {
  background: #c82333;
  transform: translateY(-1px);
}

.remove-header-btn:disabled {
  background: #6c757d;
  cursor: not-allowed;
  opacity: 0.5;
}

.add-header-btn {
  padding: 8px 16px;
  background: #28a745;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  transition: all 0.2s;
  margin-top: 5px;
}

.add-header-btn:hover {
  background: #218838;
  transform: translateY(-1px);
}

.toggle-btn.active {
  background: #007bff;
  color: white;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.toggle-btn.active:hover {
  background: #0056b3;
  color: white;
}

/* Header Collection Styles */
.header-collection-section {
  margin-bottom: 20px;
  padding: 16px;
  background: linear-gradient(135deg, #f8f9fa 0%, #ffffff 100%);
  border-radius: 12px;
  border: 1px solid #e3e6ea;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  transition: all 0.2s ease;
}





/* Collection Settings Styles */
.collection-settings {
  margin-bottom: 12px;
  padding: 0 12px;
}

.settings-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 12px;
}

.setting-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.setting-label {
  font-size: 11px;
  font-weight: 600;
  color: #6c757d;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.setting-select {
  width: 100%;
  padding: 6px 8px;
  border: 1px solid #ced4da;
  border-radius: 4px;
  background: white;
  font-size: 12px;
  color: #495057;
  cursor: pointer;
  transition: all 0.2s ease;
}

.setting-select:hover {
  border-color: #adb5bd;
}

.setting-select:focus {
  outline: none;
  border-color: #80bdff;
  box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
}

/* Header Settings Modal Styles */
.header-settings-modal {
  width: 90vw;
  max-width: 1200px;
  height: 90vh;
  max-height: 800px;
  padding: 0;
  overflow: hidden;
}

.header-settings-modal .modal-content {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.empty-icon {
  width: 24px;
  height: 24px;
  color: #ced4da;
}
</style>