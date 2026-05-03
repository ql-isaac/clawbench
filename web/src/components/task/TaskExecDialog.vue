<template>
  <ModalDialog :open="open" :title="view === 'detail' ? '' : '执行记录'" @close="handleClose">
    <template #header>
      <!-- Detail view: back arrow + time -->
      <template v-if="view === 'detail' && selectedExec">
        <button class="back-btn" @click.stop="view = 'list'" title="返回列表">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="15 18 9 12 15 6"/>
          </svg>
        </button>
        <span class="modal-title">{{ formatAbsoluteTime(selectedExec.createdAt) }}</span>
      </template>
      <!-- List view: icon + task name -->
      <template v-else>
        <svg class="modal-header-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>
        <span class="modal-title">{{ task?.name || '执行记录' }}</span>
      </template>
    </template>

    <!-- List view -->
    <div v-if="view === 'list'" class="executions-content">
      <div v-if="loading" class="dialog-empty">加载中...</div>
      <div v-else-if="executions.length === 0" class="dialog-empty">暂无执行记录</div>
      <div v-for="(exec, idx) in executions" :key="idx" class="execution-item" :class="{ unread: exec.isUnread }" @click="openDetail(exec)">
        <div class="execution-row">
          <div class="execution-info">
            <div class="execution-time-row">
              <span class="exec-absolute-time">{{ formatAbsoluteTime(exec.createdAt) }}</span>
              <span class="exec-relative-time">{{ chatRender.formatMessageTime(exec.createdAt) }}</span>
              <span v-if="exec.isUnread" class="exec-unread-dot"></span>
            </div>
            <div v-if="exec.summary" class="exec-summary">{{ exec.summary }}</div>
            <div v-else class="exec-summary empty">无文本输出</div>
          </div>
          <svg class="exec-chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="9 18 15 12 9 6"/>
          </svg>
        </div>
      </div>
    </div>

    <!-- Detail view -->
    <div v-if="view === 'detail' && selectedExec" class="detail-content">
      <ContentBlocks
        :blocks="selectedExec.blocks"
        msgId="exec-detail"
        :msgIndex="0"
        :expandedTools="expandedTools"
        :blockProposals="{}"
        :renderTextBlock="chatRender.renderTextBlock"
        :formatToolInput="chatRender.formatToolInput"
        :toolCallSummary="chatRender.toolCallSummary"
        @toggle-tool="toggleTool"
      />
    </div>

    <template #footer>
      <button class="btn btn-secondary" @click="handleClose">关闭</button>
    </template>
  </ModalDialog>
</template>

<script setup>
import { ref, watch, inject } from 'vue'
import ModalDialog from '@/components/common/ModalDialog.vue'
import ContentBlocks from '@/components/chat/ContentBlocks.vue'
import { useChatRender } from '@/composables/useChatRender.ts'

const props = defineProps({
  open: Boolean,
  task: Object,
})

const emit = defineEmits(['close'])

const loading = ref(false)
const executions = ref([])
const expandedTools = ref({})
const view = ref('list')  // 'list' | 'detail'
const selectedExec = ref(null)

// Create chatRender instance for rendering execution blocks
const renderTheme = inject('theme', ref('light'))
const chatRender = useChatRender({ messages: ref([]), theme: renderTheme, currentSessionId: ref('') })

function toggleTool(key) {
  expandedTools.value = { ...expandedTools.value, [key]: !expandedTools.value[key] }
}

function formatAbsoluteTime(createdAt) {
  const d = new Date(createdAt)
  const y = d.getFullYear()
  const mo = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const h = String(d.getHours()).padStart(2, '0')
  const mi = String(d.getMinutes()).padStart(2, '0')
  const s = String(d.getSeconds()).padStart(2, '0')
  return `${y}-${mo}-${day} ${h}:${mi}:${s}`
}

