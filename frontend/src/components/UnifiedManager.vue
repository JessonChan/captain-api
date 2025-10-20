<template>
  <div class="unified-manager-overlay" @click="handleOverlayClick">
    <div class="unified-manager" @click.stop>
      <!-- Close Button -->
      <button @click.stop="$emit('close')" class="manager-close-btn" title="Close Manager">
        <span class="close-icon">×</span>
      </button>
      
      <!-- Manager Tabs -->
    <div class="manager-tabs">
      <button 
        v-for="tab in tabs" 
        :key="tab.id"
        @click="activeTab = tab.id"
        :class="['tab-btn', { active: activeTab === tab.id }]"
      >
        <span class="tab-icon">{{ tab.icon }}</span>
        <span class="tab-text">{{ tab.label }}</span>
      </button>
    </div>

    <!-- Tab Content -->
    <div class="tab-content">
      <!-- Environment Management -->
      <div v-if="activeTab === 'environments'" class="environments-tab">
        <div class="section-header">
          <h4>Environments</h4>
          <button @click="showEnvironmentForm = true" class="add-btn">
            <span class="btn-icon">+</span>
            Add Environment
          </button>
        </div>

        <div class="items-list">
          <div v-if="environments.length === 0" class="empty-state">
            <div class="empty-icon">🌐</div>
            <div class="empty-message">No environments yet</div>
            <div class="empty-description">Create your first environment to manage different API endpoints</div>
          </div>

          <div 
            v-for="env in environments" 
            :key="env.id"
            class="item-card"
            :class="{ active: env.isActive }"
          >
            <div class="item-info">
              <div class="item-name">{{ env.name }}</div>
              <div class="item-description">{{ env.description }}</div>
              <div class="item-details">
                <span class="url-badge">{{ env.baseUrl }}</span>
                <span v-if="env.isActive" class="active-badge">Active</span>
              </div>
            </div>
            <div class="item-actions">
              <button 
                v-if="!env.isActive"
                @click="activateEnvironment(env.id)"
                class="action-btn activate"
                title="Activate environment"
              >
                <span class="btn-icon">✓</span>
              </button>
              <button 
                @click="editEnvironment(env)"
                class="action-btn edit"
                title="Edit environment"
              >
                <span class="btn-icon">✏️</span>
              </button>
              <button 
                v-if="environments.length > 1"
                @click="deleteEnvironment(env)"
                class="action-btn delete"
                title="Delete environment"
              >
                <span class="btn-icon">🗑️</span>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Header Collections Management -->
      <div v-if="activeTab === 'headers'" class="headers-tab">
        <div class="section-header">
          <h4>Header Collections</h4>
          <div class="header-actions">
            <button @click="triggerImport" class="import-btn">
              <span class="btn-icon">📥</span>
              Import
            </button>
            <button @click="showHeaderForm = true" class="add-btn">
              <span class="btn-icon">+</span>
              Add Collection
            </button>
          </div>
          <input 
            ref="fileInput"
            type="file" 
            @change="importHeaderCollection" 
            accept=".json" 
            style="display: none;"
          />
        </div>

        <div class="items-list">
          <div v-if="headerCollections.length === 0" class="empty-state">
            <div class="empty-icon">📋</div>
            <div class="empty-message">No header collections yet</div>
            <div class="empty-description">Create your first header collection to reuse headers across requests</div>
          </div>

          <div 
            v-for="collection in headerCollections" 
            :key="collection.id"
            class="item-card"
          >
            <div class="item-info">
              <div class="item-name">{{ collection.name }}</div>
              <div class="item-description">{{ collection.description }}</div>
              <div class="item-details">
                <span class="header-count-badge">
                  {{ Object.keys(collection.headers || {}).length }} headers
                </span>
                <span v-if="isActiveHeaderCollection(collection.id)" class="active-badge">Active</span>
              </div>
            </div>
            <div class="item-actions">
              <button 
                @click="setActiveHeaderCollection(collection.id)"
                class="action-btn activate"
                title="Set as active"
              >
                <span class="btn-icon">✓</span>
              </button>
              <button 
                @click="editHeaderCollection(collection)"
                class="action-btn edit"
                title="Edit collection"
              >
                <span class="btn-icon">✏️</span>
              </button>
              <button 
                @click="exportHeaderCollection(collection)"
                class="action-btn export"
                title="Export collection"
              >
                <span class="btn-icon">📤</span>
              </button>
              <button 
                @click="deleteHeaderCollection(collection)"
                class="action-btn delete"
                title="Delete collection"
              >
                <span class="btn-icon">🗑️</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Popup Forms -->
  <EnvironmentFormPopup
    :is-visible="showEnvironmentForm"
    :environment="editingEnvironment ?? undefined"
    @submit="saveEnvironment"
    @close="cancelEnvironmentForm"
  />
  
  <HeaderCollectionFormPopup
    :is-visible="showHeaderForm"
    :header-collection="editingHeaderCollection ?? undefined"
    @submit="saveHeaderCollection"
    @close="cancelHeaderForm"
  />
