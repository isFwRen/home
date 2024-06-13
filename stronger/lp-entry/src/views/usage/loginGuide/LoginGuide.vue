<template>
  <div class="login-guide">
    <ul>
      <li 
        v-for="(item, index) in imageList"
        :key="`image_${ index }`"
        class="mb-4 text-center"
      >
        <h1>{{ item.header }}</h1><br/>
        <p class="mb-2">{{ item.text }}</p>
        <img :src="item.path" >
      </li>
    </ul>
  </div>
</template>

<script>
  export default {
    name: 'LoginGuide',

    data() {
      return {
        imageList: [
          {
            header: '如何登录',
            path: require('@/assets/images/usage/loginGuide/1.jpg'),
            text: '输入手机号、密码和验证码，验证码在钉钉工作通知上。第一次登录网址需要修改密码，原密码是123456。'
          },
          {
            path: require('@/assets/images/usage/loginGuide/2.jpg')
          },
          {
            path: require('@/assets/images/usage/loginGuide/3.jpg'),
            text: '成功修改密码后，输入手机号、密码和验证码进行登录，进入系统主页。'
          },
          {
            path: require('@/assets/images/usage/loginGuide/4.jpg')
          }
        ]
      }
    }
  }
</script>

<style scoped lang="scss">
  .login-guide {
    h1{
      font-size: 22px;
      font-weight:bold;
    }
  }
</style>