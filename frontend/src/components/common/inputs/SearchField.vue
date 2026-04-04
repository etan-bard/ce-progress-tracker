<template>
  <v-text-field
    v-model="internalSearch"
    label="Search"
    prepend-inner-icon="mdi-magnify"
    variant="outlined"
    hide-details
    single-line
    clearable
    :loading="loading"
    :disabled="disabled"
    density="comfortable"
    class="search-field"
    @update:model-value="handleSearch"
  >
    <template v-if="$slots.append" v-slot:append>
      <slot name="append"></slot>
    </template>
  </v-text-field>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  modelValue: string
  loading?: boolean
  disabled?: boolean
  debounce?: number
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'search', value: string): void
}>()

const internalSearch = ref(props.modelValue)

// Handle search input with optional debounce
const handleSearch = (value: string) => {
  internalSearch.value = value
  emit('update:modelValue', value)
  emit('search', value)
}

// Watch for external changes to modelValue
watch(() => props.modelValue, (newValue) => {
  if (newValue !== internalSearch.value) {
    internalSearch.value = newValue
  }
})
</script>

<style scoped>
</style>