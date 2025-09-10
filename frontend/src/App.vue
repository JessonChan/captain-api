<script setup lang="ts">
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import RequestBuilder from './components/RequestBuilder.vue'
import ResponseViewer from './components/ResponseViewer.vue'
import CollectionSidebar from './components/CollectionSidebar.vue'
import UnifiedManager from './components/UnifiedManager.vue'
import Welcome from './components/Welcome.vue'
import { Events } from '@wailsio/runtime'
import { CollectionService } from '../bindings/captain-api'
import { main } from '../bindings/captain-api/models'

// Tabs management
const tabs = ref([
  { id: 'welcome-tab', name: 'Welcome', active: true, isWelcome: true }
])
const activeTabId = ref('welcome-tab')
const tabCounter = ref(2)

// References
const requestBuilderRefs = ref({})
const collectionSidebarRef = ref(null)
const collections = ref<main.Collection[]>([])
const currentCollectionId = ref(null)

// Unified Manager state
const showUnifiedManager = ref(false)
const unifiedManagerCollectionId = ref(null)

// Computed properties
const activeTab = computed(() => tabs.value.find(tab => tab.id === activeTabId.value))
const activeResponse = computed(() => activeTab.value?.response)
const selectedHeaders = ref(null)

const loadCollections = () => {
  CollectionService.GetAllCollections().then(allCollections => {
    collections.value = allCollections
    if (allCollections.length > 0) {
      if (!currentCollectionId.value || !allCollections.some(c => c.id === currentCollectionId.value)) {
        currentCollectionId.value = allCollections[0].id
      }
      updateSelectedHeaders()
    }
  })
}

const updateSelectedHeaders = () => {
  selectedHeaders.value = null
  const collection = collections.value.find((c: main.Collection) => c.id === currentCollectionId.value)
  if (collection) {
    const activeHeaderCollection = collection.headerCollections.find((h: main.HeaderCollection) => h.id === collection.activeHeaderCollectionId)
    if (activeHeaderCollection) {
      selectedHeaders.value = activeHeaderCollection.headers
    }
  }
}

onMounted(() => {
  loadCollections()
  Events.On('header-collection-updated', loadCollections)
  Events.On('collections-updated', loadCollections)
  Events.On('request-saved', (savedRequest) => {
    loadCollections()
    const active = activeTab.value
    if (active && !active.requestId) {
      active.requestId = savedRequest.id
      active.name = savedRequest.name
    }
  })
})

const handleCollectionSelected = (collectionId: string) => {
  console.log('Collection selected:', collectionId)
  currentCollectionId.value = collectionId
  updateSelectedHeaders()
}

const handleHeaderCollectionSelected = (headers: main.HeaderCollection) => {
  selectedHeaders.value = headers
}

const handleResponseReceived = async (responseData: main.ResponseData) => {
  // Update the response for the active tab
  const index = tabs.value.findIndex(tab => tab.id === activeTabId.value)
  if (index !== -1) {
    tabs.value[index].response = responseData
    
    // Update tab name with status code if available
    if (responseData && responseData.statusCode) {
      // Only append status code if it's not already in the name
      if (!tabs.value[index].name.includes(`[${responseData.statusCode}]`)) {
        tabs.value[index].name = `${tabs.value[index].name} [${responseData.statusCode}]`
      }
    }
  }
  
  // Refresh logs after a request is made
  if (collectionSidebarRef.value && collectionSidebarRef.value.requestLogsRef) {
    await collectionSidebarRef.value.requestLogsRef.loadLogs()
  }
}

