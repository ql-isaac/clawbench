// API utility functions
import i18n from '@/i18n'

function localeHeaders(): Record<string, string> {
    return { 'X-Locale': i18n.global.locale.value as string }
}

// Default timeout for API requests (10 seconds)
const API_TIMEOUT_MS = 10_000

export async function apiGet<T = unknown>(url: string): Promise<T> {
    const controller = new AbortController()
    const timer = setTimeout(() => controller.abort(), API_TIMEOUT_MS)
    try {
        const resp = await fetch(url, { headers: localeHeaders(), signal: controller.signal })
        if (!resp.ok) throw new Error(await resp.text())
        return resp.json()
    } finally {
        clearTimeout(timer)
    }
}

export async function apiPost<T = unknown>(url: string, body: unknown): Promise<T> {
    const controller = new AbortController()
    const timer = setTimeout(() => controller.abort(), API_TIMEOUT_MS)
    try {
        const resp = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json', ...localeHeaders() },
            body: JSON.stringify(body),
            signal: controller.signal,
        })
        const data = await resp.json().catch(() => ({})) as Record<string, unknown>
        if (!resp.ok) throw new Error(data.error ? String(data.error) : resp.statusText)
        return data as T
    } finally {
        clearTimeout(timer)
    }
}

export async function apiPut<T = unknown>(url: string, body: unknown): Promise<T> {
    const controller = new AbortController()
    const timer = setTimeout(() => controller.abort(), API_TIMEOUT_MS)
    try {
        const resp = await fetch(url, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json', ...localeHeaders() },
            body: JSON.stringify(body),
            signal: controller.signal,
        })
        const data = await resp.json().catch(() => ({})) as Record<string, unknown>
        if (!resp.ok) throw new Error(data.error ? String(data.error) : resp.statusText)
        return data as T
    } finally {
        clearTimeout(timer)
    }
}

export async function apiDelete<T = unknown>(url: string): Promise<T> {
    const controller = new AbortController()
    const timer = setTimeout(() => controller.abort(), API_TIMEOUT_MS)
    try {
        const resp = await fetch(url, { method: 'DELETE', headers: localeHeaders(), signal: controller.signal })
        if (!resp.ok) throw new Error(resp.statusText)
        return resp.json()
    } finally {
        clearTimeout(timer)
    }
}

export async function cancelChat(sessionId: string): Promise<void> {
    const controller = new AbortController()
    const timer = setTimeout(() => controller.abort(), API_TIMEOUT_MS)
    try {
        const resp = await fetch(`/api/ai/chat/cancel?session_id=${encodeURIComponent(sessionId)}`, {
            method: 'POST',
            headers: localeHeaders(),
            signal: controller.signal,
        })
        if (!resp.ok) throw new Error(resp.statusText)
    } finally {
        clearTimeout(timer)
    }
}
