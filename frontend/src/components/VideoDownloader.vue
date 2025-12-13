<template>
  <div class="downloader">
    <div class="input-group">
      <div class="input-wrapper">
        <input
          ref="urlInput"
          v-model="url"
          type="url"
          :placeholder="placeholder"
          @keyup.enter="fetchVideoInfo"
          @input="validateUrlInput"
          :disabled="loading"
          aria-label="Video URL input"
          :aria-invalid="urlError ? 'true' : 'false'"
        />
        <button
          @click="pasteFromClipboard"
          class="paste-btn"
          :disabled="loading"
          aria-label="Paste URL from clipboard"
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
            <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
          </svg>
        </button>
      </div>
      <button
        @click="fetchVideoInfo"
        :disabled="loading || !url || !!urlError"
        class="get-btn"
        aria-label="Fetch video information"
      >
        <span v-if="loading" class="btn-loading">
          <svg class="spinner-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
          </svg>
          Loading
        </span>
        <span v-else class="btn-content">
          Get Video
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M5 12h14M12 5l7 7-7 7"/>
          </svg>
        </span>
      </button>
    </div>

    <div v-if="urlError" class="url-hint">
      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="10"/>
        <path d="M12 8v4M12 16h.01"/>
      </svg>
      {{ urlError }}
    </div>

    <div v-if="error" class="error" role="alert">
      <div class="error-icon">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <path d="M15 9l-6 6M9 9l6 6"/>
        </svg>
      </div>
      <div class="error-content">
        <p class="error-message">{{ error }}</p>
        <button v-if="canRetry" @click="retryLastAction" class="retry-btn" aria-label="Retry failed operation">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M1 4v6h6M23 20v-6h-6"/>
            <path d="M20.49 9A9 9 0 0 0 5.64 5.64L1 10m22 4l-4.64 4.36A9 9 0 0 1 3.51 15"/>
          </svg>
          Retry
        </button>
      </div>
    </div>

    <div v-if="downloadProgress" class="progress-indicator" role="status" aria-live="polite">
      <div class="progress-spinner">
        <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
        </svg>
      </div>
      <div class="progress-text">
        <p>{{ downloadProgress }}</p>
        <p v-if="estimatedTimeRemaining" class="time-remaining">{{ estimatedTimeRemaining }}</p>
      </div>
    </div>

    <div v-if="loading" class="loader" role="status" aria-live="polite">
      <div class="skeleton-loader">
        <div class="skeleton-header">
          <div class="skeleton-thumbnail"></div>
          <div class="skeleton-meta">
            <div class="skeleton-title"></div>
            <div class="skeleton-subtitle"></div>
            <div class="skeleton-subtitle short"></div>
          </div>
        </div>
        <div class="skeleton-formats">
          <div class="skeleton-format"></div>
          <div class="skeleton-format"></div>
          <div class="skeleton-format"></div>
        </div>
      </div>
      <p class="loading-text">{{ loadingMessage }}</p>
    </div>

    <div v-if="videoInfo && !loading" class="video-info">
      <div v-if="videoInfo.is_multi_video" class="multi-video-selection">
        <h4>Select video to download</h4>
        <div class="video-list" role="list">
          <button
            v-for="video in videoInfo.multi_videos"
            :key="video.index"
            @click="selectVideo(video.index)"
            :disabled="downloading"
            class="video-item-btn"
            role="listitem"
          >
            <img
              v-if="video.thumbnail"
              :src="video.thumbnail"
              :alt="`Thumbnail for video ${video.index}`"
              class="video-item-thumbnail"
            />
            <div class="video-item-meta">
              <span class="video-item-number">Video {{ video.index }}</span>
              <span class="video-item-title">{{ video.title || 'Untitled' }}</span>
              <span class="video-item-duration">{{ formatDuration(video.duration) }}</span>
            </div>
            <svg class="video-item-arrow" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M5 12h14M12 5l7 7-7 7"/>
            </svg>
          </button>
        </div>
      </div>

      <div v-else>
        <div class="video-header">
          <img
            v-if="videoInfo.thumbnail"
            :src="videoInfo.thumbnail"
            :alt="`Thumbnail for ${videoInfo.title}`"
            class="thumbnail"
          />
          <div class="video-meta">
            <h3>{{ videoInfo.title }}</h3>
            <p class="video-uploader">{{ videoInfo.uploader }}</p>
            <p class="video-duration">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <path d="M12 6v6l4 2"/>
              </svg>
              {{ formatDuration(videoInfo.duration) }}
            </p>
          </div>
        </div>

        <div v-if="videoInfo.formats && videoInfo.formats.length > 0" class="formats">
          <h4 id="quality-heading">Select quality</h4>
          <div class="format-grid" role="list" aria-labelledby="quality-heading">
            <button
              v-for="format in videoInfo.formats"
              :key="format.format_id"
              @click="downloadVideo(format.format_id)"
              :disabled="downloading"
              class="format-btn"
              role="listitem"
            >
              <span class="format-quality">{{ format.quality }}</span>
              <span class="format-size">{{ formatSize(format.filesize) }}</span>
              <svg class="format-download-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                <polyline points="7 10 12 15 17 10"/>
                <line x1="12" y1="15" x2="12" y2="3"/>
              </svg>
            </button>
          </div>
        </div>

        <button
          v-else
          @click="downloadVideo('best')"
          :disabled="downloading"
          class="download-best-btn"
          aria-label="Download video in best available quality"
        >
          <span v-if="downloading" class="btn-loading">
            <svg class="spinner-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
            </svg>
            Downloading...
          </span>
          <span v-else class="btn-content">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
              <polyline points="7 10 12 15 17 10"/>
              <line x1="12" y1="15" x2="12" y2="3"/>
            </svg>
            Download Best Quality
          </span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const props = defineProps({
  placeholder: {
    type: String,
    default: 'Paste video URL here...'
  },
  filterDomain: {
    type: String,
    default: null
  }
})

