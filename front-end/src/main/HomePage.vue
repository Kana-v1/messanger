<template>
  <body>
    <div class="messangerItself d-flex w-100 h-100" v-show="isAuthorized">
      <div id = "user">
      <UserBar />
      </div>
      <div id = "chat">
      <ChatBar />
      </div>
    </div>
  </body>
</template>

<script>
import UserBar from "../userList/UsersBar.vue";
import ChatBar from "../chatList/ChatBar.vue";

export default {
  name: "HomePage",
  components: {
    UserBar,
    ChatBar,
  },
  data() {
    return {
      isAuthorized: false,
    };
  },
};
</script>
<style scoped>
body {
  height: 100vh;
  background: rgba(218, 145, 116, 0.425);
  
}
#chat {
  width: 100%; 
}
</style>
