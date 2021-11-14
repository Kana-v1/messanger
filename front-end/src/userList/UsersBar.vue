<template>
  <body class="bg-dark">
    <div class="userList bg-dark">
      <div>
        <b-list-group
          style="max-width: 300px"
          class="bg-dark"
          v-for="(user, i) in users"
          :key="i"
        >
          <b-list-group-item
            :class="
              hovered
                ? 'd-flex align-items-center border-light border-dark bg-danger'
                : 'd-flex align-items-center bg-danger border-light'
            "
            v-on:mouseover="(event) => changeColor(event, user.Id)"
            v-on:mouseleave="() => originalColor(user.Id)"
          >
            <b-avatar class="mr-3 bg-white"></b-avatar>
            <span class="mr-auto ms-2">{{ user.Name }}</span>
          </b-list-group-item>
        </b-list-group>
      </div>
    </div>
  </body>
</template>

<script>
export default {
  name: "UserBar",
  methods: {
    changeColor(e, userId) {
      e.preventDefault();
      this.isHovered.set(userId, true);
    },
    userHovered(userId) {
      return this.isHovered.get(userId);
    },
    originalColor(userId) {
      this.isHovered.set(userId, false)
    },
  },

  data() {
    return {
      isHovered: Map,
    };
  },

  computed: {
    users() {
      let users = this.$store.state.users;
      if (users != null) {
        users.forEach((user) => {
          this.isHovered.set(user.Id, false);
        });
      }
      return users;
    },
    hovered() {
      return this.isHovered.get(1)
    }
  },

  created() {
    this.isHovered = new Map();
    this.$store.dispatch("getUsers");
  },
};
</script>
