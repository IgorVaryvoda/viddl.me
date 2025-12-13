<template>
  <div class="platform-page">
    <section class="hero" :style="{ '--platform-color': platform.color }">
      <div class="hero-badge">
        <span class="hero-icon">{{ platform.icon }}</span>
      </div>
      <h1>
        <span class="hero-platform">{{ platform.name }}</span>
        Video Downloader
      </h1>
      <p class="subtitle">Download {{ platform.name }} videos in high quality. Free, fast, no signup required.</p>

      <VideoDownloader
        :placeholder="platform.placeholder"
        :filter-domain="platform.domain"
        @download-complete="handleDownloadComplete"
      />

      <div class="trust-badges">
        <span class="badge">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2v4M12 18v4M4.93 4.93l2.83 2.83M16.24 16.24l2.83 2.83M2 12h4M18 12h4M4.93 19.07l2.83-2.83M16.24 7.76l2.83-2.83"/>
          </svg>
          Free Forever
        </span>
        <span class="badge">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            <circle cx="12" cy="7" r="4"/>
          </svg>
          No Signup
        </span>
        <span class="badge">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
            <polyline points="22 4 12 14.01 9 11.01"/>
          </svg>
          High Quality
        </span>
        <span class="badge">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/>
          </svg>
          Fast Downloads
        </span>
      </div>
    </section>

    <section class="how-to">
      <div class="section-header">
        <span class="section-tag">Guide</span>
        <h2>How to Download {{ platform.name }} Videos</h2>
      </div>
      <div class="steps">
        <div class="step">
          <div class="step-number">01</div>
          <div class="step-content">
            <h3>Copy the Video URL</h3>
            <p>Go to {{ platform.name }} and find the video you want to download. Copy the URL from your browser's address bar or use the share button.</p>
          </div>
        </div>
        <div class="step">
          <div class="step-number">02</div>
          <div class="step-content">
            <h3>Paste the URL Above</h3>
            <p>Paste the {{ platform.name }} video URL into the input field above and click "Get Video" to fetch the available formats.</p>
          </div>
        </div>
        <div class="step">
          <div class="step-number">03</div>
          <div class="step-content">
            <h3>Download Your Video</h3>
            <p>Select your preferred quality and click to download. Your video will be saved as an MP4 file.</p>
          </div>
        </div>
      </div>
    </section>

    <section class="features-section">
      <div class="section-header">
        <span class="section-tag">Features</span>
        <h2>{{ platform.name }} Downloader Features</h2>
      </div>
      <div class="features-grid">
        <div v-for="feature in platform.features" :key="feature" class="feature-item">
          <svg class="feature-check" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
            <polyline points="20 6 9 17 4 12"/>
          </svg>
          <span>{{ feature }}</span>
        </div>
      </div>
    </section>

    <section class="faq-section">
      <div class="section-header">
        <span class="section-tag">FAQ</span>
        <h2>{{ platform.name }} Downloader FAQ</h2>
      </div>
      <div class="faq-list">
        <details v-for="(faq, index) in platform.faqs" :key="index" class="faq-item">
          <summary>
            <span>{{ faq.q }}</span>
            <svg class="faq-chevron" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M6 9l6 6 6-6"/>
            </svg>
          </summary>
          <p>{{ faq.a }}</p>
        </details>
      </div>
    </section>

    <section class="other-platforms">
      <div class="section-header">
        <span class="section-tag">More Platforms</span>
        <h2>Download from Other Platforms</h2>
      </div>
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
          <svg class="platform-link-arrow" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M5 12h14M12 5l7 7-7 7"/>
          </svg>
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
  padding: 2rem 0 4rem;
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 72px;
  height: 72px;
  background: linear-gradient(135deg, var(--platform-color) 0%, color-mix(in srgb, var(--platform-color) 60%, var(--accent-secondary)) 100%);
  border-radius: 20px;
  margin-bottom: 1.5rem;
  box-shadow: 0 20px 40px -20px var(--platform-color);
}

.hero-icon {
  font-size: 2.25rem;
}

.hero h1 {
  font-size: clamp(2rem, 6vw, 2.75rem);
  font-weight: 800;
  margin: 0 0 0.75rem 0;
  line-height: 1.2;
}

