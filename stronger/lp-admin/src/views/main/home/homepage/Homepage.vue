<template>
	<div class="lp-home">
		<div class="firstBlock">
			<div class="myAim">
				<HomeMinTitle text="我的目标">
					<div @click="changeMyAim" class="edit">
						<EditSVG />
					</div>
				</HomeMinTitle>
				<div class="content">
					<div class="todayYield">
						<div>
							<USVG />
						</div>
						<span>今日产量</span>
						<span>{{ yeild.nowYeild }}</span>
						<span>目标:{{ yeild.aimYeild }}</span>
						<span>{{
							sayArr[yeild.aimYeild] || sayArr[yeild.nowYeild > yeild.aimYeild]
						}}</span>
					</div>
					<div class="statistics">
						<TreeView :inputData="treeData" />
					</div>
				</div>
			</div>
			<div class="todayRanking">
				<HomeMinTitle text="今日排行榜">
					<div class="more" @click="showRanking()">
						<span>更多>></span>
					</div>
				</HomeMinTitle>
				<table>
					<tr>
						<th>名次</th>
						<th>工号</th>
						<th>姓名</th>
						<th>产量</th>
					</tr>
					<tr>
						<td>
							<div class="icon">
								{{ ranking.userYieldRanking.myOrder }}
							</div>
						</td>
						<td>{{ ranking.userYieldRanking.userCode }}</td>
						<td>{{ ranking.userYieldRanking.userName }}</td>
						<td>{{ ranking.userYieldRanking.value }}</td>
					</tr>
					<div class="ranking">
						<tr v-for="(e, i) in ranking.list">
							<td v-if="i < 3">
								<img width="20" height="30" :src="Medal[i]" alt="" />
							</td>
							<td v-else>{{ i + 1 }}</td>
							<td>{{ e.userCode }}</td>
							<td>{{ e.userName }}</td>
							<td>{{ e.value }}</td>
						</tr>
					</div>
				</table>
			</div>
		</div>
		<div class="secondBlock">
			<div class="companyNotice">
				<HomeMinTitle text="公司通告">
					<div class="more" @click="showNotice">
						<span>更多>></span>
					</div>
				</HomeMinTitle>
				<div class="content">
					<li v-for="(e, i) in ruleList" :key="i" @click="watchDetial(e)">
						<h4>{{ e.title }}</h4>
						<time>{{ e.releaseDate.substr(0, 10) }}</time>
					</li>
				</div>
			</div>
			<div class="ruleInformation">
				<HomeMinTitle text="规则信息">
					<div class="more" @click="showRule">
						<span>更多>></span>
					</div>
				</HomeMinTitle>
				<div class="content">
					<li v-for="(e, i) in NoteList" :key="i" @click="watchDetial(e)">
						<h4>{{ e.title }}</h4>
						<time>{{ e.releaseDate.substr(0, 10) }}</time>
					</li>
				</div>
			</div>
		</div>
		<div class="thirdBlock">
			<div class="imgBox">
				<a href="../PM/case">
					<img src="@/assets/images/home/1.png" alt="" /><span>案件列表</span></a
				>
			</div>
			<div class="imgBox">
				<a href="../PM/prescription">
					<img src="@/assets/images/home/2.png" alt="" /><span>时效管理</span>
				</a>
			</div>
			<div class="imgBox">
				<a href="../PM/qualities/manage">
					<img src="@/assets/images/home/3.png" alt="" /><span>质量管理</span></a
				>
			</div>
			<div class="imgBox">
				<a href="../PM/task">
					<img src="@/assets/images/home/mission.jpg" alt="" /><span>任务管理</span></a
				>
			</div>
		</div>
		<ChangeTarget ref="ChangeTarget" @submited="init" />
		<ShowRule ref="ShowRules" @showDetial="watchDetial" />
		<ShowNotice ref="ShowNotice" @showDetial="watchDetial" />
		<WatchDialogVue :content="DetialContent" ref="showWatchDialog" />
		<Ranking ref="showRanking" />
	</div>
</template>

<script>
import HomeMinTitle from "./HomeMinTitle.vue";
import EditSVG from "./Sundries/edit.svg.vue";
import USVG from "./Sundries/uLike.svg.vue";
import TreeView from "./Sundries/TreeView.vue";
import DialogMixins from "@/mixins/DialogMixins";
import ChangeTarget from "./Sundries/ChangeTarget.vue";
import moment from "moment";
import ShowRule from "./showRule/showRule.vue";
import ShowNotice from "./showNotice/showNotice.vue";
import WatchDialogVue from "../../pm/notice/updateDialog/WatchDialog.vue";
import Gold from "@/assets/images/home/goldMedal.png";
import Silver from "@/assets/images/home/SilverMedal.png";
import Bronze from "@/assets/images/home/BronzeMedal.png";
import Ranking from "./ranking/Ranking.vue";