const emit = defineEmits(['download-complete'])

const url = ref('')
const urlInput = ref(null)
const videoInfo = ref(null)
const error = ref('')
const urlError = ref('')
const loading = ref(false)
const downloading = ref(false)
const downloadProgress = ref('')
const selectedVideoIndex = ref(0)
const loadingStartTime = ref(0)
const loadingMessage = ref('Fetching video information...')
const estimatedTimeRemaining = ref('')
const downloadStartTime = ref(0)

const canRetry = ref(false)
const lastAction = ref(null)
const retryCount = ref(0)
const MAX_RETRIES = 2

const videoInfoCache = ref(new Map())
const CACHE_DURATION = 5 * 60 * 1000

const API_URL = import.meta.env.VITE_API_URL || '/api'

const SUPPORTED_DOMAINS = [
  'youtube.com', 'youtu.be', 'twitter.com', 'x.com',
  'instagram.com', 'facebook.com', 'vimeo.com',
  'reddit.com', 'twitch.tv', 'threads.com', 'sirv.com'
]

const ERROR_MESSAGES = {
  'domain not allowed': 'This website is not supported. Try YouTube, Twitter, Instagram, or Facebook.',
  'invalid URL format': 'Please enter a valid URL starting with http:// or https://',
  'URL too long': 'The URL is too long. Please check and try again.',
  'invalid protocol': 'Only http:// and https:// URLs are supported.',
  'Failed to fetch video information': 'Could not retrieve video info. The video may be private, deleted, or region-locked.',
  'Download failed or file exceeds size limit': 'Download failed. The file may be too large (max 2GB) or temporarily unavailable.',
  'Too many requests': 'Too many requests. Please wait a moment before trying again.',
  'rate limit': 'You\'re making requests too quickly. Please wait a minute.',
}

onMounted(() => {
  if (urlInput.value) {
    urlInput.value.focus()
  }
})

