const headers = [
  { text: '日期', value: 'submitTime' },
  { text: '工号', value: 'code' },
  { text: '姓名', value: 'nickName' },
  {
    text: '汇总',
    value: 'Summary',
    children: [
      {
        text: '字符总量', value: 'summaryFieldCharacter',
        icon: 'mdi-plus'
      },
      { text: '有效字符总量', value: 'summaryFieldEffectiveCharacter' },
      { text: '准确率', value: 'summaryAccuracyRate' },
      {
        text: '时间', value: 'summaryCostTime',
        icon: 'mdi-plus'
      },
      { text: '分块数量', value: 'summaryBlockNum' },
      { text: '分块效率', value: 'summaryBlockEfficiency' },
      { text: '字符效率', value: 'summaryFieldEfficiency' },
      { text: '录入?数量', value: 'summaryQuestionMarkNumber' },
      { text: '录入?比例', value: 'summaryQuestionMarkProportion' }
    ]
  },
  {
    text: '初审',
    value: 'first',
    children: [
      { text: '字符总量', value: 'op0FieldCharacter' },
      { text: '有效字符总量', value: 'op0FieldEffectiveCharacter' },
      { text: '准确率', value: 'op0AccuracyRate' },
      { text: '时间', value: 'op0CostTime' },
      { text: '发票数量', value: 'op0InvoiceNum' },
      { text: '分块数量', value: 'op0BlockNum' },
      { text: '分块效率', value: 'op0BlockEfficiency' },
      { text: '字符效率', value: 'op0FieldEfficiency' },
      { text: '录入?数量', value: 'op0QuestionMarkNumber' },
      { text: '录入?比例', value: 'op0QuestionMarkProportion' }
    ]
  },
  {
    text: '一码',
    value: 'one',
    children: [
      { text: '非报销单字符总量', value: 'op1NotExpenseAccountFieldCharacter' },
      { text: '非报销单有效字符总量', value: 'op1NotExpenseAccountFieldEffectiveCharacter' },
      { text: '报销单字符总量', value: 'op1ExpenseAccountFieldCharacter' },
      { text: '报销单有效字符总量', value: 'op1ExpenseAccountFieldEffectiveCharacter' },
      { text: '准确率', value: 'op1AccuracyRate' },
      { text: '时间', value: 'op1CostTime' },
      { text: '分块数量', value: 'op1BlockNum' },
      { text: '分块效率', value: 'op1BlockEfficiency' },
      { text: '字符效率', value: 'op1FieldEfficiency' },
      { text: '录入?数量', value: 'op1QuestionMarkNumber' },
      { text: '录入?比例', value: 'op1QuestionMarkProportion' }
    ]
  },
  {
    text: '二码',
    value: 'two',
    children: [
      { text: '非报销单字符总量', value: 'op2NotExpenseAccountFieldCharacter' },
      { text: '非报销单有效字符总量', value: 'op2NotExpenseAccountFieldEffectiveCharacter' },
      { text: '报销单字符总量', value: 'op2ExpenseAccountFieldCharacter' },
      { text: '报销单有效字符总量', value: 'op2ExpenseAccountFieldEffectiveCharacter' },
      { text: '准确率', value: 'op2AccuracyRate' },
      { text: '时间', value: 'op2CostTime' },
      { text: '分块数量', value: 'op2BlockNum' },
      { text: '分块效率', value: 'op2BlockEfficiency' },
      { text: '字符效率', value: 'op2FieldEfficiency' },
      { text: '录入?数量', value: 'op2QuestionMarkNumber' },
      { text: '录入?比例', value: 'op2QuestionMarkProportion' }
    ]
  },
  {
    text: '问题件',
    value: 'problem',
    children: [
      { text: '字符总量', value: 'opQFieldCharacter' },
      { text: '有效字符总量', value: 'opQFieldEffectiveCharacter' },
      { text: '准确率', value: 'opQAccuracyRate' },
      { text: '时间', value: 'opQCostTime' },
      { text: '分块数量', value: 'opQBlockNum' },
      { text: '分块效率', value: 'opQBlockEfficiency' },
      { text: '字符效率', value: 'opQFieldEfficiency' },
      { text: '录入?数量', value: 'opQQuestionMarkNumber' },
      { text: '录入?比例', value: 'opQQuestionMarkProportion' }
    ]
  },
]

const pheaders = [
  { text: '开始时间', value: 'applyAt' },
  { text: '结束时间', value: 'submitAt' },
  { text: '工号', value: 'code' },
  { text: '姓名', value: 'name' },
  { text: '字符总量', value: 'summaryFieldCharacter' },
  { text: '有效字符总量', value: 'summaryFieldEffectiveCharacter' },
  { text: '准确率', value: 'summaryAccuracyRate' },
  { text: '练习时间', value: 'costTime' },
  { text: '分块数量', value: 'summaryBlockNum' },
  { text: '分块效率', value: 'summaryBlockEfficiency' },
  { text: '字符效率', value: 'summaryFieldEfficiency' },
  { text: '录入?数量', value: 'summaryQuestionMarkNumber' },
  { text: '录入?比例', value: 'summaryQuestionMarkProportion' },
  { text: '错误查询', value: 'search' },
]

const prefixPath = '/main/yield/'

const tabsOptions = [
  {
    key: 'lu',
    label: '录入产量',
    path: `${prefixPath}lu`
  },

  {
    key: 'pError',
    label: '练习产量',
    path: `${prefixPath}practice`
  },
]

export default {
  headers,
  pheaders,
  tabsOptions
}