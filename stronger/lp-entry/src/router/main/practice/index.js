const PracticeRoutes = {
  path: 'practice',
  name: 'Practice',
  meta: {
    key: 'practice',
    realm: 'practice',
    title: '上岗练习',
  },
  component: () => import('@/views/main/practice'),

  children: [
    {
      path: '/main/practice',
      redirect: 'practiceChannel'
    },

    {
      path: 'practiceChannel',
      name: 'PracticeChannel',
      meta: {
        key: 'practiceChannel',
        pKey: 'practice',
        realm: 'practice',
        title: '上岗练习',
      },
      component: () => import('@/views/main/practice/channel'),
      
      children: [
        {
          path: 'opp',
          name: 'Opp',
          meta: {
            key: 'opp',
            pKey: 'practice',
            realm: 'practiceChannel',
            title: '练习',
            path: 'opp'
          },
          component: () => import('@/views/main/entry/channel/taskDialog/opp')
        }
      ]
    }
  ]
}

export default PracticeRoutes