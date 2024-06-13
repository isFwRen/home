<template>
	<div class="code">
		<ul>
			<li v-for="(item, index) in imageList" :key="`image_${index}`" class="mb-4 text-center">
				<h1>{{ item.header }}</h1>
				<br />
				<p class="mb-4">{{ item.text }}</p>
				<img :src="item.path" />
			</li>
		</ul>
	</div>
</template>

<script>
export default {
	name: "Code",

	data() {
		return {
			imageList: [
				{
					header: "无法获取验证码",
					path: require("@/assets/images/usage/code/1.jpg"),
					text: "-检查一下登陆的手机号是否和注册的手机号以及钉钉登录的手机号相同。"
				},
				{
					path: require("@/assets/images/usage/code/2.jpg"),
					text: "-检查所使用的浏览器是否为谷歌浏览器或火狐浏览器。"
				},
				{
					path: require("@/assets/images/usage/code/3.jpg"),
					text: "-若上述方法均无法获取验证码,请主动联系管理员。"
				}
			]
		};
	}
};
</script>

<style scoped lang="scss">
.code {
	h1 {
		font-size: 22px;
		font-weight: bold;
	}
}
</style>
