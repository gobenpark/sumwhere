import Vue from 'vue'
import Router from 'vue-router'
import LoginInput from '@/components/LoginInput'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'LoginInput',
      component: LoginInput
    }
  ]
})
