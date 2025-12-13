<template>
  <div class="home">
    <section class="hero">
      <h1>Download Videos from Anywhere</h1>
      <p class="subtitle">Fast, free, and easy. No signup required.</p>

      <VideoDownloader
        placeholder="Paste video URL here..."
        @download-complete="handleDownloadComplete"
      />

      <div class="supported-platforms">
        <router-link
          v-for="platform in platforms"
          :key="platform.path"
          :to="`/${platform.path}`"
          class="platform-chip"
          :style="{ '--platform-color': platform.color }"
        >
          <span class="platform-icon">{{ platform.icon }}</span>
          <span>{{ platform.name }}</span>
        </router-link>
      </div>
    </section>

    <section class="how-it-works">
      <h2>How It Works</h2>
      <div class="steps">
        <div class="step">
          <div class="step-number">1</div>
          <h3>Copy URL</h3>
          <p>Copy the video URL from YouTube, Twitter, Instagram, or any supported site</p>
        </div>
        <div class="step">
          <div class="step-number">2</div>
          <h3>Paste & Fetch</h3>
          <p>Paste the URL above and click "Get Video" to fetch available formats</p>
        </div>
        <div class="step">
          <div class="step-number">3</div>
          <h3>Download</h3>
          <p>Choose your preferred quality and download the video instantly</p>
        </div>
      </div>
    </section>

    <section class="features">
      <h2>Why Choose viddl.me?</h2>
      <div class="feature-grid">
        <div class="feature">
          <div class="feature-icon">⚡</div>
          <h3>Lightning Fast</h3>
          <p>No waiting. Downloads start instantly with our optimized servers.</p>
        </div>
        <div class="feature">
          <div class="feature-icon">🔒</div>
          <h3>Secure & Private</h3>
          <p>No data stored. No tracking. Your downloads are completely private.</p>
        </div>
        <div class="feature">
          <div class="feature-icon">💯</div>
          <h3>100% Free</h3>
          <p>No hidden fees, no subscriptions, no limits. Completely free forever.</p>
        </div>
        <div class="feature">
          <div class="feature-icon">📱</div>
          <h3>Works Everywhere</h3>
          <p>Download on any device - desktop, tablet, or mobile phone.</p>
        </div>
      </div>
    </section>

    <section class="platforms-section">
      <h2>Supported Platforms</h2>
      <div class="platform-cards">
        <router-link
          v-for="platform in platforms"
          :key="platform.path"
          :to="`/${platform.path}`"
          class="platform-card"
          :style="{ '--platform-color': platform.color }"
        >
          <div class="platform-card-icon">{{ platform.icon }}</div>
          <h3>{{ platform.name }}</h3>
          <p>Download {{ platform.name }} videos</p>
          <span class="platform-arrow">→</span>
        </router-link>
      </div>
    </section>

    <section class="faq-preview">
      <h2>Frequently Asked Questions</h2>
      <div class="faq-list">
        <details class="faq-item">
          <summary>Is viddl.me free to use?</summary>
          <p>Yes! viddl.me is completely free with no hidden fees, subscriptions, or download limits. We believe everyone should have access to download their favorite videos.</p>
        </details>
        <details class="faq-item">
          <summary>Do I need to create an account?</summary>
          <p>No account or signup required. Just paste your video URL and download. It's that simple.</p>
        </details>
        <details class="faq-item">
          <summary>What video quality can I download?</summary>
          <p>We offer multiple quality options including 360p, 480p, 720p HD, 1080p Full HD, and 4K when available from the source.</p>
        </details>
        <details class="faq-item">
          <summary>Is downloading videos legal?</summary>
          <p>Downloading videos for personal use is generally acceptable. However, please respect copyright laws and the platform's terms of service. Don't download content you don't have rights to redistribute.</p>
        </details>
      </div>
      <router-link to="/faq" class="view-all-faq">View All FAQs →</router-link>
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
  text-align: center;
  padding: 2rem 0 3rem;
}

