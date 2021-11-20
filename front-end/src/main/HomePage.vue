<template>
  <body class = "w-100 h-100">
    <div class="messangerItself d-flex w-100 h-100" v-show="isAuthorized">
      <UserBar />
      <ChatBar />
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
      isAuthorized: true,
    };
  },
};
</script>
