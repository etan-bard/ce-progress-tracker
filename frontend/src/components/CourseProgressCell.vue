<template>
  <div class="course-progress-cell">
    <div v-if="completion !== null" class="cell-content">
      <CompletedPercentageCell :value="completion" />
      <div v-if="lastAccessed" class="date-text text-body-small">
        {{ formatDate(lastAccessed) }}
      </div>
      <div v-else class="no-date text-body-small">
        N/A
      </div>
    </div>
    <div v-else class="no-data">
      —
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue'
import CompletedPercentageCell from "@/components/CompletedPercentageCell.vue";

const props = defineProps<{
  completion: number | null
  lastAccessed: string | null
}>()

const getCompletionColor = (percentage: number): string => {
  if (percentage >= 1) return 'success'
  else if (percentage > 0.1) return 'warning'
  else return 'error'
}

const formatCompletion = (value: number): string => {
  return `${Math.round(value * 100)}%`
}

const formatDate = (dateString: string): string => {
  try {
    const date = new Date(dateString)
    return date.toLocaleDateString('en-US', { 
      month: 'short', 
      day: 'numeric', 
      year: '2-digit' 
    })
  } catch (error) {
    return dateString
  }
}
</script>

<style scoped>
  .text-body-small {
    padding-top: 0.3em;
  }
</style>