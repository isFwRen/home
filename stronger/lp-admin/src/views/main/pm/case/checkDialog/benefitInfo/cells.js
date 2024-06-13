const headers = [
  { text: "受益人姓名", value: "beneficiaryName", width: 190 },
  { text: "理赔受益人类型", value: "beneficiaryType" },
  { text: "领款人姓名", value: "getMoneyerName" },
  { text: "受益金额", value: "beneficiaryMoney" },
  { text: "收益比例", value: "incomeRation" },
  { text: "支付方式", value: "payment" },
];

const form1 = []
const form2 = [
  {
    name: 'radio',
    field: 'beneficiary2',
    title: '',
    span: 24,
    circle: [
      {
        label: '自然人',
        content: '自然人'
      },
      {
        label: '法人',
        content: '法人'
      },
    ]
  },
]
const form3 = [
  {
    name: 'input',
    field: 'name3',
    title: '姓名',
    span: 8,
    placeholder: '请输入姓名'
  },
  {
    name: 'select',
    field: 'insuranceRelation3',
    title: '与被保人关系',
    span: 8,
    placeholder: '请选择关系',
    options: [
      {
        label: '夫妻',
        value: '夫妻'
      },
      {
        label: '母子',
        value: '母子'
      },
    ]
  },
  {
    name: 'select',
    field: 'getMoneyerRelation3',
    title: '与领款人关系',
    span: 8,
    placeholder: '请选择关系',
    options: [
      {
        label: '夫妻',
        value: '夫妻'
      },
      {
        label: '父子',
        value: '父子'
      },
    ]
  },
  {
    name: 'select',
    field: 'sex3',
    title: '性别',
    span: 8,
    placeholder: '请选择性别',
    options: [
      {
        label: '男',
        value: '男'
      },
      {
        label: '女',
        value: '女'
      },
    ]
  },
  {
    name: 'select',
    field: 'beneficiary3',
    title: '受益人职业',
    span: 8,
    placeholder: '请选择职业',
    options: [
      {
        label: '务农',
        value: '务农'
      },
      {
        label: '军人',
        value: '军人'
      },
    ]
  },
  {
    name: 'date',
    field: 'birthTime3',
    title: '出生日期',
    span: 8,
    placeholder: '请选择日期',
  },
  {
    name: 'select',
    field: 'credentialType3',
    title: '证件类型',
    span: 8,
    placeholder: '请选择证件类型',
    options: [
      {
        label: '身份证',
        value: '身份证'
      },
      {
        label: '港澳通行证',
        value: '港澳通行证'
      },
    ]
  },
  {
    name: 'input',
    field: 'credentialNumber3',
    title: '证件号码',
    span: 8,
    placeholder: '请输入证件号码'
  },
  {
    name: 'date',
    field: 'credentialStartDate3',
    title: '证件有效起期',
    span: 8,
    placeholder: '请选择日期',
  },
  {
    name: 'date',
    field: 'credentialEndDate3',
    title: '证件有效止期',
    span: 8,
    placeholder: '请选择日期',
  },
  {
    name: 'select',
    field: 'nation3',
    title: '国籍',
    span: 8,
    placeholder: '请选择国籍',
    options: [
      {
        label: '中国',
        value: '中国'
      },
      {
        label: '美国',
        value: '美国'
      },
    ]
  },
  {
    name: 'input',
    field: 'telephone3',
    title: '联系电话',
    span: 8,
    placeholder: '请输入联系电话'
  },
  {
    name: 'select',
    field: 'province3',
    title: '省份',
    span: 8,
    placeholder: '请选择省份/直辖市/自治区',
    options: [
      {
        label: '广东',
        value: '广东'
      },
      {
        label: '河南',
        value: '河南'
      },
    ]
  },
  {
    name: 'select',
    field: 'city3',
    title: '城市',
    span: 8,
    placeholder: '请选择城市',
    options: [
      {
        label: '珠海',
        value: '珠海'
      },
      {
        label: '深圳',
        value: '深圳'
      },
    ]
  },
  {
    name: 'select',
    field: 'district3',
    title: '区县',
    span: 8,
    placeholder: '请选择区县',
    options: [
      {
        label: '香洲',
        value: '香洲'
      },
      {
        label: '斗门',
        value: '斗门'
      },
    ]
  },
  {
    name: 'input',
    field: 'detailAddress3',
    title: '详细地址',
    span: 24,
    placeholder: '请输入详细地址',
  },
]
const form4 = [
  {
    name: 'input',
    field: 'name4',
    title: '姓名',
    span: 8,
    placeholder: '请输入姓名'
  },
  {
    name: 'select',
    field: 'sex4',
    title: '性别',
    span: 8,
    placeholder: '请选择性别',
    options: [
      {
        label: '男',
        value: '男'
      },
      {
        label: '女',
        value: '女'
      },
    ]
  },
  {
    name: 'select',
    field: 'nation4',
    title: '国籍',
    span: 8,
    placeholder: '请选择国籍',
    options: [
      {
        label: '中国',
        value: '中国'
      },
      {
        label: '美国',
        value: '美国'
      },
    ]
  },
  {
    name: 'select',
    field: 'credentialTyp4',
    title: '证件类型',
    span: 8,
    placeholder: '请选择证件类型',
    options: [
      {
        label: '身份证',
        value: '身份证'
      },
      {
        label: '港澳通行证',
        value: '港澳通行证'
      },
    ]
  },
  {
    name: 'input',
    field: 'credentialNumber4',
    title: '证件号码',
    span: 8,
    placeholder: '请输入证件号码'
  },
  {
    name: 'date',
    field: 'credentialStartDate4',
    title: '证件有效起期',
    span: 8,
    placeholder: '请选择日期',
  },
  {
    name: 'date',
    field: 'credentialEndDate4',
    title: '证件有效止期',
    span: 8,
    placeholder: '请选择日期',
  },
  {
    name: 'date',
    field: 'birthTime4',
    title: '出生日期',
    span: 8,
    placeholder: '请选择日期',
  },
  {
    name: 'input',
    field: 'payment4',
    title: '支付方式',
    span: 8,
    placeholder: '请输入支付方式'
  },
  {
    name: 'input',
    field: 'accountName4',
    title: '银行户名',
    span: 8,
    placeholder: '请输入银行户名'
  },
  {
    name: 'input',
    field: 'bankName4',
    title: '银行名称',
    span: 8,
    placeholder: '请输入银行名称'
  },
  {
    name: 'input',
    field: 'accountBank4',
    title: '银行账号',
    span: 8,
    placeholder: '请输入银行账号'
  },
]