const validateUrlInput = () => {
  if (!url.value) {
    urlError.value = ''
    return
  }

  try {
    const parsed = new URL(url.value)
    if (parsed.protocol !== 'http:' && parsed.protocol !== 'https:') {
      urlError.value = 'URL must start with http:// or https://'
      return
    }

    const hostname = parsed.hostname.replace(/^www\./, '')
    const isSupported = SUPPORTED_DOMAINS.some(domain =>
      hostname === domain || hostname.endsWith('.' + domain)
    )

    if (!isSupported) {
      urlError.value = `Unsupported site. Try: ${SUPPORTED_DOMAINS.slice(0, 4).join(', ')}...`
      return
    }

    if (props.filterDomain) {
      const matchesDomain = hostname === props.filterDomain ||
        hostname.endsWith('.' + props.filterDomain) ||
        (props.filterDomain === 'youtube.com' && hostname === 'youtu.be') ||
        (props.filterDomain === 'twitter.com' && hostname === 'x.com')

      if (!matchesDomain) {
        urlError.value = `Please enter a ${props.filterDomain} URL`
        return
      }
    }

    urlError.value = ''
  } catch {
    if (url.value.length > 5) {
      urlError.value = 'Please enter a valid URL'
    } else {
      urlError.value = ''
    }
  }
}

const getReadableError = (errorMsg) => {
  if (!errorMsg) return 'An unexpected error occurred. Please try again.'

  const lowerError = errorMsg.toLowerCase()

  for (const [key, message] of Object.entries(ERROR_MESSAGES)) {
    if (lowerError.includes(key.toLowerCase())) {
      return message
    }
  }

  if (lowerError.includes('private') || lowerError.includes('unavailable')) {
    return 'This video is private or unavailable.'
  }
  if (lowerError.includes('age') || lowerError.includes('sign in')) {
    return 'This video requires age verification or login.'
  }
  if (lowerError.includes('copyright') || lowerError.includes('blocked')) {
    return 'This video is blocked due to copyright restrictions.'
  }
  if (lowerError.includes('network') || lowerError.includes('timeout')) {
    return 'Network error. Please check your connection and try again.'
  }

  return errorMsg
}

const getCachedVideoInfo = (videoUrl) => {
  const cached = videoInfoCache.value.get(videoUrl)
  if (cached && Date.now() - cached.timestamp < CACHE_DURATION) {
    return cached.data
  }
  videoInfoCache.value.delete(videoUrl)
  return null
}

const setCachedVideoInfo = (videoUrl, data) => {
  videoInfoCache.value.set(videoUrl, {
    data,
    timestamp: Date.now()
  })
}

let loadingMessageInterval = null

const updateLoadingMessage = () => {
  const elapsed = Date.now() - loadingStartTime.value
  if (elapsed > 15000) {
    loadingMessage.value = 'Still working... Some videos take longer to process'
  } else if (elapsed > 8000) {
    loadingMessage.value = 'Processing video information...'
  } else if (elapsed > 3000) {
    loadingMessage.value = 'Connecting to video source...'
  }
}

const fetchVideoInfo = async () => {
  if (!url.value || urlError.value) return

  const cached = getCachedVideoInfo(url.value)
  if (cached) {
    videoInfo.value = cached
    return
  }

  error.value = ''
  loading.value = true
  videoInfo.value = null
  selectedVideoIndex.value = 0
  loadingStartTime.value = Date.now()
  loadingMessage.value = 'Fetching video information...'

  loadingMessageInterval = setInterval(updateLoadingMessage, 1000)

  lastAction.value = { type: 'fetch', url: url.value }
  canRetry.value = false

  try {
    const response = await axios.post(`${API_URL}/info`, {
      url: url.value
    }, {
      timeout: 60000
    })
    videoInfo.value = response.data
    setCachedVideoInfo(url.value, response.data)
    retryCount.value = 0
  } catch (err) {
    const rawError = err.response?.data?.error || err.message || 'Failed to fetch video information'
    error.value = getReadableError(rawError)
    canRetry.value = retryCount.value < MAX_RETRIES

    if (err.code === 'ECONNABORTED' || err.message?.includes('timeout')) {
      error.value = 'Request timed out. The video source may be slow. Please try again.'
      canRetry.value = true
    }
  } finally {
    loading.value = false
    clearInterval(loadingMessageInterval)
  }
}

