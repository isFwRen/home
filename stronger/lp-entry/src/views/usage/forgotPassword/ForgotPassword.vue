<template>
  <div class="forgot-password">
    <ul>
      <li 
        v-for="(item, index) in imageList"
        :key="`image_${ index }`"
        class="mb-4 text-center"
      >
        <h1>{{ item.header }}</h1><br/>
        <p class="mb-3">{{ item.text }}</p>
        <img :src="item.path" >
      </li>
    </ul>
  </div>
</template>

<script>
  export default {
    name: 'ForgotPassword',

    data() {
      return {
        imageList: [
          {
            header: '忘记密码',
            path: require('@/assets/images/usage/forgotPassword/1.jpg'),
            text: '点击忘记密码，输入手机号码和验证码，验证码在钉钉工作通知上。'
          },
          {
            path: require('@/assets/images/usage/forgotPassword/2.jpg')
          },
          {
            path: require('@/assets/images/usage/forgotPassword/3.jpg'),
            text: '成功重置密码后，点击修改密码，进入修改密码页面，重置密码即为重置成原始密码123456，需成功修改密码后，方可进行登录。 '
          },
          {
            path: require('@/assets/images/usage/forgotPassword/4.jpg')
          }
        ]
      }
    }
  }
</script>

<style scoped lang="scss">
  .forgot-password {
    h1{
      font-size: 22px;
      font-weight:bold;
    }    
  }
</style>