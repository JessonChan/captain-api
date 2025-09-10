import { ref } from 'vue'

// Simple popup service that works with existing modals
class PopupService {
  private alertDialog: any = null
  private confirmDialog: any = null

  // Set dialog references (called from App.vue)
  setDialogRefs(alertDialog: any, confirmDialog: any) {
    this.alertDialog = alertDialog
    this.confirmDialog = confirmDialog
  }

  // Show alert dialog
  async alert(message: string, config: any = {}): Promise<void> {
    return new Promise((resolve) => {
      if (this.alertDialog) {
        this.alertDialog.show(
          config.title || 'Alert',
          message,
          config.details || '',
          config.closeText || 'OK',
          config.severity || 'info',
          resolve
        )
      } else {
        // Fallback to console if dialog not available
        console.log(`Alert: ${message}`, config)
        resolve()
      }
    })
  }

  // Show confirm dialog
  async confirm(message: string, config: any = {}): Promise<boolean> {
    return new Promise((resolve) => {
      if (this.confirmDialog) {
        this.confirmDialog.show(
          config.title || 'Confirm Action',
          message,
          config.confirmText || 'Confirm',
          config.cancelText || 'Cancel',
          (result: boolean) => resolve(result),
          () => resolve(false)
        )
      } else {
        // Fallback to console if dialog not available
        console.log(`Confirm: ${message}`, config)
        resolve(false)
      }
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