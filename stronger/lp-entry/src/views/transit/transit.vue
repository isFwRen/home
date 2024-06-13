<template>
	<div class="transit">
		<div class="transit-top">
			<div class="spinner">
				<div></div>
				<div></div>
				<div></div>
				<div></div>
				<div></div>
				<div></div>
			</div>
			<div class="button">
				<button class="contactButton" @click="backToLogin">
					back to Login
					<div class="iconButton">
						<svg height="24" width="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
							<path d="M0 0h24v24H0z" fill="none"></path>
							<path
								d="M16.172 11l-5.364-5.364 1.414-1.414L20 12l-7.778 7.778-1.414-1.414L16.172 13H4v-2z"
								fill="currentColor"
							></path>
						</svg>
					</div>
				</button>
			</div>
		</div>

		<div class="transit-wrap">
			<h2 class="transit-title">🔥 请选择您的项目:</h2>
			<div class="containers">
				<div class="item" v-for="item in list" :key="item.proCode">
					<div class="item-logo" @click="linkTo(item)">
						<div class="item-title">{{ item.proCode }}</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>
<script>
import { tools as lpTools } from "@/libs/util";
import { localStorage, sessionStorage } from "vue-rocket";
const isIntranet = lpTools.isIntranet();
export default {
	data() {
		return {
			list: []
		};
	},
	created() {
		this.getList("INPUT_GET_TRANSIT_LIST");
	},
	methods: {
		linkTo(item) {
			const { innerIp, inAppPort, outIp, outAppPort } = item;

			const origin = `https://${isIntranet ? `${innerIp}:${inAppPort}` : `${outIp}:${outAppPort}`}`;

			const token = localStorage.get("token");
			const user = localStorage.get("user");
			const secret = localStorage.get("secret") || "";
			const height = localStorage.get("viewport").height;
			const width = localStorage.get("viewport").width;
			const isApp = sessionStorage.get("isApp").isApp;

			const url = `${origin}/main/entry/channel?sync=1&isApp=${isApp}&height=${height}&width=${width}&userId=${user.id}&token=${token}&secret=${secret}`;
			window.open(url);
			//window.open(`${origin}/main/entry/channel?token=${token}`);
		},
		backToLogin() {
			this.$modal({
				visible: true,
				title: "返回提示",
				content: `请确认是否要退返回登录页？`,
				confirm: async () => {
					let isApp = sessionStorage.get("isApp");
					this.$router.replace(`/login?isApp=${isApp.isApp}`);
				}
			});
		},
		async getList(dispatch) {
			const result = await this.$store.dispatch(dispatch);
			if (result.code === 200) {
				const sortArr = result.data.list.sort(
					(a, b) => Number(a.proCode.slice(2, a.length)) - Number(b.proCode.slice(2, b.length))
				);
				this.list = sortArr;
			}
			this.toasted.dynamic(result.data.msg, result.code);
		}
	}
};
</script>
<style scoped lang="scss">
.transit {
	background-color: #4b90e1;
	width: 100%;
	height: 100%;
	position: relative;

	.transit-top {
		padding: 20px 40px 40px 40px;
		display: flex;
		justify-content: space-between;
	}

	.transit-wrap {
		width: 950px;
		margin: 0 auto;
		background: rgba(255, 255, 255, 0.25);
		box-shadow: 0 8px 32px 0 rgba(31, 38, 135, 0.37);
		backdrop-filter: blur(4px);
		-webkit-backdrop-filter: blur(4px);
		border-radius: 10px;
		border: 1px solid rgba(255, 255, 255, 0.18);
		.transit-title {
			color: #fff;
			font-size: 20px;
			padding: 10px;
		}
	}
	.containers {
		padding: 40px 20px;
		box-sizing: border-box;
		display: flex;
		flex-wrap: wrap;
		backdrop-filter: blur(4px);
		.item {
			margin-left: 15px;
			margin-right: 15px;
			margin-bottom: 30px;
			cursor: pointer;
			box-sizing: border-box;
			&:hover {
				.item-logo {
					user-select: none;
					transform: translateY(-10px);
					background-color: rgba(255, 255, 255, 0.35);
					color: hsl(0, 0%, 100%);
					filter: drop-shadow(0 0 1em #646cffaa);
				}
			}
			.item-logo {
				color: #e6f4ff;
				border-radius: 20px;
				width: 120px;
				height: 120px;
				display: flex;
				justify-content: center;
				align-items: center;
				background: rgba(255, 255, 255, 0.25);
				border: 1px solid rgba(255, 255, 255, 0.18);
				backdrop-filter: blur(4px);
				transition: transform 0.2s, background-color 0.2s, color 0.2s, filter 0.2s ease;
				.item-title {
					font-weight: 600;
					font-size: 20px;
				}
			}
		}
	}
}

.spinner {
	width: 50px;
	height: 50px;
	--clr: rgb(127, 207, 255);
	--clr-alpha: rgb(127, 207, 255, 0.1);
	animation: spinner 4s infinite linear;
	transform-style: preserve-3d;
}

.spinner > div {
	background-color: var(--clr-alpha);
	height: 100%;
	position: absolute;
	width: 100%;
	border: 3px solid var(--clr);
}

.spinner div:nth-of-type(1) {
	transform: translateZ(-22px) rotateY(180deg);
}

.spinner div:nth-of-type(2) {
	transform: rotateY(-270deg) translateX(50%);
	transform-origin: top right;
}

.spinner div:nth-of-type(3) {
	transform: rotateY(270deg) translateX(-50%);
	transform-origin: center left;
}

.spinner div:nth-of-type(4) {
	transform: rotateX(90deg) translateY(-50%);
	transform-origin: top center;
}

.spinner div:nth-of-type(5) {
	transform: rotateX(-90deg) translateY(50%);
	transform-origin: bottom center;
}

.spinner div:nth-of-type(6) {
	transform: translateZ(22px);
}

@keyframes spinner {
	0% {
		transform: rotate(0deg) rotateX(0deg) rotateY(0deg);
	}

	50% {
		transform: rotate(180deg) rotateX(180deg) rotateY(180deg);
	}

	100% {
		transform: rotate(360deg) rotateX(360deg) rotateY(360deg);
	}
}

.contactButton {
	background: #3cb371;
	color: white;
	font-family: inherit;
	padding: 0.5em;
	padding-left: 1.2em;
	font-size: 18px;
	font-weight: 600;
	border-radius: 1em;
	border: none;
	letter-spacing: 0.05em;
	display: flex;
	align-items: center;
	box-shadow: inset 0 0 1.8em -0.7em #2e8b57;
	overflow: hidden;
	position: relative;
	height: 3em;
	padding-right: 3.5em;
}

.iconButton {
	margin-left: 1.2em;
	position: absolute;
	display: flex;
	align-items: center;
	justify-content: center;
	height: 2.4em;
	width: 2.4em;
	border-radius: 1.2em;
	box-shadow: 0.15em 0.15em 0.8em 0.3em #32cd32;
	right: 0.4em;
	transition: all 0.3s;
}

.contactButton:hover {
	transform: translate(-0.05em, -0.05em);
	box-shadow: 0.2em 0.2em #228b22;
}

.contactButton:active {
	transform: translate(0.05em, 0.05em);
	box-shadow: 0.1em 0.1em #228b22;
}
</style>
