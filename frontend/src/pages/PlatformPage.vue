<template>
  <div class="platform-page">
    <section class="hero" :style="{ '--platform-color': platform.color }">
      <div class="hero-badge">
        <span class="hero-icon">{{ platform.icon }}</span>
      </div>
      <h1>{{ platform.name }} Video Downloader</h1>
      <p class="subtitle">Download {{ platform.name }} videos in high quality. Free, fast, no signup required.</p>

      <VideoDownloader
        :placeholder="platform.placeholder"
        :filter-domain="platform.domain"
        @download-complete="handleDownloadComplete"
      />

      <div class="trust-badges">
        <span class="badge">Free Forever</span>
        <span class="badge">No Signup</span>
        <span class="badge">High Quality</span>
        <span class="badge">Fast Downloads</span>
      </div>
    </section>

    <section class="how-to">
      <h2>How to Download {{ platform.name }} Videos</h2>
      <div class="steps">
        <div class="step">
          <div class="step-number">1</div>
          <div class="step-content">
            <h3>Copy the Video URL</h3>
            <p>Go to {{ platform.name }} and find the video you want to download. Copy the URL from your browser's address bar or use the share button.</p>
          </div>
        </div>
        <div class="step">
          <div class="step-number">2</div>
          <div class="step-content">
            <h3>Paste the URL Above</h3>
            <p>Paste the {{ platform.name }} video URL into the input field above and click "Get Video" to fetch the available formats.</p>
          </div>
        </div>
        <div class="step">
          <div class="step-number">3</div>
          <div class="step-content">
            <h3>Download Your Video</h3>
            <p>Select your preferred quality and click to download. Your video will be saved as an MP4 file.</p>
          </div>
        </div>
      </div>
    </section>

    <section class="features-section">
      <h2>{{ platform.name }} Downloader Features</h2>
      <div class="features-grid">
        <div v-for="feature in platform.features" :key="feature" class="feature-item">
          <span class="feature-check">✓</span>
          <span>{{ feature }}</span>
        </div>
      </div>
    </section>

    <section class="faq-section">
      <h2>{{ platform.name }} Downloader FAQ</h2>
      <div class="faq-list">
        <details v-for="(faq, index) in platform.faqs" :key="index" class="faq-item">
          <summary>{{ faq.q }}</summary>
          <p>{{ faq.a }}</p>
        </details>
      </div>
    </section>

    <section class="other-platforms">
      <h2>Download from Other Platforms</h2>
      <div class="platform-links">
        <router-link
          v-for="p in otherPlatforms"
          :key="p.path"
          :to="`/${p.path}`"
          class="platform-link"
          :style="{ '--platform-color': p.color }"
        >
          <span class="platform-link-icon">{{ p.icon }}</span>
          <span>{{ p.name }}</span>
        </router-link>
      </div>
    </section>

    <script type="application/ld+json" v-html="structuredData"></script>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import VideoDownloader from '../components/VideoDownloader.vue'
import { platforms as platformsData } from '../router'

const props = defineProps({
  platform: {
    type: Object,
    required: true
  }
})

const otherPlatforms = computed(() =>
  Object.entries(platformsData)
    .filter(([path]) => path !== props.platform.path)
    .map(([path, data]) => ({ ...data, path }))
)

const structuredData = computed(() => JSON.stringify({
  "@context": "https://schema.org",
  "@type": "WebApplication",
  "name": `viddl.me ${props.platform.name} Downloader`,
  "url": `https://viddl.me/${props.platform.path}`,
  "description": props.platform.description,
  "applicationCategory": "MultimediaApplication",
  "operatingSystem": "Any",
  "offers": {
    "@type": "Offer",
    "price": "0",
    "priceCurrency": "USD"
  },
  "aggregateRating": {
    "@type": "AggregateRating",
    "ratingValue": "4.8",
    "ratingCount": "1250"
  }
}))

const handleDownloadComplete = (data) => {
  console.log('Download complete:', data)
}

