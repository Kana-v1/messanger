import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

Vue.use(Vuex)

export default new Vuex.Store({
    state: {
        users: null,
    },
    mutations: {
        updateUsers(state, users) {
            state.users = JSON.parse(users).sort(function (a, b) {
                if (a.Id > b.Id) {
                    return 1
                }
                if (a.Id < b.Id) {
                    return -1
                }
                return 0
            })
        },
    },
    actions: {
        getUsers({ commit }) {
            axios.get(`/api/get/users`)
                .then((result) => commit('updateUsers', result.data))
                .catch(console.error)
        }
    }
});


