const regex_phone = /^(?:(?:\+|00)86)?1[3-9]\d{9}$/;

const usernameField = {
  formKey: 'account',
  label: "工号",
  validation: [
    { rule: 'required', message: '请输入工号!' }
  ]
}

const phoneField = {
  formKey: 'account',
  label: "手机号",
  validation: [
    { rule: 'required', message: '请输入手机号!' },
    { regex: regex_phone, message: '手机格式不正确.' }
  ]
}

export {
  regex_phone
}

export default {
  usernameField,
  phoneField
}