<script setup>
import { ref } from 'vue'
import RequestBuilder from './components/RequestBuilder.vue'
import ResponseViewer from './components/ResponseViewer.vue'
import CollectionSidebar from './components/CollectionSidebar.vue'
import EnvironmentSelector from './components/EnvironmentSelector.vue'
import HeaderSettings from './components/HeaderSettings.vue'

const showHeaderSettings = ref(false)

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
  // Always clear the response viewer when loading a new request
  response.value = null

  if (requestBuilderRef.value) {
    // Handle both old format (direct request) and new format (with request and response)
    const requestData = logData.request || logData
    requestBuilderRef.value.loadRequest(requestData)

    // If response data is available (e.g., from logs), set it after clearing
    if (logData.response) {
      response.value = logData.response
    }
  }
}

const handleNewRequest = () => {
  if (requestBuilderRef.value) {
    requestBuilderRef.value.newRequest()
  }
  response.value = null
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
            <HeaderSettings :show="showHeaderSettings" @close="showHeaderSettings = false" @open="showHeaderSettings = true" />
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
            @new-request="handleNewRequest"
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
