const RuleRoutes = {
  path: 'rule',
  name: 'Rule',
  meta: {
    key: 'rule',
    realm: 'rule',
    title: '项目规则',
    isRequired: true
  },
  component: () => import('@/views/main/rule'),

  children: [
    {
      path: '/main/rule',
      redirect: 'business'
    }, 

    {
      path: 'business',
      name: 'Business',
      meta: {
        key: 'business',
        realm: 'rule',
        pKey: 'rule',
        title: '业务规则',
      },
      component: () => import('@/views/main/rule/business')
    },

    {
      path: 'template',
      name: 'Template',
      meta: {
        key: 'template',
        realm: 'rule',
        pKey: 'rule',
        title: '报销单模板',
      },
      component: () => import('@/views/main/rule/template')
    },

    {
      path: 'video',
      name: 'Video',
      meta: {
        key: 'video',
        realm: 'rule',
        pKey: 'rule',
        title: '教学视频',
      },
      component: () => import('@/views/main/rule/video')
    }
  ]
}

export default RuleRoutes