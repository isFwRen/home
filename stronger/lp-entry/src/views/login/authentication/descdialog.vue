<template>
    <v-dialog v-model="dialog" transition="dialog-bottom-transition" :width="dialogWidth" @keydown.enter="onClose">
        <div class="desc-dialog-wrap">
            <v-row no-gutters style="height: 50px;" align="center">
                <v-col cols="11">
                    <div class="desc-header-title">新版验证码获取功能使用说明</div>
                </v-col>
                <v-col cols="1" style="text-align: center;">
                    <v-btn icon @click="onClose">
                        <v-icon>mdi-close</v-icon>
                    </v-btn>
                </v-col>
            </v-row>
            <div class="desc-scroll-content">
                <div class="desc-scroll">
                    <v-row no-gutters>
                        <v-col>
                            <p> <span class="font-weight-bold">第一步：</span> 下载app</p>
                            <p>安卓app：<a
                                    href="https://f-droid.org/repo/com.beemdevelopment.aegis_57.apk">https://f-droid.org/repo/com.beemdevelopment.aegis_57.apk</a>
                            </p>
                            <p>苹果app：Google Authenticator</p>
                            <div class="desc-img-box">
                                <v-img :src="authImages.imgUrl1" :lazy-src="authImages.imgUrl1" />
                            </div>
                        </v-col>
                    </v-row>

                    <v-row no-gutters>
                        <v-col cols="12">
                            <p class="p-margin-top"> <span
                                    class="font-weight-bold">第二步：</span>打开需要访问的网址，点击页面右上角“身份验证”，根据输入框提示输入对应手机号及密码，点击验证
                            </p>
                            <div class="desc-img-box-bg">
                                <v-img :src="authImages.imgUrl2" :lazy-src="authImages.imgUrl2" />
                            </div>
                            <p class="p-margin-top"> <span class="font-weight-bold">第三步：</span>用下载的app扫描右侧二维码，保存人员信息。
                            </p>
                            <p class="text-red"> （手机app有该账号信息后，后续登录无需再扫二维码。）</p>
                            <div class="desc-img-box-bg">
                                <v-img :src="authImages.imgUrl3"
                                    :lazy-src="authImages.imgUrl3" />
                            </div>
                            <p class="p-margin-top"> <span class="font-weight-bold">第四步：</span>回到登录页面，输入手机号和密码，直接打开app查看
                                <span class="text-bold">六位数验证码</span>
                                ，将验证码输入后即可登录该系统。
                            </p>
                            <div class="desc-img-box-md">
                                <v-img :src="authImages.imgUrl4"
                                    :lazy-src="authImages.imgUrl4" />
                            </div>
                            <div class="text-red tips">
                                注：首次登录系统按此流程进行，APP扫二维码保存人员信息后，往后每次登录只需要录入手机号+密码+打开手机app查看6位数验证码输入后登录。（所有系统通用）
                            </div>
                        </v-col>
                    </v-row>

                </div>
            </div>
            <div class="desc-btn"><v-btn color="primary"> 关闭</v-btn>
            </div>
        </div>
    </v-dialog>
</template>

<script>
import cells from './cells'
export default {
    name: 'descdialog',
    data() {
        return {
            dialog: false,
            dialogWidth: Math.ceil(document.documentElement.clientWidth * 0.62) + 'px',
            authImages: cells.authImages
        }
    },
    props: {
        componwidth: {
            type: Number
        }
    },
    methods: {
        onOpen() {
            this.dialog = true
        },
        onClose() {
            this.dialog = false
        }
    }

}
</script>

<style  lang="scss">
.font-weight-bold {
    font-weight: 600;
}

.p-margin-top {
    margin-top: 10px;
}

.v-dialog {
    overflow-y: hidden !important;
    max-height: 700px !important;
}

.desc-dialog-wrap {
    background-color: #fff;
    padding: 10px 10px;
    height: 90%;

    .desc-header-title {
        font-size: 14px;
        font-weight: 600;
    }

    .desc-scroll-content {
        height: 500px;
        overflow-y: scroll;
        padding-bottom: 40px;

        .desc-scroll {
            font-size: 14px;

            .desc-img-box {
                width: 48%;
            }

            .desc-img-box-bg {
                width: 60%;
            }

            .desc-img-box-md {
                width: 52%;
            }

            .text-red {
                color: red;
            }
        }

    }

    .desc-btn {
        width: 100px;
        height: 50px;
        padding-top: 15px;
        margin: 0 auto;
    }
}
</style>
