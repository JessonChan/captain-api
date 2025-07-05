<template>
  <div class="welcome-container">
    <img src="/logo.png" alt="Captain API Logo" class="logo" />
    <h1 class="title">Welcome to Captain API</h1>
    <p class="subtitle">Your friendly neighborhood API client.</p>
    <button @click="createNewRequest" class="new-request-btn">
      <span class="plus-icon">+</span> Create New Request
    </button>
    <div class="tips">
      <p><strong>Pro Tip:</strong> Use Ctrl/Cmd + Enter to send a request.</p>
    </div>
  </div>
</template>

<script setup>
const emit = defineEmits(['new-request']);

const createNewRequest = () => {
  emit('new-request');
};
</script>

<style scoped>
.welcome-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  text-align: center;
  background-color: #f8f9fa;
  padding: 40px;
  border-radius: 8px;
}

.logo {
  width: 80px;
  height: 80px;
  margin-bottom: 20px;
}

.title {
  font-size: 2.5rem;
  font-weight: 600;
  color: #343a40;
  margin-bottom: 10px;
}

.subtitle {
  font-size: 1.2rem;
  color: #6c757d;
  margin-bottom: 30px;
}

.new-request-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 12px 24px;
  border: none;
  border-radius: 6px;
  background-color: #007bff;
  color: white;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 2px 8px rgba(0, 123, 255, 0.3);
}

.new-request-btn:hover {
  background-color: #0069d9;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 123, 255, 0.4);
}

.plus-icon {
  margin-right: 8px;
  font-weight: bold;
}

.tips {
  margin-top: 40px;
  font-size: 0.9rem;
  color: #6c757d;
}
</style>
