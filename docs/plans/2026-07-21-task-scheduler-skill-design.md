# Task Scheduler Skill Design

**Date**: 2026-07-21
**Status**: Approved

## Overview

Refactor the scheduled task system from passive `<schedule-proposal>` tag detection to an AI-driven Skill + CLI architecture. AI agents actively call `clawbench task` CLI commands to manage tasks, then reference them with `<scheduled-task id="..." />` tags. Frontend renders cards by querying the API with the task ID, ensuring cards always reflect the latest task state.

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| AI invocation method | `clawbench task` dedicated CLI subcommands | No auth needed, type-safe, reuses Go code |
| Old system compatibility | Complete removal, no backward compat | Clean architecture, old tags display as plain text |
| New tag format | `<scheduled-task id="task-xxx" />` (self-closing, anchor only) | Card data comes from API, not tag |
| CLI command scope | Full operations: create/update/delete/pause/resume/trigger | AI has complete control |
| Anti-recursion protection | Skill-level filtering + CLI secondary check | Fundamental block at discovery layer |
| Card data source | API real-time query | Always up-to-date |
| Card inline actions | Edit/Pause/Resume/Trigger/Delete buttons | No need to navigate to TaskDrawer |
| Old tag display | Accept raw tag text in old messages | Simplicity over cosmetics |

## Architecture

### Data Flow (New)

```
User: "帮我每天9点总结代码"
  │
  ├─ AI sees task-scheduler skill trigger match
  │   ├─ GET /api/skills/task-scheduler.md → reads CLI docs
  │   └─ Decides to create a scheduled task
  │
  ├─ AI calls Bash: clawbench task create --name "..." --cron "0 9 * * *" ...
  │   ├─ CLI outputs: {"ok":true,"task":{"id":"task-abc",...}}
  │   └─ AI gets task_id = "task-abc"
  │
  ├─ AI outputs text + <scheduled-task id="task-abc" />
  │   ├─ Tag is anchor only, no data payload
  │   └─ AI can write natural language explanation around the tag
  │
  ├─ Backend finalize phase:
  │   ├─ No longer detects schedule-proposal
  │   └─ Tag stored as plain text in DB
  │
  ├─ SSE / load history → Frontend rendering:
  │   ├─ Regex matches <scheduled-task id="task-abc" />
  │   ├─ Strips tag, renders text
  │   ├─ Async GET /api/tasks/task-abc → populates card
  │   └─ Card data = API real-time response (always latest)
  │
  └─ Subsequent: card inline action buttons call API directly
     ├─ Pause: PUT /api/tasks/task-abc (action=pause)
     ├─ Trigger: PUT /api/tasks/task-abc (action=trigger)
     └─ Delete: DELETE /api/tasks/task-abc → card shows "已删除"
```

### Data Flow (Old - Removed)

```
AI outputs <schedule-proposal>{JSON}</schedule-proposal>
  → Backend finalize: regex detect → auto-create task → inject task_id into JSON
  → Frontend: extract JSON from tag → render card from static snapshot
  → Updates via store.state.tasks polling
```

## Phase 1: Backend CLI Subcommands (Non-breaking)

### 1.1 Entry point dispatch (`cmd/server/main.go`)

Add `os.Args[1] == "task"` dispatch before server startup:

```go
if len(os.Args) > 1 && os.Args[1] == "task" {
    os.Exit(runTaskCommand(os.Args[2:]))
}
```

### 1.2 CLI implementation (`internal/cli/task.go` - new file)

6 subcommands using standard `flag` package:

```bash
clawbench task create --name NAME --cron EXPR --agent AGENT_ID --prompt PROMPT --repeat MODE [--max-runs N]
clawbench task update TASK_ID [--name NAME] [--cron EXPR] [--agent AGENT_ID] [--prompt PROMPT] [--repeat MODE] [--max-runs N]
clawbench task delete TASK_ID
clawbench task pause TASK_ID
clawbench task resume TASK_ID
clawbench task trigger TASK_ID
```

**Unified output format**:
```json
// Success
{"ok":true,"task":{"id":"task-xxx","name":"...","status":"active",...}}

// Failure
{"ok":false,"error":"task not found"}
```

**Implementation requirements**:
- Load config via `model.LoadConfig()`
- Init database via `service.InitDB()`
- Initialize `GlobalScheduler` (reuse existing constructor)
- Call corresponding Scheduler methods
- Output JSON to stdout
- Exit code: 0=success, 1=failure

**Anti-recursion secondary check**:
- `clawbench task create` checks for `CLAWBENCH_SCHEDULED=1` environment variable
- Parent `CLIBackend.ExecuteStream()` injects this env var when `ScheduledExecution=true`
- Returns `{"ok":false,"error":"scheduled execution cannot create new tasks"}`

### 1.3 Tests (`internal/cli/task_test.go` - new file)

Unit tests for all 6 subcommands covering success and error cases.

## Phase 2: Skill File (Non-breaking)

### 2.1 Skill definition (`config/skills/task-scheduler.md` - new file)

```yaml
---
name: task-scheduler
description: Create, update, delete, pause, resume, and trigger scheduled tasks via clawbench task CLI.
  Triggers on: scheduled task, cron job, recurring, periodic, timer, 定时任务, 定期执行, 周期性, 计划任务
condition: scheduler.enabled
---
```

Skill body contains:
- Full CLI reference for all 6 subcommands with flags and examples
- Tag format: `<scheduled-task id="task-xxx" />`
- Rules: never use `<schedule-proposal>`, always use CLI, always include tag after creation
- Cron expression quick reference

