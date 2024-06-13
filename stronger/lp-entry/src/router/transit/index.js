const TransitRoutes = {
  path: '/transit',
  name: 'transit',
  meta: {
    title: '登录',
    key: 'transit',
    path: 'transit'
  },
  component: () => import('@/views/transit/transit.vue')
}

export default TransitRoutes