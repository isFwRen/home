import moment from 'moment'
import { request } from '@/api/service'

const actions = {
	async GET_SALARY_PT_LIST({}, params) {
		const data = {
			pageIndex: params.pageIndex,
			pageSize: params.pageSize,
			start: moment(params.date[0]).format('YYYY-MM'),
			end : moment(params.date[1]).format('YYYY-MM'),
		}
		const result = await request({
			url: 'report-management/pt/salary-task/list',
			params: data,
		})
		return result
	},
}
export default {
    actions
}