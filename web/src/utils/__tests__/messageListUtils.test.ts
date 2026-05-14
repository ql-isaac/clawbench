import { describe, expect, it } from 'vitest'
import { computeRemainingCount, computeLastRoundIndices, isCollapsed } from '@/utils/messageListUtils.ts'

describe('computeRemainingCount', () => {
  it('returns 0 when hasMore is false', () => {
    expect(computeRemainingCount(false, 100, 20)).toBe(0)
  })
  it('returns difference when hasMore is true', () => {
    expect(computeRemainingCount(true, 100, 20)).toBe(80)
  })
  it('returns 0 when total equals loaded', () => {
    expect(computeRemainingCount(true, 20, 20)).toBe(0)
  })
  it('clamps to 0 when loaded exceeds total', () => {
    expect(computeRemainingCount(true, 10, 20)).toBe(0)
  })
  it('handles zero total', () => {
    expect(computeRemainingCount(true, 0, 0)).toBe(0)
  })
})

describe('computeLastRoundIndices', () => {
  it('returns empty set for empty array', () => {
    expect(computeLastRoundIndices([])).toEqual(new Set())
  })
  it('returns empty set for null', () => {
    expect(computeLastRoundIndices(null as any)).toEqual(new Set())
  })
  it('finds last assistant and preceding user', () => {
    const msgs = [
      { role: 'user' },
      { role: 'assistant' },
      { role: 'user' },
      { role: 'assistant' },
    ]
    const indices = computeLastRoundIndices(msgs)
    expect(indices.has(2)).toBe(true) // last user before last assistant
    expect(indices.has(3)).toBe(true) // last assistant
    expect(indices.size).toBe(2)
  })
  it('finds only last user when no assistant', () => {
    const msgs = [
      { role: 'user' },
      { role: 'user' },
    ]
    const indices = computeLastRoundIndices(msgs)
    expect(indices.has(1)).toBe(true) // last user
    expect(indices.size).toBe(1)
  })
  it('finds assistant with no preceding user', () => {
    const msgs = [
      { role: 'assistant' },
    ]
    const indices = computeLastRoundIndices(msgs)
    expect(indices.has(0)).toBe(true)
    expect(indices.size).toBe(1)
  })
  it('finds closest preceding user, not any user', () => {
    const msgs = [
      { role: 'user' },      // 0
      { role: 'assistant' },  // 1
      { role: 'user' },       // 2 - this is the one that should be paired
      { role: 'assistant' },  // 3
    ]
    const indices = computeLastRoundIndices(msgs)
    expect(indices.has(2)).toBe(true)  // user right before assistant 3
    expect(indices.has(3)).toBe(true)  // assistant 3
    expect(indices.has(0)).toBe(false) // earlier user, not paired
  })
  it('handles single user message', () => {
    const indices = computeLastRoundIndices([{ role: 'user' }])
    expect(indices.has(0)).toBe(true)
    expect(indices.size).toBe(1)
  })
  it('skips non-user/non-assistant messages when looking for user', () => {
    const msgs = [
      { role: 'system' },
      { role: 'user' },
      { role: 'system' },
      { role: 'assistant' },
    ]
    const indices = computeLastRoundIndices(msgs)
    expect(indices.has(1)).toBe(true) // user
    expect(indices.has(3)).toBe(true) // assistant
  })
})

describe('isCollapsed', () => {
  const lastRoundIndices = new Set([2, 3])

  it('returns true when in collapsedSet', () => {
    expect(isCollapsed(3, { role: 'assistant' }, new Set([3]), lastRoundIndices, new Set())).toBe(true)
  })
  it('returns false when in lastRoundIndices (and not in collapsedSet)', () => {
    expect(isCollapsed(3, { role: 'assistant' }, new Set(), lastRoundIndices, new Set())).toBe(false)
  })
  it('returns false when in expandedSet', () => {
    expect(isCollapsed(0, { role: 'user' }, new Set(), lastRoundIndices, new Set([0]))).toBe(false)
  })
  it('returns true by default (not in any set)', () => {
    expect(isCollapsed(0, { role: 'user' }, new Set(), lastRoundIndices, new Set())).toBe(true)
  })
  it('collapsedSet takes priority over lastRoundIndices', () => {
    // User explicitly collapsed a last-round message
    expect(isCollapsed(2, { role: 'user' }, new Set([2]), lastRoundIndices, new Set())).toBe(true)
  })
  it('collapsedSet takes priority over expandedSet', () => {
    // Both sets is unusual, but collapsedSet wins
    expect(isCollapsed(1, { role: 'user' }, new Set([1]), lastRoundIndices, new Set([1]))).toBe(true)
  })
})
