import { describe, expect, it } from 'vitest'
import { validateQuickSendForm } from '@/utils/quickSendValidation.ts'

describe('validateQuickSendForm', () => {
  it('returns error key when label is empty', () => {
    expect(validateQuickSendForm({ label: '', command: 'git status' })).toBe('chat.quickSend.itemRequired')
  })
  it('returns error key when command is empty', () => {
    expect(validateQuickSendForm({ label: 'Git Status', command: '' })).toBe('chat.quickSend.itemRequired')
  })
  it('returns error key when both are empty', () => {
    expect(validateQuickSendForm({ label: '', command: '' })).toBe('chat.quickSend.itemRequired')
  })
  it('returns empty string when both are valid', () => {
    expect(validateQuickSendForm({ label: 'Build', command: 'go build ./...' })).toBe('')
  })
  it('trims whitespace before validation', () => {
    expect(validateQuickSendForm({ label: '  ', command: 'test' })).toBe('chat.quickSend.itemRequired')
  })
  it('accepts label with only whitespace after trim', () => {
    expect(validateQuickSendForm({ label: 'Build', command: '   ' })).toBe('chat.quickSend.itemRequired')
  })
  it('accepts valid input with leading/trailing spaces', () => {
    expect(validateQuickSendForm({ label: '  Build  ', command: '  go build  ' })).toBe('')
  })
})