.hero-platform {
  background: linear-gradient(135deg, var(--text) 0%, var(--platform-color) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
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
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 1rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 100px;
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--text-secondary);
}

.badge svg {
  color: var(--accent);
}

.section-header {
  text-align: center;
  margin-bottom: 2.5rem;
}

.section-tag {
  display: inline-block;
  padding: 0.375rem 0.875rem;
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-secondary) 100%);
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: #09090b;
  margin-bottom: 1rem;
}

.section-header h2 {
  font-size: 1.75rem;
  margin: 0;
}

.how-to {
  padding: 4rem 0;
  border-top: 1px solid var(--border);
}

.steps {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.step {
  display: flex;
  gap: 1.5rem;
  padding: 1.75rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 16px;
  transition: all 0.3s;
}

.step:hover {
  border-color: var(--accent);
  transform: translateX(8px);
}

.step-number {
  flex-shrink: 0;
  font-family: var(--font-display);
  font-size: 2rem;
  font-weight: 800;
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-secondary) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  line-height: 1;
}

.step-content h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.125rem;
}

.step-content p {
  margin: 0;
  color: var(--text-secondary);
  line-height: 1.6;
}

.features-section {
  padding: 4rem 0;
  border-top: 1px solid var(--border);
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
  max-width: 600px;
  margin: 0 auto;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 0.875rem;
  padding: 1.125rem 1.25rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 12px;
  transition: all 0.2s;
}

.feature-item:hover {
  border-color: var(--accent);
}

.feature-check {
  color: var(--success);
  flex-shrink: 0;
}

.faq-section {
  padding: 4rem 0;
  border-top: 1px solid var(--border);
}

.faq-list {
  max-width: 700px;
  margin: 0 auto;
}

.faq-item {
  margin-bottom: 0.75rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 12px;
  overflow: hidden;
  transition: border-color 0.2s;
}

.faq-item:hover {
  border-color: var(--accent);
}

.faq-item[open] {
  border-color: var(--accent);
}

.faq-item summary {
  padding: 1.25rem 1.5rem;
  cursor: pointer;
  font-weight: 500;
  list-style: none;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  transition: background 0.2s;
}

.faq-item summary::-webkit-details-marker {
  display: none;
}

.faq-item summary:hover {
  background: var(--bg-tertiary);
}

.faq-chevron {
  flex-shrink: 0;
  color: var(--text-secondary);
  transition: transform 0.3s, color 0.2s;
}

.faq-item[open] .faq-chevron {
  transform: rotate(180deg);
  color: var(--accent);
}

.faq-item p {
  padding: 0 1.5rem 1.25rem;
  margin: 0;
  color: var(--text-secondary);
  line-height: 1.7;
}

.other-platforms {
  padding: 4rem 0;
  border-top: 1px solid var(--border);
}

.platform-links {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 0.875rem;
}

.platform-link {
  display: flex;
  align-items: center;
  gap: 0.625rem;
  padding: 0.875rem 1.25rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 12px;
  color: var(--text);
  text-decoration: none;
  font-weight: 500;
  transition: all 0.2s;
}

.platform-link:hover {
  border-color: var(--platform-color);
  transform: translateY(-2px);
  box-shadow: 0 10px 20px -10px rgba(0, 0, 0, 0.3);
}

.platform-link:hover .platform-link-arrow {
  transform: translateX(4px);
  color: var(--platform-color);
}

.platform-link-icon {
  font-size: 1.25rem;
}

.platform-link-arrow {
  color: var(--text-secondary);
  transition: all 0.2s;
}

@media (max-width: 768px) {
  .hero {
    padding: 1.5rem 0 3rem;
  }

  .hero-badge {
    width: 60px;
    height: 60px;
    border-radius: 16px;
  }

  .hero-icon {
    font-size: 1.75rem;
  }

  .features-grid {
    grid-template-columns: 1fr;
  }

  .step {
    flex-direction: column;
    gap: 1rem;
    padding: 1.5rem;
  }

  .step-number {
    font-size: 1.75rem;
  }

  .platform-links {
    gap: 0.75rem;
  }

  .platform-link {
    padding: 0.75rem 1rem;
    font-size: 0.9375rem;
  }
}
</style>
