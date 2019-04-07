import Vue from 'vue'
import Router from 'vue-router'
import LoginInput from '@/components/LoginInput'
import SignupInput from '@/components/SignupInput'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'LoginInput',
      component: LoginInput
    },
    {
      path: '/login',
      name: 'LoginInput',
      component: LoginInput
    },
    {
      path: '/signup',
      name: 'SignupInput',
      component: SignupInput
    }
  ]
})