</div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { CollectionService } from '../../bindings/captain-api'
import { EventBusService } from '../../bindings/captain-api'
import popupService from './popupService'
import EnvironmentFormPopup from './EnvironmentFormPopup.vue'
import HeaderCollectionFormPopup from './HeaderCollectionFormPopup.vue'
import { HeaderCollection as ModelHeaderCollection } from '../../bindings/captain-api/models'

// Local lightweight types to avoid depending on generated model types at build time
type CollectionEnvironment = {
  id: string
  name: string
  baseUrl: string
  description?: string
  isActive: boolean
}

type HeaderCollection = {
  id: string
  name: string
  description?: string
  headers: Record<string, string>
  createdAt?: any
  updatedAt?: any
}

type RequestItem = {
  id: string
  name: string
  method: string
  url: string
}

type Collection = {
  id: string
  name: string
  description?: string
  environments: CollectionEnvironment[]
  headerCollections: HeaderCollection[]
  activeHeaderCollectionId?: string
  requests: RequestItem[]
}

// Props
const props = defineProps({
  collectionId: {
    type: String,
    required: true
  }
})

// Emits
const emit = defineEmits(['close', 'updated'])

// State
const activeTab = ref('environments')
const environments = ref<CollectionEnvironment[]>([])
const headerCollections = ref<HeaderCollection[]>([])
const collection = ref<Collection | null>(null)

// Form states
const showEnvironmentForm = ref(false)
const showHeaderForm = ref(false)
const editingEnvironment = ref<CollectionEnvironment | null>(null)
const editingHeaderCollection = ref<HeaderCollection | null>(null)


// File input ref
const fileInput = ref<HTMLInputElement | null>(null)

// Tabs configuration
const tabs = [
  { id: 'environments', label: 'Environments', icon: '🌐' },
  { id: 'headers', label: 'Headers', icon: '📋' }
]

// Load collection data
const loadCollectionData = async () => {
  try {
    // Load collection
    const allCollections = await CollectionService.GetAllCollections()
    collection.value = (allCollections as any[]).find((c: any) => c.id === props.collectionId) as Collection | null
    
    if (collection.value) {
      // Load environments
      const envs = await CollectionService.GetCollectionEnvironments(props.collectionId)
      environments.value = (envs as any[]) as CollectionEnvironment[]
      
      // Load header collections
      headerCollections.value = (collection.value.headerCollections || []) as HeaderCollection[]
    }
  } catch (error) {
    console.error('Failed to load collection data:', error)
    await popupService.alert('Failed to load collection data', {
      severity: 'error'
    })
  }
}

// Environment management
const activateEnvironment = async (envId) => {
  try {
    await CollectionService.SetActiveCollectionEnvironment(props.collectionId, envId)
    await loadCollectionData()
    emit('updated')
  } catch (error) {
    console.error('Failed to activate environment:', error)
    await popupService.alert('Failed to activate environment', {
      severity: 'error'
    })
  }
}

const editEnvironment = (env) => {
  editingEnvironment.value = env
  showEnvironmentForm.value = true
}

const saveEnvironment = async (envData) => {
  try {
    if (editingEnvironment.value) {
      await CollectionService.UpdateCollectionEnvironment(props.collectionId, envData)
    } else {
      await CollectionService.AddCollectionEnvironment(props.collectionId, envData)
    }

    await loadCollectionData()
    cancelEnvironmentForm()
    emit('updated')
  } catch (error) {
    console.error('Failed to save environment:', error)
    await popupService.alert('Failed to save environment', {
      severity: 'error'
    })
  }
}

