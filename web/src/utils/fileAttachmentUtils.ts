/**
 * Pure functions extracted from FileAttachmentList.vue for testability.
 */

/** Normalize a file entry to { path: string } format.
 *  Backend returns string[], local push uses [{path: "..."}]. */
export function normalizeFileEntry(f: string | { path: string }): { path: string } {
  if (typeof f === 'string') return { path: f }
  return { path: f.path || '' }
}

/** Check if a path points to an uploaded file (in .clawbench/uploads/). */
export function isUploadPath(path: string): boolean {
  return path.startsWith('.clawbench/uploads/') || path.startsWith('.clawbench\\uploads\\')
}

/** Common image file extensions. */
const IMAGE_EXTENSIONS = ['.png', '.jpg', '.jpeg', '.gif', '.webp', '.svg', '.bmp', '.ico', '.tiff', '.tif', '.avif']

/** Check if a path points to an image file based on its extension. */
export function isImageFile(path: string | null | undefined): boolean {
  if (!path) return false
  const lower = path.toLowerCase()
  return IMAGE_EXTENSIONS.some(ext => lower.endsWith(ext))
}
