<template>
  <v-switch
      v-model="darkMode"
      true-icon="mdi-weather-night"
      false-icon="mdi-weather-sunny"
      inset
      persistent-hint
      hide-details
      class="theme-toggle-switch"
  ></v-switch>
</template>

<script setup lang="ts">
import { useTheme } from 'vuetify'
import { ref, watch, onMounted } from 'vue'

const theme = useTheme()

// Initialize darkMode based on current theme
const darkMode = ref<boolean>(theme.global.current.value.dark)

// Watch for theme changes and update localStorage
watch(darkMode, (newValue: boolean) => {
  localStorage.setItem('themePreference', newValue ? 'dark' : 'light')
  theme.change(newValue ? 'dark' : 'light')
})

// Initialize from localStorage if available
onMounted(() => {
  const savedPreference = localStorage.getItem('themePreference')
  if (savedPreference) {
    darkMode.value = (savedPreference === 'dark')
    theme.change(savedPreference)
  }
})
</script>

<style scoped>
</style>