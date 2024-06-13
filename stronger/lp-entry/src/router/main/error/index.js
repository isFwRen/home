const ErrorRoutes = {
  path: 'error',
  name: 'Error',
  meta: {
    key: 'error',
    realm: 'error',
    title: '错误查询',
    isRequired: true
  },
  component: () => import('@/views/main/error'),

  children: [
    {
      path: '/main/error',
      redirect: 'practice'
    },

    {
      path: 'lu',
      name: 'Lu',
      meta: {
        key: 'lu',
        realm: 'error',
        pKey: 'error',
        title: '录入错误明细',
      },
      component: () => import('@/views/main/error/luError')
    },

    {
      path: 'practice',
      name: 'Practice',
      meta: {
        key: 'practice',
        realm: 'error',
        pKey: 'error',
        title: '练习错误明细',
      },
      component: () => import('@/views/main/error/pError')
    },
  ]
}

export default ErrorRoutes