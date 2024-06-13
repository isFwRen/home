const YieldRoutes = {
  path: 'yield',
  name: 'Yield',
  meta: {
    key: 'yield',
    realm: 'yield',
    title: '产量查询',
    isRequired: true
  },
  component: () => import('@/views/main/yield'),

  children: [
    {
      path: '/main/yield',
      redirect: 'lu'
    },

    {
      path: 'lu',
      name: 'Lu',
      meta: {
        key: 'lu',
        realm: 'yield',
        pKey: 'yield',
        title: '练习产量',
      },
      component: () => import('@/views/main/yield/luYield')
    },

    {
      path: 'practice',
      name: 'Practice',
      meta: {
        key: 'practice',
        realm: 'yield',
        pKey: 'yield',
        title: '录入产量',
      },
      component: () => import('@/views/main/yield/pYield')
    },
  ]
}

export default YieldRoutes