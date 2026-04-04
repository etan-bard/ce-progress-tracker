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
        <v-menu offset-y density="compact" :close-on-content-click="true">
          <template v-slot:activator="{ props: menuProps }">
            <v-btn 
              v-bind="menuProps"
              size="small"
              variant="outlined"
              color="primary"
            >
              <v-icon left size="small" class="mr-1">{{ sortIcons[sortMode] }}</v-icon>
              Sort: {{ sortMode === 'completion' ? 'Completion %' : 'Last Accessed' }}
              <v-icon right size="small" class="ml-1">mdi-chevron-down</v-icon>
              <v-tooltip activator="parent" location="top">
                Change course column sort criteria
              </v-tooltip>
            </v-btn>
          </template>
          <v-list density="compact" nav>
            <v-list-subheader>SORT CRITERIA</v-list-subheader>
            <v-list-item 
              @click="setGlobalSortMode('completion')"
              :active="sortMode === 'completion'"
              prepend-icon="mdi-percent"
              title="Completion %"
            >
            </v-list-item>
            <v-list-item 
              @click="setGlobalSortMode('date')"
              :active="sortMode === 'date'"
              prepend-icon="mdi-calendar"
              title="Last Accessed"
            >
            </v-list-item>
          </v-list>
        </v-menu>
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

    <ErrorAlert
      v-if="props.error"
      :message="props.error"
    />

    <v-data-table
      :headers="dynamicHeaders"
      :items="transformedItems"
      v-model:sort-by="sortBy"
      :search="search"
      :loading="isLoading"
      loading-text="Preparing cross-tab view..."
      hover
      class="elevation-0 cross-tab-v-data-table mt-2"
    >
      <!-- Dynamic slots for each course column -->
      <template 
        v-for="courseId in pivotData.courses" 
        v-slot:[`item.${courseId}`]="{ item }"
      >
        <CourseProgressCell 
          :completion="item[courseId]?.completion ?? null"
          :last-accessed="item[courseId]?.lastAccessed ?? null"
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
import { ref, computed } from 'vue'
import type { Takes } from '@/types'
import CourseProgressCell from './CourseProgressCell.vue'
import SearchField from '@/components/common/inputs/SearchField.vue'
import { usePivotData, pivotDataSort } from '@/composables/usePivotData'
import HelpAlert from "@/components/common/messages/HelpAlert.vue";
import ErrorAlert from '@/components/common/messages/ErrorAlert.vue'

const props = defineProps<{
  items: Takes[]
  isLoading?: boolean
  error?: string | null
}>()

const search = ref('')
const showDataTableHelp = ref(false)
const sortBy = ref<{ key: string; order: 'asc' | 'desc' }[]>([{ key: 'participantId', order: 'asc' }])
const sortMode = ref<'completion' | 'date'>('completion')

const sortIcons = {
  completion: 'mdi-percent',
  date: 'mdi-calendar'
}

// Use the pivot data composable
const { pivotData, transformedItems } = usePivotData(
  computed(() => props.items),
  search,
  sortBy,
  sortMode
)

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
      width: '10em'
    },
    ...pivotData.value.courses.map(courseId => ({
      title: `Course ${courseId}`,
      key: courseId.toString(),
      value: (item: any) => item[`${courseId}_${sortMode.value === 'completion' ? 'completion' : 'date'}`],
      sortable: true,
      align: 'end' as const,
      minWidth: '8em',
      sort: (a: any, b: any) => {
        // Find current sort order for this header's key
        const sortItem = sortBy.value.find(s => s.key === courseId.toString());
        const isDescending = sortItem?.order === 'desc';

        return pivotDataSort(a, b, isDescending);
      }
    }))
  ]
})

const setGlobalSortMode = (mode: 'completion' | 'date') => {
  sortMode.value = mode
  
  // If we are currently sorting by a course, we need to refresh the sortBy to trigger a re-sort
  const currentSortItem = sortBy.value.find(s => s.key !== 'participantId')
  if (currentSortItem) {
    const currentOrder = currentSortItem.order
    const currentKey = currentSortItem.key
    // Find index again to be safe with reactivity if needed, though findIndex was fine
    const currentSortIndex = sortBy.value.findIndex(s => s.key === currentKey)
    if (currentSortIndex !== -1) {
      const newSortBy = [...sortBy.value]
      // Re-assign same key/order to trigger reactive update
      newSortBy[currentSortIndex] = { key: currentKey, order: currentOrder }
      sortBy.value = newSortBy
    }
  }
}
</script>

<style scoped>
</style>