const handleLoadRequest = (logData) => {
  // Handle both old format (direct request) and new format (with request and response)
  const requestData = logData.request || logData
  const requestUrl = requestData.url || ''
  const method = requestData.method || 'GET'
  
  // Update currentCollectionId if the request belongs to a different collection
  if (requestData.collectionId && requestData.collectionId !== currentCollectionId.value) {
    handleCollectionSelected(requestData.collectionId)
  }
  
  // Generate a meaningful tab name
  let tabName = ''
  if (requestData.name) {
    tabName = requestData.name
  } else {
    // Extract endpoint from URL
    const urlParts = requestUrl.split('/')
    const endpoint = urlParts.length > 0 ? urlParts[urlParts.length - 1] : ''
    
    if (endpoint) {
      tabName = `${method} ${endpoint}`
    } else if (requestUrl) {
      // Use the domain if endpoint is not available
      try {
        const url = new URL(requestUrl)
        tabName = `${method} ${url.hostname}`
      } catch {
        tabName = `${method} Request`
      }
    } else {
      tabName = `${method} Request`
    }
  }
  
  // Check if there's already a tab with this request ID
  const existingTabIndex = tabs.value.findIndex(tab => tab.requestId === requestData.id)
  
  if (existingTabIndex !== -1) {
    // If tab with this URL exists, switch to it
    switchTab(tabs.value[existingTabIndex].id)
    
    // Update the request if needed (might have different headers, body, etc.)
    const existingTabBuilder = requestBuilderRefs.value[tabs.value[existingTabIndex].id]
    if (existingTabBuilder) {
      existingTabBuilder.loadRequest(requestData, currentCollectionId.value)
      
      // If response data is available (e.g., from logs), set it
      if (logData.response) {
        tabs.value[existingTabIndex].response = logData.response
        
        // Update tab name with status code if available
        if (logData.response.statusCode) {
          tabs.value[existingTabIndex].name = `${tabName} [${logData.response.statusCode}]`
        } else {
          tabs.value[existingTabIndex].name = tabName
        }
      } else {
        tabs.value[existingTabIndex].name = tabName
      }
    }
  } else {
    // Create a new tab for this request
    addNewTab(requestData, logData.response, requestData.id)
  }
}

// Handle request updates from RequestBuilder
const handleRequestUpdated = (data) => {
  const index = tabs.value.findIndex(tab => tab.id === activeTabId.value)
  if (index !== -1) {
    // Store a deep copy of the request object in the tab
    tabs.value[index].request = JSON.parse(JSON.stringify(data.request))
    
    // Update tab name based on request data if no custom name is provided
    if (!data.name) {
      const urlParts = data.request.url?.split('/') || []
      const endpoint = urlParts.length > 0 ? urlParts[urlParts.length - 1] : ''
      const method = data.request.method || 'GET'
      
      // Create a meaningful tab name
      if (endpoint) {
        tabs.value[index].name = `${method} ${endpoint}`
      } else if (data.request.url) {
        // Use the domain if endpoint is not available
        try {
          const url = new URL(data.request.url)
          tabs.value[index].name = `${method} ${url.hostname}`
        } catch {
          tabs.value[index].name = `${method} Request`
        }
      } else {
        tabs.value[index].name = `${method} Request`
      }
      
      // If there's a status code in the tab name, preserve it
      const statusMatch = tabs.value[index].name.match(/\[(\d+)\]$/)
      if (statusMatch && tabs.value[index].response) {
        tabs.value[index].name = `${tabs.value[index].name} [${tabs.value[index].response.statusCode}]`
      }
    } else {
      // Use the custom name provided
      tabs.value[index].name = data.name
    }
  }
}

const handleNewRequest = () => {
  // Create a new tab
  addNewTab()
}


const addNewTab = (requestData = null, responseData = null, requestId = null) => {
  const newTabId = `tab-${tabCounter.value}`
  tabCounter.value++
  
  // Deactivate all tabs
  tabs.value.forEach(tab => tab.active = false)
  
  // Create tab name if request data is provided
  let tabName = 'New Request'
  if (requestData && requestData.url) {
    const urlParts = requestData.url.split('/') || []
    const endpoint = urlParts.length > 0 ? urlParts[urlParts.length - 1] : ''
    tabName = endpoint || `${requestData.method || 'GET'} Request`
  }
  
  const newTab = {
    id: newTabId,
    name: tabName,
    active: true,
    request: requestData,
    response: responseData,
    requestId: requestId
  }

  // If the only tab is the welcome tab, replace it
  if (tabs.value.length === 1 && tabs.value[0].isWelcome) {
    tabs.value = [newTab]
  } else {
    tabs.value.push(newTab)
  }
  
  // Set the new tab as active
  activeTabId.value = newTabId
  
  // Initialize the request builder in the next tick
  nextTick(() => {
    if (requestBuilderRefs.value[newTabId]) {
      if (requestData) {
        requestBuilderRefs.value[newTabId].loadRequest(requestData, currentCollectionId.value)
      } else {
        requestBuilderRefs.value[newTabId].newRequest()
      }
    }
  })
}

