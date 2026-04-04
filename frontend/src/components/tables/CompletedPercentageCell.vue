<template>
  <v-chip
    :border="`${getColor(value)} thin opacity-25`"
    :color="getColor(value)"
    :text="formatCompletion(value)"
    size="small"
  ></v-chip>
</template>

<script setup lang="ts">

const props = defineProps<{
  value: number
}>()

const getColor = (percentage: number): string => {
  if (percentage >= 1) return 'success'
  else if (percentage > 0.1) return 'warning'
  else return 'error'
}

const formatCompletion = (value: number): string => {
  return `${(value * 100).toFixed(2)}%`
}
</script>

<style scoped>
</style>