<template>
  <div class="app" :class="{ 'light-theme': !isDarkTheme }">
    <div class="grain-overlay"></div>
    <div class="grid-bg"></div>

    <header class="header">
      <div class="header-content">
        <router-link to="/" class="logo">
          <span class="logo-text">viddl</span>
          <span class="logo-dot">.</span>
          <span class="logo-text">me</span>
        </router-link>
        <nav class="nav">
          <router-link to="/faq" class="nav-link">FAQ</router-link>
          <a href="https://github.com/IgorVaryvoda/viddl.me" target="_blank" rel="noopener" class="nav-link nav-link-github">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
            </svg>
          </a>
          <button
            @click="toggleTheme"
            class="theme-btn"
            :aria-label="isDarkTheme ? 'Switch to light theme' : 'Switch to dark theme'"
          >
            <span class="theme-icon" :class="{ 'is-light': !isDarkTheme }">
              <svg v-if="isDarkTheme" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="5"/>
                <path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/>
              </svg>
              <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
              </svg>
            </span>
          </button>
        </nav>
      </div>
    </header>

    <main class="main">
      <router-view />
    </main>

    <footer class="footer">
      <div class="footer-content">
        <div class="footer-grid">
          <div class="footer-brand">
            <span class="footer-logo">viddl<span class="footer-dot">.</span>me</span>
            <p>Free & open source video downloader</p>
          </div>
          <div class="footer-links-group">
            <h4>Platforms</h4>
            <router-link to="/youtube-downloader">YouTube</router-link>
            <router-link to="/twitter-downloader">Twitter/X</router-link>
            <router-link to="/instagram-downloader">Instagram</router-link>
          </div>
          <div class="footer-links-group">
            <h4>Resources</h4>
            <router-link to="/">Home</router-link>
            <router-link to="/faq">FAQ</router-link>
            <a href="https://github.com/IgorVaryvoda/viddl.me" target="_blank" rel="noopener">GitHub</a>
          </div>
        </div>
        <div class="footer-bottom">
          <p>&copy; {{ new Date().getFullYear() }} viddl.me</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const isDarkTheme = ref(true)

onMounted(() => {
  const savedTheme = localStorage.getItem('viddl_theme')
  if (savedTheme) {
    isDarkTheme.value = savedTheme === 'dark'
  } else {
    isDarkTheme.value = window.matchMedia('(prefers-color-scheme: dark)').matches
  }
  applyTheme()
})

const toggleTheme = () => {
  isDarkTheme.value = !isDarkTheme.value
  localStorage.setItem('viddl_theme', isDarkTheme.value ? 'dark' : 'light')
  applyTheme()
}

const applyTheme = () => {
  if (isDarkTheme.value) {
    document.documentElement.style.setProperty('--bg', '#09090b')
    document.documentElement.style.setProperty('--bg-secondary', '#18181b')
    document.documentElement.style.setProperty('--bg-tertiary', '#27272a')
    document.documentElement.style.setProperty('--text', '#fafafa')
    document.documentElement.style.setProperty('--text-secondary', '#a1a1aa')
    document.documentElement.style.setProperty('--border', '#27272a')
    document.documentElement.style.setProperty('--glow-opacity', '0.6')
  } else {
    document.documentElement.style.setProperty('--bg', '#fafafa')
    document.documentElement.style.setProperty('--bg-secondary', '#f4f4f5')
    document.documentElement.style.setProperty('--bg-tertiary', '#e4e4e7')
    document.documentElement.style.setProperty('--text', '#09090b')
    document.documentElement.style.setProperty('--text-secondary', '#52525b')
    document.documentElement.style.setProperty('--border', '#d4d4d8')
    document.documentElement.style.setProperty('--glow-opacity', '0.3')
  }
}
</script>