const formRule = {
  name2: [{ required: true, message: "必填项" }],
  sex2: [{ required: true, message: "必填项" }],
  birthTime2: [{ required: true, message: "必填项" }],
  career2: [{ required: true, message: "必填项" }],
  nation2: [{ required: true, message: "必填项" }],
  credentialType2: [{ required: true, message: "必填项" }],
  credentialNumber2: [{ required: true, message: "必填项" }],
  credentialStartDate2: [{ required: true, message: "必填项" }],
  credentialEndDate2: [{ required: true, message: "必填项" }],
  telephone2: [{ required: true, message: "必填项" }],
  insuranceRelation2: [{ required: true, message: "必填项" }],
  policyholderRelation2: [{ required: true, message: "必填项" }],
  province2: [{ required: true, message: "必填项" }],
  city2: [{ required: true, message: "必填项" }],
  district2: [{ required: true, message: "必填项" }],
  detailAddress2: [{ required: true, message: "必填项" }],

  name3: [{ required: true, message: "必填项" }],
  sex3: [{ required: true, message: "必填项" }],
  birthTime3: [{ required: true, message: "必填项" }],
  career3: [{ required: true, message: "必填项" }],
  nation3: [{ required: true, message: "必填项" }],
  credentialType3: [{ required: true, message: "必填项" }],
  credentialNumber3: [{ required: true, message: "必填项" }],
  credentialStartDate3: [{ required: true, message: "必填项" }],
  credentialEndDate3: [{ required: true, message: "必填项" }],
  telephone3: [{ required: true, message: "必填项" }],
  insuranceRelation3: [{ required: true, message: "必填项" }],
  policyholderRelation3: [{ required: true, message: "必填项" }],
  province3: [{ required: true, message: "必填项" }],
  city3: [{ required: true, message: "必填项" }],
  district3: [{ required: true, message: "必填项" }],
  detailAddress3: [{ required: true, message: "必填项" }],
  beneficiary3: [{ required: true, message: "必填项" }],
  getMoneyerRelation3: [{ required: true, message: "必填项" }],

  principalType4: [{ required: true, message: "必填项" }],
  baileeName4: [{ required: true, message: "必填项" }],
  sex4: [{ required: true, message: "必填项" }],
  credentialType4: [{ required: true, message: "必填项" }],
  credentialNumber4: [{ required: true, message: "必填项" }],
  telephone4: [{ required: true, message: "必填项" }],
  payment4: [{ required: true, message: "必填项" }],
  accountName4: [{ required: true, message: "必填项" }],
  bankName4: [{ required: true, message: "必填项" }],
  accountBank4: [{ required: true, message: "必填项" }],
  credentialStartDate4: [{ required: true, message: "必填项" }],
  credentialEndDate4: [{ required: true, message: "必填项" }],
  nation4: [{ required: true, message: "必填项" }],
  name4: [{ required: true, message: "必填项" }],
  birthTime4: [{ required: true, message: "必填项" }],
  credentialTyp4: [{ required: true, message: "必填项" }],
}

export default {
  form1,
  form2,
  form3,
  form4,
  formRule,
  headers
}