const selectVideo = async (index) => {
  selectedVideoIndex.value = index
  downloadVideo('best', index)
}

const calculateETA = (loaded, total, startTime) => {
  const elapsed = (Date.now() - startTime) / 1000
  if (elapsed < 2 || loaded === 0) return ''

  const speed = loaded / elapsed
  const remaining = total - loaded
  const etaSeconds = Math.round(remaining / speed)

  if (etaSeconds < 60) {
    return `~${etaSeconds}s remaining`
  } else {
    const mins = Math.floor(etaSeconds / 60)
    const secs = etaSeconds % 60
    return `~${mins}m ${secs}s remaining`
  }
}

const downloadVideo = async (format, videoIndex = 0) => {
  downloading.value = true
  error.value = ''
  downloadProgress.value = 'Preparing download...'
  estimatedTimeRemaining.value = ''
  downloadStartTime.value = Date.now()

  lastAction.value = { type: 'download', format, videoIndex }
  canRetry.value = false

  try {
    downloadProgress.value = 'Fetching video from server...'

    const requestData = {
      url: url.value,
      format: format
    }

    if (videoIndex > 0) {
      requestData.video_index = videoIndex
    }

    const response = await axios.post(
      `${API_URL}/download`,
      requestData,
      {
        responseType: 'blob',
        timeout: 600000,
        onDownloadProgress: (progressEvent) => {
          if (progressEvent.total) {
            const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total)
            downloadProgress.value = `Downloading: ${percentCompleted}%`
            estimatedTimeRemaining.value = calculateETA(
              progressEvent.loaded,
              progressEvent.total,
              downloadStartTime.value
            )
          } else {
            downloadProgress.value = `Downloading: ${(progressEvent.loaded / 1024 / 1024).toFixed(1)} MB`
            estimatedTimeRemaining.value = ''
          }
        }
      }
    )

    downloadProgress.value = 'Processing file...'
    estimatedTimeRemaining.value = ''

    const blob = new Blob([response.data])
    const downloadUrl = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = downloadUrl

    const contentDisposition = response.headers['content-disposition']
    let filename = 'video'
    if (contentDisposition) {
      const filenameMatch = contentDisposition.match(/filename=(.+)/)
      if (filenameMatch) {
        filename = filenameMatch[1]
      }
    }

    link.download = filename
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(downloadUrl)

    downloadProgress.value = 'Download complete!'
    retryCount.value = 0

    emit('download-complete', {
      url: url.value,
      title: videoInfo.value?.title,
      thumbnail: videoInfo.value?.thumbnail
    })

    setTimeout(() => {
      downloadProgress.value = ''
    }, 2000)
  } catch (err) {
    const rawError = err.response?.data?.error || err.message || 'Download failed'
    error.value = getReadableError(rawError)
    downloadProgress.value = ''
    estimatedTimeRemaining.value = ''
    canRetry.value = retryCount.value < MAX_RETRIES

    if (err.code === 'ECONNABORTED' || err.message?.includes('timeout')) {
      error.value = 'Download timed out. The file may be too large or the server is slow.'
      canRetry.value = true
    }
  } finally {
    downloading.value = false
  }
}

const retryLastAction = async () => {
  if (!lastAction.value || retryCount.value >= MAX_RETRIES) return

  retryCount.value++
  error.value = ''
  canRetry.value = false

  if (lastAction.value.type === 'fetch') {
    await fetchVideoInfo()
  } else if (lastAction.value.type === 'download') {
    await downloadVideo(lastAction.value.format, lastAction.value.videoIndex)
  }
}

