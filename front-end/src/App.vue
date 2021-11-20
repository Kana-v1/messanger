<template>
<body class = "bg-dark">
  <div id="app">
    <HomePage/>
  </div>
</body>
</template>

<script>
import HomePage from './main/HomePage.vue'


export default {
  name: 'App',
  components: {
    HomePage
  },
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
}
body {
  background: black;
}
</style>
