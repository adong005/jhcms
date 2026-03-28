<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { RouterLink, RouterView } from 'vue-router';
import { loadSiteMetaByHost } from './api/site';

const siteTitle = ref('企业官网');

function updateMeta(name: string, content: string) {
  let tag = document.querySelector(`meta[name="${name}"]`) as HTMLMetaElement | null;
  if (!tag) {
    tag = document.createElement('meta');
    tag.setAttribute('name', name);
    document.head.appendChild(tag);
  }
  tag.setAttribute('content', content);
}

onMounted(async () => {
  const meta = await loadSiteMetaByHost();
  if (!meta) return;
  siteTitle.value = meta.title || siteTitle.value;
  document.title = siteTitle.value;
  updateMeta('keywords', meta.keywords || '');
  updateMeta('description', meta.description || '');
});
</script>

<template>
  <div class="site-shell">
    <header class="site-header">
      <div class="container nav-wrap">
        <div class="brand">{{ siteTitle }}</div>
        <nav class="nav">
          <RouterLink to="/">首页</RouterLink>
          <RouterLink to="/service">服务中心</RouterLink>
          <RouterLink to="/culture">企业文化</RouterLink>
          <RouterLink to="/contact">联系我们</RouterLink>
        </nav>
      </div>
    </header>

    <main>
      <section class="hero">
        <div class="container">
          <h1>专注数字化与品牌建设</h1>
          <p>一个通用、专业、可快速落地的企业官网模板。</p>
        </div>
      </section>
      <RouterView />
    </main>

    <footer class="site-footer">
      <div class="container">© {{ new Date().getFullYear() }} {{ siteTitle }}</div>
    </footer>
  </div>
</template>