### 2.2 Skill condition filtering (`internal/model/agent.go`)

`RemoveSkillsByCondition()` gains `scheduled bool` parameter:
- When `scheduled=true`, additionally removes `task-scheduler` skill
- Called from `buildCommonPrompt()` when `ScheduledExecution=true`
- AI in scheduled execution cannot discover the skill → fundamental block

## Phase 3: Remove Old System (Breaking)

### 3.1 Rules cleanup (`config/rules.md`)

Remove entire "Scheduled Tasks (Highest Priority)" section including:
- `<schedule-proposal>` format specification
- Forbidden behaviors list (CronCreate, crontab, systemctl, etc.)
- Cron expression examples

### 3.2 Handler cleanup (`internal/handler/chat.go`)

Remove:
- `scheduleProposalRe` regex constant
- `detectAndCreateScheduleProposals()` function
- `injectTaskIDsIntoProposals()` function
- Schedule-proposal detection block in `finalizeStreamRun()` (lines ~602-637)

Keep:
- `ChatRequest.ScheduledExecution` field (anti-recursion still needed)
- `finalizeStreamRun()` function itself (still handles other finalization)

### 3.3 Accumulator cleanup (`internal/ai/accumulate.go`)

Remove schedule-proposal related comments (lines 39-41).

## Phase 4: Frontend New Tag Rendering (Replaces Old Logic)

### 4.1 Tag extraction (`web/src/composables/useChatRender.ts`)

**Remove**:
- `blockProposals` reactive object and all related logic
- `extractScheduleProposals()` method
- Old `<schedule-proposal>` regex
- Watch on `store.state.tasks` for proposal sync

**Add**:
- `blockTasks` reactive object: `{ taskId, task, loading, deleted }`
- New regex: `/<scheduled-task\s+id="([^"]+)"\s*\/>/g`
- `extractScheduledTasks()` method: extract task IDs, strip tags from text
- Async API fetch: on first render of each tag, `GET /api/tasks/{id}`
- `blockTasks` watch on `store.state.tasks` for state updates (pause/resume/delete)

### 4.2 Card component (`web/src/components/chat/ContentBlocks.vue`)

**Remove**:
- Old proposal card template (`hasProposals`, `proposalKeys`, `blockProposals` references)
- `.schedule-proposal-card` styling

**Add**:
- New card template with real-time data from `blockTasks`:
  - Header: ⏰ icon + task name + edit button
  - Body: frequency, executor, repeat mode, status
  - Status line: running indicator + run count + last/next run time
  - Action buttons: pause/resume, trigger, delete
- Loading state: skeleton card while API fetch pending
- Deleted state: 🗑️ icon + "Task Deleted" (same UX as old)
- `.scheduled-task-card` styling

### 4.3 Event handling (`web/src/components/chat/ChatPanel.vue`)

**Remove**:
- Old `@edit-task="openTaskEdit"` handler
- `TaskFormDialog` import for proposal card editing

**Add**:
- Card inline action handlers: onPause, onResume, onTrigger, onDelete
- Each action calls the existing API endpoint directly
- Card re-fetches task data after action to reflect updated state
- Edit button still opens `TaskFormDialog` (keep existing dialog)

## Edge Cases

### AI creation failure
- CLI returns `{"ok":false,"error":"..."}`, AI sees error and can retry or inform user
- No `<scheduled-task />` tag rendered → no empty card

### Multiple tasks in one message
- AI calls `clawbench task create` multiple times, outputs multiple tags
- Frontend renders independent cards, each queries API separately

### AI updates/deletes existing tasks
- After `clawbench task update/pause/resume/delete`, existing cards auto-reflect changes on next refresh
- No need to output new tags — task_id unchanged, API returns updated data

### Historical messages with old tags
- Old `<schedule-proposal>` tags display as plain text (accepted)
- TaskDrawer still works for managing old tasks

### Task not found on API query
- `GET /api/tasks/{id}` returns 404 → card shows "已删除" state
- Same UX as old system but triggered by API response, not store polling miss

## Anti-Recursion Protection (Dual Layer)

1. **Skill layer**: `RemoveSkillsByCondition(scheduled=true)` removes `task-scheduler` from AI's skill table → AI cannot discover the capability
2. **CLI layer**: `clawbench task create` checks `CLAWBENCH_SCHEDULED=1` env var → hard block even if AI hardcodes the command

Both layers are needed: skill layer prevents accidental attempts, CLI layer prevents deliberate bypass.

## File Change Summary

| File | Action | Phase |
|------|--------|-------|
| `cmd/server/main.go` | Modify: add task dispatch | 1 |
| `internal/cli/task.go` | **New**: CLI subcommand implementation | 1 |
| `internal/cli/task_test.go` | **New**: CLI tests | 1 |
| `config/skills/task-scheduler.md` | **New**: Skill definition | 2 |
| `internal/model/agent.go` | Modify: `RemoveSkillsByCondition()` add `scheduled` param | 2 |
| `config/rules.md` | Modify: remove schedule-proposal section | 3 |
| `internal/handler/chat.go` | Modify: remove proposal detection/injection | 3 |
| `internal/ai/accumulate.go` | Modify: remove proposal comment | 3 |
| `web/src/composables/useChatRender.ts` | Modify: replace blockProposals with blockTasks | 4 |
| `web/src/components/chat/ContentBlocks.vue` | Modify: replace proposal card with task card | 4 |
| `web/src/components/chat/ChatPanel.vue` | Modify: replace edit-task with inline actions | 4 |
