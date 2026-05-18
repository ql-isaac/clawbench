import { ref } from 'vue'
import { getCachedCommitInfo } from '@/composables/useCommitHashAnnotation.ts'

// Module-level pending SHA — set by chat click, consumed by Git history components
const pendingSha = ref<string | null>(null)

/**
 * Set a pending commit navigation request.
 * Called from App.vue's handleNavigateToCommit.
 */
export function setPendingCommitNavigation(sha: string) {
    pendingSha.value = sha
}

/**
 * Check if there's a pending commit navigation and consume it.
 * Returns the SHA or null.
 */
export function consumePendingCommitNavigation(): string | null {
    const sha = pendingSha.value
    pendingSha.value = null
    return sha
}

/**
 * Check if there's a pending commit navigation without consuming it.
 */
export function hasPendingCommitNavigation(): boolean {
    return pendingSha.value !== null
}

/**
 * Shared commit navigation logic for GitHistory components.
 * Takes the component's reactive state and functions as parameters.
 */
export function useCommitNavigation(options: {
    commits: any              // ref([])
    selectedSHA: any          // ref(null)
    currentView: any          // ref('commits')
    loadCommitFiles: (sha: string) => Promise<void>
    loadProjectHistory?: () => Promise<void>
}) {
    const { commits, selectedSHA, currentView, loadCommitFiles, loadProjectHistory } = options

    /**
     * Fetch a single commit's info via verify-commits API.
     */
    async function fetchCommitInfo(sha: string): Promise<any | null> {
        try {
            const resp = await fetch('/api/git/verify-commits', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ shas: [sha] }),
            })
            if (!resp.ok) return null
            const data = await resp.json()
            const info = data.results?.[sha]
            return (info && info.sha) ? info : null
        } catch {
            return null
        }
    }

    /**
     * Navigate directly to a specific commit's files view.
     * Ensures the commit info is in the commits array so breadcrumbs work.
     */
    async function navigateToCommit(sha: string) {
        selectedSHA.value = sha
        currentView.value = 'files'

        // Ensure the commit exists in the commits array for selectedCommit computed
        const existing = commits.value.find(c => c.sha === sha)
        if (!existing) {
            // Try annotation cache first
            const info = getCachedCommitInfo(sha)
            if (info && info.sha) {
                commits.value.unshift(info)
            } else {
                // Fallback: fetch commit info via API
                const fetched = await fetchCommitInfo(sha)
                if (fetched && !commits.value.find(c => c.sha === sha)) {
                    commits.value.unshift(fetched)
                }
            }
        }

        loadCommitFiles(sha).catch(() => {})
    }

    /**
     * Handle drill-back to commits list when arriving from a deep-linked commit.
     * Loads the full project history if the commits array only has the one commit.
     */
    function handleDrillBackToCommits() {
        if (commits.value.length <= 1 && loadProjectHistory) {
            loadProjectHistory()
        }
    }

    return {
        navigateToCommit,
        handleDrillBackToCommits,
        fetchCommitInfo,
    }
}
