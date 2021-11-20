import VueRouter from 'vue-router'
import Vue from 'vue'
import HomePage from '../main/HomePage.vue'
import AuthPage from '../authentication/AuthenticationPage.vue'

Vue.use(VueRouter)

export default new VueRouter({
    mode: 'history',
    routes: [{
        path: '/main',
        component: HomePage
    },
    {
        path: '/auth',
        component: AuthPage,
    }]
})