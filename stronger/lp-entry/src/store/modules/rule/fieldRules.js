import { request } from '@/api/service'

const actions = {
  // 列表
  async GET_PM_TEACHING_FIELD_RULES_LIST({}, params) {
    const data = {
      pageSize: 1,
      pageIndex: 1,
      proCode: params.proCode,
      fieldsName: params.fieldsName,
      rule: '有'
    }

    const result = await request({
      url: 'pro-manager/fieldsRule/list',
      params: data
    })

    return result
  }
}

export default {
  actions
}