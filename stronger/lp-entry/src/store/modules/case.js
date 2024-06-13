import { request } from '@/api/service'

const actions = {
  // 删除
  async DELETE_CASE_ITEM({}, form) {
    const data = {
      proCode: form.proCode,
      id: form.id,
      delRemarks: form.delRemarks
    }

    const result = await request({
      method: 'DELETE',
      url: 'pro-manager/bill-list/delete',
      data
    })

    return result
  }
}

export default {
  actions
}