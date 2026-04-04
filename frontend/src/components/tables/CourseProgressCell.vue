<template>
  <div class="course-progress-cell d-flex flex-column align-end">
    <template v-if="completion !== null">
      <CompletedPercentageCell :value="completion" />
      <div v-if="lastAccessed" class="text-body-small text-medium-emphasis">
        {{ formatDate(lastAccessed) }}
      </div>
      <div v-else class="text-body-small text-medium-emphasis">
        N/A
      </div>
    </template>
    <div v-else class="no-data text-medium-emphasis">
      —
    </div>
  </div>
</template>

<script setup lang="ts">
import CompletedPercentageCell from "@/components/tables/CompletedPercentageCell.vue";

const props = defineProps<{
  completion: number | null
  lastAccessed: string | null
}>()

const formatDate = (dateString: string): string => {
  try {
    const date = new Date(dateString)
    // Treat dates before 1971 as "N/A"
    if (date.getTime() < 31536000000) {
      return 'N/A'
    }
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