function extractSummary(blocks) {
  for (const block of blocks) {
    if (block.type === 'text' && block.text) {
      // Strip schedule-proposal tags and markdown
      const clean = block.text
        .replace(/<schedule-proposal>[\s\S]*?<\/schedule-proposal>/g, '')
        .replace(/[#*`_~\[\]()]/g, '')
        .trim()
      if (clean) {
        return clean.length > 80 ? clean.substring(0, 80) + '...' : clean
      }
    }
  }
  return ''
}

function openDetail(exec) {
  selectedExec.value = exec
  expandedTools.value = {}
  view.value = 'detail'
}

function handleClose() {
  view.value = 'list'
  selectedExec.value = null
  emit('close')
}

async function loadExecutions() {
  if (!props.task?.id) return
  loading.value = true
  try {
    const resp = await fetch(`/api/tasks/${props.task.id}/executions`)
    const data = await resp.json()
    const rawExecutions = data.executions || []
    executions.value = rawExecutions.map(exec => {
      const { blocks } = chatRender.parseAssistantContent(exec.content)
      const summary = extractSummary(blocks)
      return { ...exec, blocks, summary }
    })
  } catch (err) {
    console.error('Failed to load executions:', err)
  } finally {
    loading.value = false
  }
}

async function markTaskRead() {
  if (!props.task?.id) return
  try {
    await fetch(`/api/tasks/${props.task.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ action: 'read' }),
    })
  } catch (err) {
    console.error('Failed to mark task as read:', err)
  }
}

watch(() => props.open, (isOpen) => {
  if (!isOpen) return
  view.value = 'list'
  selectedExec.value = null
  expandedTools.value = {}
  loadExecutions()
  markTaskRead()
})
</script>

<style scoped>
.executions-content {
  flex: 1;
  overflow-y: auto;
  padding: 2px 0;
}

.detail-content {
  flex: 1;
  overflow-y: auto;
  padding: 6px;
}

.execution-item {
  border-bottom: 1px solid var(--border-color, #e5e5e5);
}

.execution-item:last-child {
  border-bottom: none;
}

.execution-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  cursor: pointer;
  transition: background 0.15s;
}

.execution-row:hover {
  background: var(--bg-secondary);
}

.execution-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.execution-time-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.exec-absolute-time {
  font-size: 12px;
  color: var(--text-primary);
  font-weight: 500;
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
}

.exec-relative-time {
  font-size: 11px;
  color: var(--text-muted, #999);
  white-space: nowrap;
}

.exec-unread-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--accent-color, #0066cc);
  flex-shrink: 0;
  animation: exec-unread-pulse 1.2s ease-in-out infinite;
}

@keyframes exec-unread-pulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.5; transform: scale(0.7); }
}

.exec-summary {
  font-size: 12px;
  color: var(--text-secondary, #666);
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.exec-summary.empty {
  color: var(--text-muted, #999);
  font-style: italic;
}

.execution-item.unread .exec-absolute-time {
  color: var(--accent-color, #0066cc);
}

.execution-item.unread {
  animation: exec-unread-flash 0.8s ease-in-out infinite;
}

@keyframes exec-unread-flash {
  0%, 100% {
    background: transparent;
  }
  50% {
    background: color-mix(in srgb, var(--accent-color, #0066cc) 6%, transparent);
  }
}

.exec-chevron {
  flex-shrink: 0;
  color: var(--text-muted, #ccc);
}

.back-btn {
  width: 22px;
  height: 22px;
  border: none;
  background: none;
  color: var(--accent-color, #0066cc);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: background 0.15s;
}

.back-btn:hover {
  background: rgba(0, 102, 204, 0.1);
}

.dialog-empty {
  text-align: center;
  padding: 20px 12px;
  color: var(--text-muted, #999);
  font-size: 13px;
}

/* Buttons */
.btn {
  padding: 5px 14px;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.15s, opacity 0.15s;
}

.btn-secondary {
  background: var(--bg-tertiary, #f0f0f0);
  color: var(--text-primary, #1a1a1a);
}

.btn-secondary:hover { background: #e0e0e0; }
</style>
