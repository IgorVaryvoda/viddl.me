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
          <span class="paste-icon"></span>
        </button>
      </div>
      <button
        @click="fetchVideoInfo"
        :disabled="loading || !url || !!urlError"
        class="get-btn"
        aria-label="Fetch video information"
      >
        {{ loading ? 'Loading...' : 'Get Video' }}
      </button>
    </div>

    <div v-if="urlError" class="url-hint">
      {{ urlError }}
    </div>

    <div v-if="error" class="error" role="alert">
      <span class="error-icon"></span>
      <div class="error-content">
        <p class="error-message">{{ error }}</p>
        <button v-if="canRetry" @click="retryLastAction" class="retry-btn" aria-label="Retry failed operation">
          Retry
        </button>
      </div>
    </div>

    <div v-if="downloadProgress" class="progress-indicator" role="status" aria-live="polite">
      <div class="spinner" aria-hidden="true"></div>
      <p>{{ downloadProgress }}</p>
      <p v-if="estimatedTimeRemaining" class="time-remaining">{{ estimatedTimeRemaining }}</p>
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
        <h4>Select Video to Download</h4>
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
            <p>{{ videoInfo.uploader }}</p>
            <p>{{ formatDuration(videoInfo.duration) }}</p>
          </div>
        </div>

        <div v-if="videoInfo.formats && videoInfo.formats.length > 0" class="formats">
          <h4 id="quality-heading">Select Quality (MP4)</h4>
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
          {{ downloading ? 'Downloading...' : 'Download Best Quality' }}
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
  'reddit.com', 'twitch.tv', 'threads.net', 'sirv.com'
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
  max-width: 600px;
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
  padding: 1rem 3rem 1rem 1rem;
  border: 2px solid var(--border);
  border-radius: 12px;
  background: var(--bg-secondary);
  color: var(--text);
  font-size: 1rem;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.input-wrapper input:focus {
  outline: none;
  border-color: var(--accent);
  box-shadow: 0 0 0 3px rgba(255, 107, 53, 0.2);
}

.paste-btn {
  position: absolute;
  right: 0.75rem;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  cursor: pointer;
  padding: 0.5rem;
  opacity: 0.6;
  transition: opacity 0.2s;
  font-size: 1.2rem;
}

.paste-btn:hover {
  opacity: 1;
}

.paste-icon::before {
  content: '📋';
}

.get-btn {
  padding: 1rem 1.5rem;
  background: var(--accent);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s, transform 0.1s;
  white-space: nowrap;
}

.get-btn:hover:not(:disabled) {
  background: var(--accent-hover);
  transform: translateY(-1px);
}

.get-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.url-hint {
  color: var(--accent);
  font-size: 0.875rem;
  margin-bottom: 1rem;
  text-align: center;
}

.error {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  padding: 1rem;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 12px;
  margin-bottom: 1rem;
}

.error-icon::before {
  content: '⚠️';
}

.error-content {
  flex: 1;
}

.error-message {
  color: #ef4444;
  margin: 0;
}

.retry-btn {
  margin-top: 0.5rem;
  padding: 0.5rem 1rem;
  background: rgba(239, 68, 68, 0.2);
  border: 1px solid rgba(239, 68, 68, 0.4);
  border-radius: 8px;
  color: #ef4444;
  cursor: pointer;
  font-size: 0.875rem;
}

.retry-btn:hover {
  background: rgba(239, 68, 68, 0.3);
}

.progress-indicator {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
  padding: 1.5rem;
  background: var(--bg-secondary);
  border-radius: 12px;
  margin-bottom: 1rem;
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--border);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.time-remaining {
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.loader {
  padding: 1.5rem;
  background: var(--bg-secondary);
  border-radius: 12px;
  margin-bottom: 1rem;
}

.skeleton-loader {
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.skeleton-header {
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
}

.skeleton-thumbnail {
  width: 120px;
  height: 68px;
  background: var(--border);
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
  background: var(--border);
  border-radius: 4px;
  width: 80%;
}

.skeleton-subtitle {
  height: 1rem;
  background: var(--border);
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
  height: 3rem;
  background: var(--border);
  border-radius: 8px;
}

.loading-text {
  text-align: center;
  color: var(--text-secondary);
  margin-top: 1rem;
}

.video-info {
  background: var(--bg-secondary);
  border-radius: 12px;
  padding: 1.5rem;
}

.video-header {
  display: flex;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.thumbnail {
  width: 160px;
  height: 90px;
  object-fit: cover;
  border-radius: 8px;
}

.video-meta {
  flex: 1;
}

.video-meta h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1rem;
  line-height: 1.4;
}

.video-meta p {
  margin: 0;
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.formats h4 {
  margin: 0 0 1rem 0;
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.format-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 0.75rem;
}

.format-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.25rem;
  padding: 0.75rem;
  background: var(--bg);
  border: 2px solid var(--border);
  border-radius: 8px;
  cursor: pointer;
  transition: border-color 0.2s, background 0.2s;
}

.format-btn:hover:not(:disabled) {
  border-color: var(--accent);
  background: rgba(255, 107, 53, 0.1);
}

.format-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.format-quality {
  font-weight: 600;
  color: var(--text);
}

.format-size {
  font-size: 0.75rem;
  color: var(--text-secondary);
}

.download-best-btn {
  width: 100%;
  padding: 1rem;
  background: var(--accent);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}

.download-best-btn:hover:not(:disabled) {
  background: var(--accent-hover);
}

.download-best-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.multi-video-selection h4 {
  margin: 0 0 1rem 0;
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
  padding: 0.75rem;
  background: var(--bg);
  border: 2px solid var(--border);
  border-radius: 8px;
  cursor: pointer;
  text-align: left;
  transition: border-color 0.2s;
}

.video-item-btn:hover:not(:disabled) {
  border-color: var(--accent);
}

.video-item-thumbnail {
  width: 80px;
  height: 45px;
  object-fit: cover;
  border-radius: 4px;
}

.video-item-meta {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.video-item-number {
  font-size: 0.75rem;
  color: var(--text-secondary);
}

.video-item-title {
  font-weight: 500;
  color: var(--text);
}

.video-item-duration {
  font-size: 0.875rem;
  color: var(--text-secondary);
}

@media (max-width: 640px) {
  .input-group {
    flex-direction: column;
  }

  .get-btn {
    width: 100%;
  }

  .video-header {
    flex-direction: column;
  }

  .thumbnail {
    width: 100%;
    height: auto;
    aspect-ratio: 16/9;
  }
}
</style>
