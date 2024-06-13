const ExamRoutes = {
  path: 'exam',
  name: 'Exam',
  meta: {
    key: 'exam',
    realm: 'exam',
    title: '上岗考核',
    isRequired: true
  },
  component: () => import('@/views/main/exam'),

  children: [
    {
      path: '/main/exam',
      redirect: 'examChannel'
    },

    {
      path: 'examChannel',
      name: 'ExamChannel',
      meta: {
        key: 'examChannel',
        pKey: 'exam',
        realm: 'exam',
        title: '上岗考核',
      },
      component: () => import('@/views/main/exam/channel'),
      
      children: [
        {
          path: 'ope',
          name: 'Ope',
          meta: {
            key: 'ope',
            pKey: 'exam',
            realm: 'examChannel',
            title: '考核',
            path: 'ope'
          },
          component: () => import('@/views/main/entry/channel/taskDialog/ope')
        }
      ]
    }
  ]
}

export default ExamRoutes