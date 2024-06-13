const form1 = [
  {
    name: 'input',
    field: 'caseNum1',
    title: '立案号',
    span: 8,
    placeholder: '请输入立案号'
  },
  {
    name: 'date',
    field: 'applyTime1',
    title: '申请日期',
    span: 8,
    placeholder: '请选择日期',
  },
  {
    name: 'input',
    field: 'manageAgency1',
    title: '管辖机构',
    span: 8,
    placeholder: '请输入管辖机构',
  },
  {
    name: 'checkbox',
    field: 'compensationType1',
    title: '理赔类型',
    span: 24,
    rect: [
      {
        label: '身故',
        content: '身故'
      },
      {
        label: '疾病',
        content: '疾病'
      },
      {
        label: '医疗',
        content: '医疗'
      },
      {
        label: '伤残/高残/全残',
        content: '伤残/高残/全残'
      },
    ]
  },
  {
    name: 'radio',
    field: 'caseType1',
    title: '案件类型',
    span: 24,
    circle: [
      {
        label: '普通案件',
        content: '普通案件'
      },
      {
        label: '简易案件',
        content: '简易案件'
      },
    ]
  },
]
const form2 = [
  {
    name: 'input',
    field: 'name2',
    title: '姓名',
    span: 8,
    placeholder: '请输入姓名'
  },
  {
    name: 'select',
    field: 'sex2',
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
    name: 'date',
    field: 'birthTime2',
    title: '出生日期',
    span: 8,
    placeholder: '请选择日期',
  },
  {
    name: 'select',
    field: 'career2',
    title: '职业',
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
    name: 'select',
    field: 'nation2',
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
    field: 'credentialType2',
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
    field: 'credentialNumber2',
    title: '证件号码',
    span: 8,
    placeholder: '请输入证件号码'
  },
  {
    name: 'date',
    field: 'credentialStartDate2',
    title: '证件有效起期',
    span: 8,
    placeholder: '请选择日期',
  },
  {
    name: 'date',
    field: 'credentialEndDate2',
    title: '证件有效止期',
    span: 8,
    placeholder: '请选择日期',
  },
  {
    name: 'input',
    field: 'telephone2',
    title: '联系电话',
    span: 8,
    placeholder: '请输入联系电话'
  },
  {
    name: 'select',
    field: 'insuranceRelation2',
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
    field: 'policyholderRelation2',
    title: '与投保人关系',
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
    field: 'province2',
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
    field: 'city2',
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
    field: 'district2',
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
    field: 'detailAddress2',
    title: '详细地址',
    span: 24,
    placeholder: '请输入详细地址',
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
    name: 'date',
    field: 'birthTime3',
    title: '出生日期',
    span: 8,
    placeholder: '请选择日期',
  },
  {
    name: 'select',
    field: 'career3',
    title: '职业',
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
    name: 'input',
    field: 'telephone3',
    title: '联系电话',
    span: 8,
    placeholder: '请输入联系电话'
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
    field: 'policyholderRelation3',
    title: '与投保人关系',
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
    name: 'select',
    field: 'principalType4',
    title: '委托人类型',
    span: 8,
    placeholder: '请选择委托人类型',
    options: [
      {
        label: '银行职员',
        value: '银行职员'
      },
      {
        label: '银行经理',
        value: '银行经理'
      },
    ]
  },
  {
    name: 'input',
    field: 'baileeName4',
    title: '受托人姓名',
    span: 8,
    placeholder: '请输入受托人姓名'
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
    name: 'input',
    field: 'telephone4',
    title: '联系电话',
    span: 8,
    placeholder: '请输入联系电话'
  },
]

const formRule = {
  caseNum1: [{ required: true, message: "请输入立案号" }],
  applyTime1: [{ required: true, message: "请选择申请日期" }],
  manageAgency1: [{ required: true, message: "请输入管辖机构" }],
  compensationType1: [{ required: true, message: "请选择理赔类型" }],
  caseType1: [{ required: true, message: "请选择案件类型" }],

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

  principalType4: [{ required: true, message: "必填项" }],
  baileeName4: [{ required: true, message: "必填项" }],
  sex4: [{ required: true, message: "必填项" }],
  credentialType4: [{ required: true, message: "必填项" }],
  credentialNumber4: [{ required: true, message: "必填项" }],
  telephone4: [{ required: true, message: "必填项" }],
  credentialTyp4: [{ required: true, message: "必填项" }],
}

export default {
  form1,
  form2,
  form3,
  form4,
  formRule
}