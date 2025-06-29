<template>
  <div class="header-collection-form">
    <div class="form-header">
      <h4>{{ isEditing ? 'Edit Header Collection' : 'Add New Header Collection' }}</h4>
      <div v-if="isEditing" class="editing-badge">Editing</div>
    </div>
    <form @submit.prevent="handleSubmit">
      <div class="form-group">
        <div class="form-label-row">
          <label for="collection-name">Collection Name:</label>
          <span class="char-counter" :class="{ 'near-limit': localFormData.name.length > 40, 'at-limit': localFormData.name.length >= 50 }">
            {{ localFormData.name.length }}/50
          </span>
        </div>
        <input 
          id="collection-name"
          :value="localFormData.name"
          @input="e => updateField('name', e.target.value)"
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
          <span class="char-counter" :class="{ 'near-limit': localFormData.description.length > 80, 'at-limit': localFormData.description.length >= 100 }">
            {{ localFormData.description.length }}/100
          </span>
        </div>
        <input 
          id="collection-description"
          :value="localFormData.description"
          @input="e => updateField('description', e.target.value)"
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
          <div class="header-count">{{ localFormData.headers.filter(h => h.key && h.value).length }} valid headers</div>
        </div>
        
        <div class="header-table">
          <div class="header-table-head">
            <div class="header-name-col">Header Name</div>
            <div class="header-value-col">Value</div>
            <div class="header-action-col"></div>
          </div>
          
          <div class="header-rows">
            <div v-if="localFormData.headers.length === 0" class="empty-headers">
              <div class="empty-headers-message">No headers added yet</div>
              <div class="empty-headers-description">Click "Add Header" below to add your first header</div>
            </div>
            
            <div 
              class="header-row" 
              v-for="(header, index) in localFormData.headers" 
              :key="index"
              :class="{'is-valid': header.key && header.value, 'is-invalid': (!header.key && header.value) || (header.key && !header.value)}"
            >
              <div class="input-wrapper">
                <input 
                  :value="header.key" 
                  type="text" 
                  class="form-input header-input" 
                  placeholder="Header name"
                  @input="e => updateHeaders(index, 'key', e.target.value)"
                />
                <div v-if="!header.key && header.value" class="validation-message">Name required</div>
              </div>
              <div class="input-wrapper">
                <input 
                  :value="header.value" 
                  type="text" 
                  class="form-input header-input" 
                  placeholder="Header value"
                  @input="e => updateHeaders(index, 'value', e.target.value)"
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
        <button type="button" @click="$emit('cancel')" class="cancel-btn">
          <span class="btn-icon">✕</span>
          <span class="btn-text">Cancel</span>
        </button>
        <button type="submit" class="save-btn">
          <span class="btn-icon">{{ isEditing ? '✓' : '+' }}</span>
          <span class="btn-text">{{ isEditing ? 'Update Collection' : 'Add Collection' }}</span>
        </button>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

// Simple debounce implementation with TypeScript types
const debounce = <T extends (...args: any[]) => void>(
  fn: T,
  delay: number
): ((...args: Parameters<T>) => void) => {
  let timeoutId: ReturnType<typeof setTimeout> | null = null
  
  return function(this: any, ...args: Parameters<T>) {
    if (timeoutId !== null) {
      clearTimeout(timeoutId)
    }
    timeoutId = setTimeout(() => {
      fn.apply(this, args)
      timeoutId = null
    }, delay)
  }
}

interface Header {
  key: string
  value: string
}

interface HeaderCollection {
  id: string
  name: string
  description: string
  headers: Header[]
}

const props = withDefaults(defineProps<{
  collection: Partial<HeaderCollection>
  isEditing?: boolean
}>(), {
  collection: () => ({
    id: '',
    name: '',
    description: '',
    headers: []
  }),
  isEditing: false
})

// Local copy of the form data
const localFormData = ref<HeaderCollection>({
  id: '',
  name: '',
  description: '',
  headers: []
})

// Initialize local form data when props change
watch(() => props.collection, (newCollection) => {
  localFormData.value = {
    id: newCollection.id || '',
    name: newCollection.name || '',
    description: newCollection.description || '',
    headers: Array.isArray(newCollection.headers) 
      ? newCollection.headers.map(h => ({ ...h }))
      : newCollection.headers && typeof newCollection.headers === 'object'
        ? Object.entries(newCollection.headers).map(([key, value]) => ({ key, value }))
        : [{ key: '', value: '' }]
  }
  
  // Ensure at least one header row
  if (localFormData.value.headers.length === 0) {
    localFormData.value.headers = [{ key: '', value: '' }]
  }
}, { immediate: true, deep: true })

const emit = defineEmits<{
  (e: 'submit', value: HeaderCollection): void
  (e: 'cancel'): void
  (e: 'update:collection', value: HeaderCollection): void
}>()

const formData = ref<HeaderCollection>({
  id: '',
  name: '',
  description: '',
})

// Watch for changes in the collection prop
watch(() => props.collection, (newVal) => {
  if (newVal) {
    formData.value = {
      id: newVal.id || '',
      name: newVal.name || '',
      description: newVal.description || '',
      headers: Array.isArray(newVal.headers) 
        ? [...newVal.headers] 
        : newVal.headers && typeof newVal.headers === 'object' && !Array.isArray(newVal.headers)
          ? Object.entries(newVal.headers).map(([key, value]) => ({ key, value }))
          : [{ key: '', value: '' }]
    }
  }
}, { immediate: true, deep: true })

// Handle input changes with debounce
const updateField = (field: string, value: any) => {
  localFormData.value[field] = value
  debouncedEmit()
}

