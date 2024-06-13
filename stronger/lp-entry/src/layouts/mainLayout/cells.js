export const menus = [
  {
    pId: '-1',
    id: 'home',
    key: 'home',
    icon: 'mdi-home-outline',
    realm: 'home',
    title: '首页',
    link: '/main/home',
    leaf: false,
    visible: true
  },

  {
    pId: '-1',
    id: 'rule',
    key: 'rule',
    icon: 'mdi-pencil-ruler',
    realm: 'rule',
    title: '项目规则',
    link: '/main/rule',
    leaf: false,
    visible: true
  },

  {
    pId: '-1',
    id: 'entry',
    key: 'entry',
    icon: 'mdi-circle-edit-outline',
    realm: 'entry',
    title: '录入通道',
    link: '/main/entry',
    leaf: false,
    visible: true
  },

  {
    pId: '-1',
    id: 'train',
    key: 'train',
    icon: 'mdi-file-tree-outline',
    realm: 'train',
    title: '培训流程',
    link: '/main/train',
    leaf: false,
    visible: true
  },

  {
    pId: '-1',
    id: 'practice',
    key: 'practice',
    icon: 'mdi-pencil-outline',
    realm: 'practice',
    title: '上岗练习',
    link: '/main/practice',
    leaf: false,
    visible: true
  },

  {
    pId: '-1',
    id: 'exam',
    key: 'exam',
    icon: 'mdi-file-document-edit-outline',
    realm: 'exam',
    title: '上岗考核',
    link: '/main/exam',
    leaf: false,
    visible: true
  },

  {
    pId: '-1',
    id: 'yield',
    key: 'yield',
    icon: 'mdi-chart-line',
    realm: 'yield',
    title: '产量查询',
    link: '/main/yield',
    leaf: false,
    visible: true
  },

  {
    pId: '-1',
    id: 'error',
    key: 'error',
    icon: 'mdi-alert-circle-outline',
    realm: 'error',
    title: '错误查询',
    link: '/main/error',
    leaf: false,
    visible: true
  },

  {
    pId: '-1',
    id: 'complaint',
    key: 'complaint',
    icon: 'mdi-phone-in-talk-outline',
    realm: 'complaint',
    title: '客户投诉',
    link: '/main/complaint',
    leaf: false,
    visible: true
  },

  {
    pId: '-1',
    id: 'salary',
    key: 'salary',
    icon: 'mdi-cash',
    realm: 'salary',
    title: '我的工资',
    link: '/main/salary',
    leaf: false,
    visible: true
  }
]

export const menuAuth = ['train','rule','practice','exam','exam','yield','error']