const pasteFromClipboard = async () => {
  try {
    const text = await navigator.clipboard.readText()
    url.value = text.trim()
    validateUrlInput()
  } catch (err) {
    error.value = 'Failed to read clipboard. Please paste manually or grant clipboard permission.'
    canRetry.value = false
    setTimeout(() => {
      error.value = ''
    }, 3000)
  }
}

const formatDuration = (seconds) => {
  if (!seconds) return 'Unknown duration'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

const formatSize = (bytes) => {
  if (!bytes) return 'Unknown size'
  const mb = bytes / (1024 * 1024)
  if (mb >= 1) {
    return `${mb.toFixed(1)} MB`
  }
  return `${(bytes / 1024).toFixed(1)} KB`
}
</script>

<style scoped>
.downloader {
  width: 100%;
  max-width: 640px;
  margin: 0 auto;
}

.input-group {
  display: flex;
  gap: 0.75rem;
  margin-bottom: 1rem;
}

.input-wrapper {
  flex: 1;
  position: relative;
}

.input-wrapper input {
  width: 100%;
  padding: 1rem 3.5rem 1rem 1.25rem;
  border: 2px solid var(--border);
  border-radius: 12px;
  background: var(--bg-secondary);
  color: var(--text);
  font-family: var(--font-body);
  font-size: 1rem;
  transition: all 0.2s;
}

.input-wrapper input::placeholder {
  color: var(--text-secondary);
}

.input-wrapper input:focus {
  outline: none;
  border-color: var(--accent);
  box-shadow: 0 0 0 4px var(--accent-glow), 0 0 30px var(--accent-glow);
}

.paste-btn {
  position: absolute;
  right: 0.75rem;
  top: 50%;
  transform: translateY(-50%);
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 0.5rem;
  color: var(--text-secondary);
  transition: all 0.2s;
  border-radius: 6px;
}

.paste-btn:hover {
  color: var(--accent);
  background: var(--bg-tertiary);
}

.get-btn {
  padding: 1rem 1.5rem;
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-secondary) 100%);
  color: #fff;
  border: none;
  border-radius: 12px;
  font-family: var(--font-display);
  font-size: 0.9375rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.get-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 10px 30px -10px var(--accent-glow);
}

.get-btn:active:not(:disabled) {
  transform: translateY(0);
}

.get-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-content,
.btn-loading {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.spinner-icon {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.url-hint {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--accent-secondary);
  font-size: 0.875rem;
  margin-bottom: 1rem;
  padding: 0 0.25rem;
}

.error {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  padding: 1rem 1.25rem;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.2);
  border-radius: 12px;
  margin-bottom: 1rem;
}

.error-icon {
  color: var(--error);
  flex-shrink: 0;
  margin-top: 0.125rem;
}

.error-content {
  flex: 1;
}

.error-message {
  color: var(--error);
  margin: 0;
  font-size: 0.9375rem;
}

.retry-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  margin-top: 0.75rem;
  padding: 0.5rem 0.875rem;
  background: transparent;
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 6px;
  color: var(--error);
  cursor: pointer;
  font-size: 0.8125rem;
  font-weight: 500;
  transition: all 0.2s;
}

.retry-btn:hover {
  background: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.5);
}

.progress-indicator {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.25rem 1.5rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 12px;
  margin-bottom: 1rem;
}

.progress-spinner svg {
  animation: spin 1s linear infinite;
  color: var(--accent);
}

.progress-text {
  flex: 1;
}

.progress-text p {
  margin: 0;
  font-weight: 500;
}

.time-remaining {
  font-size: 0.875rem;
  color: var(--text-secondary);
  font-weight: 400 !important;
  margin-top: 0.25rem !important;
}

.loader {
  padding: 1.5rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 12px;
  margin-bottom: 1rem;
}

