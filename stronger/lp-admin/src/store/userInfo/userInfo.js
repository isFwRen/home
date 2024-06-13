import { request } from "@/api/service";

const state = {

}

const actions = {
  // 修改用户头像
  async UPLOAD_AVATAR({ }, data) {
    const result = await request({
      method: "post",
      url: 'sys-user/upload-user/avatar',
      data,
    })
    return result
  },


  // 注销账号
  async DELETE_ACCOUNT({ }, data) {
    const result = await request({
      method: "post",
      url: 'sys-user-leave/user-leave/resignation',
      data
    })
    return result
  }


}



export default {
  state,
  actions
}
