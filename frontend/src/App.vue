<template>
  <v-app>
    <AppNavbar @reload="fetchData" />
    <v-main>
      <v-container fluid>
        <DataTable 
          :items="takesList" 
          :is-loading="isLoading"
          :error="error"
          class="mb-8"
        />

        <CrossTabTable 
          :items="takesList" 
          :is-loading="isLoading"
          :error="error"
          class="mt-6"
        />
      </v-container>
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
import AppNavbar from './components/common/navigation/AppNavbar.vue'
import { useTheme } from 'vuetify'
import { onMounted, ref } from 'vue'
import DataTable from "@/components/tables/DataTable.vue"
import CrossTabTable from "@/components/tables/CrossTabTable.vue"
import type { Takes } from '@/types'
import { fetchParticipantCourses } from '@/services/api'

const theme = useTheme()

onMounted(() => {
  const savedPreference = localStorage.getItem('themePreference')
  if (savedPreference) {
    theme.change(savedPreference)
  } else if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
    theme.change('dark')
  }
})

// API data fetching with loading and error states
const takesList = ref<Takes[]>([])
const isLoading = ref<boolean>(false)
const error = ref<string | null>(null)

// Extract data fetching into a reusable method
const fetchData = async () => {
  isLoading.value = true
  error.value = null

  try {
    takesList.value = await fetchParticipantCourses()

    if (takesList.value.length == 0) {
      error.value = 'No data was found.'
    }
  } catch (err) {
    error.value = 'Failed to load participant data. Please try again later.'
  } finally {
    isLoading.value = false
  }
}

// Initial data load
onMounted(fetchData)
</script>

<style scoped>
</style>