const deleteEnvironment = async (env) => {
  try {
    const confirmed = await popupService.confirm(
      `Are you sure you want to delete the "${env.name}" environment?`,
      {
        details: 'This action cannot be undone.',
        title: 'Delete Environment'
      }
    )

    if (confirmed) {
      await CollectionService.DeleteCollectionEnvironment(props.collectionId, env.id)
      await loadCollectionData()
      emit('updated')
    }
  } catch (error) {
    console.error('Failed to delete environment:', error)
    await popupService.alert('Failed to delete environment', {
      severity: 'error'
    })
  }
}

const cancelEnvironmentForm = () => {
  showEnvironmentForm.value = false
  editingEnvironment.value = null
}

// Header collection management
const isActiveHeaderCollection = (collectionId: string) => {
  return collection.value?.activeHeaderCollectionId === collectionId
}

const setActiveHeaderCollection = async (collectionId: string) => {
  try {
    await CollectionService.SetActiveHeaderCollection(props.collectionId, collectionId)
    await loadCollectionData()
    emit('updated')
  } catch (error) {
    console.error('Failed to set active header collection:', error)
    await popupService.alert('Failed to set active header collection', {
      severity: 'error'
    })
  }
}

const editHeaderCollection = (collection: HeaderCollection) => {
  editingHeaderCollection.value = collection
  showHeaderForm.value = true
}

const saveHeaderCollection = async (collectionData: HeaderCollection & { headers: { key?: string; value?: string }[] }) => {
  try {
    // Convert headers array to object
    const headersObject = {}
    ;(collectionData.headers as any[]).forEach((header: any) => {
      if (header.key) headersObject[header.key] = header.value
    })

    const payload = ModelHeaderCollection.createFrom({
      id: collectionData.id || `header_col_${Date.now()}`,
      collectionId: props.collectionId,
      name: collectionData.name,
      description: (collectionData.description || '').toString(),
      headers: headersObject,
      createdAt: editingHeaderCollection.value?.createdAt ?? null,
      updatedAt: new Date().toISOString()
    })

    if (editingHeaderCollection.value) {
      await CollectionService.UpdateHeaderCollection(props.collectionId, payload)
    } else {
      await CollectionService.CreateHeaderCollection(props.collectionId, payload)
    }

    await loadCollectionData()
    cancelHeaderForm()
    emit('updated')
  } catch (error) {
    console.error('Failed to save header collection:', error)
    await popupService.alert('Failed to save header collection', {
      severity: 'error'
    })
  }
}

const deleteHeaderCollection = async (collection: HeaderCollection) => {
  try {
    const confirmed = await popupService.confirm(
      `Are you sure you want to delete the "${collection.name}" header collection?`,
      {
        details: 'This action cannot be undone.',
        title: 'Delete Header Collection'
      }
    )

    if (confirmed) {
      await CollectionService.DeleteHeaderCollection(props.collectionId, collection.id)
      await loadCollectionData()
      emit('updated')
    }
  } catch (error) {
    console.error('Failed to delete header collection:', error)
    await popupService.alert('Failed to delete header collection', {
      severity: 'error'
    })
  }
}

