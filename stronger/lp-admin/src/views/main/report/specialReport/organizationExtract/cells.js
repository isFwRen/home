import moment from 'moment'

const headers = [
  { text: '机构', value: 'institution' },
  { text: '核心立案号', value: 'billCode' },
  { text: '疾病', value: 'disease' },
  { text: '是否匹配', value: 'isMatch' },
  { text: '票据类型', value: 'ticketType' },
  { text: '医院名称', value: 'hospitalName' },
  { text: '总金额', value: 'totalMoney' },
  { text: '统筹金额', value: 'sumMoney' },
  { text: '范围外金额', value: 'outMoney' },
  { text: '范围内金额', value: 'inMoney' },
]

const fields = [
  {
    cols: 2,
    formKey: 'proCode',
    inputType: 'select',
    hideDetails: true,
    label: '项目',
    options: [],
    defaultValue: undefined
  },

  {
    cols: 3,
    formKey: 'time',
    inputType: 'date',
    hideDetails: true,
    label: '日期',
    clearable: true,
    range: true,
    defaultValue: [moment().format('YYYY-MM-DD'), moment().format('YYYY-MM-DD')]
  },
]


export default {
  headers,
  fields,
}