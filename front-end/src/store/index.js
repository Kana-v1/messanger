import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

Vue.use(Vuex)

export default new Vuex.Store({
    state: {
        users: null,
        chats: null,
        signedId: -1,
        signUpMessage: '',
        signInMessage: '',
        isAuthorized: false,
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
        updateChats(state, chats) {
            state.chats = JSON.parse(chats)
            state.chats.forEach(chat => chat.Messages.forEach(message => message.Message = new Buffer.from(message.Message, 'base64').toString('binary')))
        },
        signInAccount(state, accId, message) {
            state.signedIn.push(accId)
            state.signInMessage = message,
                state.isAuthorized = true
        },
        signUpAccount(state, accId, message) {
            console.log(message);
            state.signedId = accId
            state.signUpMessage = message
        }

    },
    actions: {
        getUsers({ commit }) {
            axios.get(`/api/get/users`)
                .then((result) => commit('updateUsers', result.data))
                .catch(console.error)
        },
        getChats({ commit }) {
            axios.get('api/get/chats')
                .then(result => commit('updateChats', result.data))
                .catch(console.error)
        },
        signIn({ commit }) {
            return function (accId, log, password) {
                axios.post('api/post/signIn', {
                    accountId: accId,
                    log: log,
                    password: password
                }).then(result => commit('signInAccount', result.data, 'Succesfully signed up'))
                    .catch(error => this.state.signInMessage = error.reponse.data)
            }
        },
        signUp({ commit }, logData) {
            console.log(commit);
            //commit('signUpAccount', result.data, 'Succesfully signed in'))
            async () => axios.post('/api/post/signUp', {
                log: logData.log,
                password: logData.password
            }).then(response => {
                console.log(response)
            })
                .catch(error => {
                    console.log(error.response)
                });//.then(() => console.log(125125125125)).catch(err => { console.log('Call Failed!'); console.log(err); });//error => this.state.signUpMessage = error.reponse.data)
        },
    }
});


