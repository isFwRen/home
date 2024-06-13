<template>
    <transition enter-active-class="animate__animated animate__fadeIn"
        leave-active-class="animate__animated animate__fadeOut">
        <div class="authen-wrap" v-if="dialog">
            <div class="authen-wrap-bgcolo" @click="onClose"></div>
            <!---验证框--->
            <div class="authen-content">
                <v-row no-gutters style="height: 50px;" align="center">
                    <v-col cols="11">
                        <div class="header-title">
                            <span class="header-title-icon">
                                <v-icon>mdi-account-key</v-icon>
                            </span>
                            <h5 class="text-lg-h6" style="font-size: 16px !important;">
                                身份认证
                            </h5>
                        </div>
                    </v-col>
                    <v-col cols="1" style="text-align: center;">
                        <v-btn icon @click="onClose">
                            <v-icon>mdi-close</v-icon>
                        </v-btn>
                    </v-col>
                </v-row>
                <v-row align="center" justify="center" no-gutters style="margin-top: 15px;">
                    <v-col cols="8">
                        <div class="authen-forms">
                            <v-form ref="entryForm">
                                <div class="forms-item">
                                    <v-text-field :label="userField.label" v-model="forms.phone" :rules="userField.rules">
                                    </v-text-field>
                                </div>
                                <div class="forms-item">
                                    <v-text-field :label="passwordField.label" @click:append="eye = !eye"
                                        :type="eye ? 'text' : 'password'" :append-icon="eye ? 'mdi-eye' : 'mdi-eye-off'"
                                        v-model="forms.password" :rules="passwordField.rules">
                                    </v-text-field>
                                </div>
                                <!-- <div class="forms-item flex-item">
                                    <v-text-field label="验证码" v-model="forms.captcha" :rules="captchaField.rules">
                                    </v-text-field>
                                    <z-btn class="pb-5 ml-4" color="primary" :disabled="!validAccount || counting"
                                        @click="sendCode">{{ text }}</z-btn>
                                </div> -->
                                <div class="forms-item">
                                    <z-btn block color="primary" @click="authFun">验证</z-btn>
                                </div>
                            </v-form>
                        </div>
                    </v-col>
                    <v-col>
                        <div class="authen-img-box" v-if="isverify">
                            <div class="authen-img">
                                <div class="au-img-box" style="background-color: #fff; height: 120px;z-index:10; ">
                                    <v-img :src="logoImg" :lazy-src="logoImg" width="120" />
                                </div>
                                <div class="rects-logo">
                                    <div class="top-left rect"></div>
                                    <div class="top-right rect"></div>
                                    <div class="bottom-left rect"></div>
                                    <div class="bottom-right rect"></div>
                                </div>

                            </div>
                            <div class="desc" @click="descDialot">身份验证流程介绍</div>
                        </div>
                        <div class="code-img-box" v-else>
                            <div class="authen-img">
                                <div class="v-img-box">
                                    <div style="display:none">
                                        <canvas ref="canvas" id="canvas"></canvas>
                                    </div>
                                    <v-img :src="qrImage" :lazy-src="qrImage" width="160" height="160" />
                                </div>
                                <div class="rects-qr">
                                    <div class="top-left rect"></div>
                                    <div class="top-right rect"></div>
                                    <div class="bottom-left rect"></div>
                                    <div class="bottom-right rect"></div>
                                </div>
                            </div>
                            <p class="code-desc">使用验证码app扫描上方的二维码获取登录验证码后关闭弹窗，使用该验证码进行登录</p>
                        </div>
                    </v-col>
                </v-row>
            </div>
            <!---身份验证浏览--->
            <descdialog ref="descDiaDom"></descdialog>
        </div>
    </transition>
