import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import Status from '../views/Status.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Status',
    component: Status
  },
  {
    path: '/edit',
    name: 'Edit',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "edit" */ '../views/Edit.vue')
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
