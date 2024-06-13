<template>
    <div class="minTitle">
        <h2>{{text}}</h2>
        <slot></slot>
    </div>
</template>

<script>
  export default {
    props:["text"]
}
</script>

<style lang="scss" scoped>
.minTitle {
    font-family: SourceHanSansSC;
    font-weight: 700;
    font-size: 24px;
    color: rgba(6, 48, 121, 0.8);
    font-style: normal;
    letter-spacing: 0px;
    line-height: 36px;
    text-decoration: none;
    border-bottom: 2px solid #999;
    position: relative;
    h2{
        display: inline-block;
    }
}
</style>