</template>
<script>
import cells from './cells'
import { regex_phone } from './cells'
import { getQRcode } from '../../../plugins/qrcode'
import ButtonMixins from '@/mixins/ButtonMixins'
export default {
    name: 'authentication',
    mixins: [ButtonMixins],
    data() {
        return {
            formId: 'authentication',
            dialog: false,
            userField: cells.phoneField,
            passwordField: cells.passwordField,
            captchaField: cells.captchaField,
            eye: false,
            isverify: true,
            qrImage: '',
            validAccount: false,
            text: '获取验证码',
            capRules: [],
            logoImg: require('@/assets/logo.png'),
            forms: {
                phone: '',
                password: '',
                // captcha: ''
            }
        }
    },
    watch: {
        'forms.phone': {
            handler() {
                if (this.forms['phone']) {
                    const { phone } = this.forms
                    this.validAccount = this.networkIdetify ? !!phone : phone && regex_phone.test(phone)
                }
            },
            immediate: true
        }
    },
    props: ['networkIdetify'],
    methods: {
        async authFun() {
            for (let i in this.forms) {
                if (this.forms[i] === '') {
                    this.captchaField.rules = [value => !!value || '请输入验证码.']
                    this.$refs.entryForm.validate()
                    return
                }
            }

            const form = {
                isIntranet: this.networkIdetify,
                ...this.forms
            }
            const result = await this.$store.dispatch(
                'VERIFY_IDENTITY',
                form
            )
            if (result.code === 200) {
                this.toasted.dynamic(result.msg, result.code)
                this.qrImage = getQRcode(result.data, 160)
                this.isverify = false;
            } else {
                this.toasted.dynamic(result.msg, result.code);
            }
        },
        async sendCode() {
            this.clock = setInterval(this.countdown, 1000)
            const form = {
                isIntranet: this.networkIdetify,
                account: this.forms['phone'],
                accountKey: 'phone'
            }

            const result = await this.$store.dispatch(
                'GET_DING_CODE',
                form
            )
            if (result.code === 200) {
                this.toasted.dynamic(result.msg, result.code)
            }
        },
        descDialot() {
            this.$refs.descDiaDom.onOpen()
        },
        onOpen() {
            this.dialog = true
        },
        onClose() {
            this.forms = {
                phone: '',
                password: ''
            }
            this.qrImage = '';
            this.isverify = true;
            this.dialog = false;
        }
    },
    components: {
        'descdialog': () => import('./descdialog.vue')
    }
}
</script>

<style  lang="scss">
@mixin position_box ($position: null, $left: null, $top: null, $right: null, $bottom: null, $index: null) {
    position: $position;

    @if $left {
        left: $left;
    }

    @if $right {
        right: $right;
    }

    @if $bottom {
        bottom: $bottom;
    }

    @if $top {
        top: $top;
    }

    @if $index {
        z-index: $index
    }
}

@mixin w_h($w: null, $h: null) {
    @if $w {
        width: $w;
    }

    @if $h {
        height: $h
    }
}







.authen-wrap {
    @include position_box (absolute, 0, 0, 0, 0, 10);

    .authen-content {
        width: 600px;
        padding-bottom: 50px;
        @include position_box (absolute, 50%, 50%, null, null, 10);
        margin-top: -200px;
        margin-left: -300px;
        background-color: #fff;
        border-radius: 2px;

        .header-title {
            display: flex;
            align-items: center;
            height: 50px;

            .header-title-icon {
                margin-left: 10px;
            }

            .text-lg-h6 {
                margin-left: 10px;
            }
        }
    }

    .authen-forms {
        padding: 0 22px;

        .forms-item {
            padding-top: 10px;
        }

        .flex-item {

            display: flex;
            align-items: center;
        }
    }

    .authen-img-box {
        padding: 5px 5px;
        width: 150px;
        position: relative;
        margin: 0 auto;

        .authen-img {
            @include w_h(125px, 125px);
            box-sizing: border-box;
            position: relative;
            box-sizing: border-box;

            .v-image__image--cover {
                background-size: 80% !important;
            }

            .rects-logo {
                width: 129px;
                height: 124px;
                @include position_box (absolute, -2px, -2px, null, null, -1);


                .rect {
                    position: absolute;
                    @include w_h(25px, 25px);
                    border: 2px solid #6a6a6a;
                }

                .top-left {
                    top: 0;
                    left: 0;
                }

                .top-right {
                    top: 0;
                    right: 0
                }

                .bottom-left {
                    bottom: 0;
                    left: 0
                }

                .bottom-right {
                    bottom: 0;
                    right: 0
                }
            }
        }

        .desc {
            width: 127px;
            color: #2440b3;
            cursor: pointer;
            text-align: center;
            font-size: 12px;
            margin-top: 20px;
            margin-right: 15px;
            letter-spacing: 2px;
        }
    }

    .code-img-box {
        width: 160px;
        position: relative;
        margin: 0 auto;
        font-size: 14px;
        color: #6c6c6c;
        transform: translateY(30px);

        .authen-img {
            position: relative;
            @include w_h(160px, 160px);
            box-sizing: border-box;
            box-shadow: 0 3px 1px -2px rgba(0, 0, 0, .2), 0 2px 2px 0 rgba(0, 0, 0, .14), 0 1px 5px 0 rgba(0, 0, 0, .12) !important;
        }

        .rects-qr {
            @include w_h(166px, 166px);
            @include position_box (absolute, -3px, -3px, null, null, -1);

            .rect {
                position: absolute;
                @include w_h(25px, 25px);
                border: 2px solid #6a6a6a;
            }

            .top-left {
                top: 0;
                left: 0;
            }

            .top-right {
                top: 0;
                right: 0
            }

            .bottom-left {
                bottom: 0;
                left: 0
            }

            .bottom-right {
                bottom: 0;
                right: 0
            }
        }

        .v-img-box {
            background-color: #fff;
        }



        .code-desc {
            padding-top: 5px;
        }
    }

    .authen-wrap-bgcolo {
        background-color: rgba(0, 0, 0, 0.3);
        @include position_box (absolute, 0, 0, 0, 0, 9);
    }
}
</style>