<template>
  <v-app-bar
    app
    color="primary"
    dark
    elevated
    height="64"
  >
    <v-toolbar-title class="font-weight-bold">
      CE Progress Tracker
    </v-toolbar-title>
    
    <v-spacer></v-spacer>

    <ThemeSwitch class="mr-4" />

    <!-- Reload Data Button -->
    <v-btn
      variant="text"
      @click="reloadData"
      class="mr-2"
    >
      <v-icon left>mdi-refresh</v-icon>
      Reload Data
    </v-btn>
  </v-app-bar>
</template>

<script setup>
import ThemeSwitch from "@/components/ThemeSwitch.vue";

const emit = defineEmits(['reload'])

const reloadData = () => {
  // Implement data reload logic here
  console.log('Reloading data...');
  emit('reload');
  // You would typically call an API here to refresh the data
}
</script>

<style scoped>

</style>