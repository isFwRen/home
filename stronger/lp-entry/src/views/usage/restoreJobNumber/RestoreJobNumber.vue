<template>
  <div class="restore-job-number">
    <ul>
      <li 
        v-for="(item, index) in imageList"
        :key="`image_${ index }`"
        class="mb-4 text-center"
      >
        <h1>{{ item.header }}</h1><br/>
        <p class="mb-6">{{ item.text }}</p>
        <img :src="item.path" >
      </li>
    </ul>
  </div>
</template>

<script>
  export default {
    name: 'RestoreJobNumber',

    data() {
      return {
        imageList: [
          {
            header: '恢复工号',
            path: require('@/assets/images/usage/restoreJobNumber/1.jpg'),
            text: '- 在微信查找公众号“珠海汇流HR”，点击“兼职招聘”，选择“恢复工号权限申请”，按要求填写申请资料并提交。'
          },
          {
            path: require('@/assets/images/usage/restoreJobNumber/2.jpg')
          },
          {
            path: require('@/assets/images/usage/restoreJobNumber/3.jpg')
          },
          {
            path: require('@/assets/images/usage/restoreJobNumber/4.jpg'),
            text: '- 提交后添加管理员微信用于收取审核通知，未添加或未正确填写微信号无法收取通知。'
          }
        ]
      }
    }
  }
</script>

<style scoped lang="scss">
  .restore-job-number {
    h1{
      font-size: 22px;
      font-weight:bold;
    }
  }
</style>