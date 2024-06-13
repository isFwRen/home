import {
  request
} from "@/api/service";
import moment from "moment";
import { localStorage } from "vue-rocket";

const state = {
  yieldInfo: {
    staffYield: {},
  },
  pro: {
    proOptions: [],
    proMap: {}
  }
};

const getters = {
  yieldInfo: () => state.yieldInfo,
  pro: () => state.pro,
};

const mutations = {
  UPDATE_PRO(state, data) {
    const {
      pro
    } = state;
    state.pro = Object.assign(pro, data);
  },
};

const actions = {

  async GET_STAFF_TOTAL({ }, param) {
    const data = {
      pageIndex: param.pageIndex,
      pageSize: param.pageSize,
      proCode: param.proCode,
      startTime: (param.date ? param.date[0] : moment().format("yyyy-MM-DD")) + " 00:00:00",
      endTime: (param.date ? param.date[1] : moment().format("yyyy-MM-DD")) + " 23:59:59",
      isCheckAll: 2,
    };

    const result = await request({
      //baseURL: 'http://127.0.0.1:9999/',
      url: "report-management/output-Statistics-task/list",
      params: data,
    });
    return result
  },

  // 练习错误
  async SUM_GET_LIST({ }, params) {
    const data = {
      pageIndex: params.pageIndex,
      pageSize: params.pageSize,
      proCode: params.proCode,
      name: params.name,
      code: localStorage.get('user').code,
      startTime: (params.date ? params.date[0] : moment().format('YYYY-MM-DD')) + "T00:00:00.000Z",
      endTime: (params.date ? params.date[1] : moment().format('YYYY-MM-DD')) + "T23:59:59.000Z",
    }

    const result = await request({
      url: '/practice/sum',
      params: data
    })

    return result
  },

};

export default {
  state,
  getters,
  mutations,
  actions,
};