<template>
  <div class="ops-empty">
    <div class="ops-empty__art" aria-hidden="true">
      <svg viewBox="0 0 132 96" fill="none" xmlns="http://www.w3.org/2000/svg">
        <defs>
          <linearGradient :id="gid" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0" :stop-color="accent" stop-opacity="0.22" />
            <stop offset="1" :stop-color="accent" stop-opacity="0.04" />
          </linearGradient>
        </defs>
        <rect x="18" y="30" width="96" height="54" rx="10" :fill="`url(#${gid})`" :stroke="accent" stroke-opacity="0.35" />
        <path d="M18 42 H114" :stroke="accent" stroke-opacity="0.25" />
        <rect x="30" y="20" width="72" height="14" rx="7" :fill="surface" :stroke="border" />
        <path d="M52 58 H80" :stroke="accent" stroke-opacity="0.5" stroke-linecap="round" />
        <path d="M56 68 H76" :stroke="muted" stroke-opacity="0.5" stroke-linecap="round" />
        <circle cx="104" cy="24" r="13" :fill="surface" :stroke="accent" stroke-opacity="0.5" />
        <path d="M104 19 V29 M99 24 H109" :stroke="accent" stroke-opacity="0.7" stroke-linecap="round" />
      </svg>
    </div>
    <div class="ops-empty__title">{{ title }}</div>
    <div v-if="desc" class="ops-empty__desc">{{ desc }}</div>
    <div v-if="$slots.default" class="ops-empty__action"><slot /></div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
const props = defineProps({
  title: { type: String, required: true },
  desc: { type: String, default: '' }
})
const accent = 'var(--accent)'
const surface = 'var(--surface)'
const border = 'var(--border)'
const muted = 'var(--text-muted)'
const gid = computed(() => 'ops-empty-grad-' + Math.random().toString(36).slice(2, 8))
</script>

<style scoped>
.ops-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 40px 20px;
  gap: 4px;
}
.ops-empty__art {
  width: 132px;
  height: 96px;
  margin-bottom: 12px;
  opacity: 0.95;
  filter: drop-shadow(0 8px 20px rgba(0, 0, 0, 0.12));
}
[data-theme='dark'] .ops-empty__art { filter: drop-shadow(0 8px 24px var(--accent-glow)); }
.ops-empty__title {
  font-family: var(--font-display);
  font-size: 14px;
  letter-spacing: 0.04em;
  color: var(--text);
}
.ops-empty__desc {
  font-size: 12.5px;
  color: var(--text-muted);
  max-width: 360px;
  line-height: 1.6;
}
.ops-empty__action { margin-top: 14px; }
</style>
