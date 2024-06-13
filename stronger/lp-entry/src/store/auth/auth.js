import { request } from '@/api/service'
import { tools, localStorage } from 'vue-rocket'

const defaultAuth = {
  menus: [],
  perm: [],
  mapPro: new Map(),
  proItems: []
}

const state = {
  project: {
    id: '',
    code: ''
  },

  auth: tools.deepClone(defaultAuth)
}

const getters = {
  project: () => state.project,
  auth: () => state.auth
}

const mutations = {
  SET_PROJECT_INFO(state, project) {
    if (project) {
      state.project = project
    }
    else {
      state.project = {
        id: '',
        code: ''
      }
    }

    localStorage.set('project', state.project)
  },

  UPDATE_AUTH(state, auth) {
    if (tools.isYummy(auth)) {
      state.auth = auth
      localStorage.set('auth', state.auth)
    }
    else {
      const storageAuth = localStorage.get('auth')

      const auth = tools.isYummy(storageAuth) ? storageAuth : state.auth
      state.auth = auth
    }

    // console.log(state.auth) 
  }
}

const actions = {
  // 用户信息
  async GET_USER_INFO(data) {
    const result = await request({
      url: "/sys-base/sys-login/user-info",
      method: 'POST',
      data
    })

    return result

  },
  // 角色
  async GET_ROLE_SYS_MENU() {
    const result = await request({
      url: 'sys-menu/role/get'
    })

    return result
  }
}

export default {
  state,
  getters,
  mutations,
  actions
}