const switchTab = (tabId) => {
  // Deactivate all tabs
  tabs.value.forEach(tab => tab.active = false)
  
  // Activate the selected tab
  const index = tabs.value.findIndex(tab => tab.id === tabId)
  if (index !== -1) {
    tabs.value[index].active = true
    activeTabId.value = tabId
  }
}

const closeTab = (tabId, event) => {
  // Prevent the click from triggering the tab switch
  if (event) {
    event.stopPropagation()
  }
  
  const index = tabs.value.findIndex(tab => tab.id === tabId)
  if (index !== -1) {
    // If closing the active tab
    if (tabId === activeTabId.value) {
      // If it's the last request tab, replace it with the welcome tab
      if (tabs.value.length === 1) {
        tabs.value = [{ id: 'welcome-tab', name: 'Welcome', active: true, isWelcome: true }]
        activeTabId.value = 'welcome-tab'
        return
      }
      
      // Otherwise, activate the next tab or the previous one if it's the last
      const newActiveIndex = index === tabs.value.length - 1 ? index - 1 : index + 1
      tabs.value[newActiveIndex].active = true
      activeTabId.value = tabs.value[newActiveIndex].id
    }
    
    // Remove the tab
    tabs.value.splice(index, 1)
  }
}

const setRequestBuilderRef = (el, tabId) => {
  if (el) {
    requestBuilderRefs.value[tabId] = el
  }
}

// Unified Manager methods
const openUnifiedManager = (collectionId) => {
  unifiedManagerCollectionId.value = collectionId
  showUnifiedManager.value = true
}

const closeUnifiedManager = () => {
  showUnifiedManager.value = false
  unifiedManagerCollectionId.value = null
}

const handleManagerUpdated = () => {
  loadCollections()
  updateSelectedHeaders()
}
</script>

<template>
  <div class="app">
    <header class="app-header">
      <div class="header-content">
        <div class="logo-section">
          <img src="/logo.png" class="app-logo" alt="Captain API" />
          <h1 class="app-title">Captain API</h1>
        </div>
        <div class="header-subtitle">
          A Postman-like API client built with Wails3 + Vue
        </div>
      </div>
    </header>

    <main class="app-main">
      <div class="main-layout">
        <!-- Sidebar -->
        <aside class="sidebar">
          <CollectionSidebar
            v-if="currentCollectionId"
            ref="collectionSidebarRef"
            :collection-id="currentCollectionId"
            @load-request="handleLoadRequest"
            @new-request="handleNewRequest"
            @collection-selected="handleCollectionSelected"
            @header-collection-selected="handleHeaderCollectionSelected"
            @open-unified-manager="openUnifiedManager"
          />
        </aside>

        <!-- Main Content -->
        <div class="main-content">
          <!-- Tab Bar -->
          <div class="tab-bar">
            <div 
              v-for="tab in tabs" 
              :key="tab.id"
              class="tab"
              :class="{ active: tab.id === activeTabId }"
              @click="switchTab(tab.id)"
            >
              <span class="tab-name">{{ tab.name }}</span>
              <button v-if="!tab.isWelcome" class="tab-close" @click="closeTab(tab.id, $event)">&times;</button>
            </div>
            <button class="new-tab-button" @click="handleNewRequest">+</button>
          </div>
          
          <Welcome v-if="activeTab.isWelcome" @new-request="handleNewRequest" />
          <template v-else>
            <!-- Request Builders - One per tab, only show active one -->
            <div v-for="tab in tabs" :key="`builder-${tab.id}`" v-show="tab.id === activeTabId" class="request-section">
              <RequestBuilder
                :ref="el => setRequestBuilderRef(el, tab.id)"
                :collectionId="currentCollectionId"
                :virtualHeaders="selectedHeaders"
                @response-received="handleResponseReceived"
                @request-updated="handleRequestUpdated"
                @load-request="handleLoadRequest"
              />
            </div>
            
            <!-- Response Viewer - Shows active tab's response -->
            <section class="response-section">
              <ResponseViewer :response="activeResponse" />
            </section>
          </template>
        </div>
      </div>
    </main>
    
    <!-- Unified Manager Modal -->
    <UnifiedManager 
      v-if="showUnifiedManager" 
      :collection-id="unifiedManagerCollectionId" 
      @close="closeUnifiedManager" 
      @updated="handleManagerUpdated"
    />
  </div>
