<template>
	<div id="lpLogin" class="lp-login">
		<div class="container">
			<div class="z-flex justify-between align-center">
				<div class="logo">
					<img src="@/assets/logo.png" />
				</div>
				<div class="btn_wrap">
					<a class="btn" href="https://www.i-confluence.com" target="_blank">关于汇流</a>
				</div>
			</div>
			<div class="z-flex justify-around content_wrap pt-16">
				<div class="l_content">
					<h2 class="title2">珠海汇流理赔数据处理平台2.0</h2>
					<p class="desc">基于智能化数据采集系统和成熟的数据服务体系，珠海汇流提供专业的数据采集与应用服务....</p>
					<div class="l_btn_wrap mt-10">
						<a class="l_btn" :href="downloadUrl64" target="_blank">Windows 64位下载</a>
						<a class="l_btn ml-2" :href="downloadUrl32" target="_blank">Windows 32位下载</a>
					</div>
				</div>
				<div class="r_img">
					<div class="img_wrap elevation-6">
						<img src="@/assets/images/download.png" alt="" />
					</div>
					<div class="img_desc z-flex justify-between">
						<div class="i_left">联系方式：0756-3630987</div>
						<div class="i_right">在线客服：关注公众号“汇流珠海HR”</div>
					</div>
				</div>
			</div>
		</div>

		<bubble-background></bubble-background>
	</div>
</template>

<script>
export default {
	name: "download",
	data() {
		return {
			downloadUrl: ""
		};
	},
	mounted() {
		this.downloadUrl64 = location.origin + "/app/理赔录入-beta-win32-x64.rar";
		this.downloadUrl32 = location.origin + "/app/理赔录入-beta-win32-ia32.rar";
	},
	methods: {},
	components: {
		"bubble-background": () => import("../login/bubbleBackground")
	}
};
</script>

<style lang="scss" scoped>
.lp-login {
	position: relative;
	width: 100%;
	height: 100%;
	background-color: #4e54c8;
	overflow: hidden;

	.container {
		position: relative;
		z-index: 10;

		.logo {
			width: 80px;
			height: 80px;
			filter: invert(100%) grayscale(100);

			img {
				width: 100%;
				height: 100%;
			}
		}

		.btn_wrap {
			.btn {
				color: #fff;
				padding: 5px 10px;
				border-radius: 10px;
				cursor: pointer;
				display: inline-block;
				transition: all 0.2s ease;
				border: 1px solid #fff;
				text-decoration: none;

				&:focus {
					color: #4e54c8;
					background: #fff;
				}

				&:hover {
					color: #4e54c8;
					background: #fff;
				}
			}
		}

		@media (min-width: 1920px) {
			.content_wrap {
				.l_content {
					padding: 20px 20px 20px 40px;

					//width: 50%;
					.title2 {
						color: #fff;
						font-size: 33px;
						min-width: 600px;
					}

					.desc {
						width: 600px;
						padding-top: 30px;
						color: #fff;
						font-size: 27px;
					}

					.l_btn_wrap {
						padding-top: 100px;

						.l_btn {
							text-decoration: none;
							color: rgba(0, 0, 0, 0.85);
							background: rgba(255, 255, 255, 0.85);
							padding: 12px 25px;
							border-radius: 10px;
							cursor: pointer;
							display: inline-block;
							transition: all 0.2s ease;

							&:focus {
								color: #fff;
								background: #5ab55e;
							}

							&:hover {
								color: #fff;
								background: #5ab55e;
							}
						}
					}
				}

				.r_img {
					width: 40%;

					.img_wrap {
						width: 550px;
						height: 450px;
						border-radius: 10px;
						overflow: hidden;

						img {
							width: 100%;
							height: 100%;
						}
					}

					.img_desc {
						width: 550px;
						padding: 20px 5px;
						color: #fff;
					}
				}
			}
		}

		@media (min-width: 1200px) and (max-width: 1919px) {
			.content_wrap {
				.l_content {
					padding: 5px 20px 20px 40px;

					//width: 50%;
					.title2 {
						color: #fff;
						font-size: 24px;
						min-width: 550px;
					}

					.desc {
						width: 550px;
						padding-top: 30px;
						color: #fff;
						font-size: 22px;
					}

					.l_btn_wrap {
						padding-top: 100px;

						.l_btn {
							color: rgba(0, 0, 0, 0.85);
							background: rgba(255, 255, 255, 0.85);
							padding: 10px 25px;
							border-radius: 10px;
							cursor: pointer;
							display: inline-block;
							transition: all 0.2s ease;

							&:focus {
								color: #fff;
								background: #5ab55e;
							}
						}
					}
				}

				.r_img {
					width: 40%;

					.img_wrap {
						width: 440px;
						height: 360px;
						border-radius: 10px;
						overflow: hidden;

						img {
							width: 100%;
							height: 100%;
						}
					}

					.img_desc {
						color: #fff;
						width: 440px;
						padding: 20px 3px;
						font-size: 14px;

						.i_left,
						.l_right {
							white-space: nowrap;
						}
					}
				}
			}
		}

		@media all and (max-width: 1200px) {
			.content_wrap {
				.l_content {
					padding: 20px 20px 20px 40px;

					//width: 50%;
					.title2 {
						color: #fff;
						font-size: 22px;
						min-width: 400px;
					}

					.desc {
						width: 400px;
						padding-top: 30px;
						color: #fff;
						font-size: 20px;
					}

					.l_btn_wrap {
						.l_btn {
							color: rgba(0, 0, 0, 0.85);
							background: rgba(255, 255, 255, 0.85);
							padding: 10px 25px;
							border-radius: 10px;
							cursor: pointer;
							display: inline-block;
							transition: all 0.2s ease;

							&:focus {
								color: #fff;
								background: #5ab55e;
							}
						}
					}
				}

				.r_img {
					width: 40%;

					.img_wrap {
						width: 380px;
						height: 280px;
						border-radius: 10px;
						overflow: hidden;

						img {
							width: 100%;
							height: 100%;
						}
					}

					.img_desc {
						width: 350px;
						padding: 20px 5px;
						color: #fff;
					}
				}
			}
		}
	}
}
</style>