<style>
:root {
  --bg: #09090b;
  --bg-secondary: #18181b;
  --bg-tertiary: #27272a;
  --text: #fafafa;
  --text-secondary: #a1a1aa;
  --border: #27272a;
  --accent: #ff6b35;
  --accent-secondary: #ff9500;
  --accent-glow: rgba(255, 107, 53, 0.4);
  --accent-secondary-glow: rgba(255, 149, 0, 0.3);
  --glow-opacity: 0.6;
  --success: #22c55e;
  --error: #ef4444;
  --font-display: 'Syne', sans-serif;
  --font-body: 'Outfit', sans-serif;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html {
  scroll-behavior: smooth;
}

body {
  font-family: var(--font-body);
  background: var(--bg);
  color: var(--text);
  line-height: 1.6;
  min-height: 100vh;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

h1, h2, h3, h4, h5, h6 {
  font-family: var(--font-display);
  font-weight: 700;
  line-height: 1.2;
}

.app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  position: relative;
  overflow-x: hidden;
}

.grain-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 1000;
  opacity: 0.03;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noise)'/%3E%3C/svg%3E");
}

.light-theme .grain-overlay {
  opacity: 0.02;
}

.grid-bg {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: -1;
  background-image:
    linear-gradient(var(--border) 1px, transparent 1px),
    linear-gradient(90deg, var(--border) 1px, transparent 1px);
  background-size: 60px 60px;
  opacity: 0.3;
  mask-image: radial-gradient(ellipse 80% 50% at 50% 0%, black 40%, transparent 100%);
  -webkit-mask-image: radial-gradient(ellipse 80% 50% at 50% 0%, black 40%, transparent 100%);
}

.header {
  position: sticky;
  top: 0;
  background: rgba(9, 9, 11, 0.8);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid var(--border);
  z-index: 100;
}

.light-theme .header {
  background: rgba(250, 250, 250, 0.8);
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 1rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo {
  font-family: var(--font-display);
  font-size: 1.75rem;
  font-weight: 800;
  text-decoration: none;
  display: flex;
  align-items: center;
  transition: transform 0.2s;
}

.logo:hover {
  transform: scale(1.02);
}

.logo-text {
  color: var(--text);
}

.logo-dot {
  color: var(--accent);
  text-shadow: 0 0 20px var(--accent-glow);
}

.nav {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.nav-link {
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 0.9375rem;
  font-weight: 500;
  padding: 0.625rem 1rem;
  border-radius: 8px;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.nav-link:hover,
.nav-link.router-link-active {
  color: var(--text);
  background: var(--bg-secondary);
}

.nav-link-github {
  padding: 0.625rem;
}

.theme-btn {
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  cursor: pointer;
  padding: 0.625rem;
  border-radius: 8px;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
}

.theme-btn:hover {
  background: var(--bg-tertiary);
  color: var(--accent);
  border-color: var(--accent);
  box-shadow: 0 0 20px var(--accent-glow);
}

.theme-icon {
  display: flex;
  transition: transform 0.3s;
}

.theme-icon.is-light {
  transform: rotate(-20deg);
}

.main {
  flex: 1;
  padding: 2rem;
  position: relative;
  z-index: 1;
}

.footer {
  border-top: 1px solid var(--border);
  padding: 4rem 2rem 2rem;
  margin-top: auto;
  background: var(--bg-secondary);
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
}

.footer-grid {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr;
  gap: 3rem;
  margin-bottom: 3rem;
}

.footer-brand {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.footer-logo {
  font-family: var(--font-display);
  font-size: 1.5rem;
  font-weight: 800;
  color: var(--text);
}

.footer-dot {
  color: var(--accent);
}

.footer-brand p {
  color: var(--text-secondary);
  font-size: 0.9375rem;
  max-width: 280px;
}

.footer-links-group {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.footer-links-group h4 {
  font-family: var(--font-display);
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 0.5rem;
}

.footer-links-group a {
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 0.9375rem;
  transition: color 0.2s;
}

.footer-links-group a:hover {
  color: var(--accent);
}

.footer-bottom {
  padding-top: 2rem;
  border-top: 1px solid var(--border);
  text-align: center;
}

.footer-bottom p {
  color: var(--text-secondary);
  font-size: 0.875rem;
}

@media (max-width: 768px) {
  .header-content {
    padding: 0.875rem 1.25rem;
  }

  .logo {
    font-size: 1.5rem;
  }

  .main {
    padding: 1.5rem 1.25rem;
  }

  .footer {
    padding: 3rem 1.25rem 1.5rem;
  }

  .footer-grid {
    grid-template-columns: 1fr;
    gap: 2rem;
  }

  .footer-brand {
    text-align: center;
    align-items: center;
  }

  .footer-brand p {
    max-width: none;
  }

  .footer-links-group {
    text-align: center;
  }
}
</style>