</template>

<style lang="scss">
// Variables
$primary-gradient: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
$background-color: #f5f6fa;
$max-width: 100vw;
$sidebar-width: 300px;
$border-radius: 10px;
$box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

// Font stack
$font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
  'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue', sans-serif;

// Global styles
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: $font-family;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background: $background-color;
}

.app {
  height: 100vh;
  display: flex;
  flex-direction: column;

  &-header {
    background: $primary-gradient;
    color: white;
    padding: 20px;
    box-shadow: $box-shadow;
  }

  &-logo {
    width: 40px;
    height: 40px;
    border-radius: $border-radius;
  }

  &-title {
    font-size: 28px;
    font-weight: 600;
    margin: 0;
  }

  &-main {
    flex: 1;
    overflow: hidden;
  }
}

.header {
  &-content {
    max-width: $max-width;
    margin: 0 auto;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 20px;
    height: 100%;
    flex-wrap: wrap;
  }

  &-subtitle {
    font-size: 16px;
    opacity: 0.9;
    font-weight: 300;
  }

  &-nav {
    display: flex;
    align-items: center;
    margin-top: 10px;
  }
}

.header-content{
  gap: 20px;
  padding: 0 100px;
  display: flex;
  justify-content: initial;
}

.logo-section {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 8px;
}

.environment-container {
  margin-left: auto;
  min-width: 300px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.main {
  &-layout {
    display: flex;
    height: 100%;
    max-width: $max-width;
    margin: 0 auto;
    padding: 20px;
    gap: 20px;
  }

  &-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 20px;
    min-width: 0;
  }
}

.sidebar {
  max-width: $sidebar-width;
  flex: 0 0 auto;
}

.request-section {
  flex-shrink: 0;
}

.response-section {
  flex: 1;
  min-height: 0;
}

// Tab bar styles
.tab-bar {
  display: flex;
  background-color: #f0f0f0;
  border-bottom: 1px solid #ddd;
  overflow-x: auto;
  height: 40px;
  align-items: center;
  
  &::-webkit-scrollbar {
    height: 4px;
  }
  
  .tab {
    display: flex;
    align-items: center;
    padding: 0 15px;
    height: 100%;
    min-width: 120px;
    max-width: 200px;
    background-color: #f8f8f8;
    border-right: 1px solid #ddd;
    cursor: pointer;
    position: relative;
    user-select: none;
    transition: background-color 0.2s;
    
    &:hover {
      background-color: #fff;
    }
    
    &.active {
      background-color: #fff;
      border-bottom: 2px solid #667eea;
    }
    
    .tab-name {
      flex: 1;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      font-size: 14px;
    }
    
    .tab-close {
      background: transparent;
      border: none;
      color: #999;
      font-size: 16px;
      cursor: pointer;
      padding: 0 5px;
      margin-left: 5px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      
      &:hover {
        background-color: rgba(0,0,0,0.05);
        color: #333;
      }
    }
  }
  
  .new-tab-button {
    background: transparent;
    border: none;
    color: #667eea;
    font-size: 20px;
    cursor: pointer;
    padding: 0 15px;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 0.2s;
    
    &:hover {
      background-color: rgba(102, 126, 234, 0.1);
    }
  }
}

// Scrollbar styling
::-webkit-scrollbar {
  width: 8px;

  &-track {
    background: #f1f1f1;
    border-radius: 4px;
  }

  &-thumb {
    background: #c1c1c1;
    border-radius: 4px;

    &:hover {
      background: #a8a8a8;
    }
  }
}

// Responsive design
@media (max-width: 1024px) {
  .main-layout {
    flex-direction: column;
    padding: 15px;
  }
  
  .sidebar {
    width: 100%;
    max-height: 300px;
  }
  
  .main-content {
    min-height: 0;
  }
}

@media (max-width: 768px) {
  .app-header {
    padding: 15px;
  }
  
  .app-title {
    font-size: 24px;
  }
  
  .header-subtitle {
    font-size: 14px;
  }
  
  .main-layout {
    padding: 10px;
    gap: 15px;
  }
}
</style>
