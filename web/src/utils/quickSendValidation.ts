/**
 * QuickSend form validation logic, extracted for testability.
 */

export interface QuickSendFormData {
  label: string
  command: string
}

/**
 * Validate a QuickSend form.
 * Returns an error message string, or empty string if valid.
 */
export function validateQuickSendForm(form: QuickSendFormData): string {
  const label = form.label.trim()
  const command = form.command.trim()
  if (!label || !command) {
    return 'chat.quickSend.itemRequired'
  }
  return ''
}
