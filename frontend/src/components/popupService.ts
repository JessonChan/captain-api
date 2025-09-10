import { ref } from 'vue'

// Simple popup service that works with existing modals
class PopupService {
  // Show alert dialog
  async alert(message: string, config: any = {}): Promise<void> {
    return new Promise((resolve) => {
      // Use browser alert as fallback for now
      // In a real implementation, this would use a proper modal system
      console.log(`Alert: ${message}`, config)
      alert(message)
      resolve()
    })
  }

  // Show confirm dialog
  async confirm(message: string, config: any = {}): Promise<boolean> {
    return new Promise((resolve) => {
      // Use browser confirm as fallback for now
      console.log(`Confirm: ${message}`, config)
      const result = confirm(message)
      resolve(result)
    })
  }

  // Show custom dialog
  async show(config: any): Promise<any> {
    return new Promise((resolve) => {
      // For now, return a simple result
      console.log('Show custom dialog:', config)
      resolve(null)
    })
  }

  // Show form dialog
  async showForm(config: any): Promise<any> {
    return this.show(config)
  }

  // Close current popup
  close(): void {
    console.log('Close popup')
  }

  // Check if popup is visible
  isVisible(): boolean {
    return false
  }
}

// Create singleton instance
const popupService = new PopupService()

// Export service
export default popupService