<template>
  <div id="app">
     <!-- <vue-particles
        color="#fff"
        :particleOpacity="0.7"
        :particlesNumber="60"
        shapeType="circle"
        :particleSize="4"
        linesColor="#fff"
        :linesWidth="1"
        :lineLinked="true"
        :lineOpacity="0.4"
        :linesDistance="150"
        :moveSpeed="2"
        :hoverEffect="true"
        hoverMode="grab"
        :clickEffect="true"
        clickMode="push"
        class="lizi"
      >
      </vue-particles> -->
    <!-- <img src="./assets/logo.png"> -->
    <router-view/>
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
  /* 大屏幕 ：大于等于1200px*/
  @media (min-width: 1200px) {
    body{
        min-width: auto;
        /*overflow-x:hidden;*/
      }
    }

  /*默认*/
  @media (min-width: 980px){
    body{
      min-width: auto;
    }
  }

  /* 平板电脑和小屏电脑之间的分辨率 */
  @media (min-width: 768px) and (max-width: 979px) {
    body{
      min-width: 1500px;
    }
  }

  /* 横向放置的手机和竖向放置的平板之间的分辨率 */
  @media (max-width: 767px) {
    body{
      min-width: 1500px;
    }
  }

  /* 横向放置的手机及分辨率更小的设备 */
  @media (max-width: 480px) {
    body{
      min-width: 1500px;
    }
  }
/* #app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
} */
.el-table thead.is-group th {
  background: #F5F7FA;
}
body::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}
body::-webkit-scrollbar-thumb {
  border-radius: 10px;
  -webkit-box-shadow: inset 0 0 5px rgba(0, 0, 0, 0.2);
  background: #8cc1f3;
}
body::-webkit-scrollbar-track {
  -webkit-box-shadow: inset 0 0 5px rgba(0, 0, 0, 0.2);
  border-radius: 10px;
  background: #ededed;
}
</style>
