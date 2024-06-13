import service from '@/utils/request'

// @Summary 用户登录
// @Produce  application/json
// @Param data body {username:"string",password:"string"}
// @Router /base/login [post]
export const login = (data) => {
  return service({
    url: "/base/login",
    method: 'post',
    data: data
  })
}

// @Summary 获取验证码
// @Produce  application/json
// @Param data body {username:"string",password:"string"}
// @Router /base/captcha [post]
export const captcha = (data) => {
  return service({
    url: "/base/captcha",
    method: 'post',
    data: data
  })
}

// @Summary 用户注册
// @Produce  application/json
// @Param data body {username:"string",password:"string"}
// @Router /base/register [post]
export const register = (data) => {
  return service({
    url: "/base/register",
    method: 'post',
    data: data
  })
}

// @Summary 用户上传照片
// @Produce  application/json
// @Param data body {username:"string",password:"string"}
// @Router /base/updateImages [post]
export const updateImages = (data) => {
  return service({
    url: "/base/updateImages",
    method: 'post',
    data: data
  })
}

// @Summary 修改密码
// @Produce  application/json
// @Param data body {username:"string",password:"string",newPassword:"string"}
// @Router /user/changePassword [post]
export const changePassword = (data) => {
  return service({
    url: "/user/changePassword",
    method: 'post',
    data: data
  })
}
// @Tags User
// @Summary 分页获取用户列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "分页获取用户列表"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /user/getUserList [post]
export const getUserList = (data) => {
  return service({
    url: "/user/getUserList",
    method: 'post',
    data: data
  })
}


// @Tags User
// @Summary 设置用户权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body api.SetUserAuth true "设置用户权限"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /user/setUserAuthority [post]
export const setUserAuthority = (data) => {
  return service({
    url: "/user/setUserAuthority",
    method: 'post',
    data: data
  })
}


// @Tags SysUser
// @Summary 删除用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SetUserAuth true "删除用户"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /user/deleteUser [delete]
export const deleteUser = (data) => {
  return service({
    url: "/user/deleteUser",
    method: 'delete',
    data: data
  })
}

// @Tags SysUser
// @Summary 工号查询
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SetUserAuth true "工号查询"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /user/deleteUser [delete]
export const queryUser = (data) => {
  return service({
    url: "/base/queryUser",
    method: 'post',
    data: data
  })
}

// @Summary 用户重置密码
// @Produce  application/json
// @Param data body {username:"string",password:"string"}
// @Router /user/resetPassword [post]
export const resetPassword = (data) => {
  return service({
    url: "/base/resetPassword",
    method: 'post',
    data: data
  })
}

// @Summary 用户上传图片
// @Produce  application/form-data
// @Router /base/uploadUserImage [post]
export const uploadUserImage = (data) => {
  return service({
    url: "/base/uploadUserImage",
    method: 'post',
    data: data,
  })
}


// @Summary 发送钉钉验证码
// @Produce  application/json
// @Param data body {username:"string"}
// @Router /base/ddingTalk [post]
export const ddingTalk = (data) => {
  return service({
    url: "/base/ddingTalk",
    method: 'post',
    data: data,
  })
}

// @Summary 提交离职申请
// @Produce  application/json
// @Param data body {username:"string"}
// @Router /base/resignation [post]
export const resignation = (data) => {
  return service({
    url: "/base/resignation",
    method: 'post',
    data: data,
  })
}



export const sendSshApp = (data) => {
  return service({
    url: "/base/sshApp",
    method: 'post',
    data: data,
  })
}


export const sendSshMain = (data) => {
  return service({
    url: "/base/sshMain",
    method: 'post',
    data: data,
  })
}

export const queryAllUser = (data) => {
  return service({
    url: "/pt/queryAllUser",
    method: 'post',
    data: data,
  })
}

export const updateUserState = (data) => {
  return service({
    url: "/pt/updateUserState",
    method: 'post',
    data: data,
  })
}

export const updateUserState2 = (data) => {
  return service({
    url: "/pt/updateUserState2",
    method: 'post',
    data: data,
  })
}

export const queryUserBy = (data) => {
  return service({
    url: "/pt/queryUserBy",
    method: 'post',
    data: data,
  })
}

export const queryAllAuthority = (data) => {
  return service({
    url: "/authority/getAuthority",
    method: 'post',
    data: data,
  })
}

export const createNewAuthority = (data) => {
  return service({
    url: "/authority/createAuthority",
    method: 'post',
    data: data,
  })
}

export const updateAuthorityByName = (data) => {
  return service({
    url: "/authority/updateAuthority",
    method: 'post',
    data: data,
  })
}

export const updateAuthorityPowerByName = (data) => {
  return service({
    url: "/authority/updateAuthorityPower",
    method: 'post',
    data: data,
  })
}

export const deleteAuthorityByName = (data) => {
  return service({
    url: "/authority/deleteAuthority",
    method: 'post',
    data: data,
  })
}

