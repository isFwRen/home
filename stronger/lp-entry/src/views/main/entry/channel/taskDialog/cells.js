// 顶部导航
const opTabs = [
  {
    authKey: 'hasOp0',
    label: '初审',
    name: 'Op0',
    value: 'op0',
  },

  {
    authKey: 'hasOp1',
    label: '一码',
    name: 'Op1',
    value: 'op1',
  },

  {
    authKey: 'hasOp2',
    label: '二码',
    name: 'Op2',
    value: 'op2',
  },

  {
    authKey: 'hasOpq',
    label: '问题件',
    name: 'Opq',
    value: 'opq',
  },

  {
    authKey: 'hasOpp',
    label: '练习',
    name: 'Opp',
    value: 'opp',
  },

  {
    authKey: 'hasOpe',
    label: '考核',
    name: 'Ope',
    value: 'ope',
  }
]

// 删单
const deleteFields = [
  {
    formKey: 'password',
    inputType: 'text',
    autofocus: true,
    label: '删单密码',
    prependOuter: '*',
    prependOuterClass: 'error--text',
    type: 'password',
    validation: [
      { rule: 'required', message: '删单密码不能为空.' }
    ]
  },

  {
    formKey: 'delRemarks',
    inputType: 'textarea',
    hideDetails: false,
    label: '备注',
    prependOuter: '*',
    prependOuterClass: 'error--text',
    validation: [
      { rule: 'required', message: '备注不能为空.' }
    ]
  }
]

// 提示
const toastedOptions = {
  duration: 2000,
  position: 'top-center'
}

export {
  opTabs,
  deleteFields,
  toastedOptions
}

export default {
  opTabs,
  deleteFields,
  toastedOptions
}