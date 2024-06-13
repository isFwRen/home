const headers = [
  { text: '日期', value: 'submitDay' },
  { text: '工号', value: 'code' },
  { text: '姓名', value: 'nickName' },
  { text: '案件号', value: 'billNum' },
  { text: '机构号', value: 'agency' },
  { text: '字段', value: 'fieldName' },
  { text: '错误数据', value: 'wrong' },
  { text: '正确数据', value: 'right' },
  { text: '解析', value: 'analysis' },
  { text: '申诉', value: 'isComplain', width: 120 },
  { text: '差错审核', value: 'isWrongConfirm' }
]

const pheaders = [
  { text: '日期', value: 'submitDay' },
  { text: '工号', value: 'code' },
  { text: '姓名', value: 'nickName' },
  { text: '案件号', value: 'billNum' },
  { text: '机构号', value: 'agency' },
  { text: '字段', value: 'fieldName' },
  { text: '错误数据', value: 'wrong' },
  { text: '正确数据', value: 'right' },
]

const complaintOptions = [
  {
    value: 'true',
    label: '已申诉',
  },
  {
    value: 'false',
    label: '待申诉',
  }
]

const btns = [
  {
    class: 'pr-3',
    color: 'primary',
    text: '批量申诉'
  },
]

const reviewOptions = [
  {
    value: true,
    label: '通过'
  },

  {
    value: false,
    label: '不通过'
  }
]
const prefixPath = '/main/error/'

const tabsOptions = [
  {
    key: 'lu',
    label: '录入错误明细',
    path: `${prefixPath}lu`
  },

  {
    key: 'pError',
    label: '练习错误明细',
    path: `${prefixPath}practice`
  },
]


export default {
  headers,
  pheaders,
  btns,
  complaintOptions,
  reviewOptions,
  tabsOptions
}