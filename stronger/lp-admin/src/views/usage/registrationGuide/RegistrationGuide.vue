<template>
	<div class="registration-guide">
		<ul>
			<li v-for="(item, index) in imageList" :key="`image_${index}`" class="mb-4 text-center">
				<h1>{{ item.header }}</h1>
				<br />
				<p class="mb-1">{{ item.text }}</p>
				<img :src="item.path" />
			</li>
		</ul>
	</div>
</template>

<script>
export default {
	name: "RegistrationGuide",

	data() {
		return {
			imageList: [
				{
					header: "如何注册",
					path: require("@/assets/images/usage/registrationGuide/1.jpg"),
					text: "- 进入网址 https://www.i-confluence.com:11111/  ，点击右上角的“用户注册”，同意协议后输入个人信息并提交，即注册完毕。"
				},
				{
					path: require("@/assets/images/usage/registrationGuide/2.jpg")
				},
				{
					path: require("@/assets/images/usage/registrationGuide/3.jpg"),
					text: "- 注册成功后主动添加管理员微信用于收取审核通知，未添加或未正确填写微信号无法收取通知；可在注册网址点击“注册查询”查询审核状态。"
				},
				{
					path: require("@/assets/images/usage/registrationGuide/4.jpg")
				}
			]
		};
	}
};
</script>

<style scoped lang="scss">
.registration-guide {
	h1 {
		font-size: 22px;
		font-weight: bold;
	}
}
</style>