// Handle header changes
const updateHeaders = (index: number, field: 'key' | 'value', value: string) => {
  localFormData.value.headers[index][field] = value
  debouncedEmit()
}

// Debounce emit to prevent too many updates
const debouncedEmit = debounce(() => {
  emit('update:collection', { ...localFormData.value })
}, 300)

// Add a new header field
const addHeader = (): void => {
  localFormData.value.headers.push({ key: '', value: '' })
  debouncedEmit()
}

// Remove a header field
const removeHeader = (index: number): void => {
  localFormData.value.headers.splice(index, 1)
  debouncedEmit()
}
// Handle form submission
const handleSubmit = () => {
  emit('submit', {
    id: localFormData.value.id,
    name: localFormData.value.name,
    description: localFormData.value.description,
    headers: localFormData.value.headers.filter(h => h.key || h.value)
  })
}
</script>

<style lang="scss" scoped>
/* All the styles from the original component can be copied here */
.header-collection-form {
  background-color: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;

  .form-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    h4 {
      margin: 0;
      font-size: 1.1rem;
      color: #333;
    }
  }


  .editing-badge {
    background-color: #e6f7ff;
    color: #1890ff;
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 0.75rem;
    font-weight: 500;
  }


  .form-group {
    margin-bottom: 16px;

    .form-label-row {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 6px;

      label {
        font-weight: 500;
        font-size: 0.9rem;
        color: #555;
      }
    }

  }

  .char-counter {
    font-size: 0.75rem;
    color: #888;

    &.near-limit {
      color: #f39c12;
    }


    &.at-limit {
      color: #e74c3c;
      font-weight: bold;
    }
  }

  .form-input {
    width: 100%;
    padding: 10px 12px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 0.9rem;
    transition: border-color 0.2s, box-shadow 0.2s;

    &:focus {
      border-color: #1890ff;
      box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
      outline: none;
    }
  }
}

// Headers Editor
.headers-editor {
  margin-top: 24px;
  border: 1px solid #eee;
  border-radius: 6px;
  overflow: hidden;

  &-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    background-color: #f9f9f9;
    border-bottom: 1px solid #eee;

    h5 {
      margin: 0;
      font-size: 0.9rem;
      color: #333;
    }
  }


  .header-count {
    font-size: 0.75rem;
    color: #888;
  }


  .header-table {
    width: 100%;

    &-head {
      display: flex;
      background-color: #f5f5f5;
      border-bottom: 1px solid #eee;
      font-size: 0.8rem;
      font-weight: 600;
      color: #666;
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }
  }

  .header-name-col, .header-value-col {
    padding: 10px 12px;
  }

  .header-name-col {
    width: 40%;
    border-right: 1px solid #eee;
  }

  .header-value-col {
    width: 55%;
  }

  .header-action-col {
    width: 5%;
    min-width: 40px;
  }

  .header-rows {
    max-height: 300px;
    overflow-y: auto;
  }

  .empty-headers {
    padding: 24px;
    text-align: center;
    color: #999;

    &-message {
      font-weight: 500;
      margin-bottom: 4px;
    }

    &-description {
      font-size: 0.85rem;
    }
  }

  .header-row {
    display: flex;
    align-items: flex-start;
    padding: 12px;
    border-bottom: 1px solid #f0f0f0;
    transition: background-color 0.2s;

    &:hover {
      background-color: #fafafa;
    }

    &.is-valid {
      background-color: #f6ffed;
    }


    &.is-invalid {
      background-color: #fff1f0;
    }
  }

  .input-wrapper {
    position: relative;
    flex: 1;
    margin-right: 8px;
  }

  .header-input {
    width: 100%;
    padding: 8px 12px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 0.9rem;
    transition: border-color 0.2s, box-shadow 0.2s;

    &:focus {
      border-color: #1890ff;
      box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
      outline: none;
    }

    &.invalid-input {
      border-color: #ff4d4f;
    }
  }

  .validation-message {
    position: absolute;
    bottom: -18px;
    left: 0;
    color: #ff4d4f;
    font-size: 0.75rem;
    line-height: 1;
  }

  .remove-header-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    border: none;
    background: none;
    color: #999;
    cursor: pointer;
    border-radius: 4px;
    transition: background-color 0.2s, color 0.2s;

    &:hover {
      background-color: #f5f5f5;
      color: #ff4d4f;
    }
  }
  .remove-icon {
    font-size: 1.2rem;
    line-height: 1;
  }

  .add-header-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    padding: 10px;
    background-color: #f9f9f9;
    border: 1px dashed #d9d9d9;
    color: #666;
    cursor: pointer;
    font-size: 0.9rem;
    transition: all 0.2s;
    border-radius: 0 0 6px 6px;

    &:hover {
      background-color: #f0f0f0;
      border-color: #bbb;
      color: #333;
    }
  }
  .add-icon {
    margin-right: 6px;
    font-weight: bold;
  }
}

// Form Actions
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;

  .cancel-btn, .save-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 8px 16px;
    border: 1px solid #d9d9d9;
    border-radius: 4px;
    background-color: #fff;
    color: #333;
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }

  .cancel-btn {
    &:hover {
      background-color: #f5f5f5;
      border-color: #bbb;
    }
  }

  .save-btn {
    background-color: #1890ff;
    border-color: #1890ff;
    color: #fff;

    &:hover {
      background-color: #40a9ff;
      border-color: #40a9ff;
    }
  }

  .btn-icon {
    margin-right: 6px;
    line-height: 1;
  }

  .btn-text {
    line-height: 1.5;
  }
}
</style>