const exportHeaderCollection = async (collection: HeaderCollection) => {
  try {
    const jsonData = JSON.stringify(collection, null, 2)
    const blob = new Blob([jsonData], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    
    const a = document.createElement('a')
    a.href = url
    a.download = `${collection.name.replace(/\s+/g, '_')}_headers.json`
    
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
    
    await popupService.alert(`Header collection "${collection.name}" exported successfully`, {
      severity: 'success'
    })
  } catch (error) {
    console.error('Failed to export header collection:', error)
    await popupService.alert('Failed to export header collection', {
      severity: 'error'
    })
  }
}

const triggerImport = () => {
  fileInput.value?.click()
}

const importHeaderCollection = async (event: Event) => {
  const target = event.target as HTMLInputElement | null
  const file = target?.files?.[0]
  if (!file) return
  
  const reader = new FileReader()
  reader.onload = async (e) => {
    try {
      const result = e.target ? (e.target as FileReader).result : null
      const importedData = JSON.parse((result as string) || '{}')
      
      if (!importedData.name || !importedData.headers) {
        throw new Error('Invalid header collection format')
      }
      
      const payload = ModelHeaderCollection.createFrom({
        id: 'header_col_' + Date.now(),
        collectionId: props.collectionId,
        name: importedData.name,
        description: (importedData.description || '').toString(),
        headers: importedData.headers,
        createdAt: null,
        updatedAt: new Date().toISOString()
      })
      
      // Check for duplicate name
      const existingNames = headerCollections.value.map(c => (c.name || '').toLowerCase())
      if (existingNames.includes((payload.name || '').toLowerCase())) {
        payload.name = `${payload.name} (Imported)`
      }
      
      await CollectionService.CreateHeaderCollection(props.collectionId, payload)
      await loadCollectionData()
      emit('updated')
      
      await popupService.alert(`Successfully imported "${payload.name}" collection`, {
        severity: 'success'
      })
    } catch (error) {
      console.error('Failed to import header collection:', error)
      await popupService.alert('Failed to import header collection: Invalid format', {
        severity: 'error'
      })
    }
    
    if (target) target.value = ''
  }
  
  reader.onerror = () => {
    popupService.alert('Error reading file', {
      severity: 'error'
    })
    if (target) target.value = ''
  }
  
  reader.readAsText(file)
}

const cancelHeaderForm = () => {
  showHeaderForm.value = false
  editingHeaderCollection.value = null
}


// Methods
const handleOverlayClick = () => {
  emit('close')
}

// Lifecycle
onMounted(() => {
  loadCollectionData()
})

watch(() => props.collectionId, () => {
  loadCollectionData()
})
</script>

<style scoped>
.unified-manager-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 2000;
  backdrop-filter: blur(4px);
  animation: fadeIn 0.3s ease-out;
}

.unified-manager {
  display: flex;
  flex-direction: column;
  height: 90vh;
  width: 90vw;
  max-width: 1200px;
  background-color: #fff;
  border-radius: 12px;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.3);
  animation: slideUp 0.4s ease-out;
  position: relative;
  overflow: hidden;
}

.manager-close-btn {
  position: absolute;
  top: 15px;
  right: 15px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid #e9ecef;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 10;
  transition: all 0.3s ease;
  font-size: 20px;
  color: #495057;
  backdrop-filter: blur(10px);
}

.manager-close-btn:hover {
  background: rgba(220, 53, 69, 0.1);
  border-color: #dc3545;
  color: #dc3545;
  transform: scale(1.1);
}

.close-icon {
  line-height: 1;
  font-weight: 300;
}

.manager-tabs {
  display: flex;
  border-bottom: 1px solid #e9ecef;
  background-color: #f8f9fa;
  border-radius: 8px 8px 0 0;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  border: none;
  background: none;
  color: #6c757d;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all 0.2s;
  border-bottom: 2px solid transparent;
}

.tab-btn:hover {
  color: #495057;
  background-color: rgba(0, 123, 255, 0.05);
}

.tab-btn.active {
  color: #007bff;
  background-color: white;
  border-bottom-color: #007bff;
}

.tab-icon {
  font-size: 1rem;
}

.tab-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-header h4 {
  margin: 0;
  color: #343a40;
  font-size: 1.1rem;
  font-weight: 600;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.add-btn, .import-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.85rem;
  font-weight: 500;
  transition: all 0.2s;
}

.add-btn {
  background-color: #0d6efd;
  color: white;
}

.add-btn:hover {
  background-color: #0b5ed7;
  transform: translateY(-1px);
}

.import-btn {
  background-color: #17a2b8;
  color: white;
}

.import-btn:hover {
  background-color: #138496;
  transform: translateY(-1px);
}

.btn-icon {
  font-size: 0.9rem;
}

.items-list {
  margin-bottom: 20px;
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: #6c757d;
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 12px;
}

.empty-message {
  font-size: 1.1rem;
  font-weight: 600;
  margin-bottom: 8px;
}

.empty-description {
  font-size: 0.9rem;
  max-width: 300px;
  margin: 0 auto;
  line-height: 1.5;
}

.item-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  margin-bottom: 12px;
  background: white;
  transition: all 0.2s;
}

.item-card:hover {
  border-color: #007bff;
  box-shadow: 0 2px 8px rgba(0, 123, 255, 0.1);
}