.skeleton-loader {
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.skeleton-header {
  display: flex;
  gap: 1rem;
  margin-bottom: 1.25rem;
}

.skeleton-thumbnail {
  width: 140px;
  height: 80px;
  background: var(--bg-tertiary);
  border-radius: 8px;
}

.skeleton-meta {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.skeleton-title {
  height: 1.25rem;
  background: var(--bg-tertiary);
  border-radius: 4px;
  width: 85%;
}

.skeleton-subtitle {
  height: 1rem;
  background: var(--bg-tertiary);
  border-radius: 4px;
  width: 60%;
}

.skeleton-subtitle.short {
  width: 40%;
}

.skeleton-formats {
  display: flex;
  gap: 0.75rem;
}

.skeleton-format {
  flex: 1;
  height: 3.5rem;
  background: var(--bg-tertiary);
  border-radius: 8px;
}

.loading-text {
  text-align: center;
  color: var(--text-secondary);
  margin: 1rem 0 0;
  font-size: 0.9375rem;
}

.video-info {
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 1.5rem;
  animation: slideUp 0.3s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.video-header {
  display: flex;
  gap: 1.25rem;
  margin-bottom: 1.5rem;
}

.thumbnail {
  width: 180px;
  height: 100px;
  object-fit: cover;
  border-radius: 10px;
  flex-shrink: 0;
}

.video-meta {
  flex: 1;
  min-width: 0;
}

.video-meta h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.0625rem;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.video-uploader {
  margin: 0 0 0.375rem 0;
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.video-duration {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  margin: 0;
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.formats h4 {
  margin: 0 0 1rem 0;
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.format-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(110px, 1fr));
  gap: 0.75rem;
}

.format-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.375rem;
  padding: 1rem 0.75rem;
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
  overflow: hidden;
}

.format-btn::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-secondary) 100%);
  opacity: 0;
  transition: opacity 0.2s;
}

.format-btn:hover:not(:disabled) {
  border-color: var(--accent);
  transform: translateY(-2px);
}

.format-btn:hover:not(:disabled)::before {
  opacity: 0.1;
}

.format-btn:hover:not(:disabled) .format-download-icon {
  opacity: 1;
  transform: translateY(0);
}

.format-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.format-quality {
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 1rem;
  color: var(--text);
  position: relative;
}

.format-size {
  font-size: 0.75rem;
  color: var(--text-secondary);
  position: relative;
}

.format-download-icon {
  position: absolute;
  bottom: 0.5rem;
  right: 0.5rem;
  opacity: 0;
  transform: translateY(4px);
  transition: all 0.2s;
  color: var(--accent);
}

.download-best-btn {
  width: 100%;
  padding: 1.125rem;
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-secondary) 100%);
  color: #fff;
  border: none;
  border-radius: 10px;
  font-family: var(--font-display);
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.download-best-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 10px 30px -10px var(--accent-glow);
}

.download-best-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.multi-video-selection h4 {
  margin: 0 0 1rem 0;
  font-size: 1rem;
}

.video-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.video-item-btn {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.875rem;
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 10px;
  cursor: pointer;
  text-align: left;
  transition: all 0.2s;
}

.video-item-btn:hover:not(:disabled) {
  border-color: var(--accent);
  transform: translateX(4px);
}

.video-item-btn:hover:not(:disabled) .video-item-arrow {
  opacity: 1;
  color: var(--accent);
}

.video-item-thumbnail {
  width: 80px;
  height: 45px;
  object-fit: cover;
  border-radius: 6px;
  flex-shrink: 0;
}

.video-item-meta {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
  min-width: 0;
}

.video-item-number {
  font-size: 0.75rem;
  color: var(--text-secondary);
}

.video-item-title {
  font-weight: 500;
  color: var(--text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.video-item-duration {
  font-size: 0.8125rem;
  color: var(--text-secondary);
}

.video-item-arrow {
  flex-shrink: 0;
  opacity: 0;
  color: var(--text-secondary);
  transition: all 0.2s;
}

@media (max-width: 640px) {
  .input-group {
    flex-direction: column;
  }

  .get-btn {
    width: 100%;
    justify-content: center;
  }

  .video-header {
    flex-direction: column;
    gap: 1rem;
  }

  .thumbnail {
    width: 100%;
    height: auto;
    aspect-ratio: 16/9;
  }

  .format-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}
</style>
