<script setup>
import { ref } from 'vue'
import RequestBuilder from './components/RequestBuilder.vue'
import ResponseViewer from './components/ResponseViewer.vue'
import CollectionSidebar from './components/CollectionSidebar.vue'
import EnvironmentSelector from './components/EnvironmentSelector.vue'
import HeaderSettings from './components/HeaderSettings.vue'

const response = ref(null)
const requestBuilderRef = ref(null)
const collectionSidebarRef = ref(null)
const currentCollectionId = ref('default')

const handleResponseReceived = async (responseData) => {
  response.value = responseData
  // Refresh logs after a request is made
  if (collectionSidebarRef.value && collectionSidebarRef.value.requestLogsRef) {
    await collectionSidebarRef.value.requestLogsRef.loadLogs()
  }
}

const handleLoadRequest = (logData) => {
  if (requestBuilderRef.value) {
    // Handle both old format (direct request) and new format (with request and response)
    const requestData = logData.request || logData
    requestBuilderRef.value.loadRequest(requestData)
    
    // If response data is available, set it
    if (logData.response) {
      response.value = logData.response
    }
  }
}

const handleRequestSaved = () => {
  if (collectionSidebarRef.value) {
    collectionSidebarRef.value.refresh()
  }
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
        <div class="header-nav">
          <div class="environment-container">
            <EnvironmentSelector />
            <HeaderSettings />
          </div>
        </div>
      </div>
    </header>

    <main class="app-main">
      <div class="main-layout">
        <!-- Sidebar -->
        <aside class="sidebar">
          <CollectionSidebar 
            ref="collectionSidebarRef"
            @load-request="handleLoadRequest"
          />
        </aside>

        <!-- Main Content -->
        <div class="main-content">
          <!-- Request Builder -->
          <section class="request-section">
            <RequestBuilder 
              ref="requestBuilderRef"
              :collectionId="currentCollectionId"
              @response-received="handleResponseReceived"
              @request-saved="handleRequestSaved"
            />
          </section>

          <!-- Response Viewer -->
          <section class="response-section">
            <ResponseViewer :response="response" />
          </section>
        </div>
      </div>
    </main>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
    'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
    sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background: #f5f6fa;
}

.app {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  height: 100%;
  flex-wrap: wrap;
}

.logo-section {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 8px;
}

.app-logo {
  width: 40px;
  height: 40px;
  border-radius: 10px;
}

.app-title {
  font-size: 28px;
  font-weight: 600;
  margin: 0;
}

.header-subtitle {
  font-size: 16px;
  opacity: 0.9;
  font-weight: 300;
}

.environment-container {
  margin-left: auto;
  min-width: 300px;
}

.app-main {
  flex: 1;
  overflow: hidden;
}

.main-layout {
  display: flex;
  height: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  gap: 20px;
}

.sidebar {
  width: 300px;
  flex-shrink: 0;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;
  min-width: 0;
}

.request-section {
  flex-shrink: 0;
}

.response-section {
  flex: 1;
  min-height: 0;
}

/* Header Navigation */
.header-nav {
  display: flex;
  align-items: center;
  margin-top: 10px;
}

.environment-container {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 10px;
}



/* Responsive Design */
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

/* Scrollbar Styling */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>