onMounted(() => {
  const ogTitle = document.querySelector('meta[property="og:title"]')
  const ogDesc = document.querySelector('meta[property="og:description"]')
  const ogUrl = document.querySelector('meta[property="og:url"]')

  if (ogTitle) ogTitle.setAttribute('content', props.platform.title)
  if (ogDesc) ogDesc.setAttribute('content', props.platform.description)
  if (ogUrl) ogUrl.setAttribute('content', `https://viddl.me/${props.platform.path}`)
})
</script>

<style scoped>
.platform-page {
  max-width: 900px;
  margin: 0 auto;
}

.hero {
  text-align: center;
  padding: 2rem 0 3rem;
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  background: color-mix(in srgb, var(--platform-color) 15%, transparent);
  border: 2px solid var(--platform-color);
  border-radius: 16px;
  margin-bottom: 1.5rem;
}

.hero-icon {
  font-size: 2rem;
}

.hero h1 {
  font-size: 2.25rem;
  font-weight: 700;
  margin: 0 0 0.75rem 0;
  color: var(--text);
}

.subtitle {
  font-size: 1.125rem;
  color: var(--text-secondary);
  margin: 0 0 2rem 0;
  max-width: 500px;
  margin-left: auto;
  margin-right: auto;
}

.trust-badges {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 0.75rem;
  margin-top: 2rem;
}

.badge {
  padding: 0.5rem 1rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 100px;
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.how-to {
  padding: 3rem 0;
  border-top: 1px solid var(--border);
}

.how-to h2 {
  text-align: center;
  font-size: 1.5rem;
  margin: 0 0 2rem 0;
}

.steps {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.step {
  display: flex;
  gap: 1.25rem;
  padding: 1.5rem;
  background: var(--bg-secondary);
  border-radius: 12px;
  border: 1px solid var(--border);
}

.step-number {
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--accent);
  color: white;
  font-weight: 700;
  border-radius: 50%;
}

.step-content h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.125rem;
}

.step-content p {
  margin: 0;
  color: var(--text-secondary);
  line-height: 1.5;
}

.features-section {
  padding: 3rem 0;
  border-top: 1px solid var(--border);
}

.features-section h2 {
  text-align: center;
  font-size: 1.5rem;
  margin: 0 0 2rem 0;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
  max-width: 500px;
  margin: 0 auto;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem;
  background: var(--bg-secondary);
  border-radius: 8px;
  border: 1px solid var(--border);
}

.feature-check {
  color: #22c55e;
  font-weight: 700;
}

.faq-section {
  padding: 3rem 0;
  border-top: 1px solid var(--border);
}

.faq-section h2 {
  text-align: center;
  font-size: 1.5rem;
  margin: 0 0 2rem 0;
}

.faq-list {
  max-width: 700px;
  margin: 0 auto;
}

.faq-item {
  margin-bottom: 0.75rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 8px;
  overflow: hidden;
}

.faq-item summary {
  padding: 1rem 1.25rem;
  cursor: pointer;
  font-weight: 500;
  list-style: none;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.faq-item summary::after {
  content: '+';
  font-size: 1.25rem;
  color: var(--text-secondary);
}

.faq-item[open] summary::after {
  content: '−';
}

.faq-item p {
  padding: 0 1.25rem 1rem;
  margin: 0;
  color: var(--text-secondary);
  line-height: 1.6;
}

.other-platforms {
  padding: 3rem 0;
  border-top: 1px solid var(--border);
}

.other-platforms h2 {
  text-align: center;
  font-size: 1.5rem;
  margin: 0 0 2rem 0;
}

.platform-links {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 0.75rem;
}

.platform-link {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.25rem;
  background: var(--bg-secondary);
  border: 2px solid var(--border);
  border-radius: 100px;
  color: var(--text);
  text-decoration: none;
  transition: border-color 0.2s, background 0.2s;
}

.platform-link:hover {
  border-color: var(--platform-color);
  background: color-mix(in srgb, var(--platform-color) 10%, transparent);
}

.platform-link-icon {
  font-size: 1.125rem;
}

@media (max-width: 768px) {
  .hero h1 {
    font-size: 1.75rem;
  }

  .features-grid {
    grid-template-columns: 1fr;
  }

  .step {
    flex-direction: column;
    text-align: center;
  }

  .step-number {
    margin: 0 auto;
  }
}
</style>
