<template>
  <v-app-bar
    app
    color="primary"
    dark
    elevated
    height="64"
    class="overflow-visible"
  >
    <div class="cat-logo-container">
      <img
        src="@/assets/logo.svg"
        alt="CE Progress Tracker Logo"
        class="cat-logo"
      />
    </div>
    <v-toolbar-title class="font-weight-bold ml-16">
      CE Progress Tracker
    </v-toolbar-title>
    
    <v-spacer></v-spacer>

    <ThemeSwitch class="mr-4" />

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
import ThemeSwitch from "@/components/common/inputs/ThemeSwitch.vue";

const emit = defineEmits(['reload'])

const reloadData = () => {
  // Implement data reload logic here
  console.log('Reloading data...');
  emit('reload');
  // You would typically call an API here to refresh the data
}
</script>

<style scoped>
:deep(.v-toolbar__content) {
  overflow: visible !important;
}

.cat-logo-container {
  position: relative;
  width: 64px;
  height: 64px;
}

.cat-logo {
  position: absolute;
  bottom: -35px;
  left: 6px;
  height: 80px;
  width: auto;
  z-index: 1000;
}
</style>