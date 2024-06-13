<template>
	<div :class="['time-brief', expand ? 'toLeft' : 'toRight']" :style="moveStyle" ref="dragEle">
		<div class="time-brief-container" v-if="expand">
			<div
				class="time-brief-title font-weight-black z-flex justify-space-between align-center"
			>
				<v-btn class="point_auto" icon color="primary" @click="toggleBox(false)">
					<svg
						t="1693213018385"
						class="icon"
						viewBox="0 0 1024 1024"
						version="1.1"
						xmlns="http://www.w3.org/2000/svg"
						p-id="14209"
						width="22"
						height="22"
					>
						<path
							d="M170.667 149.333h682.666a42.667 42.667 0 0 1 0 85.334H170.667a42.667 42.667 0 1 1 0-85.334z m0 640h682.666a42.667 42.667 0 0 1 0 85.334H170.667a42.667 42.667 0 0 1 0-85.334z m256-213.333h426.666a42.667 42.667 0 0 1 0 85.333H426.667a42.667 42.667 0 0 1 0-85.333z m0-213.333h426.666a42.667 42.667 0 0 1 0 85.333H426.667a42.667 42.667 0 0 1 0-85.333zM289.707 535.04L180.139 651.776a29.227 29.227 0 0 1-43.179 0 33.707 33.707 0 0 1-8.96-23.04V395.264c0-17.963 13.653-32.555 30.55-32.555 8.106 0 15.871 3.414 21.589 9.558L289.707 488.96a34.133 34.133 0 0 1 0 46.08z"
							p-id="14210"
							fill="#1976d2"
						></path>
					</svg>
				</v-btn>
				<h6 class="point_auto" ref="dragTarget">{{ proCode }}时效简报</h6>
				<v-btn class="point_auto" icon color="green" @click="getTimeBrief">
					<v-icon>mdi-cached</v-icon>
				</v-btn>
			</div>
			<div class="residue red--text">
				超时30分钟未回传:
				<span class="point_auto user_select" @click="showDetail('not30ReturnMap')">{{
					timeBriefs.not30ReturnMapLen
				}}</span>
			</div>
			<div class="residue amber--text">
				时效剩余5分钟以内:
				<span class="point_auto user_select" @click="showDetail('in5MinutesMap')">{{
					timeBriefs.in5MinutesMapLen
				}}</span>
			</div>
			<div class="residue lime--text">
				时效剩余15分钟以内:
				<span class="point_auto user_select" @click="showDetail('in15MinutesMap')">{{
					timeBriefs.in15MinutesMapLen
				}}</span>
			</div>
			<div class="residue green--text">
				时效剩余20分钟以内:
				<span class="point_auto user_select" @click="showDetail('in20MinutesMap')">{{
					timeBriefs.in20MinutesMapLen
				}}</span>
			</div>
		</div>
		<div class="expand" v-else>
			<v-btn class="point_auto" icon color="primary" @click="toggleBox(true)">
				<svg
					t="1693212884338"
					class="icon"
					viewBox="0 0 1024 1024"
					version="1.1"
					xmlns="http://www.w3.org/2000/svg"
					p-id="14035"
					width="20"
					height="20"
				>
					<path
						d="M393.142857 432h548.571429c5.028571 0 9.142857-4.114286 9.142857-9.142857v-64c0-5.028571-4.114286-9.142857-9.142857-9.142857H393.142857c-5.028571 0-9.142857 4.114286-9.142857 9.142857v64c0 5.028571 4.114286 9.142857 9.142857 9.142857z m-9.142857 233.142857c0 5.028571 4.114286 9.142857 9.142857 9.142857h548.571429c5.028571 0 9.142857-4.114286 9.142857-9.142857v-64c0-5.028571-4.114286-9.142857-9.142857-9.142857H393.142857c-5.028571 0-9.142857 4.114286-9.142857 9.142857v64z m576-555.428571H64c-5.028571 0-9.142857 4.114286-9.142857 9.142857v64c0 5.028571 4.114286 9.142857 9.142857 9.142857h896c5.028571 0 9.142857-4.114286 9.142857-9.142857v-64c0-5.028571-4.114286-9.142857-9.142857-9.142857z m0 722.285714H64c-5.028571 0-9.142857 4.114286-9.142857 9.142857v64c0 5.028571 4.114286 9.142857 9.142857 9.142857h896c5.028571 0 9.142857-4.114286 9.142857-9.142857v-64c0-5.028571-4.114286-9.142857-9.142857-9.142857zM58.742857 519.885714L237.371429 660.571429c6.628571 5.257143 16.457143 0.571429 16.457142-7.885715V371.314286c0-8.457143-9.714286-13.142857-16.457142-7.885715L58.742857 504.114286a9.988571 9.988571 0 0 0 0 15.771428z"
						p-id="14036"
						fill="#1976d2"
					></path>
				</svg>
			</v-btn>
		</div>

		<div class="time-brief-dialog point_auto" v-if="timeBriefList.length && brief_dialog">
			<v-card>
				<v-system-bar color="#fff">
					<v-spacer></v-spacer>
					<v-icon @click="brief_dialog = false">mdi-close</v-icon>
				</v-system-bar>
				<v-card-text>
					<div class="list">
						<span v-for="bill in timeBriefList" :key="bill">
							{{ bill }}
						</span>
					</div>
				</v-card-text>
			</v-card>
		</div>
	</div>
