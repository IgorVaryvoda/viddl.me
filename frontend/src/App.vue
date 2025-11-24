<template>
  <div class="container">
    <header :class="{ 'header-minimal': videoInfo }">
      <h1>viddl.me</h1>
      <p class="subtitle" v-if="!videoInfo">Fast & free video downloader</p>

      <div class="supported-sites-top" v-if="!videoInfo">
        <span class="site-tag">YouTube</span>
        <span class="site-tag">Twitter/X</span>
        <span class="site-tag">Instagram</span>
        <span class="site-tag">Facebook</span>
        <span class="site-tag">+4 more</span>
      </div>
    </header>

    <div class="input-group">
      <div class="input-wrapper">
        <input
          ref="urlInput"
          v-model="url"
          type="url"
          placeholder="Paste video URL here..."
          @keyup.enter="fetchVideoInfo"
          @input="validateUrlInput"
          :disabled="loading"
          aria-label="Video URL input"
          aria-describedby="url-hint"
          :aria-invalid="urlError ? 'true' : 'false'"
        />
        <span id="url-hint" class="sr-only">Enter a video URL from YouTube, Twitter, Instagram, or other supported sites</span>
        <button
          @click="pasteFromClipboard"
          class="paste-btn"
          :disabled="loading"
          aria-label="Paste URL from clipboard"
        >
          📋
        </button>
      </div>
      <button
        @click="fetchVideoInfo"
        :disabled="loading || !url || !!urlError"
        aria-label="Fetch video information"
      >
        {{ loading ? 'Loading...' : 'Get Video' }}
      </button>
    </div>

    <div v-if="urlError" class="url-hint">
      {{ urlError }}
    </div>

    <div v-if="error" class="error" role="alert">
      <span class="error-icon">⚠️</span>
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
      <!-- Multi-video selection -->
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
            :aria-label="`Download video ${video.index}: ${video.title || 'Untitled'}`"
          >
            <img
              v-if="video.thumbnail"
              :src="video.thumbnail"
              :alt="`Thumbnail for video ${video.index}: ${video.title || 'Untitled'}`"
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

      <!-- Single video display -->
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
              :aria-label="`Download ${format.quality} quality, ${formatSize(format.filesize)}`"
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

    <footer class="footer">
      <p>Free & open source • <a href="https://github.com/IgorVaryvoda/viddl.me" target="_blank" rel="noopener">GitHub</a></p>
    </footer>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'

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
const downloadedBytes = ref(0)
const totalBytes = ref(0)

const canRetry = ref(false)
const lastAction = ref(null)
const retryCount = ref(0)
const MAX_RETRIES = 2

const videoInfoCache = ref(new Map())
const CACHE_DURATION = 5 * 60 * 1000

const API_URL = import.meta.env.VITE_API_URL || '/api'

const SUPPORTED_DOMAINS = [
  'youtube.com', 'youtu.be', 'twitter.com', 'x.com',
  'instagram.com', 'facebook.com', 'tiktok.com',
  'vimeo.com', 'reddit.com', 'twitch.tv'
]

const ERROR_MESSAGES = {
  'domain not allowed': 'This website is not supported. Try YouTube, Twitter, Instagram, TikTok, or Facebook.',
  'invalid URL format': 'Please enter a valid URL starting with http:// or https://',
  'URL too long': 'The URL is too long. Please check and try again.',
  'invalid protocol': 'Only http:// and https:// URLs are supported.',
  'Failed to fetch video information': 'Could not retrieve video info. The video may be private, deleted, or region-locked.',
  'Download failed or file exceeds size limit': 'Download failed. The file may be too large (max 2GB) or temporarily unavailable.',
  'Too many requests': 'Too many requests. Please wait a moment before trying again.',
  'rate limit': 'You\'re making requests too quickly. Please wait a minute.',
}

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
  downloadedBytes.value = 0
  totalBytes.value = 0

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
          downloadedBytes.value = progressEvent.loaded
          totalBytes.value = progressEvent.total || 0

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
