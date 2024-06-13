const headers = [
  { text: '项目编码', value: 'proCode' },
  { text: '规则名称', value: 'ruleName' },
  { text: '规则类型', value: 'ruleType' },
  { text: '更新日期', value: 'UpdatedAt' },
  { text: '操作', value: 'options', width: 260 }
]

const ruleTypes = [
  { label: '项目规则', value: '项目规则' },
  { label: '易错规则', value: '易错规则' }
]

const projectCode = [
  { label: 'B0102', value: 'B0102' },
  { label: 'B0103', value: 'B0103' },
  { label: 'B0106', value: 'B0106' },
  { label: 'B0108', value: 'B0108' },
  { label: 'B0110', value: 'B0110' },
  { label: 'B0113', value: 'B0113' },
  { label: 'B0114', value: 'B0114' },
  { label: 'B0116', value: 'B0116' },
  { label: 'B0118', value: 'B0118' },
  { label: 'B0121', value: 'B0121' },
  { label: 'B0122', value: 'B0122' },
]

const rule = [
  { label: '项目规则', value: '项目规则' },
  { label: '易错规则', value: '易错规则' },
]

export default {
  headers,
  ruleTypes,
  projectCode,
  rule
}