</template>
<script>
import { mapGetters } from "vuex";
import dragMixin from "@/mixins/dragMixin";
export default {
	name: "timeBrief",
	mixins: [dragMixin],
	data() {
		return {
			expand: false,
			brief_dialog: false,
			timeBriefList: []
		};
	},
	props: ["proCode"],
	created() {
		this.getTimeBrief();
	},
	computed: {
		...mapGetters(["timeBriefs"])
	},
	methods: {
		async getTimeBrief() {
			this.timeBriefList = [];
			if (!this.proCode) {
				return;
			}
			const params = {
				proCode: this.proCode
			};
			const result = await this.$store.dispatch("GET_TIME_BRIEF", params);
			this.toasted.dynamic(result.msg, result.code);
		},
		toggleBox(bool) {
			this.dragBox.style.transform = "";
			this.expand = bool;
			if (bool) {
				this.$nextTick(() => {
					this.position.x = 0;
					this.position.y = 0;
					this.unbindEventListener(this.$refs.dragTarget, {});
					this.bindEventListenr(this.$refs.dragTarget, {});
					this.dragBox = this.$refs.dragEle;
				});
			} else {
				this.brief_dialog = false;
			}
		},
		showDetail(key) {
			if (this.timeBriefs[key] && this.timeBriefs[key].length > 0) {
				this.brief_dialog = true;
				this.timeBriefList = this.timeBriefs[key];
			} else {
				this.timeBriefList = [];
			}
		}
	},
	watch: {
		proCode(newVal, oldValue) {
			if (newVal !== oldValue) {
				this.getTimeBrief();
			}
		}
	}
};
</script>
<style scoped lang="scss">
.time-brief {
	width: 235px;
	height: 220px;
	border-radius: 5px;
	position: fixed;
	z-index: 100;
	pointer-events: none;
	// top: 270px;
	bottom: 20px;
	left: 20px;
	transition: right 0.25s ease-in;
	&-title {
		margin-bottom: 5px;
		cursor: move;
		user-select: none;
	}
	&-container {
		background: rgba(255, 255, 255, 0.25);
		box-shadow: 0 5px 20px 0 rgba(31, 38, 135, 0.37);
		backdrop-filter: blur(1px);
		-webkit-backdrop-filter: blur(1px);
		border-radius: 10px;
		border: 1px solid rgba(255, 255, 255, 0.18);
	}
	.expand {
		border: 1px solid #eee;
		position: absolute;
		border-radius: 50%;
		top: 5px;
		left: 0;
	}
}
.residue {
	padding-left: 35px;
	box-sizing: border-box;
	text-align: left;
	margin-bottom: 20px;
	cursor: pointer;
}
.point_auto {
	pointer-events: auto;
}
.toRight {
	right: -185px;
}
.toLeft {
	right: 20px;
}
.user_select {
	user-select: none;
}
.time-brief-dialog {
	position: absolute;
	top: 0;
	left: -210px;
	width: 200px;
	height: 220px;
	border-radius: 10px;
	background-color: #fff;
	.list {
		height: 200px;
		overflow: auto;
	}
}
</style>
