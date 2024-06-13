import { request } from "@/api/service";
import moment from "moment";
import { R } from "vue-rocket";
const TODAY = moment().format("YYYY-MM-DD");

const actions = {
  async GET_DESTRUCTION_REPORT_LIST({ }, params) {
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
    };

    const result = await request({
      url: "pro-manager/bill-del-list/page",
      params: data
    });

    return result;
  },

  // 导出
  async EXPORT_DESTRUCTION_REPORT_EXCEL({ }, params) {
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
      // pageIndex: 1,
      // pageSize: 10,
    };

    const result = await request({
      url: "pro-manager/bill-del-list/export",
      params: data,
      responseType: "blob"
    });

    return result;
  },

  async EXPORT_DESTRUCTION_REPORT_WORD({ }, params) {
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
      // pageIndex: 1,
      // pageSize: 10,
    };

    const result = await request({
      url: "pro-manager/bill-del-list/sum",
      params: data,
      responseType: "blob"
    });
    return result;
  },

  async GET_HISTORY_DETAIL({ }, params) {
    const data = {
      proCode: params.proCode,
    };

    const result = await request({
      url: "/pro-manager/bill-del-list/history",
      params: data,
      responseType: "blob"
    });
    return result;
  }
};

export default {
  actions
};
