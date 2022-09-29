import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter)

const Login = () => import('../base/login.vue')
const Index = () => import('../base/index.vue')
const Files = () => import('../components/Files.vue')


const routes = [
    // base
    {
        path: '/login',
        component: Login
    },
    {
        path: '/index',
        component: Index,
        redirect: '/files',
        children: [
        { path: '/files', component: Files },
    ]
    },
]

const router = new VueRouter({
    mode: 'history',
    routes
})

export default router
