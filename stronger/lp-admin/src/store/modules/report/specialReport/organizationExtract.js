import { request } from '@/api/service'
import moment from 'moment'
import { R } from 'vue-rocket'
const TODAY = moment().format('YYYY-MM-DD')

const actions = {
  async GET_ORGANIZATION_EXTRACT_LIST({ }, params) {
    const startTime = moment(R.isYummy(params.time) ? params.time[0] : TODAY).format('YYYY-MM-DD HH:mm:ss')
    const endTime = moment(R.isYummy(params.time) ? params.time[1] : TODAY).format('YYYY-MM-DD HH:mm:ss')

    const data = {
      pageIndex: params.pageIndex,
      pageSize: params.pageSize,
      proCode: params.proCode,
      startTime,
      endTime,
      types: params.types
    }

    const result = await request({
      url: 'report-management/abnormal-bill/list',
      params: data
    })

    return result
  },

  // 导出
  async EXPORT_ABNORMAL_PART_BILL({ }, params) {
    const startTime = `${params.time ? params.time[0] : TODAY} 00:00:00`
    const endTime = `${params.time ? params.time[1] : TODAY} 23:59:59`

    const data = {
      startTime,
      endTime,
      proCode: params.proCode
    }

    const result = await request({
      url: 'report-management/abnormal-bill/export',
      params: data,
      responseType: "blob"
    })

    return result
  }
}

export default {
  actions
}