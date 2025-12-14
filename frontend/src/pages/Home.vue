<template>
  <div class="home">
    <section class="hero">
      <div class="hero-content">
        <div class="hero-badge">
          <span class="badge-dot"></span>
          Free & Open Source
        </div>
        <h1>
          <span class="hero-line">Download videos</span>
          <span class="hero-line hero-gradient">from anywhere</span>
        </h1>
        <p class="hero-subtitle">
          Grab videos from YouTube, Twitter, Instagram, and more.
          No signup. No limits. Just paste and download.
        </p>
      </div>

      <div class="downloader-wrapper">
        <VideoDownloader
          placeholder="Paste any video URL..."
          @download-complete="handleDownloadComplete"
        />
      </div>

      <div class="platforms-row">
        <span class="platforms-label">Supported:</span>
        <div class="platform-tags">
          <router-link
            v-for="platform in platforms"
            :key="platform.path"
            :to="`/${platform.path}`"
            class="platform-tag"
          >
            <span class="platform-icon">{{ platform.icon }}</span>
            {{ platform.name }}
          </router-link>
        </div>
      </div>
    </section>

    <section class="steps-section">
      <div class="section-header">
        <span class="section-tag">How it works</span>
        <h2>Three simple steps</h2>
      </div>
      <div class="steps">
        <div class="step" v-for="(step, index) in steps" :key="index">
          <div class="step-number">{{ String(index + 1).padStart(2, '0') }}</div>
          <h3>{{ step.title }}</h3>
          <p>{{ step.desc }}</p>
        </div>
      </div>
    </section>

    <section class="features-section">
      <div class="section-header">
        <span class="section-tag">Why viddl.me</span>
        <h2>Built different</h2>
      </div>
      <div class="features-grid">
        <div class="feature-card" v-for="feature in features" :key="feature.title">
          <div class="feature-icon-wrapper">
            <div class="feature-icon">{{ feature.icon }}</div>
          </div>
          <h3>{{ feature.title }}</h3>
          <p>{{ feature.desc }}</p>
        </div>
      </div>
    </section>

    <section class="platforms-section">
      <div class="section-header">
        <span class="section-tag">Platforms</span>
        <h2>Download from your favorites</h2>
      </div>
      <div class="platform-cards">
        <router-link
          v-for="platform in platforms"
          :key="platform.path"
          :to="`/${platform.path}`"
          class="platform-card"
          :style="{ '--platform-color': platform.color }"
        >
          <div class="platform-card-icon">{{ platform.icon }}</div>
          <div class="platform-card-content">
            <h3>{{ platform.name }}</h3>
            <p>Download {{ platform.name }} videos</p>
          </div>
          <svg class="platform-arrow" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M5 12h14M12 5l7 7-7 7"/>
          </svg>
        </router-link>
      </div>
    </section>

    <section class="faq-section">
      <div class="section-header">
        <span class="section-tag">FAQ</span>
        <h2>Common questions</h2>
      </div>
      <div class="faq-list">
        <details class="faq-item" v-for="(faq, index) in faqs" :key="index">
          <summary>
            <span>{{ faq.q }}</span>
            <svg class="faq-chevron" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M6 9l6 6 6-6"/>
            </svg>
          </summary>
          <p>{{ faq.a }}</p>
        </details>
      </div>
      <router-link to="/faq" class="faq-link">
        View all FAQs
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M5 12h14M12 5l7 7-7 7"/>
        </svg>
      </router-link>
    </section>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import VideoDownloader from '../components/VideoDownloader.vue'
import { platforms as platformsData } from '../router'

const platforms = computed(() =>
  Object.entries(platformsData).map(([path, data]) => ({
    ...data,
    path
  }))
)

const steps = [
  { title: 'Copy URL', desc: 'Copy the video link from YouTube, Twitter, Instagram, or any supported platform.' },
  { title: 'Paste & fetch', desc: 'Paste the URL above and click "Get Video" to see available formats.' },
  { title: 'Download', desc: 'Choose your preferred quality and download the video instantly.' }
]

const features = [
  { icon: '\u26A1', title: 'Lightning fast', desc: 'Optimized servers deliver your downloads instantly. No waiting.' },
  { icon: '\uD83D\uDD12', title: 'Private & secure', desc: 'No tracking, no data stored. Your downloads stay private.' },
  { icon: '\uD83D\uDCB0', title: '100% free', desc: 'No hidden fees, no subscriptions. Free forever.' },
  { icon: '\uD83C\uDF10', title: 'Works everywhere', desc: 'Desktop, tablet, or phone. Download on any device.' }
]

const faqs = [
  { q: 'Is viddl.me free to use?', a: 'Yes! Completely free with no hidden fees, subscriptions, or download limits.' },
  { q: 'Do I need to create an account?', a: 'No account required. Just paste your URL and download instantly.' },
  { q: 'What video quality can I download?', a: 'We offer multiple quality options including 360p, 720p HD, 1080p, and 4K when available.' },
  { q: 'Is downloading videos legal?', a: 'Downloading for personal use is generally acceptable. Please respect copyright laws.' }
]

const handleDownloadComplete = (data) => {
  console.log('Download complete:', data)
}
</script>

<style scoped>
.home {
  max-width: 1000px;
  margin: 0 auto;
}

.hero {
  padding: 3rem 0 4rem;
  text-align: center;
}

.hero-content {
  margin-bottom: 2.5rem;
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 100px;
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 1.5rem;
}

