import { request } from "@/api/service";

const actions = {
  async HOME_GET_USER_YIELD({}) {
    const result = await request({
      url: "homepage/home/user-yield",
    });
    return result;
  },
  async HOME_GET_RANK_YIELD({}, params) {
    const data = {
      rankingType: params.rankingType || 0,
      pageIndex: params.pageIndex || 2,
      pageSize: params.pageSize || 10,
    };
    const result = await request({
      url: "homepage/home/ranking-yield",
      method: "GET",
      params: data,
    });

    return result;
  },
  async HOME_GET_ANNOUNCEMENT({}, params) {
    const data = {
      releaseType: params.releaseType || 1,
      pageIndex: params.pageIndex || 1,
      pageSize: params.pageSize || 10,
      startTime:
        (params.dayRange && params.dayRange[0] + "T16:00:00.000Z") || "",
      endTime: (params.dayRange && params.dayRange[1] + "T16:00:00.000Z") || "",
      title: params.title || "",
      proCode: params.proCode || "",
    };
    const result = await request({
      url: "homepage/home/announcement",
      method: "GET",
      params: data,
    });
    return result;
  },
  async HOME_SET_USER_YIELD({}, params) {
    const data = {
      target: +params.target || 0,
    };
    const result = await request({
      url: "homepage/home/set-target",
      method: "POST",
      data,
    });
    return result;
  },
  async WHATCH_NOTICE_VIEW_NUMBER({}, params) {
    const data = {
      id: params.id,
    };
    const result = await request({
      url: "homepage/home/announcement-view",
      method: "GET",
      params: data,
    });
    return result;
  },
};

export default {
  actions,
};
