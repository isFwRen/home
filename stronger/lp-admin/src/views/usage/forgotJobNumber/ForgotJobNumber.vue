<template>
	<div class="forgot-job-number">
		<ul>
			<li v-for="(item, index) in imageList" :key="`image_${index}`" class="mb-4 text-center">
				<h1>{{ item.header }}</h1>
				<br />
				<p class="mb-5">{{ item.text }}</p>
				<img :src="item.path" />
			</li>
		</ul>
	</div>
</template>

<script>
export default {
	name: "ForgetJobNumber",

	data() {
		return {
			imageList: [
				{
					header: "忘记工号",
					path: require("@/assets/images/usage/forgetJobNumber/1.jpg"),
					text: "方法1：进入网址 https://www.i-confluence.com:11111/  ，点击右上角的“注册查询”查询。"
				},
				{
					path: require("@/assets/images/usage/forgetJobNumber/2.jpg")
				},
				{
					path: require("@/assets/images/usage/forgetJobNumber/3.jpg"),
					text: "方法2：点击钉钉头像→点击呢称→可看到姓名前方存在P开头的工号。"
				},
				{
					path: require("@/assets/images/usage/forgetJobNumber/4.jpg")
				}
			]
		};
	}
};
</script>

<style scoped lang="scss">
.forgot-job-number {
	h1 {
		font-size: 22px;
		font-weight: bold;
	}
}
</style>