.badge-dot {
  width: 8px;
  height: 8px;
  background: var(--accent);
  border-radius: 50%;
  animation: pulse-dot 2s ease-in-out infinite;
}

@keyframes pulse-dot {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.6; transform: scale(0.9); }
}

.hero h1 {
  font-size: clamp(2rem, 5vw, 2.75rem);
  font-weight: 700;
  line-height: 1.15;
  margin: 0 0 1.5rem 0;
}

.hero-line {
  display: block;
}

.hero-gradient {
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-secondary) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.hero-subtitle {
  font-size: 1.25rem;
  color: var(--text-secondary);
  max-width: 520px;
  margin: 0 auto;
  line-height: 1.6;
}

.downloader-wrapper {
  margin-bottom: 2rem;
}

.platforms-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.platforms-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
  font-weight: 500;
}

.platform-tags {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 0.5rem;
}

.platform-tag {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.375rem 0.75rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 6px;
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 0.8125rem;
  font-weight: 500;
  transition: all 0.2s;
}

.platform-tag:hover {
  background: var(--bg-tertiary);
  color: var(--text);
  border-color: var(--accent);
}

.platform-icon {
  font-size: 0.875rem;
}

.section-header {
  text-align: center;
  margin-bottom: 3rem;
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
  color: #fff;
  margin-bottom: 1rem;
}

.section-header h2 {
  font-size: 2rem;
  margin: 0;
}

.steps-section {
  padding: 4rem 0;
  border-top: 1px solid var(--border);
}

.steps {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 2rem;
}

.step {
  text-align: center;
  padding: 2rem 1.5rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 16px;
  transition: all 0.3s;
}

.step:hover {
  transform: translateY(-4px);
  border-color: var(--accent);
  box-shadow: 0 20px 40px -20px rgba(255, 107, 53, 0.25);
}

.step-number {
  font-family: var(--font-display);
  font-size: 3rem;
  font-weight: 800;
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-secondary) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  line-height: 1;
  margin-bottom: 1rem;
}

.step h3 {
  font-size: 1.25rem;
  margin: 0 0 0.75rem 0;
}

.step p {
  margin: 0;
  color: var(--text-secondary);
  font-size: 0.9375rem;
  line-height: 1.5;
}

.features-section {
  padding: 4rem 0;
  border-top: 1px solid var(--border);
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1.5rem;
}

.feature-card {
  padding: 2rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 16px;
  transition: all 0.3s;
}

.feature-card:hover {
  border-color: var(--border);
  background: var(--bg-tertiary);
}

.feature-icon-wrapper {
  margin-bottom: 1.25rem;
}

.feature-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-secondary) 100%);
  border-radius: 12px;
  font-size: 1.5rem;
}

.feature-card h3 {
  font-size: 1.125rem;
  margin: 0 0 0.5rem 0;
}

.feature-card p {
  margin: 0;
  color: var(--text-secondary);
  font-size: 0.9375rem;
  line-height: 1.5;
}

.platforms-section {
  padding: 4rem 0;
  border-top: 1px solid var(--border);
}

.platform-cards {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1rem;
}

.platform-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.25rem 1.5rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 12px;
  text-decoration: none;
  color: var(--text);
  transition: all 0.3s;
}

.platform-card:hover {
  transform: translateY(-2px);
  border-color: var(--platform-color, var(--accent));
  box-shadow: 0 10px 30px -10px rgba(0, 0, 0, 0.3);
}

.platform-card:hover .platform-arrow {
  transform: translateX(4px);
  color: var(--platform-color, var(--accent));
}

.platform-card-icon {
  font-size: 1.75rem;
  flex-shrink: 0;
}

.platform-card-content {
  flex: 1;
  min-width: 0;
}

.platform-card h3 {
  font-size: 1rem;
  margin: 0 0 0.125rem 0;
}

.platform-card p {
  margin: 0;
  color: var(--text-secondary);
  font-size: 0.8125rem;
}

.platform-arrow {
  flex-shrink: 0;
  color: var(--text-secondary);
  transition: all 0.2s;
}

.faq-section {
  padding: 4rem 0;
  border-top: 1px solid var(--border);
}

.faq-list {
  max-width: 700px;
  margin: 0 auto 2rem;
}

.faq-item {
  margin-bottom: 0.75rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 12px;
  overflow: hidden;
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
  transition: transform 0.3s;
}

.faq-item[open] .faq-chevron {
  transform: rotate(180deg);
}

.faq-item p {
  padding: 0 1.5rem 1.25rem;
  margin: 0;
  color: var(--text-secondary);
  line-height: 1.6;
}

.faq-link {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--accent);
  text-decoration: none;
  font-weight: 500;
  margin: 0 auto;
  display: flex;
  justify-content: center;
  transition: gap 0.2s;
}

.faq-link:hover {
  gap: 0.75rem;
}

@media (max-width: 768px) {
  .hero {
    padding: 2rem 0 3rem;
  }

  .hero-subtitle {
    font-size: 1.0625rem;
  }

  .platforms-row {
    flex-direction: column;
    gap: 0.75rem;
  }

  .steps {
    grid-template-columns: 1fr;
    gap: 1rem;
  }

  .step {
    padding: 1.5rem;
  }

  .step-number {
    font-size: 2.5rem;
  }

  .features-grid {
    grid-template-columns: 1fr;
    gap: 1rem;
  }

  .feature-card {
    padding: 1.5rem;
  }

  .platform-cards {
    grid-template-columns: 1fr;
  }

  .section-header h2 {
    font-size: 1.75rem;
  }
}
</style>
