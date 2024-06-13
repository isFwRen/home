const headers1 = [
  { text: "发票数量", value: "invoiceCount", width: 190 },
  { text: "原始金额总和", value: "originalAmount" },
  { text: "调整金额总和", value: "adjustAmount" },
  { text: "扣费金额总和", value: "deductionAmount" },
  { text: "报销金额总和", value: "reimbursementAmount" },
];
const headers2 = [
  { text: "账单号", value: "billName2", width: 130 },
  { text: "诊疗方式", value: "treatmentType2", width: 90 },
  { text: "医院名称", value: "hospitalName2", width: 190 },
  { text: "开始时间", value: "startTime2", width: 150 },
  { text: "结束时间", value: "endTime2", width: 150 },
  { text: "天数", value: "days2", width: 80 },
  { text: "原始金额", value: "originalMoney", width: 90 },
  { text: "调整金额", value: "adjustMoney", width: 90 },
  { text: "自费金额", value: "selfMoney", width: 90 },
  { text: "自付金额", value: "selfPayMoney", width: 90 },
  { text: "其它金额", value: "otherMoney", width: 90 },
  { text: "报销金额", value: "reimbursementMoney", width: 90 },
];

const headers3 = [
  { text: "项目名称", value: "projectName2", width: 130 },
  { text: "总金额", value: "totalMoney2", width: 90 },
  { text: "甲类", value: "jia2", width: 90 },
  { text: "乙类", value: "yi2", width: 90 },
  { text: "部分自费", value: "partCost2", width: 90 },
  { text: "全自费", value: "allCost2", width: 90 },
  { text: "自费合计", value: "selfTotalCost2", width: 90 },
  { text: "其他", value: "other2", width: 90 },
  { text: "不合理费用", value: "reasonableCost2", width: 90 },
  { text: "社保内费用", value: "inCost2", width: 90 },
  { text: "社保外费用", value: "outCost2", width: 90 },
];

const form2 = [
  {
    name: 'input',
    field: 'billName2',
    title: '账单号',
    span: 8,
    placeholder: '请输入账单号'
  },
  {
    name: 'input',
    field: 'hospital2',
    title: '治疗医院',
    span: 8,
    placeholder: '请输入治疗医院',
  },
  {
    name: 'select',
    field: 'hospitalLevel2',
    title: '医院等级',
    span: 8,
    placeholder: '请选择医院等级',
    options: [
      {
        label: '三级甲等',
        value: '三级甲等'
      },
      {
        label: '二级乙等',
        value: '二级乙等'
      },
    ]
  },
  {
    name: 'select',
    field: 'hospitalProperty2',
    title: '医院性质',
    span: 8,
    placeholder: '请选择医院性质',
    options: [
      {
        label: '公立',
        value: '公立'
      },
      {
        label: '私立',
        value: '私立'
      },
    ]
  },
  {
    name: 'select',
    field: 'isPosition2',
    title: '是否定点',
    span: 8,
    placeholder: '请选择是否定点',
    options: [
      {
        label: '是',
        value: '是'
      },
      {
        label: '否',
        value: '否'
      },
    ]
  },
  {
    name: 'select',
    field: 'reimbursementType2',
    title: '报销类型',
    span: 8,
    placeholder: '请选择报销类型',
    options: [
      {
        label: '社保统筹支付',
        value: '社保统筹支付'
      },
      {
        label: '农和支付',
        value: '农和支付'
      },
      {
        label: '城镇职工',
        value: '城镇职工'
      },
      {
        label: '城镇居民',
        value: '城镇居民'
      },
    ]
  },
  {
    name: 'select',
    field: 'medicalType2',
    title: '诊疗类型',
    span: 8,
    placeholder: '请选择诊疗类型',
    options: [
      {
        label: '门诊',
        value: '门诊'
      },
      {
        label: '住院',
        value: '住院'
      },
    ]
  },
  {
    name: 'select',
    field: 'ticketType2',
    title: '票据类型',
    span: 8,
    placeholder: '请选择票据类型',
    options: [
      {
        label: '财政电子票据',
        value: '财政电子票据'
      },
      {
        label: '普通发票',
        value: '普通发票'
      },
    ]
  },
  {
    name: 'input',
    field: 'ticketCode2',
    title: '票据代码',
    span: 8,
    placeholder: '请输入票据代码',
  },
  {
    name: 'input',
    field: 'ticketNumber2',
    title: '票据号码',
    span: 8,
    placeholder: '请输入票据号码',
  },
  {
    name: 'input',
    field: 'checkCode2',
    title: '校验码',
    span: 8,
    placeholder: '请输入校验码'
  },
  {
    name: 'date',
    field: 'credentialEndDate2',
    title: '开票日期',
    span: 8,
    placeholder: '请选择日期',
  },
  {
    name: 'select',
    field: 'ticketState2',
    title: '票据状态',
    span: 8,
    placeholder: '请选择票据状态',
    options: [
      {
        label: '票据正常',
        value: '票据正常'
      },
      {
        label: '票据异常',
        value: '票据异常'
      },
    ]
  },
  {
    name: 'input',
    field: 'ticketRemark2',
    title: '电票查询备注',
    span: 16,
    placeholder: '请输入电票查询备注',
  },
  {
    name: 'date',
    field: 'ticketStartDate2',
    title: '票据开始日期',
    span: 8,
    placeholder: '请选择日期',
  },
  {
    name: 'date',
    field: 'ticketEndDate2',
    title: '票据结束日期',
    span: 8,
    placeholder: '请选择日期',
  },
  {
    name: 'input',
    field: 'treatMentDays2',
    title: '就诊天数',
    span: 8,
    placeholder: '请输入就诊天数',
  },
  {
    name: 'select',
    field: 'ticketInputType2',
    title: '票据录入方式',
    span: 8,
    placeholder: '请选择票据录入方式',
    options: [
      {
        label: '明细录入',
        value: '明细录入'
      },
      {
        label: '汇总录入',
        value: '汇总录入'
      },
    ]
  },
  {
    name: 'input',
    field: 'ticketBalance2',
    title: '票据余额',
    span: 8,
    placeholder: '请输入票据余额',
  },
  {
    name: 'input',
    field: 'socialInsurancePayment2',
    title: '社保支付',
    span: 8,
    placeholder: '请输入社保支付',
  },
  {
    name: 'input',
    field: 'thirdPartyAmount2',
    title: '第三方支付总额',
    span: 8,
    placeholder: '请输入第三方支付总额',
  },
  {
    name: 'input',
    field: 'thirdPartyAgency2',
    title: '第三方支付机构',
    span: 8,
    placeholder: '请输入第三方支付机构',
  },
]

