<template>
  <PopupMenu :show="show" @update:show="emit('update:show', $event)" :target-element="targetElement" :max-width="180" :menu-items-count="3" anchor="right">
    <button class="tab-menu-item" @click="handleClose">
      {{ t('terminal.closeTab') }}
    </button>
    <button class="tab-menu-item" @click="handleCopyPath">
      {{ t('terminal.copyPath') }}
    </button>
    <button class="tab-menu-item" @click="handleNewTabHere">
      {{ t('terminal.newTabHere') }}
    </button>
  </PopupMenu>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useToast } from '@/composables/useToast'
import PopupMenu from '@/components/common/PopupMenu.vue'

const props = defineProps<{
  show: boolean
  targetElement: HTMLElement | null
  cwd: string
}>()

const emit = defineEmits<{
  'update:show': [value: boolean]
  close: []
  copyPath: []
  newTabHere: []
}>()

const { t } = useI18n()
const toast = useToast()

function handleClose() {
  emit('update:show', false)
  emit('close')
}

function handleCopyPath() {
  emit('update:show', false)
  navigator.clipboard.writeText(props.cwd).catch(() => {})
  toast.show(t('common.copied'), { type: 'success', duration: 1500 })
  emit('copyPath')
}

function handleNewTabHere() {
  emit('update:show', false)
  emit('newTabHere')
}
</script>

<style>
.tab-menu-item {
  display: block;
  width: 100%;
  padding: 8px 14px;
  border: none;
  background: none;
  color: var(--text-primary);
  font-size: 13px;
  cursor: pointer;
  text-align: left;
  transition: background 0.12s, color 0.12s;
  position: relative;
  overflow: hidden;
}

.tab-menu-item:hover {
  background: var(--accent-color, #0066cc);
  color: #fff;
}
</style>
