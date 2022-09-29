<template>
    <div class="login">
      <div class="main">
        <div class="logoContainer">
          <div class="logo"><img src="../assets/img/logo.png" alt="" /></div>
          <div class="name">666盘</div>
        </div>
        <div
          class="mainBox"
          :class="activeName === 'first' ? '' : 'mainBoxRegistered'"
        >
          <el-tabs
            v-model="activeName"
            type="card"
            stretch
          >
            <el-tab-pane label="登录" name="first">
              <div class="loginInput">
                <el-form ref="form" :model="login" label-width="80px">
                  <el-form-item>
                    <el-input
                      v-model="login.name"
                      placeholder="请输入姓名"
                    ></el-input>
                  </el-form-item>
                  <el-form-item>
                    <el-input
                      v-model="login.password"
                      type="password"
                      placeholder="请输入密码"
                    ></el-input>
                  </el-form-item>
                  <el-form-item>
                    <el-button type="primary" @click="onSubmit">登录</el-button>
                  </el-form-item>
                </el-form>
              </div>
            </el-tab-pane>
            <el-tab-pane label="注册" name="second">
              <div class="registeredInput">
                <el-form ref="form" :model="login" label-width="80px">
                  <el-form-item>
                    <el-input
                      v-model="registered.email"
                      placeholder="请输入邮箱"
                    ></el-input>
                  </el-form-item>
                  <el-form-item
                    >
                    <el-input
                      v-model="registered.password"
                      placeholder="请输入密码"
                      type="password"
                    ></el-input>
                  </el-form-item>
                  <el-form-item>
                    <el-input
                      v-model="registered.name"
                      placeholder="请输入名称"
                    ></el-input>
                  </el-form-item>
                  <el-form-item class="codeContainer">
                    <el-input
                      v-model="registered.code"
                      placeholder="请输入验证码"
                    ></el-input>
                    <div class="codeButtonContainer">
                      <el-button
                        size="mini"
                        class="getcode"
                        v-if="!isCountDownShow"
                        @click="getCode"
                        >获取验证码</el-button
                      >
                      <div class="countDown" v-else>{{ countDownSecond }} s</div>
                    </div>
                  </el-form-item>
                  <el-form-item>
                    <el-button type="primary" @click="clickRegistered"
                      >注册</el-button
                    >
                  </el-form-item>
                </el-form>
              </div>
            </el-tab-pane>
          </el-tabs>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  // 倒计时名称
  let timer;
  
  export default {
    name: "Login",
    data() {
      return {
        login: {
          name: "",
          password: "",
        },
        registered: {
          email: "",
          password: "",
          code: "",
          name: "",
        },
        activeName: 'first',
        // 倒计时秒数
        countDownSecond: 60,
        // 是否显示秒数
        isCountDownShow: false,
      };
    },
    methods: {
      onSubmit() {
        this.$axios({
          url: "/user/login",
          method: "post",
          data: {
            name: this.login.name,
            password: this.login.password,
          },
        }).then((res) => {
          console.log(res.data);
          if (res.data.success) {
            this.$message.success("登录成功！");
          // 将返回的用户信息保存至localstorage中
          window.localStorage.setItem(
            "userInfo",
            JSON.stringify(this.login.name)
          );
          // 将token存入本地
          window.localStorage.setItem("token", res.data.token);
          //   跳转至主界面
          this.$router.push("/index");
          } else {
            this.$message.warning("登录失败,账号或密码错误!");
          }
        });
      },
  
      // 获取验证码
      async getCode() {
        this.isCountDownShow = true;
        this.$axios({
          url: "/user/code/send/register",
          method: "post",
          data: {
            email: this.registered.email
          },
        }).then((res) => {
          console.log(res.data);
          if (res.data.success) {
            this.startCountDown();
          }
        });
      },
  
      // 倒计时
      startCountDown() {
        this.countDownSecond = 60;
        timer = setInterval(() => {
          this.countDownSecond--;
          if (this.countDownSecond == 0) {
            clearInterval(timer);
            this.isCountDownShow = false;
          }
        }, 1000);
      },
  
      clickRegistered() {
        this.$axios({
          url: "/user/register",
          method: "post",
          data: {
            name: this.registered.name,
            password: this.registered.password,
            code: this.registered.code,
            email: this.registered.email,
          },
        }).then((res) => {
          console.log(res.data);
          if (res.data.success) {
            this.$message.success("注册成功!");
            this.login.name = this.registered.name;
            this.activeName = 'first';
            this.registered = {
              email: "",
              password: "",
              code: "",
              name: "",
            };
          } else {
            this.$message.error("注册失败,请稍后重试!");
          }
        });
      },
    },
  };
  </script>
  
  <style scoped>
  .login {
    background-color: #ecefff;
    height: 100vh;
  }
  
  .logoContainer {
    display: flex;
    align-items: center;
    margin-bottom: 10px;
    justify-content: center;
  }
  
  .logo {
    width: 50px;
  }
  
  .logo img {
    width: 100%;
  }
  
  .name {
    color: #25262b;
    font-size: 20px;
    letter-spacing: 4px;
    font-weight: bold;
    margin-left: 7px;
  }
  
  .main {
    width: 350px;
    height: 400px;
    position: absolute;
    right: 10vw;
    top: 15vh;
  }
  
  .mainBox {
    width: 350px;
    background-color: #fff;
    height: 330px;
    border-radius: 10px;
    overflow: hidden;
  }
  
  .mainBoxRegistered {
    height: 430px;
  }
  
  .el-form /deep/ .el-form-item__content {
    margin: 0 !important;
    padding: 0 20px;
  }
  
  .el-input /deep/ input {
    border-radius: 10px;
  }
  
  .loginInput {
    margin-top: 20px;
  }
  
  .el-tabs /deep/ .is-active,
  .el-tabs /deep/ div:hover {
    color: #595bb3;
  }
  
  .el-tabs /deep/ .is-active {
    background-color: #fff;
  }
  
  .el-tabs /deep/ .el-tabs__item {
    border: none !important;
    font-size: 18px;
    height: 50px;
    line-height: 50px;
  }
  
  .el-tabs /deep/ .el-tabs__nav {
    border: none;
  }
  
  .el-tabs /deep/ .el-tabs__nav-scroll {
    background-color: #f5f5f6;
  }
  
  .el-input /deep/ .el-input__inner {
    height: 48px;
    font-size: 15px;
  }
  
  .el-button {
    width: 100%;
    background-color: #6c6dbb;
    border: none;
    border-radius: 10px;
    margin-top: 10px;
    height: 45px;
    font-size: 15px;
  }
  
  .el-button:hover {
    background-color: #595bb3;
  }
  
  .codeContainer {
    position: relative;
  }
  
  .codeButtonContainer {
    position: absolute;
    top: 50%;
    right: 30px;
    transform: translateY(-50%);
  }
  
  .getcode {
    background-color: #6c6dbb;
    color: white;
    height: 35px;
    margin: 0;
  }
  
  .countDown {
    color: rgb(141, 141, 141);
  }
  </style>
  