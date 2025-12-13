<template>
  <div class="faq-page">
    <header class="page-header">
      <h1>Frequently Asked Questions</h1>
      <p>Everything you need to know about viddl.me video downloader</p>
    </header>

    <div class="faq-categories">
      <section class="faq-category">
        <h2>General Questions</h2>
        <div class="faq-list">
          <details class="faq-item" v-for="(faq, index) in generalFaqs" :key="'general-' + index">
            <summary>{{ faq.q }}</summary>
            <p>{{ faq.a }}</p>
          </details>
        </div>
      </section>

      <section class="faq-category">
        <h2>Supported Platforms</h2>
        <div class="faq-list">
          <details class="faq-item" v-for="(faq, index) in platformFaqs" :key="'platform-' + index">
            <summary>{{ faq.q }}</summary>
            <p>{{ faq.a }}</p>
          </details>
        </div>
      </section>

      <section class="faq-category">
        <h2>Technical Questions</h2>
        <div class="faq-list">
          <details class="faq-item" v-for="(faq, index) in technicalFaqs" :key="'technical-' + index">
            <summary>{{ faq.q }}</summary>
            <p>{{ faq.a }}</p>
          </details>
        </div>
      </section>

      <section class="faq-category">
        <h2>Troubleshooting</h2>
        <div class="faq-list">
          <details class="faq-item" v-for="(faq, index) in troubleshootingFaqs" :key="'troubleshooting-' + index">
            <summary>{{ faq.q }}</summary>
            <p>{{ faq.a }}</p>
          </details>
        </div>
      </section>
    </div>

    <section class="cta-section">
      <h2>Ready to Download?</h2>
      <p>Start downloading videos from your favorite platforms now.</p>
      <router-link to="/" class="cta-button">Go to Downloader</router-link>
    </section>

    <script type="application/ld+json" v-html="structuredData"></script>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const generalFaqs = [
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

const platformFaqs = [
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

const technicalFaqs = [
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

const troubleshootingFaqs = [
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

const structuredData = computed(() => JSON.stringify({
  "@context": "https://schema.org",
  "@type": "FAQPage",
  "mainEntity": [
    ...generalFaqs,
    ...platformFaqs,
    ...technicalFaqs,
    ...troubleshootingFaqs
  ].map(faq => ({
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
  padding: 2rem 0 3rem;
}

.page-header h1 {
  font-size: 2.25rem;
  font-weight: 700;
  margin: 0 0 0.75rem 0;
}

.page-header p {
  color: var(--text-secondary);
  font-size: 1.125rem;
  margin: 0;
}

.faq-categories {
  display: flex;
  flex-direction: column;
  gap: 2.5rem;
}

.faq-category h2 {
  font-size: 1.25rem;
  margin: 0 0 1rem 0;
  padding-bottom: 0.75rem;
  border-bottom: 2px solid var(--accent);
}

.faq-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.faq-item {
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
  transition: background 0.2s;
}

.faq-item summary:hover {
  background: var(--bg);
}

.faq-item summary::after {
  content: '+';
  font-size: 1.25rem;
  color: var(--text-secondary);
  flex-shrink: 0;
  margin-left: 1rem;
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

.cta-section {
  text-align: center;
  padding: 3rem 0;
  margin-top: 2rem;
  border-top: 1px solid var(--border);
}

.cta-section h2 {
  font-size: 1.5rem;
  margin: 0 0 0.5rem 0;
}

.cta-section p {
  color: var(--text-secondary);
  margin: 0 0 1.5rem 0;
}

.cta-button {
  display: inline-block;
  padding: 1rem 2rem;
  background: var(--accent);
  color: white;
  text-decoration: none;
  border-radius: 8px;
  font-weight: 600;
  transition: background 0.2s, transform 0.1s;
}

.cta-button:hover {
  background: var(--accent-hover);
  transform: translateY(-1px);
}

@media (max-width: 768px) {
  .page-header h1 {
    font-size: 1.75rem;
  }
}
</style>
