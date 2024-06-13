import { request } from '@/api/service'
import { sessionStorage, localStorage } from 'vue-rocket'
import { tools as lpTools } from "@/libs/util";
const isIntranet = lpTools.isIntranet();
const state = {
  task: {
    baseURL: '',
    topInfo: {},
    rowInfo: {},
    config: {},
    prompt: '',
    mainHeight: 0,
    character: 0,
    accuracyRate: 0,
  },

  f3state: true
}

const getters = {
  task: () => state.task
}

const mutations = {
  UPDATE_CHANNEL(state, data) {
    const localTask = sessionStorage.get('task') || {}
    state.task = { ...state.task, ...localTask, ...data }

    sessionStorage.set('task', state.task)
    localStorage.set('task', state.task)
  },

  UPDATE_F3STATE(state, data) {
    state.f3state = data
  },

  UPDATE_PRACTICE(state, data) {
    state.task = { ...state.task, mainHeight: 0, displayPrompt: true, displayTop: false, displayRight: true, ...data }
    sessionStorage.set('task', state.task)
    localStorage.set('task', state.task)
  },

  UPDATE_EXAM(state, data) {
    state.task = { ...state.task, mainHeight: 0, displayPrompt: false, displayTop: false, prompt: "", displayRight: false, ...data }
    sessionStorage.set('task', state.task)
    localStorage.set('task', state.task)
  },
}

const actions = {
  // 录入通道列表
  async INPUT_GET_LIST() {
    const result = await request({
      url: 'data-entry/channel/list'
    })

    return result
  },

  // transit图标
  async INPUT_GET_TRANSIT_LIST() {
    const baseUrl = isIntranet ? process.env.VUE_APP_ADMIN_INNER_URL : process.env.VUE_APP_ADMIN_OUTER_URL
    const result = await request({
      baseURL: baseUrl,
      url: '/api/data-entry/channel/list'
    })

    return result
  },

  // 顶部信息
  async INPUT_GET_OPS_INFO({ }, params) {
    const data = {
      code: params.code
    }
    if (sessionStorage.get('task')?.baseURL == '') return
    const result = await request({
      baseURL: sessionStorage.get('task')?.baseURL,
      url: 'task/opNum',
      params: data
    })

    return result
  },

  // 理赔领单
  async INPUT_GET_TASK({ }, params) {
    const data = {
      code: params.code,
      op: params.op
    }

    const result = await request({
      baseURL: sessionStorage.get('task')?.baseURL,
      url: 'task/op',
      params: data,
      headers: {
        process: params.op
      }
    })

    return result
  },

  // 获取图片大小
  async INPUT_GET_IMAGE_SIZE({ }, params) {
    const data = {
      path: params.path
    }

    const result = await request({
      baseURL: sessionStorage.get('task')?.baseURL,
      url: 'task/getImageSize',
      params: data,
      headers: {
        process: params.op
      }
    })

    return result
  },

  // 返回重录
  async INPUT_GET_PREVIOUS_TASK({ }, params) {
    const data = {
      code: params.code,
      op: params.op,
      num: params.prevNums
    }

    const result = await request({
      baseURL: sessionStorage.get('task')?.baseURL,
      url: 'task/modifyBlock',
      params: data,
      headers: {
        process: params.op
      }
    })

    return result
  },

  // 配置
  async INPUT_GET_TASK_CONFIG() {
    const result = await request({
      baseURL: sessionStorage.get('task')?.baseURL || localStorage.get('task')?.baseURL,
      url: 'task/conf'
    })

    return result
  },

  // 提交操作后的图片
  async INPUT_SUBMIT_MODIFIED_IMAGE({ }, body) {
    const formData = new FormData()

    for (let key in body) {
      if (key !== 'op')
        formData.append(key, body[key])
    }

    const result = await request({
      baseURL: sessionStorage.get('task')?.baseURL,
      method: 'POST',
      url: 'task/uploadImage',
      data: formData,
      headers: {
        process: body.op
      }
    })

    return result
  },

  // 提交
  async INPUT_SUBMIT_TASK({ }, body) {
    const data = {
      bill: body.bill,
      block: body.block,
      fields: body.fields,
      op: body.op
    }
    console.log('F8提交Data--------', data);
    const result = await request({
      baseURL: sessionStorage.get('task')?.baseURL,
      method: 'POST',
      url: 'task/submit',
      data,
      headers: {
        process: body.op
      }
    })

    return result
  },

  // KeyLog 录入日志
  async INPUT_SUBMIT_TASK_KEY_OPERATION({ }, body) {
    const data = {
      billNum: body.billNum,
      log: body.log,
    }
    console.log('F8录入操作--------', data);
    const result = await request({
      baseURL: sessionStorage.get('task')?.baseURL,
      method: 'POST',
      url: 'task/keyLog',
      data,
      headers: {
        process: 'op0'
      }
    })

    return result
  },

  // 释放任务
  async TASK_ALLOCATION_TASK({ }, body) {
    const data = {
      id: body.id,
      op: body.op,
      code: body.code
    }

    const result = await request({
      method: 'POST',
      baseURL: sessionStorage.get('task')?.baseURL,
      url: 'task/releaseBill',
      data,
      headers: {
        process: ''
      },
      effect: { task: true }
    })

    return result
  },

  async ALLOCATION_ALL_TASK({ }, body) {
    const data = {
      code: body.code
    }

    const result = await request({
      method: 'POST',
      baseURL: localStorage.get('task')?.baseURL,
      url: 'task/releaseExitBlock',
      data,
      headers: {
        process: body.op,
      },
      effect: { task: true }
    })

    return result
  },

  async INPUT_SUBMIT_TASK({ }, body) {
    const data = {
      bill: body.bill,
      block: body.block,
      fields: body.fields,
      op: body.op
    }
    console.log('F8提交Data--------', data);
    const result = await request({
      baseURL: sessionStorage.get('task')?.baseURL,
      method: 'POST',
      url: 'task/submit',
      data,
      headers: {
        process: body.op
      }
    })

    return result
  },

  // 练习列表
  async PRACTICE_CODE_LIST() {
    const result = await request({
      url: '/practice/ask-list',
    })

    return result
  },

  // 练习
  async PRACTICE_LIST({ }, params) {
    const data = {
      code: params,
      name: localStorage.get('user').name
    }
    const result = await request({
      url: '/practice/get',
      params: data,
    })

    return result
  },

  // 提交练习
  async SUBMIT_PRACTICE_TASK({ }, params) {
    const data = params

    const result = await request({
      method: 'post',
      url: '/practice/submit',
      data,
    })

    return result
  },

  // 退出练习
  async EXIT_PRACTICE_TASK({ }, params) {
    const data = params

    const result = await request({
      method: 'post',
      url: '/practice/exit',
      data,
    })

    return result
  },

  // 考核列表
  async EXAM_GET_LIST() {
    const result = await request({
      url: '/assessment-exam/test-procedure/get-project-list'
    })

    return result
  },

  // 考核题目
  async GET_EXAM_TASK({ }, params) {
    const data = {
      projectCode: params,
    }

    const result = await request({
      method: 'post',
      url: '/assessment-exam/test-procedure/start-exam',
      data,
    })

    return result
  },

  // 提交考核内容
  async SUBMIT_EXAM_TASK({ }, data) {
    const result = await request({
      method: 'post',
      url: '/assessment-exam/test-procedure/end-exam',
      data,
    })

    return result
  },

  // 练习
}

export default {
  state,
  getters,
  mutations,
  actions
}