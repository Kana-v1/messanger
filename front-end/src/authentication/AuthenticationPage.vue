<template>
  <body class="bg-dark w-100 h-100">
    <div>
      <b-form>
        <b-form-input
          v-model="form.log"
          type="email"
          placeholder="Enter email"
          required
        ></b-form-input>

        <b-form-input
          v-model="form.password"
          class="mt-3"
          type="password"
          id="text-password"
          placeholder="Enter password"
          aria-describedby="password-help-block"
        ></b-form-input>
        <!-- <span>{{signInMessage()}} {{signUpMessage()}}</span> -->
        <b-button
          type="signIn"
          variant="outline-success"
          class="m-2 ms-4"
          @click="signIn()"
          >Sign In</b-button
        >
        <b-button
          type="signUp"
          variant="outline-secondary"
          class="ms-3"
          @click="signUp()"
          >Sign up</b-button
        >
      </b-form>
    </div>
  </body>
</template>

<script>
export default {
  name: "AuthPage",
  data() {
    return {
      form: {
        log: "",
        password: "",
      },
    };
  },
  methods: {
    signIn() {
      let func = this.$store.dispatch("signIn"); //(this.form.log, this.form.password);
      let result = func.then((res) =>
        res(123, this.form.log, this.form.password)
      );
      console.log(result);
    },
    signUp() {
       this.$store.dispatch("signUp", {log: this.form.log, password: this.form.password})
    },
    signInMessage() {
      return this.$store.signInMessage;
    },
    signUpMessage() {
      return this.$store.signUpMessage;
    },
  },
};
</script>

<style scoped>
div {
  position: absolute;
  top: 38%;
  left: 40%;
}
</style>
