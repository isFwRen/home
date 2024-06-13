import { tools, sessionStorage } from 'vue-rocket'
import moment from 'moment'
import { ignoreFreeValue } from '../tools'
import { MessageBox, Notification } from 'element-ui';


const B0102 = {
  op0: {
    // 记录最后一次存储的合法field
    memoFields: [],

    // 记录相同 code 的 field 的值
    memoFieldValues: [],

    // fields 的值从 targets 里的值选择
    dropdownFields: [
      {
        targets: [],
        fields: []
      }
    ],


    // 校验规则
    rules: [],

    // 提示文本
    hints: [],

    // 工序完成初始化
    init: {
      methods: {

      }
    },

    // 字段已生成
    updateFields: {
      methods: {

      }
    },

    // 回车
    enter: {
      methods: {
      }
    },

    // F8(提交前校验)
    beforeSubmit: {
      methods: {},
    }
  },

  op1op2opq: {
    // 校验规则
    rules: [],

    // 提示文本
    hints: [

    ],

    // 字段已生成
    updateFields: {
      methods: {}
    },

    // 回车
    enter: {
      methods: {}
    },

    // 临时保存
    sessionSave: {
      methods: {

      }
    },

    // 提交前
    beforeSubmit: {
      methods: {

      }
    }
  }
}

export default B0102
