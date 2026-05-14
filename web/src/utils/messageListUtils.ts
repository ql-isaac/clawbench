/**
 * Pure functions extracted from ChatMessageList.vue for testability.
 */

/**
 * Compute how many older messages are not yet loaded.
 */
export function computeRemainingCount(hasMore: boolean, totalMessages: number, loadedMessages: number): number {
  if (!hasMore) return 0
  return Math.max(0, totalMessages - loadedMessages)
}

/**
 * Compute the indices of the last "round" of messages.
 * A round consists of the last assistant message and its preceding user message.
 * If no assistant message, the last user message alone is the round.
 */
export function computeLastRoundIndices(messages: { role: string }[]): Set<number> {
  if (!messages || messages.length === 0) return new Set()

  // Find last assistant message index
  let lastAssistantIdx = -1
  for (let i = messages.length - 1; i >= 0; i--) {
    if (messages[i].role === 'assistant') {
      lastAssistantIdx = i
      break
    }
  }

  const indices = new Set<number>()
  if (lastAssistantIdx >= 0) {
    indices.add(lastAssistantIdx)
    // Find the preceding user message
    for (let i = lastAssistantIdx - 1; i >= 0; i--) {
      if (messages[i].role === 'user') {
        indices.add(i)
        break
      }
    }
  } else {
    // No assistant message — expand last user message
    for (let i = messages.length - 1; i >= 0; i--) {
      if (messages[i].role === 'user') {
        indices.add(i)
        break
      }
    }
  }

  return indices
}

/**
 * Determine whether a message should be collapsed.
 * Returns true if the message should be suggested for collapse.
 * - If user explicitly collapsed it → suggest collapse
 * - If it's in the last round → don't suggest collapse (unless user collapsed)
 * - If user explicitly expanded → don't suggest collapse
 * - Otherwise → suggest collapse
 */
export function isCollapsed(
  index: number,
  msg: { role: string },
  collapsedSet: Set<number>,
  lastRoundIndices: Set<number>,
  expandedSet: Set<number>
): boolean {
  // User explicitly collapsed this message — suggest collapse
  if (collapsedSet.has(index)) return true
  // Last round is always fully expanded unless user collapsed it
  if (lastRoundIndices.has(index)) return false
  // Manually expanded — don't suggest collapse
  if (expandedSet.has(index)) return false
  // Everything else: suggest collapse
  return true
}
