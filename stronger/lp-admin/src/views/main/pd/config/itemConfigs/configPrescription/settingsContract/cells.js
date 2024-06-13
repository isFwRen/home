const claimTypes = [
  {
    label: '',
    value: -1
  },
  {
    label: '未定义',
    value: 0
  },
  {
    label: '无发票',
    value: 3
  },
  {
    label: '无报销',
    value: 4
  },
  {
    label: '有报销',
    value: 5
  },
  {
    label: '混合型',
    value: 6
  },
  {
    label: '简易',
    value: 7
  }
]

const fields = [
  {
    formKey: "contractStartTime",
    inputType: "date",
    dateFormat: "HH:mm:ss",
    hideDetails: false,
    hint: "格式为hh:mm:ss",
    label: "时效起始时间",
    mode: "time",
    prependOuter: "*",
    prependOuterClass: "error--text",
    timeFormat: "24hr",
    timeUseSeconds: true,
    validation: [{ rule: "required", message: "时效起始时间不能为空." }]
  },

  {
    formKey: "contractEndTime",
    inputType: "date",
    dateFormat: "HH:mm:ss",
    hideDetails: false,
    hint: "格式为hh:mm:ss",
    label: "时效结束时间",
    mode: "time",
    prependOuter: "*",
    prependOuterClass: "error--text",
    timeFormat: "24hr",
    timeUseSeconds: true,
    validation: [{ rule: "required", message: "时效结束时间不能为空." }]
  },

  {
    formKey: "claimType",
    inputType: "select",
    label: "单据类型",
    hideDetails: false,
    prependOuter: "*",
    prependOuterClass: "error--text",
    options: claimTypes,
    validation: [{ rule: "required", message: "单据类型不能为空." }]
  },

  {
    formKey: "contractOutsideStartTime",
    inputType: "date",
    dateFormat: "HH:mm:ss",
    hideDetails: false,
    hint: "格式为hh:mm:ss",
    label: "时效外开始时间",
    prependOuter: " ",
    mode: "time",
    timeFormat: "24hr",
    timeUseSeconds: true
  },

  {
    formKey: "contractOutsideEndTime",
    inputType: "date",
    dateFormat: "HH:mm:ss",
    hideDetails: false,
    prependOuter: " ",
    hint: "格式为hh:mm:ss",
    label: "时效外最晚时间",
    mode: "time",
    timeFormat: "24hr",
    timeUseSeconds: true
  },

  {
    formKey: "requirementsTime",
    inputType: "text",
    hideDetails: false,
    hint: "请输入数字",
    label: "考核要求(min)",
    prependOuter: "*",
    prependOuterClass: "error--text",
    validation: [
      { rule: "required", message: "考核要求不能为空." },
      { rule: "numeric", message: "考核要求为数字." }
    ]
  }
];





export default {
  fields,
  claimTypes
}