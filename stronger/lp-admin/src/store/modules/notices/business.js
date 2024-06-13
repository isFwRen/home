import Vue from "vue";
import { request } from "@/api/service";

const state = {
  timeBriefs: {},
  businessNotice: {
    count: 0,
    proCode: "",
    title: "",
    CreatedAt: "",
    msg: "",
    ID: '',
    pushId: ''
  },
};

const getters = {
  businessNotice: () => state.businessNotice,
  timeBriefs: () => state.timeBriefs
};

const mutations = {
  BUSINESS_NOTIFICATION_UPDATE_ITEM(state, data) {
    const newData = { ...state.businessNotice, ...data };
    Vue.set(state, "businessNotice", newData);
  },
  INCREAT_COUNT(state) {
    state.businessNotice.count++
  },
  CLEAR_NOTIFICATION_PUSHID() {
    state.businessNotice.pushId = ''
  },
  NOTIFICATION_TIMEBRIEFS_UPDATE(state, data) {
    state.timeBriefs = data
  }
};

const actions = {
  // 更新已读通知
  async UPDATE_NOTIFICATION_READ({ }, data) {
    const result = await request({
      url: "msg-manager/business-push/read",
      method: 'post',
      data
    });

    return result;
  },
  // 获取消息列表
  async GET_BUSINESS_LIST({ }, params) {

    let timeArr = {}

    if (params.firstTime) {
      timeArr = {
        startTime: (params.time && params.time[0] + "T00:00:00.000Z") || "",
        endTime: (params.time && params.time[1] + "T12:00:00.000Z") || "",
      }
    } else {
      timeArr = {
        startTime: (params.time && params.time[0] + "T00:00:00.000Z") || "",
        endTime: (params.time && params.time[1] + "T23:59:59.000Z") || "",
      }
    }


    const data = {
      startTime: timeArr.startTime,
      endTime: timeArr.endTime,
      proCode: params.proCode,
      pageSize: params.pageSize,
      pageIndex: params.pageIndex,
    }

    if (params.msgType) {
      data.msgType = params.msgType
    }

    const result = await request({
      url: "msg-manager/business-push/page",
      params: data
    });

    return result;
  },

  // 获取常量
  async GET_DISCONECT() {
    const result = await request({
      url: "pro-manager/bill-list/dict-const",
    });

    return result;
  }

};

export default {
  state,
  getters,
  mutations,
  actions
};