export default {
	name: "LPHome",
	components: {
		HomeMinTitle,
		EditSVG,
		USVG,
		TreeView,
		ChangeTarget,
		ShowRule,
		ShowNotice,
		WatchDialogVue,
		Ranking
	},
	data() {
		return {
			mixins: [DialogMixins],
			yeild: {
				nowYeild: 0,
				aimYeild: 0
			},
			treeData: {
				X: [],
				Y: [],
				data: []
			},
			NoteList: [],
			ruleList: [],
			DetialContent: "",
			sayArr: {
				0: "目标暂未设定，快点设定吧",
				false: "目标还未达成,继续加油哦",
				true: "恭喜达成今日目标成绩"
			},
			ranking: {
				list: [],
				userYieldRanking: {
					value: 0,
					userName: "",
					userCode: "",
					myOrder: 10
				}
			},
			Medal: [Gold, Silver, Bronze]
		};
	},
	created() {
		this.init();
	},
	methods: {
		changeMyAim() {
			this.$refs.ChangeTarget.onOpen(1);
		},
		async init() {
			//获取目标和近期产量
			this.getRankYeild();
			this.getAnnouncement(2, result => {
				this.NoteList = result.data.list;
			});
			this.getAnnouncement(1, result => {
				this.ruleList = result.data.list;
			});
			const result = await this.$store.dispatch("HOME_GET_USER_YIELD");
			this.yeild.nowYeild = result.data.yield;
			this.yeild.aimYeild = result.data.target;
			this.setTreeData(result.data.list);
		},
		async getRankYeild() {
			const data = {
				pageIndex: 1,
				pageSize: 3,
				rankingType: 0
			};
			const result = await this.$store.dispatch("HOME_GET_RANK_YIELD", data);
			this.ranking = result.data;
		},
		async getAnnouncement(releaseType, callback) {
			const data = {
				pageIndex: 1,
				pageSize: 3,
				releaseType: releaseType
			};
			const result = await this.$store.dispatch("HOME_GET_ANNOUNCEMENT", data);
			callback(result);
		},
		setTreeData(list) {
			let treeData = {
					X: [],
					Y: [],
					data: []
				},
				Number = [];
			const date = ["周日", "周一", "周二", "周三", "周四", "周五", "周六"];
			list.forEach(e => {
				treeData.data.unshift(e.value);
				treeData.X.unshift(date[moment(e.yieldDate).format("d")]);
				Number.push(e.value);
			});
			Number.sort((a, b) => b - a);
			treeData.Y = getNumberList(Number[0], 5);
			function getNumberList(num, paragraph) {
				let bigerNum = num,
					divisor = 10,
					step = 0,
					result = [];
				bigerNum = divisor * Math.ceil(bigerNum / divisor);

				step = bigerNum / paragraph;
				for (let i = 0; i <= paragraph + 1; i++) {
					result.push(i * step);
				}
				return result;
			}
			this.treeData = treeData;
		},
		routeChange(path) {
			this.$router.push({ path });
		},
		showRule() {
			this.$refs.ShowRules.onOpen();
		},
		showNotice() {
			this.$refs.ShowNotice.onOpen();
		},
		watchDetial(e) {
			this.$store.dispatch("WHATCH_NOTICE_VIEW_NUMBER", { id: e.ID });
			this.DetialContent = e.content;
			this.$refs.showWatchDialog.onOpen();
		},
		showRanking() {
			this.$refs.showRanking.onOpen();
		}
	}
};
</script>