.hero h1 {
  font-size: 2.5rem;
  font-weight: 700;
  margin: 0 0 0.75rem 0;
  background: linear-gradient(135deg, var(--text) 0%, var(--accent) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.subtitle {
  font-size: 1.25rem;
  color: var(--text-secondary);
  margin: 0 0 2rem 0;
}

.supported-platforms {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 0.75rem;
  margin-top: 2rem;
}

.platform-chip {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 100px;
  color: var(--text);
  text-decoration: none;
  font-size: 0.875rem;
  transition: border-color 0.2s, background 0.2s;
}

.platform-chip:hover {
  border-color: var(--platform-color);
  background: color-mix(in srgb, var(--platform-color) 10%, transparent);
}

.platform-icon {
  font-size: 1rem;
}

.how-it-works {
  padding: 3rem 0;
  border-top: 1px solid var(--border);
}

.how-it-works h2 {
  text-align: center;
  font-size: 1.75rem;
  margin: 0 0 2rem 0;
}

.steps {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 2rem;
}

.step {
  text-align: center;
  padding: 1.5rem;
}

.step-number {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 1rem;
  background: var(--accent);
  color: white;
  font-size: 1.25rem;
  font-weight: 700;
  border-radius: 50%;
}

.step h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.125rem;
}

.step p {
  margin: 0;
  color: var(--text-secondary);
  font-size: 0.9375rem;
  line-height: 1.5;
}

.features {
  padding: 3rem 0;
  border-top: 1px solid var(--border);
}

.features h2 {
  text-align: center;
  font-size: 1.75rem;
  margin: 0 0 2rem 0;
}

.feature-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1.5rem;
}

.feature {
  padding: 1.5rem;
  background: var(--bg-secondary);
  border-radius: 12px;
  border: 1px solid var(--border);
}

.feature-icon {
  font-size: 2rem;
  margin-bottom: 0.75rem;
}

.feature h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.125rem;
}

.feature p {
  margin: 0;
  color: var(--text-secondary);
  font-size: 0.9375rem;
  line-height: 1.5;
}

.platforms-section {
  padding: 3rem 0;
  border-top: 1px solid var(--border);
}

.platforms-section h2 {
  text-align: center;
  font-size: 1.75rem;
  margin: 0 0 2rem 0;
}

.platform-cards {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1rem;
}

.platform-card {
  position: relative;
  padding: 1.5rem;
  background: var(--bg-secondary);
  border: 2px solid var(--border);
  border-radius: 12px;
  text-decoration: none;
  color: var(--text);
  transition: border-color 0.2s, transform 0.2s;
}

.platform-card:hover {
  border-color: var(--platform-color);
  transform: translateY(-2px);
}

.platform-card-icon {
  font-size: 2rem;
  margin-bottom: 0.75rem;
}

.platform-card h3 {
  margin: 0 0 0.25rem 0;
  font-size: 1.125rem;
}

.platform-card p {
  margin: 0;
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.platform-arrow {
  position: absolute;
  top: 1.5rem;
  right: 1.5rem;
  color: var(--text-secondary);
  transition: color 0.2s, transform 0.2s;
}

.platform-card:hover .platform-arrow {
  color: var(--platform-color);
  transform: translateX(4px);
}

.faq-preview {
  padding: 3rem 0;
  border-top: 1px solid var(--border);
}

.faq-preview h2 {
  text-align: center;
  font-size: 1.75rem;
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

.view-all-faq {
  display: block;
  text-align: center;
  margin-top: 1.5rem;
  color: var(--accent);
  text-decoration: none;
  font-weight: 500;
}

.view-all-faq:hover {
  text-decoration: underline;
}

@media (max-width: 768px) {
  .hero h1 {
    font-size: 1.75rem;
  }

  .subtitle {
    font-size: 1rem;
  }

  .steps {
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }

  .feature-grid {
    grid-template-columns: 1fr;
  }

  .platform-cards {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 480px) {
  .platform-cards {
    grid-template-columns: 1fr;
  }
}
</style>
