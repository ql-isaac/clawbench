<template>
  <BottomSheet ref="bottomSheetRef" :open="open" compact :title="t('task.title')" @close="$emit('close')">
    <template #header>
      <Clock :size="16" class="bs-header-icon" />
      <span class="bs-header-title">{{ t('task.title') }}</span>
      <button class="create-btn" @click="openCreateDialog" :title="t('task.form.createTitle')">
        <Plus :size="16" />
      </button>
    </template>

    <div class="task-list">
      <div v-if="loading" class="task-loading">{{ t('common.loading') }}</div>
      <div v-else-if="tasks.length === 0" class="task-empty">{{ t('task.noTasks') }}</div>
      <div v-for="task in tasks" :key="task.id" class="task-item" :class="[task.status, { 'has-unread': task.unreadCount > 0 }]">
        <div class="task-item-main" @click="openEditDialog(task)">
          <div class="task-item-info">
            <div class="task-item-header">
              <span class="task-item-icon">{{ getAgentIcon(task.agentId) }}</span>
              <span class="task-item-name">{{ task.name }}</span>
              <span v-if="task.runningCount > 0" class="task-item-running-dot" :title="t('task.exec.running')"></span>
              <span v-if="task.unreadCount > 0" class="task-item-unread">{{ task.unreadCount }}</span>
              <span class="task-item-status" :class="task.status">{{ statusLabel(task.status) }}</span>
            </div>
            <div class="task-item-meta">
              <span class="task-item-cron">{{ humanizeCron(task.cronExpr) }}</span>
              <span class="task-item-repeat">{{ repeatLabel(task.repeatMode, task.maxRuns) }}</span>
              <span v-if="task.repeatMode !== 'unlimited'" class="task-item-progress">{{ task.runCount }}/{{ task.maxRuns || 1 }}</span>
            </div>
            <div v-if="task.nextRunAt" class="task-item-next">
              {{ t('task.nextRun', { time: formatDateTime(task.nextRunAt) }) }}
            </div>
          </div>
          <div class="task-item-actions">
            <button class="task-action-btn history" @click.stop="openExecDialog(task)" :title="t('task.exec.title')">
              <History :size="14" />
            </button>
            <button class="task-action-btn delete" @click.stop="deleteTask(task.id)" :title="t('common.delete')">
              <Trash2 :size="14" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <TaskFormDialog
      :open="formDialogOpen"
      :mode="formMode"
      :task="selectedTask"
      @close="formDialogOpen = false"
      @saved="onFormSaved"
    />
    <TaskExecDialog
      :open="execDialogOpen"
      :task="selectedTask"
      @close="execDialogOpen = false"
    />
  </BottomSheet>
</template>

<script setup>
import { Clock, Plus, History, Trash2 } from 'lucide-vue-next'
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import BottomSheet from '@/components/common/BottomSheet.vue'
import TaskFormDialog from '@/components/task/TaskFormDialog.vue'
import TaskExecDialog from '@/components/task/TaskExecDialog.vue'
import { useAgents } from '@/composables/useAgents.ts'
import { humanizeCron, repeatLabel, statusLabel, formatDateTime } from '@/utils/format.ts'
import { store } from '@/stores/app.ts'
import { useDialog } from '@/composables/useDialog.ts'

const { t } = useI18n()
const dialog = useDialog()

const props = defineProps({
  open: Boolean,
})

const emit = defineEmits(['close'])

const bottomSheetRef = ref(null)
const tasks = computed(() => store.state.tasks)
const loading = ref(false)
const formDialogOpen = ref(false)
const formMode = ref('create')
const execDialogOpen = ref(false)
const selectedTask = ref(null)
const { agents, loadAgents, getAgentIcon } = useAgents()

defineExpose({ loadTasks })

async function loadTasks() {
  loading.value = true
  try {
    const resp = await fetch('/api/tasks')
    const data = await resp.json()
    store.state.tasks = data.tasks || []
  } catch (err) {
    console.error('Failed to load tasks:', err)
  } finally {
    loading.value = false
  }
}

