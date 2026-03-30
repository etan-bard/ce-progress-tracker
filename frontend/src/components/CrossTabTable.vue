<template>
  <v-card 
    title="Participant Progress Cross-Tab"
    flat
    class="cross-tab-table"
  >
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
        title="Cross-tab Data View Help"
        type="info"
        icon="mdi-table"
        class="mx-4 mb-4"
        @close="showDataTableHelp = false"
    >
      This table shows participant ID along the Y-axis and course ID along the X-axis. Each cell displays the completion percentage and last accessed date for that participant-course combination.
    </HelpAlert>

    <v-alert
      v-if="error"
      type="error"
      variant="tonal"
      class="mb-4"
      dismissible
    >
      {{ error }}
    </v-alert>

    <v-data-table
      :headers="dynamicHeaders"
      :items="filteredParticipants"
      :search="search"
      :loading="isLoading"
      :sort-by="sortBy"
      loading-text="Preparing cross-tab view..."
      hover
      class="elevation-0"
    >
      <!-- Dynamic slots for each course column -->
      <template 
        v-for="courseId in pivotData.courses" 
        v-slot:[`item.${courseId}`]="{ item }"
      >
        <CourseProgressCell 
          :completion="getCellValue(pivotData, item, courseId).completion"
          :last-accessed="getCellValue(pivotData, item, courseId).lastAccessed"
        />
      </template>

      <!-- Empty state -->
      <template v-slot:no-data>
        <div>
          No data available
        </div>
      </template>
    </v-data-table>
  </v-card>
</template>

<script setup lang="ts">
import { ref, computed, watch, toRef } from 'vue'
import type { Takes } from '@/types'
import CourseProgressCell from './CourseProgressCell.vue'
import SearchField from './SearchField.vue'
import { usePivotData, getCellValue } from '@/composables/usePivotData'
import HelpAlert from "@/components/HelpAlert.vue";

const props = defineProps<{
  items: Takes[]
  isLoading?: boolean
  error?: string | null
}>()

const search = ref('')
const showDataTableHelp = ref(false)
const sortBy = ref<{ key: string; order: boolean }[]>([{ key: 'participantId', order: true }])

// Use the pivot data composable
const { pivotData } = usePivotData(toRef(props, 'items'))

// Generate dynamic headers based on courses
const dynamicHeaders = computed(() => {
  if (!pivotData.value.courses.length) {
    return [
      { 
        title: 'Participant ID', 
        key: 'participantId', 
        sortable: true, 
        align: 'start' as const
      }
    ]
  }

  return [
    { 
      title: 'Participant ID', 
      key: 'participantId', 
      sortable: true, 
      align: 'start' as const,
      width: '120px'
    },
    ...pivotData.value.courses.map(courseId => ({
      title: `Course ${courseId}`,
      key: courseId.toString(),
      sortable: true,
      align: 'end' as const
    }))
  ]
})

// Filter participants based on search
const filteredParticipants = computed(() => {
  if (!search.value) return pivotData.value.participants
  
  return pivotData.value.participants.filter(participantId => 
    participantId.toString().includes(search.value)
  )
})

// Watch for sorting changes
watch(sortBy, (newSort) => {
  console.log('Sorting changed:', newSort)
  // Sorting is handled automatically by v-data-table
})
</script>

<style scoped>
</style>