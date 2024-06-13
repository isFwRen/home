import { request } from "@/api/service";
import moment from "moment";
import { R } from "vue-rocket";
const TODAY = moment().format("YYYY-MM-DD");

const actions = {
  async GET_DEDUCTION_DETAILS_LIST({ }, params) {
    const startTime = moment(R.isYummy(params.time) ? params.time[0] : TODAY).format(
      "YYYY-MM-DD"
    );
    const endTime = moment(R.isYummy(params.time) ? params.time[1] : TODAY).format(
      "YYYY-MM-DD"
    );

    const data = {
      pageIndex: params.pageIndex,
      pageSize: params.pageSize,
      proCode: params.proCode,
      timeStart: startTime + 'T00:00:00Z',
      timeEnd: endTime + 'T23:59:59Z',
      billCode: params.bill_num,
      name: params.list_name,
    };

    const result = await request({
      url: "pro-manager/bill-deduction_details/page",
      params: data
    });

    return result;
  },

  // 导出
  async EXPORT_DEDUCTION_DETAILS_EXCEL({ }, params) {
    const startTime = moment(R.isYummy(params.time) ? params.time[0] : TODAY).format(
      "YYYY-MM-DD"
    );
    const endTime = moment(R.isYummy(params.time) ? params.time[1] : TODAY).format(
      "YYYY-MM-DD"
    );

    const data = {
      timeStart: startTime + 'T00:00:00Z',
      timeEnd: endTime + 'T23:59:59Z',
      proCode: params.proCode,
      billCode: params.proCode,
      name: params.name,
    };

    const result = await request({
      url: "pro-manager/bill-deduction_details/export",
      params: data,
      responseType: "blob"
    });

    return result;
  },
};

export default {
  actions
};