.item-card.active {
  border-color: #28a745;
  background: linear-gradient(135deg, #f8fff9 0%, #ffffff 100%);
}

.item-info {
  flex: 1;
}

.item-name {
  font-weight: 600;
  color: #212529;
  margin-bottom: 4px;
}

.item-description {
  font-size: 0.9rem;
  color: #6c757d;
  margin-bottom: 8px;
}

.item-details {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.url-badge, .header-count-badge {
  background-color: #e9ecef;
  color: #495057;
  font-size: 0.8rem;
  padding: 2px 8px;
  border-radius: 12px;
  font-weight: 500;
}

.active-badge {
  background-color: #28a745;
  color: white;
  font-size: 0.8rem;
  padding: 2px 8px;
  border-radius: 12px;
  font-weight: 500;
}

.item-actions {
  display: flex;
  gap: 6px;
}

.action-btn {
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.9rem;
  transition: all 0.2s;
}

.action-btn.activate {
  background-color: #28a745;
  color: white;
}

.action-btn.activate:hover {
  background-color: #218838;
  transform: translateY(-1px);
}

.action-btn.edit {
  background-color: #ffc107;
  color: #212529;
}

.action-btn.edit:hover {
  background-color: #e0a800;
  transform: translateY(-1px);
}

.action-btn.export {
  background-color: #17a2b8;
  color: white;
}

.action-btn.export:hover {
  background-color: #138496;
  transform: translateY(-1px);
}

.action-btn.delete {
  background-color: #dc3545;
  color: white;
}

.action-btn.delete:hover {
  background-color: #c82333;
  transform: translateY(-1px);
}

.form-section {
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 20px;
  background-color: #f8f9fa;
  margin-top: 20px;
}

.form-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.form-header h5 {
  margin: 0;
  color: #343a40;
  font-size: 1rem;
  font-weight: 600;
}

.close-form-btn {
  width: 28px;
  height: 28px;
  border: none;
  background: none;
  font-size: 1.2rem;
  cursor: pointer;
  color: #6c757d;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.close-form-btn:hover {
  background-color: #e9ecef;
  color: #343a40;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-weight: 500;
  color: #495057;
  margin-bottom: 6px;
  font-size: 0.9rem;
}

.form-input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ced4da;
  border-radius: 6px;
  font-size: 0.9rem;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.form-input:focus {
  outline: none;
  border-color: #86b7fe;
  box-shadow: 0 0 0 3px rgba(13, 110, 253, 0.1);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
}

.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all 0.2s;
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

.headers-editor {
  border: 1px solid #e9ecef;
  border-radius: 6px;
  overflow: hidden;
}

.headers-header {
  display: grid;
  grid-template-columns: 1fr 1fr 40px;
  gap: 10px;
  padding: 12px 15px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
  font-weight: 600;
  font-size: 0.85rem;
  color: #495057;
}

.headers-list {
  max-height: 200px;
  overflow-y: auto;
}

.header-row {
  display: grid;
  grid-template-columns: 1fr 1fr 40px;
  gap: 10px;
  padding: 10px 15px;
  border-bottom: 1px solid #f1f3f5;
  align-items: center;
}

.header-row:last-child {
  border-bottom: none;
}

.header-input {
  width: 100%;
  padding: 6px 10px;
  border: 1px solid #ced4da;
  border-radius: 4px;
  font-size: 0.85rem;
}

.header-input:focus {
  outline: none;
  border-color: #86b7fe;
  box-shadow: 0 0 0 2px rgba(13, 110, 253, 0.1);
}

.remove-header-btn {
  width: 24px;
  height: 24px;
  border: none;
  background-color: #dc3545;
  color: white;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.8rem;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.remove-header-btn:hover {
  background-color: #c82333;
  transform: translateY(-1px);
}

.add-header-btn {
  width: 100%;
  padding: 10px;
  border: 1px dashed #ced4da;
  background-color: #f8f9fa;
  color: #495057;
  border-radius: 0 0 6px 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-size: 0.85rem;
  font-weight: 500;
  transition: all 0.2s;
}

.add-header-btn:hover {
  background-color: #e9ecef;
  border-color: #adb5bd;
}

/* Animations */
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideUp {
  from { 
    transform: translateY(30px); 
    opacity: 0; 
  }
  to { 
    transform: translateY(0); 
    opacity: 1; 
  }
}
</style>