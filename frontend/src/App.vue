<template>
  <v-app>
    <AppNavbar />
    <v-main>
      <v-container fluid>
        <!-- Linear Data Table (Original) -->
        <DataTable 
          :items="takesList" 
          :is-loading="isLoading"
          :error="error"
          class="mb-8"
        />

        <!-- Cross-Tab Table (New) -->
        <v-card flat class="mt-6">
          <v-card-text>
            <CrossTabTable 
              :items="takesList" 
              :is-loading="isLoading"
              :error="error"
            />
          </v-card-text>
        </v-card>
      </v-container>
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
import AppNavbar from './components/AppNavbar.vue'
import { useTheme } from 'vuetify'
import { onMounted, ref } from 'vue'
import DataTable from "@/components/DataTable.vue"
import CrossTabTable from "@/components/CrossTabTable.vue"
import type { Takes } from '@/types'
import { fetchParticipantCourses } from '@/services/api'

const theme = useTheme()

onMounted(() => {
  const savedPreference = localStorage.getItem('themePreference')
  if (savedPreference) {
    theme.global.name.value = savedPreference
  } else if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
    theme.global.name.value = 'dark'
  }
})

// API data fetching with loading and error states
const takesList = ref<Takes[]>([])
const isLoading = ref<boolean>(false)
const error = ref<string | null>(null)

onMounted(async () => {
  try {
    isLoading.value = true
    error.value = null
    takesList.value = await fetchParticipantCourses()
  } catch (err) {
    error.value = 'Failed to load participant data. Please try again later.'
  } finally {
    isLoading.value = false
  }
})
</script>

<style scoped>
</style>