<style scoped lang="scss">
.lp-home {
	min-width: 900px;

	.secondBlock,
	.firstBlock,
	.thirdBlock {
		width: 100%;
		display: grid;
		grid-template-columns: repeat(10, 1fr);
		box-shadow: rgb(0 0 0 / 24%) 0px 2px 2px 0px;
		padding: 20px;
		gap: 30px 30px;
		margin-bottom: 40px;
	}

	.firstBlock {
		.myAim {
			grid-column: 1/7;

			.edit {
				display: inline-block;
				position: relative;
				left: 10px;
				top: 5px;
			}

			.content {
				width: 100%;
				height: 200px;
				display: grid;
				grid-template-columns: repeat(10, 1fr);
				gap: 30px 30px;

				.todayYield {
					display: grid;
					grid-column: 1/5;
					margin-top: 20px;
					position: relative;

					> div {
						margin: 0 auto;
					}

					span {
						position: absolute;
						left: 50%;
						transform: translateX(-50%);
						font-family: SourceHanSansSC;
						font-weight: 700;
						color: rgb(16, 16, 16);
						font-style: normal;
						letter-spacing: 0px;
						text-decoration: none;
						display: block;
						width: 100%;
						text-align: center;
					}

					span:nth-of-type(1) {
						top: 10px;
						font-weight: 400;
						font-size: 1em;
					}

					span:nth-of-type(2) {
						top: 50px;
						font-size: 2em;
					}

					span:nth-of-type(3) {
						top: 100px;
					}

					span:nth-of-type(4) {
						bottom: 0px;
						font-family: SourceHanSansSC;
						font-weight: 400;
						font-size: 0.8em;
						color: rgb(16, 16, 16);
						font-style: normal;
						letter-spacing: 0px;
						line-height: 20px;
						text-decoration: none;
					}
				}

				.statistics {
					grid-column: 5/11;
					border: 1px solid rgb(187, 187, 187);
					margin-top: 10px;
				}
			}
		}

		.todayRanking {
			grid-column: 7/11;

			table {
				display: grid;
				grid-template-rows: 30px 1fr 10px 3fr;
				margin-top: 10px;
				height: calc(100% - 50px);

				div.ranking {
					display: grid;
					grid-row-start: 4;
					grid-template-rows: repeat(3, 1fr);
					border-color: rgb(193, 205, 209);
					border-width: 0px;
					border-style: solid;
					box-shadow: rgb(0 0 0 / 24%) 0px 8px 8px 0px;
					border-radius: 2px;
					font-size: 1em;
					padding: 0px;
					text-align: left;
					background: rgb(255, 255, 255);
					box-sizing: border-box;

					tr {
						position: relative;
					}

					tr:nth-child(3):after {
						border: none;
					}

					tr:after {
						content: " ";
						width: 92%;
						position: absolute;
						border-bottom: 1px solid #bbb;
						bottom: 0px;
						left: 5%;
					}
				}

				tr {
					display: grid;
					grid-template-columns: repeat(4, 1fr);

					td:last-child {
						color: rgba(255, 37, 37, 1);
					}

					th,
					td {
						vertical-align: middle;
						display: flex;
						align-items: center;
						justify-content: center;
						font-family: SourceHanSansSC;
						font-weight: 400;
						font-size: 1em;
						color: rgb(16, 16, 16);
						font-style: normal;
						letter-spacing: 0px;
						text-decoration: none;
					}
				}

				> tr:nth-child(2) {
					border-color: rgb(193, 205, 209);
					border-width: 0px;
					border-style: solid;
					box-shadow: rgb(0 0 0 / 24%) 0px 8px 8px 0px;
					border-radius: 2px;
					font-size: 1em;
					padding: 0px;
					text-align: left;
					line-height: 29px;
					font-weight: normal;
					font-style: normal;
					background: rgb(255, 255, 255);
					box-sizing: border-box;

					td .icon {
						color: rgba(252, 202, 0, 1);
						width: 30px;
						height: 30px;
						border: 1px solid rgba(252, 202, 0, 1);
						text-align: center;
						border-radius: 50%;
					}
				}
			}
		}
	}

	.secondBlock {
		.companyNotice {
			grid-column: 1/6;
		}

		.ruleInformation {
			grid-column: 6/11;
		}

		.companyNotice,
		.ruleInformation {
			.content {
				font-size: 0.8em;
				display: grid;
				grid-template-rows: repeat(3, 1fr);

				li::before {
					content: "";
					width: 10px;
					height: 10px;
					display: inline-block;
					background-color: black;
					border-radius: 50%;
					top: 50%;
					position: absolute;
					left: 1em;
					transform: translate(-50%, -50%);
				}

				li {
					list-style: none;
					position: relative;
					display: grid;
					grid-template-columns: 2em 1fr 6em;
					grid-template-rows: 2fr;
					margin: 10px 0;
					cursor: pointer;

					h4 {
						grid-column: 2/2;
						grid-row: 1/1;
						font-size: 1.3em;
					}

					span {
						grid-column: 2/2;
						grid-row: 2/2;
						color: #999;
						white-space: nowrap;
						text-overflow: ellipsis;
						overflow: hidden;
						word-break: break-all;
					}

					time {
						top: 0;
						right: 0;
						grid-column: 3/3;
						grid-row: 1/2;
						color: #999;
						font-weight: 600;
					}
				}
			}
		}
	}

	.thirdBlock {
		display: flex;
		justify-content: space-between;
		width: 100%;

		.imgBox:nth-child(1) span {
			color: red;
			left: 5%;
			top: 15%;
			width: 1em;
			font-size: 1.2em;
			font-weight: 600;
		}

		.imgBox:nth-child(2) span {
			color: #fff;
			font-size: 1.2em;
			font-weight: 600;
			bottom: 9%;
			right: 8%;
			letter-spacing: 5px;
		}
		.imgBox:nth-child(4) span {
			color: #fff;
			font-size: 1.2em;
			font-weight: 600;
			bottom: 9%;
			right: 8%;
			letter-spacing: 2px;
			bottom: 30%;
		}

		.imgBox:nth-child(3) span {
			font-size: 1.2em;
			font-weight: 600;
			color: rgb(78, 78, 151);
			bottom: 30%;
			right: 6%;
		}

		.imgBox {
			height: 100%;
			position: relative;

			span {
				position: absolute;
			}

			img {
				height: 180px;
				width: 100%;
				object-fit: contain;
			}
		}
	}

	.more {
		display: inline-block;
		position: absolute;
		color: rgb(91, 107, 115);
		font-size: 0.4em;
		right: 0;
		cursor: pointer;
		font-weight: 400;
	}
}
</style>