const form3 = [
  {
    name: 'text',
    field: 'jiaDetail3',
    title: '甲类明细',
    span: 12,
    placeholder: '请输入甲类明细'
  },
  {
    name: 'text',
    field: 'jiaRemark3',
    title: '甲类备注',
    span: 12,
    placeholder: '请输入甲类备注'
  },
  {
    name: 'text',
    field: 'yiDetail3',
    title: '乙类明细',
    span: 12,
    placeholder: '请输入乙类明细'
  },
  {
    name: 'text',
    field: 'yiRemark3',
    title: '乙类备注',
    span: 12,
    placeholder: '请输入乙类备注'
  },
  {
    name: 'text',
    field: 'partCostDetail3',
    title: '部分自费明细',
    span: 12,
    placeholder: '请输入部分自费明细'
  },
  {
    name: 'text',
    field: 'partCostRemark3',
    title: '部分自费备注',
    span: 12,
    placeholder: '请输入部分自费备注'
  },
  {
    name: 'text',
    field: 'allCostDetail3',
    title: '全自费明细',
    span: 12,
    placeholder: '请输入全自费明细'
  },
  {
    name: 'text',
    field: 'allCostRemark3',
    title: '全自费备注',
    span: 12,
    placeholder: '请输入全自费备注'
  },
  {
    name: 'text',
    field: 'selfCostDetail3',
    title: '自费合计明细',
    span: 12,
    placeholder: '请输入自费合计明细'
  },
  {
    name: 'text',
    field: 'selfCostRemark3',
    title: '自费合计备注',
    span: 12,
    placeholder: '请输入自费合计备注'
  },
  {
    name: 'text',
    field: 'otherDetail3',
    title: '其它明细',
    span: 12,
    placeholder: '请输入其它明细'
  },
  {
    name: 'text',
    field: 'otherRemark3',
    title: '其它备注',
    span: 12,
    placeholder: '请输入其它备注'
  },
  {
    name: 'text',
    field: 'noReasonable3',
    title: '不合理费用明细',
    span: 12,
    placeholder: '请输入不合理费用明细'
  },
  {
    name: 'text',
    field: 'boReasonableRemark3',
    title: '不合理费用备注',
    span: 12,
    placeholder: '请输入不合理费用备注'
  },
  {
    name: 'text',
    field: 'inCostDetail3',
    title: '社保内费用明细',
    span: 12,
    placeholder: '请输入社保内费用明细'
  },
  {
    name: 'text',
    field: 'inCostRemark3',
    title: '社保内费用备注',
    span: 12,
    placeholder: '请输入社保内费用备注'
  },
  {
    name: 'text',
    field: 'outCostDetail3',
    title: '社保外费用明细',
    span: 12,
    placeholder: '请输入社保外费用明细'
  },
  {
    name: 'text',
    field: 'outCostRemark3',
    title: '社保外费用备注',
    span: 12,
    placeholder: '请输入社保外费用备注'
  },
]

