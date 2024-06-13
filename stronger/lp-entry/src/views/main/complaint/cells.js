  import moment from 'moment'

  const headers = [
  { text: '月份', value: 'month', width: 65 },
  { text: '项目编码', value: 'proCode', width: 85 },
  { text: '案件号', value: 'billName', width: 120 },
  { text: '反馈日期', value: 'feedbackDate', width: 85 },
  { text: '录入日期', value: 'entryDate', width: 85 },
  { text: '错误字段', value: 'wrongFieldName', width: 85 },
  { text: '正确值', value: 'right', width: 65 },
  { text: '错误值', value: 'wrong', width: 65 },
  { text: '初审责任人工号', value: 'op0ResponsibleCode', width: 65 },
  { text: '初审责任人姓名', value: 'op0ResponsibleName', width: 65 },
  { text: '一码责任人工号', value: 'op1ResponsibleCode', width: 65 },
  { text: '一码责任人姓名', value: 'op1ResponsibleName', width: 65 },
  { text: '二码责任人工号', value: 'op2ResponsibleCode', width: 65 },
  { text: '二码责任人姓名', value: 'op2ResponsibleName', width: 65 },
  { text: '问题件责任人工号', value: 'opqResponsibleCode', width: 65 },
  { text: '问题件责任人姓名', value: 'opqResponsibleName', width: 65 },
  { text: '影像', value: 'imagePath', width: 125 },
]

const DEFAULT_MONTH = (() => {
  const date = new Date()

  const [year, month] = [
    date.getFullYear(),
    date.getMonth() + 1
  ]

  return moment(`${ year }-${ month }`).format('YYYY-MM')
})()
  
export default {
  headers,
  DEFAULT_MONTH
}