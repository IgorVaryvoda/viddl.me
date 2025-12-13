import { createRouter, createWebHistory } from 'vue-router'
import Home from '../pages/Home.vue'
import PlatformPage from '../pages/PlatformPage.vue'
import FAQ from '../pages/FAQ.vue'

const platforms = {
  'youtube-downloader': {
    name: 'YouTube',
    color: '#FF0000',
    icon: '▶️',
    domain: 'youtube.com',
    title: 'YouTube Video Downloader - Download YouTube Videos Free | viddl.me',
    description: 'Download YouTube videos in MP4, HD, and 4K quality. Free YouTube downloader - no signup required. Save YouTube videos instantly.',
    placeholder: 'Paste YouTube video URL here...',
    features: ['HD & 4K quality', 'YouTube Shorts support', 'Playlists supported', 'Fast downloads'],
    faqs: [
      { q: 'Can I download YouTube Shorts?', a: 'Yes! Simply paste the YouTube Shorts URL and download it like any other video.' },
      { q: 'What quality options are available?', a: 'We offer multiple quality options including 360p, 480p, 720p HD, 1080p Full HD, and 4K when available.' },
      { q: 'Can I download YouTube playlists?', a: 'Currently, you can download videos one at a time. Paste each video URL separately.' },
      { q: 'Is it free to download YouTube videos?', a: 'Yes, viddl.me is completely free with no hidden fees or subscriptions.' }
    ]
  },
  'twitter-downloader': {
    name: 'Twitter/X',
    color: '#1DA1F2',
    icon: '𝕏',
    domain: 'twitter.com',
    title: 'Twitter Video Downloader - Download X Videos Free | viddl.me',
    description: 'Download Twitter/X videos easily. Free Twitter video downloader - save tweets with video, GIFs, and media instantly.',
    placeholder: 'Paste Twitter/X video URL here...',
    features: ['Download tweets with video', 'GIF support', 'High quality downloads', 'Works with x.com'],
    faqs: [
      { q: 'How do I get the Twitter video URL?', a: 'Click on the tweet, then copy the URL from your browser address bar or use the share button.' },
      { q: 'Can I download Twitter GIFs?', a: 'Yes! GIFs are downloaded as MP4 video files for better quality and smaller size.' },
      { q: 'Does it work with X.com?', a: 'Yes, we support both twitter.com and x.com URLs.' },
      { q: 'Can I download private tweets?', a: 'No, only public tweets can be downloaded.' }
    ]
  },
  'instagram-downloader': {
    name: 'Instagram',
    color: '#E4405F',
    icon: '📷',
    domain: 'instagram.com',
    title: 'Instagram Video Downloader - Download Reels & Stories Free | viddl.me',
    description: 'Download Instagram Reels, videos, and stories. Free Instagram downloader - save IG content in high quality.',
    placeholder: 'Paste Instagram video URL here...',
    features: ['Reels support', 'IGTV videos', 'Story downloads', 'High quality'],
    faqs: [
      { q: 'Can I download Instagram Reels?', a: 'Yes! Paste the Reel URL and download it in high quality.' },
      { q: 'How do I download Instagram Stories?', a: 'Copy the story link and paste it here. Note: only public stories can be downloaded.' },
      { q: 'Can I download from private accounts?', a: 'No, only content from public accounts is accessible.' },
      { q: 'What about carousel posts?', a: 'For posts with multiple videos, you can select which one to download.' }
    ]
  },
  'facebook-downloader': {
    name: 'Facebook',
    color: '#1877F2',
    icon: '📘',
    domain: 'facebook.com',
    title: 'Facebook Video Downloader - Download FB Videos Free | viddl.me',
    description: 'Download Facebook videos in HD quality. Free FB video downloader - save Facebook videos, reels, and stories.',
    placeholder: 'Paste Facebook video URL here...',
    features: ['HD quality', 'Facebook Reels', 'Watch videos', 'Fast processing'],
    faqs: [
      { q: 'How do I get the Facebook video URL?', a: 'Right-click on the video and select "Copy video URL" or copy from the address bar.' },
      { q: 'Can I download Facebook Reels?', a: 'Yes, Facebook Reels are fully supported.' },
      { q: 'What about private videos?', a: 'Only public videos can be downloaded. Private or friends-only videos are not accessible.' },
      { q: 'Are Facebook Watch videos supported?', a: 'Yes, you can download videos from Facebook Watch.' }
    ]
  },
  'reddit-downloader': {
    name: 'Reddit',
    color: '#FF4500',
    icon: '🤖',
    domain: 'reddit.com',
    title: 'Reddit Video Downloader - Download Reddit Videos with Audio | viddl.me',
    description: 'Download Reddit videos with audio. Free Reddit video downloader - save Reddit posts, GIFs, and videos in high quality.',
    placeholder: 'Paste Reddit video URL here...',
    features: ['Videos with audio', 'GIF support', 'High quality', 'All subreddits'],
    faqs: [
      { q: 'Why do Reddit videos sometimes have no audio?', a: 'Reddit stores video and audio separately. Our tool combines them automatically for you.' },
      { q: 'Can I download Reddit GIFs?', a: 'Yes! GIFs are converted to MP4 for better quality.' },
      { q: 'Does it work with NSFW content?', a: 'Yes, as long as the content is publicly accessible.' },
      { q: 'What URL format should I use?', a: 'Use the full post URL or the direct video link.' }
    ]
  },
  'threads-downloader': {
    name: 'Threads',
    color: '#000000',
    icon: '🧵',
    domain: 'threads.net',
    title: 'Threads Video Downloader - Download Threads Videos Free | viddl.me',
    description: 'Download Threads videos easily. Free Threads downloader - save videos from threads.net in high quality.',
    placeholder: 'Paste Threads video URL here...',
    features: ['High quality', 'Fast downloads', 'Easy to use', 'No signup'],
    faqs: [
      { q: 'How do I get the Threads video URL?', a: 'Open the post in Threads app or web, tap share, and copy the link.' },
      { q: 'Can I download carousel posts?', a: 'For posts with multiple videos, you can select which one to download.' },
      { q: 'Is it free?', a: 'Yes, completely free with no limits.' },
      { q: 'What quality will I get?', a: 'Videos are downloaded in the best available quality.' }
    ]
  }
}

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: {
      title: 'viddl.me - Free Video Downloader for YouTube, Twitter, Instagram & More',
      description: 'Download videos from YouTube, Twitter/X, Instagram, Facebook, Reddit, and Threads. Fast, free, no signup required.'
    }
  },
  ...Object.entries(platforms).map(([path, platform]) => ({
    path: `/${path}`,
    name: platform.name,
    component: PlatformPage,
    props: { platform: { ...platform, path } },
    meta: {
      title: platform.title,
      description: platform.description
    }
  })),
  {
    path: '/faq',
    name: 'FAQ',
    component: FAQ,
    meta: {
      title: 'FAQ - Frequently Asked Questions | viddl.me',
      description: 'Common questions about viddl.me video downloader. Learn how to download videos from YouTube, Twitter, Instagram, and more.'
    }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

router.beforeEach((to, from, next) => {
  document.title = to.meta.title || 'viddl.me - Free Video Downloader'

  const descriptionMeta = document.querySelector('meta[name="description"]')
  if (descriptionMeta && to.meta.description) {
    descriptionMeta.setAttribute('content', to.meta.description)
  }

  const canonicalLink = document.querySelector('link[rel="canonical"]')
  if (canonicalLink) {
    canonicalLink.setAttribute('href', `https://viddl.me${to.path}`)
  }

  next()
})

export { platforms }
export default router
