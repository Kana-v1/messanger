<template>
  <body class="w-100">
    <div class="chatList">
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
              ><div class = "chatTitle">
                {{ chatTitle(chat.Id) }}
              </div>
              <i>{{messageSender(chat.Id)}}</i>&#8594; {{ lastMessage(chat.Id).Message }}
              <div class = "timeFormat">
                {{lastMessageTime(chat.Id)}}
                </div>
            </span>
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
        return curChatMessages[curChatMessages.length - 1];
      }
    },
    chatTitle(chatId) {
      let users = this.$store.state.chats.find(
        (chat) => chat.Id == chatId
      ).Users;
      let name = "";
      if (users != null) {
        users.forEach((user) => {
          name += user.Name;
          name += "-";
        });
        return name.substring(0, name.length - 1);
      }
    },
    messageSender(chatId) {
      let message = this.lastMessage(chatId)
      return this.$store.state.users.find(user => user.Id === message.Sender).Name
    },
    lastMessageTime(chatId) {
      return this.lastMessage(chatId).Time.substring(0, 19) //19 - numbers of chars for date and time only
    }
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

.chatTitle {
  font-family: Courier New;
  font-size: smaller;
  margin-bottom: 20px;
  color: grey;
}

.timeFormat {
  position:absolute; 
  right:0;
  margin-right: 10px;
  font-family: Courier New;
  font-size: smaller;
  color: grey;
}

body {
  height: 100vh;
  background:black;
}

</style>
