<template>
  <div class="faq-page">
    <header class="page-header">
      <div class="header-badge">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/>
          <path d="M12 17h.01"/>
        </svg>
        Help Center
      </div>
      <h1>Frequently Asked Questions</h1>
      <p>Everything you need to know about viddl.me</p>
    </header>

    <div class="faq-categories">
      <section class="faq-category" v-for="category in categories" :key="category.name">
        <div class="category-header">
          <div class="category-icon">{{ category.icon }}</div>
          <h2>{{ category.name }}</h2>
        </div>
        <div class="faq-list">
          <details class="faq-item" v-for="(faq, index) in category.items" :key="index">
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
    </div>

    <section class="cta-section">
      <div class="cta-content">
        <h2>Ready to download?</h2>
        <p>Start downloading videos from your favorite platforms now.</p>
        <router-link to="/" class="cta-button">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
            <polyline points="7 10 12 15 17 10"/>
            <line x1="12" y1="15" x2="12" y2="3"/>
          </svg>
          Go to Downloader
        </router-link>
      </div>
    </section>

    <script type="application/ld+json" v-html="structuredData"></script>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const categories = [
  {
    name: 'General Questions',
    icon: '\uD83D\uDCDD',
    items: [
      {
        q: 'Is viddl.me free to use?',
        a: 'Yes! viddl.me is completely free with no hidden fees, subscriptions, or download limits. We believe everyone should have access to download their favorite videos.'
      },
      {
        q: 'Do I need to create an account?',
        a: 'No account or signup required. Just paste your video URL and download. We don\'t collect any personal information.'
      },
      {
        q: 'Is downloading videos legal?',
        a: 'Downloading videos for personal use is generally acceptable in most jurisdictions. However, please respect copyright laws and the platform\'s terms of service. Don\'t download content you don\'t have rights to redistribute or use commercially.'
      },
      {
        q: 'Do you store my downloaded videos?',
        a: 'No. Videos are processed in real-time and immediately deleted from our servers after download. We don\'t keep any copies of your downloaded content.'
      },
      {
        q: 'Is viddl.me safe to use?',
        a: 'Yes! We don\'t require any downloads, installations, or account creation. Everything happens in your browser securely over HTTPS.'
      }
    ]
  },
  {
    name: 'Supported Platforms',
    icon: '\uD83C\uDF10',
    items: [
      {
        q: 'Which platforms are supported?',
        a: 'We support YouTube, Twitter/X, Instagram, Facebook, Reddit, Threads, Vimeo, and Twitch. We\'re constantly adding support for more platforms.'
      },
      {
        q: 'Can I download YouTube Shorts?',
        a: 'Yes! YouTube Shorts are fully supported. Just paste the Shorts URL and download like any other video.'
      },
      {
        q: 'Do you support Instagram Reels?',
        a: 'Yes! Instagram Reels, IGTV, and regular video posts are all supported.'
      },
      {
        q: 'Can I download private or restricted videos?',
        a: 'No, we can only download publicly available videos. Private, age-restricted, or login-required videos cannot be downloaded.'
      },
      {
        q: 'What about Twitter/X videos?',
        a: 'Both twitter.com and x.com URLs are supported. You can download videos, GIFs, and media from public tweets.'
      }
    ]
  },
  {
    name: 'Technical Questions',
    icon: '\u2699\uFE0F',
    items: [
      {
        q: 'What video quality options are available?',
        a: 'We offer multiple quality options depending on what\'s available from the source, including 360p, 480p, 720p HD, 1080p Full HD, and 4K when available.'
      },
      {
        q: 'What format are videos downloaded in?',
        a: 'All videos are downloaded in MP4 format, which is compatible with virtually all devices and media players.'
      },
      {
        q: 'Is there a file size limit?',
        a: 'Yes, the maximum file size is 2GB. Most videos are well under this limit, but very long or high-quality videos may exceed it.'
      },
      {
        q: 'Why are some videos converted?',
        a: 'Some platforms use different video formats or stream video and audio separately. We automatically combine and convert them to MP4 for compatibility.'
      },
      {
        q: 'Can I download audio only?',
        a: 'Currently, we only support video downloads. Audio-only downloads may be added in a future update.'
      }
    ]
  },
  {
    name: 'Troubleshooting',
    icon: '\uD83D\uDEE0\uFE0F',
    items: [
      {
        q: 'Why is my download failing?',
        a: 'Common reasons include: the video is private or deleted, it\'s region-locked, the URL is incorrect, or the platform is experiencing issues. Try refreshing the page and trying again.'
      },
      {
        q: 'The video has no sound. What\'s wrong?',
        a: 'Some sources store video and audio separately. Our system combines them automatically, but occasionally there may be issues. Try downloading again or choose a different quality.'
      },
      {
        q: 'Why is the download so slow?',
        a: 'Download speed depends on your internet connection and the source server. Large files or popular videos may take longer to process.'
      },
      {
        q: 'I\'m getting a "Too many requests" error. What should I do?',
        a: 'This means you\'ve made too many requests in a short time. Please wait a minute before trying again. This limit helps keep the service free for everyone.'
      },
      {
        q: 'The website says my URL is not supported. Why?',
        a: 'Make sure you\'re copying the full URL including https://. Only public videos from supported platforms (YouTube, Twitter, Instagram, etc.) can be downloaded.'
      }
    ]
  }
]

