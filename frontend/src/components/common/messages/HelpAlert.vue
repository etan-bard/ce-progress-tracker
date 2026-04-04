<template>
  <v-expand-transition>
    <v-alert
      v-if="show"
      :type="type"
      variant="tonal"
      :class="['help-alert', `help-alert--${type}`]"
      :dismissible="dismissible"
      @click:close="$emit('close')"
    >
      <template v-if="icon || $slots.prepend" v-slot:prepend>
        <slot name="prepend">
          <v-icon v-if="icon">{{ icon }}</v-icon>
        </slot>
      </template>

      <div class="help-alert__content">
        <strong v-if="title" class="help-alert__title">{{ title }}</strong>
        <slot>
          <div v-if="text" class="help-alert__text" v-html="text"></div>
        </slot>
      </div>

      <template v-if="actionText" v-slot:append>
        <v-btn
          :size="actionSize"
          :variant="actionVariant"
          @click="$emit('action')"
        >
          {{ actionText }}
        </v-btn>
      </template>
    </v-alert>
  </v-expand-transition>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    show?: boolean
    title?: string
    text?: string
    type?: 'info' | 'success' | 'warning' | 'error'
    icon?: string
    dismissible?: boolean
    actionText?: string
    actionSize?: 'x-small' | 'small' | 'default' | 'large'
    actionVariant?: 'text' | 'outlined' | 'tonal' | 'flat' | 'elevated'
  }>(),
  {
    show: true,
    type: 'info',
    icon: 'mdi-lightbulb-on',
    dismissible: true,
    actionSize: 'small',
    actionVariant: 'text'
  }
)

defineEmits<{
  (e: 'close'): void
  (e: 'action'): void
}>()
</script>

<style scoped>
</style>