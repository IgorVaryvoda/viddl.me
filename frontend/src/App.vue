<template>
  <div class="app" :class="{ 'light-theme': !isDarkTheme }">
    <header class="header">
      <div class="header-content">
        <router-link to="/" class="logo">viddl.me</router-link>
        <nav class="nav">
          <router-link to="/faq" class="nav-link">FAQ</router-link>
          <button
            @click="toggleTheme"
            class="theme-btn"
            :aria-label="isDarkTheme ? 'Switch to light theme' : 'Switch to dark theme'"
          >
            {{ isDarkTheme ? '☀️' : '🌙' }}
          </button>
        </nav>
      </div>
    </header>

    <main class="main">
      <router-view />
    </main>

    <footer class="footer">
      <div class="footer-content">
        <div class="footer-links">
          <router-link to="/">Home</router-link>
          <router-link to="/youtube-downloader">YouTube</router-link>
          <router-link to="/twitter-downloader">Twitter</router-link>
          <router-link to="/instagram-downloader">Instagram</router-link>
          <router-link to="/faq">FAQ</router-link>
        </div>
        <p class="footer-copy">
          Free & open source •
          <a href="https://github.com/IgorVaryvoda/viddl.me" target="_blank" rel="noopener">GitHub</a>
        </p>
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
  document.documentElement.style.setProperty('--bg', isDarkTheme.value ? '#0f172a' : '#f8fafc')
  document.documentElement.style.setProperty('--bg-secondary', isDarkTheme.value ? '#1e293b' : '#e2e8f0')
  document.documentElement.style.setProperty('--text', isDarkTheme.value ? '#f1f5f9' : '#0f172a')
  document.documentElement.style.setProperty('--text-secondary', isDarkTheme.value ? '#94a3b8' : '#64748b')
  document.documentElement.style.setProperty('--border', isDarkTheme.value ? '#334155' : '#cbd5e1')
}
</script>

<style>
:root {
  --bg: #0f172a;
  --bg-secondary: #1e293b;
  --text: #f1f5f9;
  --text-secondary: #94a3b8;
  --border: #334155;
  --accent: #ff6b35;
  --accent-hover: #ff8555;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
  background: var(--bg);
  color: var(--text);
  line-height: 1.5;
  min-height: 100vh;
}

.app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.header {
  position: sticky;
  top: 0;
  background: var(--bg);
  border-bottom: 1px solid var(--border);
  z-index: 100;
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 1rem 1.5rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--accent);
  text-decoration: none;
}

.nav {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.nav-link {
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 0.9375rem;
  padding: 0.5rem 0.75rem;
  border-radius: 6px;
  transition: color 0.2s, background 0.2s;
}

.nav-link:hover,
.nav-link.router-link-active {
  color: var(--text);
  background: var(--bg-secondary);
}

.theme-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1.25rem;
  padding: 0.5rem;
  border-radius: 6px;
  transition: background 0.2s;
}

.theme-btn:hover {
  background: var(--bg-secondary);
}

.main {
  flex: 1;
  padding: 1.5rem;
}

.footer {
  border-top: 1px solid var(--border);
  padding: 2rem 1.5rem;
  margin-top: auto;
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
  text-align: center;
}

.footer-links {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 1.5rem;
  margin-bottom: 1rem;
}

.footer-links a {
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 0.875rem;
  transition: color 0.2s;
}

.footer-links a:hover {
  color: var(--text);
}

.footer-copy {
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.footer-copy a {
  color: var(--accent);
  text-decoration: none;
}

.footer-copy a:hover {
  text-decoration: underline;
}

@media (max-width: 640px) {
  .header-content {
    padding: 0.75rem 1rem;
  }

  .main {
    padding: 1rem;
  }

  .footer-links {
    gap: 1rem;
  }
}
</style>
