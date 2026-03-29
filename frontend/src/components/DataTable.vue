<template>
  <v-card
    title="Course Progress"
    flat
  >
    <template v-slot:text>
      <v-text-field
        v-model="search"
        label="Search"
        prepend-inner-icon="mdi-magnify"
        variant="outlined"
        hide-details
        single-line
      ></v-text-field>
    </template>

    <v-alert
      v-if="props.error"
      type="error"
      variant="tonal"
      class="mb-4"
    >
      {{ props.error }}
    </v-alert>

    <v-data-table
      :headers="headers"
      :items="props.items"
      :search="search"
      hover
      :loading="props.isLoading"
      loading-text="Loading participant data..."
    >
      <template v-slot:item.courseCompletion="{ value }">
        <CompletedPercentageCell :value="value" />
      </template>
    </v-data-table>
  </v-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { Takes } from "@/types"
import CompletedPercentageCell from './CompletedPercentageCell.vue'

const search = ref('')

// Define props
const props = defineProps<{
  items: Takes[]
  isLoading?: boolean
  error?: string | null
}>()

const headers = [
  { key: 'participantId', title: 'Participant ID', sortable: true },
  { key: 'courseId', title: 'Course ID', sortable: true },
  { key: 'dateFirstAccessed', title: 'First Accessed', sortable: true },
  { key: 'dateLastAccessed', title: 'Last Accessed', sortable: true },
  { key: 'courseCompletion', title: 'Completion (%)', sortable: true, align: 'end' },
]
</script>

<style scoped>
</style>