<template>
  <div class="container">
    <header>
      <h1>viddl.me</h1>
      <p class="subtitle">Download videos from YouTube, Twitter, Instagram, and more</p>
    </header>

    <div class="input-group">
      <input
        v-model="url"
        type="text"
        placeholder="Paste video URL here..."
        @keyup.enter="fetchVideoInfo"
        :disabled="loading"
      />
      <button @click="fetchVideoInfo" :disabled="loading || !url">
        {{ loading ? 'Loading...' : 'Get Video' }}
      </button>
    </div>

    <div v-if="error" class="error">
      {{ error }}
    </div>

    <div v-if="downloadProgress" class="progress-indicator">
      <div class="spinner"></div>
      <p>{{ downloadProgress }}</p>
    </div>

    <div v-if="loading" class="loader">
      <div class="spinner"></div>
      <p>Fetching video information...</p>
    </div>

    <div v-if="videoInfo && !loading" class="video-info">
      <div class="video-header">
        <img
          v-if="videoInfo.thumbnail"
          :src="videoInfo.thumbnail"
          :alt="videoInfo.title"
          class="thumbnail"
        />
        <div class="video-meta">
          <h3>{{ videoInfo.title }}</h3>
          <p>{{ videoInfo.uploader }}</p>
          <p>{{ formatDuration(videoInfo.duration) }}</p>
        </div>
      </div>

      <div v-if="videoInfo.formats && videoInfo.formats.length > 0" class="formats">
        <h4>Select Quality (MP4)</h4>
        <div class="format-grid">
          <button
            v-for="format in videoInfo.formats"
            :key="format.format_id"
            @click="downloadVideo(format.format_id)"
            :disabled="downloading"
            class="format-btn"
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
        style="width: 100%; margin-top: 1rem;"
      >
        {{ downloading ? 'Downloading...' : 'Download Best Quality' }}
      </button>
    </div>

    <div class="supported-sites">
      <p>Supported platforms:</p>
      <div class="sites">
        <span class="site-tag">YouTube</span>
        <span class="site-tag">Twitter/X</span>
        <span class="site-tag">Instagram</span>
        <span class="site-tag">TikTok</span>
        <span class="site-tag">Facebook</span>
        <span class="site-tag">Vimeo</span>
        <span class="site-tag">Reddit</span>
        <span class="site-tag">Twitch</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'

const url = ref('')
const videoInfo = ref(null)
const error = ref('')
const loading = ref(false)
const downloading = ref(false)
const downloadProgress = ref('')

const API_URL = import.meta.env.VITE_API_URL || '/api'

const fetchVideoInfo = async () => {
  if (!url.value) return

  error.value = ''
  loading.value = true
  videoInfo.value = null

  try {
    const response = await axios.post(`${API_URL}/info`, {
      url: url.value
    })
    videoInfo.value = response.data
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to fetch video information'
  } finally {
    loading.value = false
  }
}

const downloadVideo = async (format) => {
  downloading.value = true
  error.value = ''
  downloadProgress.value = 'Preparing download...'

  try {
    downloadProgress.value = 'Fetching video from server...'

    const response = await axios.post(
      `${API_URL}/download`,
      {
        url: url.value,
        format: format
      },
      {
        responseType: 'blob',
        onDownloadProgress: (progressEvent) => {
          if (progressEvent.total) {
            const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total)
            downloadProgress.value = `Downloading: ${percentCompleted}%`
          } else {
            downloadProgress.value = `Downloading: ${(progressEvent.loaded / 1024 / 1024).toFixed(1)} MB`
          }
        }
      }
    )

    downloadProgress.value = 'Processing file...'

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
    setTimeout(() => {
      downloadProgress.value = ''
    }, 2000)
  } catch (err) {
    error.value = err.response?.data?.error || 'Download failed'
    downloadProgress.value = ''
  } finally {
    downloading.value = false
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