const allFaqs = categories.flatMap(c => c.items)

const structuredData = computed(() => JSON.stringify({
  "@context": "https://schema.org",
  "@type": "FAQPage",
  "mainEntity": allFaqs.map(faq => ({
    "@type": "Question",
    "name": faq.q,
    "acceptedAnswer": {
      "@type": "Answer",
      "text": faq.a
    }
  }))
}))
</script>

<style scoped>
.faq-page {
  max-width: 800px;
  margin: 0 auto;
}

.page-header {
  text-align: center;
  padding: 2rem 0 4rem;
}

.header-badge {
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

.header-badge svg {
  color: var(--accent);
}

.page-header h1 {
  font-size: clamp(2rem, 6vw, 2.75rem);
  font-weight: 800;
  margin: 0 0 0.75rem 0;
  background: linear-gradient(135deg, var(--text) 0%, var(--accent) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.page-header p {
  color: var(--text-secondary);
  font-size: 1.125rem;
  margin: 0;
}

.faq-categories {
  display: flex;
  flex-direction: column;
  gap: 3rem;
}

.faq-category {
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 20px;
  padding: 2rem;
}

.category-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.category-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-secondary) 100%);
  border-radius: 12px;
  font-size: 1.5rem;
}

.category-header h2 {
  font-size: 1.375rem;
  margin: 0;
}

.faq-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.faq-item {
  background: var(--bg);
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
  background: var(--bg-secondary);
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

.cta-section {
  margin-top: 4rem;
  padding: 3rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 20px;
  text-align: center;
  position: relative;
  overflow: hidden;
}

.cta-section::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle at center, var(--accent-glow) 0%, transparent 50%);
  opacity: 0.3;
  pointer-events: none;
}

.cta-content {
  position: relative;
}

.cta-section h2 {
  font-size: 1.75rem;
  margin: 0 0 0.5rem 0;
}

.cta-section p {
  color: var(--text-secondary);
  margin: 0 0 1.5rem 0;
  font-size: 1.0625rem;
}

.cta-button {
  display: inline-flex;
  align-items: center;
  gap: 0.625rem;
  padding: 1rem 2rem;
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-secondary) 100%);
  color: #fff;
  text-decoration: none;
  border-radius: 12px;
  font-family: var(--font-display);
  font-weight: 600;
  font-size: 1rem;
  transition: all 0.2s;
}

.cta-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 30px -10px var(--accent-glow);
}

@media (max-width: 768px) {
  .page-header {
    padding: 1.5rem 0 3rem;
  }

  .faq-category {
    padding: 1.5rem;
    border-radius: 16px;
  }

  .category-header {
    gap: 0.75rem;
  }

  .category-icon {
    width: 40px;
    height: 40px;
    font-size: 1.25rem;
    border-radius: 10px;
  }

  .category-header h2 {
    font-size: 1.125rem;
  }

  .faq-item summary {
    padding: 1rem 1.25rem;
    font-size: 0.9375rem;
  }

  .faq-item p {
    padding: 0 1.25rem 1rem;
    font-size: 0.9375rem;
  }

  .cta-section {
    padding: 2rem 1.5rem;
    margin-top: 3rem;
  }

  .cta-section h2 {
    font-size: 1.5rem;
  }
}
</style>
