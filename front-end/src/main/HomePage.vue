<template>
<body>

  <div class = "messangerItself" v-show = 'isAuthorized'>
  <UserBar/>
  </div>
</body>
</template>

<script>
import UserBar from '../userList/UsersBar.vue'

export default {
  name: 'HomePage',
  components: {
    UserBar,
  },
  data() {
    return {
      isAuthorized: true,
    }
  },

}
</script>

