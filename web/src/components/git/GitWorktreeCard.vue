<template>
  <div
    class="git-worktree-card"
    :class="{ current: worktree.isCurrent, locked: worktree.locked, missing: worktree.missing }"
    @click="!worktree.isCurrent && !worktree.missing && $emit('switch', worktree)"
  >
    <div class="wt-card-main">
      <div class="wt-card-name">
        <GitBranch :size="14" />
        <span>{{ worktree.branch || '—' }}</span>
      </div>
      <div class="wt-card-path">{{ worktree.path }}</div>
    </div>
    <div class="wt-card-status">
      <span v-if="worktree.dirty" class="wt-badge wt-badge-dirty">{{ t('git.manage.dirty', { count: worktree.untrackedCount }) }}</span>
      <span v-else class="wt-badge wt-badge-clean">{{ t('git.manage.clean') }}</span>
      <span v-if="worktree.locked" class="wt-badge wt-badge-locked">{{ t('git.manage.locked') }}</span>
      <span v-if="worktree.missing" class="wt-badge wt-badge-missing">{{ t('git.manage.pathMissing') }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { GitBranch } from 'lucide-vue-next'

const { t } = useI18n()

defineProps({
  worktree: { type: Object, required: true },
})

defineEmits(['switch'])
</script>

<style scoped>
.git-worktree-card {
  border: 1px solid var(--border-color, #dee2e6);
  border-radius: 8px;
  padding: 10px 12px;
  min-height: 44px;
  cursor: pointer;
  transition: border-color 0.15s, background 0.15s;
}

@media (hover: hover) {
  .git-worktree-card:hover {
    border-color: var(--accent-color, #4a90d9);
  }
}

.git-worktree-card.current {
  background: color-mix(in srgb, var(--accent-color) 12%, transparent);
  border-color: color-mix(in srgb, var(--accent-color) 30%, transparent);
  cursor: default;
}

.git-worktree-card.current .wt-card-name {
  color: var(--accent-color, #4a90d9);
}

.git-worktree-card.current .wt-card-path {
  color: var(--text-muted, #999);
}

.git-worktree-card.current .wt-badge-clean {
  background: color-mix(in srgb, var(--accent-color) 18%, transparent);
  color: var(--accent-color, #4a90d9);
}

.git-worktree-card.current .wt-badge-dirty {
  background: color-mix(in srgb, var(--accent-color) 18%, transparent);
  color: var(--accent-color, #4a90d9);
}

.git-worktree-card.missing {
  opacity: 0.6;
}

.git-worktree-card.locked {
  opacity: 0.8;
}

.wt-card-main {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.wt-card-name {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 15px;
  font-weight: 600;
  color: var(--accent-color, #4a90d9);
}

.wt-card-path {
  font-size: 11px;
  color: var(--text-muted, #999);
  word-break: break-all;
  line-height: 1.4;
}

.wt-card-status {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: 6px;
}

.wt-badge {
  font-size: 10px;
  font-weight: 600;
  padding: 1px 6px;
  border-radius: 4px;
  white-space: nowrap;
}

.wt-badge-dirty {
  background: var(--warning-bg, rgba(255, 159, 64, 0.15));
  color: var(--warning-color, #e67e22);
}

.wt-badge-clean {
  background: var(--success-bg, rgba(40, 167, 69, 0.12));
  color: var(--success-color, #28a745);
}

.wt-badge-locked {
  background: var(--bg-secondary, #e9ecef);
  color: var(--text-muted, #999);
}

.wt-badge-missing {
  background: var(--danger-bg, rgba(220, 53, 69, 0.12));
  color: var(--danger-color, #dc3545);
}
</style>