const formRule = {
  billName2: [{ required: true, message: "必填项" }],
  hospital2: [{ required: true, message: "必填项" }],
  hospitalLevel2: [{ required: true, message: "必填项" }],
  hospitalProperty2: [{ required: true, message: "必填项" }],
  isPosition2: [{ required: true, message: "必填项" }],
  reimbursementType2: [{ required: true, message: "必填项" }],
  medicalType2: [{ required: true, message: "必填项" }],
  ticketType2: [{ required: true, message: "必填项" }],
  ticketCode2: [{ required: true, message: "必填项" }],
  ticketNumber2: [{ required: true, message: "必填项" }],
  checkCode2: [{ required: true, message: "必填项" }],
  credentialEndDate2: [{ required: true, message: "必填项" }],
  ticketState2: [{ required: true, message: "必填项" }],
  ticketRemark2: [{ required: true, message: "必填项" }],
  ticketStartDate2: [{ required: true, message: "必填项" }],
  ticketEndDate2: [{ required: true, message: "必填项" }],
  ticketInputType2: [{ required: true, message: "必填项" }],
  ticketBalance2: [{ required: true, message: "必填项" }],
  socialInsurancePayment2: [{ required: true, message: "必填项" }],
  // thirdPartyAmount2: [{ required: true, message: "必填项" }],
  // thirdPartyAgency2: [{ required: true, message: "必填项" }],
  treatMentDays2: [{ required: true, message: "必填项" }],
}

const desserts2 = [
  {
    billName2: "",
    treatmentType2: "",
    hospitalName2: "",
    startTime2: "",
    endTime2: "",
    days2: "",
    originalMoney: "",
    adjustMoney: "",
    selfMoney: "",
    selfPayMoney: "",
    otherMoney: "",
    reimbursementMoney: ""
  },
  {
    billName2: "",
    treatmentType2: "",
    hospitalName2: "",
    startTime2: "",
    endTime2: "",
    days2: "",
    originalMoney: "",
    adjustMoney: "",
    selfMoney: "",
    selfPayMoney: "",
    otherMoney: "",
    reimbursementMoney: ""
  }
]

const desserts3 = [
  {
    projectName2: "床位费",
    totalMoney2: "",
    jia2: "",
    yi2: "",
    partCost2: "",
    allCost2: "",
    selfTotalCost2: "",
    other2: "",
    reasonableCost2: "",
    inCost2: "",
    outCost2: ""
  },
  {
    projectName2: "治疗费",
    totalMoney2: "",
    jia2: "",
    yi2: "",
    partCost2: "",
    allCost2: "",
    selfTotalCost2: "",
    other2: "",
    reasonableCost2: "",
    inCost2: "",
    outCost2: ""
  },
  {
    projectName2: "手术费",
    totalMoney2: "",
    jia2: "",
    yi2: "",
    partCost2: "",
    allCost2: "",
    selfTotalCost2: "",
    other2: "",
    reasonableCost2: "",
    inCost2: "",
    outCost2: ""
  },
  {
    projectName2: "护理费",
    totalMoney2: "",
    jia2: "",
    yi2: "",
    partCost2: "",
    allCost2: "",
    selfTotalCost2: "",
    other2: "",
    reasonableCost2: "",
    inCost2: "",
    outCost2: ""
  },
  {
    projectName2: "卫生材料费",
    totalMoney2: "",
    jia2: "",
    yi2: "",
    partCost2: "",
    allCost2: "",
    selfTotalCost2: "",
    other2: "",
    reasonableCost2: "",
    inCost2: "",
    outCost2: ""
  },
  {
    projectName2: "西药费",
    totalMoney2: "",
    jia2: "",
    yi2: "",
    partCost2: "",
    allCost2: "",
    selfTotalCost2: "",
    other2: "",
    reasonableCost2: "",
    inCost2: "",
    outCost2: ""
  },
  {
    projectName2: "中药饮片费",
    totalMoney2: "",
    jia2: "",
    yi2: "",
    partCost2: "",
    allCost2: "",
    selfTotalCost2: "",
    other2: "",
    reasonableCost2: "",
    inCost2: "",
    outCost2: ""
  }
]

export default {
  form2,
  form3,
  formRule,
  headers1,
  headers2,
  headers3,
  desserts2,
  desserts3,
}