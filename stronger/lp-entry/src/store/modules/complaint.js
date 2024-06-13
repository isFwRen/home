import { request } from '@/api/service'

const actions = {
  // 列表
  async GET_COMPLAINT_LIST({}, params) {
    const data = {
      pageIndex: params.pageIndex,
      pageSize: params.pageSize,  
      proCode : params.proCode,  
      month: params.month,  
      billName: params.billName,  
      wrongFieldName: params.wrongFieldName
    }

    const result = await request({
      url: 'task/customer_complaints/list',
      params: data
    })

    return result
  }
}

export default {
  actions
}