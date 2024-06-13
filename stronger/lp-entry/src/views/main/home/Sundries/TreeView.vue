<template>
  <div class="treeView">
    <div class="trees">
      <div class="Y">
        <li v-for="(num, i) in trees.Y" :key="i">
          {{ num }}
        </li>
      </div>
      <div class="X">
        <li v-for="(day, i) in trees.X" :key="i">
          {{ day }}
        </li>
      </div>
      <div
        :title="data"
        class="tree"
        v-for="(data, i) in trees.data"
        :key="i"
        :style="{
          height:
            (data / trees.Y[trees.Y.length - 1]) * 100 + '%'
        }"
      ></div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      trees: {
        //Y轴数据只能传数字
        Y: [],
        X: [],
        data: []
      }
    }
  },
  props: ['inputData'],
  watch: {
    'inputData': {
       handler (newData) {
        this.trees = newData
      },
      deep:true
    }
  },
}
</script>

<style lang="scss" scoped>
.treeView {
  width: 100%;
  height: 100%;
  box-sizing: border-box;
  padding: 10px 10px 2em 3em;
  position: relative;
  overflow: hidden;

  .trees {
    width: 100%;
    height: 100%;
    border-left: 1px solid #bbb;
    border-bottom: 1px solid #bbb;
    position: relative;
    display: flex;
    flex-direction: row;
    justify-content: space-around;
    align-items: flex-end;

    .tree {
      background-color: rgb(50, 145, 248);
      width: 1em;
    }
  }
  .X,
  .Y {
    font-size: 0.7em;
  }

  .X {
    position: absolute;
    display: flex;
    justify-content: space-around;
    bottom: 0;
    transform: translateY(110%);
    width: 100%;

    li {
      list-style: none;
    }
  }

  .Y {
    position: absolute;
    top: -10px;
    left: 0em;
    height: 100%;
    display: flex;
    justify-content: space-between;
    flex-direction: column-reverse;
    transform: translateX(-110%);
    height: 110%;

    li {
      list-style: none;
      text-align: right;
    }
  }
}
</style>
