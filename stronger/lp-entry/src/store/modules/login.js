import { request } from '@/api/service'

const actions = {
  // 获取钉钉验证码
  async GET_DING_CODE({ }, form) {
    // const data = {
    //   isIntranet: form.isIntranet,
    //   [form.accountKey]: form.account
    // }

    // const result = await request({
    //   method: 'POST',
    //   url: 'dinging/captcha',
    //   data
    // })
    // return result
    console.log(form, "forms");
    let forms = {}
    forms.phone = form.account
    const result = await request({
      method: "POST",
      url: "/sys-base/send-code",
      data: forms
    });
    return result;
  },

  // 登录
  async LOGIN({ }, form) {
    // const data = {
    //   isIntranet: form.isIntranet,
    //   [form.accountKey]: form.account,
    //   password: form.password,
    //   captcha: form.captcha,
    //   captchaId: form.captchaId
    // }

    // const result = await request({
    //   method: 'POST',
    //   url: 'sys-base/sys-login/login',
    //   data
    // })
    // return result

    const data = {
      //[form.accountKey]: form.account,
      password: form.password,
      // captcha: form.captcha,
      code: form.account
      //captchaId: form.captchaId
    };
    console.log(data, "data");
    const result = await request({
      method: "POST",
      url: "sys-base/login",
      data
    });
    return result
  },

  // 重置密码
  async RESET_PASSWORD({ }, form) {
    const data = {
      isIntranet: form.isIntranet,
      [form.accountKey]: form.account,
      captcha: form.captcha,
      captchaId: form.captchaId
    }

    const result = await request({
      method: 'POST',
      url: 'sys-base/sys-user-operation-pw/resetPassword',
      data
    })
    return result
  },

  // 修改密码
  async CHANGE_PASSWORD({ }, form) {
    // const data = {
    //   isIntranet: form.isIntranet,
    //   [form.accountKey]: form.account,
    //   oldpass: form.oldpass,
    //   newPass: form.checkPass
    // }

    // if (!form.isIntranet) {
    //   data.captchaId = form.captchaId
    //   data.captcha = form.captcha
    // } else {
    //   delete data.captchaId
    //   delete data.captcha
    // }

    // const result = await request({
    //   method: 'POST',
    //   url: 'sys-user/sys-user-operation-pw/changePassword',
    //   data
    // })
      const data = {
      isIntranet: form.isIntranet,
      [form.accountKey]: form.account,
      oldpass: form.oldpass,
      newPass: form.checkPass
    }
    console.log(data,'data')
  const result = await request({
      method: 'POST',
      url: 'sys-base/change-pwd-common',
      data
    })

    return result
  },

  async VERIFY_IDENTITY({ }, form) {
    const result = await request({
      method: 'POST',
      url: 'sys-base/sys-login/get-user-qrCode',
      data: form
    })
    return result
  }
}

export default {
  actions
}