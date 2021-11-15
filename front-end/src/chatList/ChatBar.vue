<template>
  <body class="bg-dark w-100 h-100">
    <div class="chatList bg-dark">
      <div>
        <b-list-group class="bg-dark">
          <b-list-group-item
            v-for="(chat, i) in chats"
            :key="i"
            v-on:mouseover="Hover = i"
            v-on:mouseleave="Hover = -1"
            :class="
              Hover === i
                ? 'd-flex align-items-center hover '
                : 'd-flex align-items-center unhover '
            "
            style="min-height: 100px"
          >
            <span class="mr-auto ms-2"
              >{{ lastMessage(chat.Id) }} {{ chatName(chat.Id) }}</span
            >
          </b-list-group-item>
        </b-list-group>
      </div>
    </div>
  </body>
</template>

<script>
export default {
  name: "ChatBar",

  data() {
    return {
      isHovered: [],
      Hover: -1,
    };
  },
  methods: {
    lastMessage(chatId) {
      let curChatMessages = this.$store.state.chats.find(
        (chat) => chat.Id == chatId
      ).Messages;
      if (curChatMessages.length > 0) {
        return curChatMessages[curChatMessages.length - 1].Message;
      }
    },
    chatName(chatId) {
      let chat = this.$store.state.chats.find((chat) => chat.Id == chatId);
      if (chat.Users != null) {
        let name = "";
        chat.Users.forEach((user) => {
          name += user.Name;
          name += "-";
        });
        return (name -= "-");
      }
    },
  },

  computed: {
    chats() {
      return this.$store.state.chats;
    },
  },
  created() {
    this.$store.dispatch("getChats");
  },
};
</script>

<style scoped>
.unhover {
  background: rgba(218, 145, 116, 0.425);
}
.hover {
  background: rgba(218, 145, 116, 0.2);
}
</style>
