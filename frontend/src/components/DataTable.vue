<template>
  <v-card title="Course Progress" flat>
    <template v-slot:text>
      <div class="d-flex align-center ga-4">
        <SearchField
          v-model="search"
          label="Search"
          icon="mdi-magnify"
          :loading="props.isLoading"
          class="flex-grow-1"
        />
        <v-btn 
          size="small"
          variant="text"
          @click="showDataTableHelp = !showDataTableHelp"
        >
          <v-icon left small>{{ showDataTableHelp ? 'mdi-information-off' : 'mdi-information' }}</v-icon>
          {{ showDataTableHelp ? 'Hide' : 'Show' }} Help
        </v-btn>
      </div>
    </template>

    <HelpAlert
      v-model:show="showDataTableHelp"
      title="Linear Data View Help"
      type="info"
      icon="mdi-table"
      class="mx-4 mb-4"
      @close="showDataTableHelp = false"
    >
      This table shows each participant-course combination as a separate row with detailed information:
    </HelpAlert>

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
import SearchField from './SearchField.vue'
import HelpAlert from './HelpAlert.vue'

const search = ref('')
const showDataTableHelp = ref(false)

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