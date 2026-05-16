// Cross-platform path utilities

// Split a path into segments, handling both / and \ separators
export function splitPath(path: string): string[] {
    return path.split(/[/\\]/)
}

// Get the last segment of a path (filename or directory name)
export function baseName(path: string): string {
    const segments = splitPath(path)
    // Walk backwards to find the last non-empty segment
    // This handles trailing slashes correctly: /home/user/ → "user"
    for (let i = segments.length - 1; i >= 0; i--) {
        if (segments[i] !== '') return segments[i]
    }
    return path
}

// Get the parent directory of a path
export function dirName(path: string): string {
    const parts = splitPath(path)
    parts.pop()
    if (parts.length === 0) return ''
    // Rejoin with original separator style
    const useBackslash = path.includes('\\') && !path.includes('/')
    const result = useBackslash ? parts.join('\\') : parts.join('/')
    // On Windows, a lone "C:" should be "C:\" (drive root)
    if (/^[A-Za-z]:$/.test(result)) return result + '\\'
    return result
}

/**
 * Convert an absolute path to a relative path based on a base path.
 * Returns the original path if base is empty or absPath does not start with basePath.
 * Returns '/' if the result would be empty (i.e., the path equals the base).
 */
export function toRelativePath(absPath: string, basePath: string): string {
    if (!basePath) return absPath
    if (!absPath.startsWith(basePath)) return absPath
    const rel = absPath.slice(basePath.length).replace(/^\//, '')
    return rel || '/'
}
