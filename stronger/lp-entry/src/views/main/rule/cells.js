const prefixPath = '/main/rule/'

const tabsOptions = [
  {
    key: 'rules',
    label: '业务规则',
    path: `${ prefixPath }business`,
    class: 'popover-rules'
  },

  {
    key: 'template',
    label: '报销单模板',
    path: `${ prefixPath }template`,
    class: 'popover-template'
  },

  {
    key: 'video',
    label: '教学视频',
    path: `${ prefixPath }video`,
    class: 'popover-video'
  }
]

export default {
  tabsOptions
}