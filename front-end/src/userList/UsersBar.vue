<template>
  <body class="bg-dark h-100">
    <div class="userList bg-dark">
      <b-list-group style="width: 200px" class="bg-dark">
        <b-list-group-item
          v-for="(user, i) in users"
          :key="i"
          v-on:mouseover="Hover = i"
          v-on:mouseleave="Hover = -1"
          :class="
            Hover === i
              ? 'd-flex align-items-center hover'
              : 'd-flex align-items-center unhover'
          "
          style = "height: 75px"
        >
          <b-avatar class="mr-3 bg-white"></b-avatar>
          <span class="mr-auto ms-2">{{ user.Name }}</span>
        </b-list-group-item>
      </b-list-group>
    </div>
  </body>
</template>

<script>
export default {
  name: "UserBar",

  data() {
    return {
      isHovered: [],
      Hover: -1,
    };
  },

  computed: {
    users() {
      return this.$store.state.users;
    },
  },

  created() {
    this.$store.dispatch("getUsers");
  },
};
</script>

<style scoped>
  .unhover {
    background: rgba(255, 72, 0, 0.548);
  }
  .hover {
    background: rgba(255, 72, 0, 0.3);
  }
</style>