async function markAllTasksRead() {
  // Mark all tasks with unread executions as read
  const unreadTasks = tasks.value.filter(t => t.unreadCount > 0)
  if (unreadTasks.length === 0) return
  await Promise.all(unreadTasks.map(t =>
    fetch(`/api/tasks/${t.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ action: 'read' }),
    }).catch(() => {})
  ))
  // Clear the global unread indicator
  store.state.taskUnread = false
}

function openCreateDialog() {
  formMode.value = 'create'
  selectedTask.value = null
  formDialogOpen.value = true
}

function openEditDialog(task) {
  formMode.value = 'edit'
  selectedTask.value = task
  formDialogOpen.value = true
}

function openExecDialog(task) {
  selectedTask.value = task
  execDialogOpen.value = true
}

async function onFormSaved() {
  formDialogOpen.value = false
  loadTasks()
}

async function deleteTask(id) {
  if (!await dialog.confirm(t('task.confirmDelete'), { dangerous: true })) return
  try {
    await fetch(`/api/tasks/${id}`, { method: 'DELETE' })
    await loadTasks()
  } catch (err) {
    console.error('Failed to delete task:', err)
  }
}

watch(() => props.open, async (val) => {
  if (val) {
    await Promise.all([loadTasks(), loadAgents()])
    markAllTasksRead()
  }
})
</script>

<style scoped>
.create-btn {
  margin-left: auto;
  width: 24px;
  height: 24px;
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

.create-btn:hover {
  background: rgba(0, 102, 204, 0.1);
}

.task-list {
  display: flex;
  flex-direction: column;
  gap: 0;
  padding: 4px 0;
  min-height: 0;
  overflow-y: auto;
  flex: 1;
}

.task-loading,
.task-empty {
  padding: 32px 12px;
  text-align: center;
  color: var(--text-muted, #999);
  font-size: 13px;
}

.task-item {
  position: relative;
}

.task-item.completed {
  opacity: 0.5;
}

.task-item-main {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 10px 16px;
  cursor: pointer;
  transition: background 0.15s;
  gap: 8px;
}

@media (hover: hover) {
  .task-item-main:hover {
    background: var(--bg-tertiary, rgba(0, 0, 0, 0.03));
  }
}

.task-item-main:active {
  background: var(--bg-tertiary, rgba(0, 0, 0, 0.06));
}

.task-item:not(:last-child)::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 16px;
  right: 16px;
  height: 1px;
  background: var(--border-color, #e5e5e5);
  opacity: 0.5;
}

.task-item-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 3px;
  min-width: 0;
  overflow: hidden;
}

.task-item-header {
  display: flex;
  align-items: center;
  gap: 5px;
  min-width: 0;
}

.task-item-icon {
  font-size: 14px;
  flex-shrink: 0;
}

.task-item-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary, #1a1a1a);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  min-width: 0;
}

.task-item-unread {
  font-size: 9px;
  padding: 1px 5px;
  border-radius: 8px;
  font-weight: 600;
  background: #ef4444;
  color: #fff;
  flex-shrink: 0;
  min-width: 14px;
  text-align: center;
  line-height: 1.3;
}

.task-item.has-unread .task-item-icon {
  animation: task-unread-flash 0.8s ease-in-out infinite;
}

@keyframes task-unread-flash {
  0%, 100% {
    opacity: 1;
    text-shadow: 0 0 0 transparent;
  }
  50% {
    opacity: 0.7;
    text-shadow: 0 0 8px color-mix(in srgb, var(--accent-color, #0066cc) 40%, transparent);
  }
}

.task-item-status {
  font-size: 9px;
  padding: 1px 5px;
  border-radius: 3px;
  font-weight: 500;
  flex-shrink: 0;
  line-height: 1.4;
}

.task-item-status.active {
  background: rgba(34, 197, 94, 0.12);
  color: #22c55e;
}

.task-item-status.paused {
  background: rgba(234, 179, 8, 0.12);
  color: #eab308;
}

.task-item-status.completed {
  background: var(--bg-tertiary, #e9ecef);
  color: var(--text-muted, #999);
}

.task-item-running-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--success-color, #22c55e);
  flex-shrink: 0;
  animation: task-running-pulse 1.5s ease-in-out infinite;
}

@keyframes task-running-pulse {
  0%, 100% { opacity: 1; box-shadow: 0 0 0 0 rgba(34, 197, 94, 0.4); }
  50% { opacity: 0.7; box-shadow: 0 0 6px 2px rgba(34, 197, 94, 0.2); }
}

.task-item-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  color: var(--text-muted, #999);
  min-width: 0;
  overflow: hidden;
  flex-wrap: wrap;
}

.task-item-cron {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 60%;
}

.task-item-next {
  font-size: 10px;
  color: var(--text-muted, #999);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.task-item-progress {
  font-weight: 500;
  color: var(--accent-color, #0066cc);
  flex-shrink: 0;
}

.task-item-actions {
  display: flex;
  gap: 2px;
  flex-shrink: 0;
  align-self: center;
}

.task-action-btn {
  width: 24px;
  height: 24px;
  border: none;
  background: none;
  color: var(--text-muted, #999);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.15s;
  opacity: 0;
}

.task-item-main:hover .task-action-btn,
.task-action-btn:focus-visible {
  opacity: 1;
}

@media (hover: none) {
  .task-action-btn {
    opacity: 0.6;
  }
}

.task-action-btn:hover {
  color: var(--text-secondary, #666);
  background: var(--bg-tertiary, rgba(0, 0, 0, 0.06));
}

.task-action-btn.delete:hover {
  color: #dc3545;
  background: rgba(220, 53, 69, 0.08);
}

.task-action-btn.exec.has-unread {
  color: var(--accent-color, #0066cc);
  animation: exec-btn-flash 0.8s ease-in-out infinite;
}

@keyframes exec-btn-flash {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
    background: color-mix(in srgb, var(--accent-color, #0066cc) 12%, transparent);
  }
}
</style>
