<template>
  <lp-dialog
    ref="dialog"
    title="修改密码"
    width="600"
    @dialog="handleDialog"
  >
    <div slot="main">
      <ul class="mt-4">
        <li class="mb-1">
          <z-text-field
            :formId="formId"
            :formKey="userField.formKey"
            :label="userField.label"
            :validation="userField.validation"
          >
          </z-text-field>
        </li>

        <!-- <li class="z-flex align-end mb-4" v-if="!isIntranet">
          <z-text-field
            :formId="formId"
            formKey="captcha"
            class="flex-grow-1"
            label="验证码"
            :validation="[{ rule: 'required', message: '请输入验证码!' }]"
          >
          </z-text-field>
          <z-btn 
            class="pb-5 ml-4" 
            :color="color"
            :disabled="!validAccount || counting"
            :lockedTime="2500"
            @click="sendCode"
          >{{ text }}</z-btn>
        </li> -->

        <li class="mb-1">
          <z-text-field
            :formId="formId"
            formKey="oldpass"
            :append-icon="originEye ? 'mdi-eye' : 'mdi-eye-off'"
            label="原始密码"
            :type="originEye ? 'text' : 'password'"
            :validation="[{ rule: 'required', message: '请输入原始密码!' }]"
            @click:append="originEye = !originEye"
          >
          </z-text-field>
        </li>

        <li class="mb-1">
          <z-text-field
            :formId="formId"
            formKey="checkPass"
            :append-icon="newEye ? 'mdi-eye' : 'mdi-eye-off'"
            label="新密码"
            :type="newEye ? 'text' : 'password'"
            :validation="[
              { rule: 'required', message: '请输入新密码!' },
              { regex: regex_newpassword_digital, message: '输入不符合密码为数字，字母和符号的组合，长度为9-16位要求的密码'},
              //{ regex: regex_newpassword_letter, message: '输入不符合密码为数字，字母和符号的组合，长度为9-16位要求的密码'},
              //{ regex: regex_newpassword_length, message: '输入不符合密码为数字，字母和符号的组合，长度为9-16位要求的密码'},
            ]"
            @click:append="newEye = !newEye"
          >
          </z-text-field>
        </li>

        <li class="mb-1">
          <z-text-field
            :formId="formId"
            formKey="rep"
            :append-icon="confirmEye ? 'mdi-eye' : 'mdi-eye-off'"
            label="确认密码"
            :type="confirmEye ? 'text' : 'password'"
            :validation="[
              { rule: 'required', message: '请输入确认密码!' },
              { rule: `is:${ checkPass }`, message: '两次输入的密码不一致!' }
            ]"
            @click:append="confirmEye = !confirmEye"
          >
          </z-text-field>
        </li>
      </ul>
    </div>

    <div slot="actions">
      <z-btn
        :formId="formId"
        btnType="validate"
        class="pb-4"
        :color="color"
        block
        @click="onConfirm"
      >修改密码</z-btn>
    </div>
  </lp-dialog>
</template>

<script>
  import ButtonMixins from '@/mixins/ButtonMixins'
  import DialogMixins from '@/mixins/DialogMixins'
  import LoginMixins from '../LoginMixins'

  export default {
    name: 'ChangePassword',
    mixins: [ButtonMixins, DialogMixins, LoginMixins],

    data() {
      return {
        formId: 'ChangePassword',
        dispatchForm: 'CHANGE_PASSWORD',
        originEye: false,
        newEye: false,
        confirmEye: false,
        checkPass: undefined,
        //最短9位，最长16位 {9,16}
        //必须包含1个数字
        //必须包含1个小/大写字母
        //必须包含1个特殊字符
        regex_newpassword_digital: /^(?=.*[0-9])(?=.*[a-zA-z])(?=.*[!@#$%^&*()`~.,;'！￥。，；‘——_|+-/《》、])(.{9,16})$/,
        //regex_newpassword_letter: /^$/,
        //regex_newpassword_length: /^.{9,16}$/,
        //regex_newpassword: /^[0-9a-zA-Z!@#$%^&*?_~`.。·;:《》<>\[\]\{\}\/\^\(\)]{9,16}$/
      }
    },

    watch: {
      'forms.ChangePassword': {
        handler(form) {
          this.checkPass = form.checkPass
        }
      }
    }
  }
</script>

