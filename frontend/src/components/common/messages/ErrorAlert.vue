<template>
  <v-alert
    type="error"
    variant="tonal"
    class="mb-4"
  >
    <slot>{{ message }}</slot>
  </v-alert>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    message?: string
  }>(),
  {
    message: ''
  }
)
</script>

<style scoped>
</style>