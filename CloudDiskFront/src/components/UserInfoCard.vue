<template>
  <div class="UserInfoCard" v-if="userInfo">
    <div class="storageProgressContainer">
      <div class="">{{ storageSize }}MB/1G</div>
      <el-progress
        :percentage="
          storageProgress.toFixed(0) * 1 ? storageProgress.toFixed(0) * 1 : 0
        "
        :format="() => storageProgress + '%'"
        color="#696bcc"
        class="storageProgress"
      ></el-progress>
    </div>

    <div class="avatar">
      <img :src="userInfo.avatar" alt="" v-if="userInfo.avatar" />
      <img src="../assets/img/avatar.png" alt="" v-else />
    </div>
    <div class="userName">{{ userInfo.name }}</div>
    <div class="menuContainer">
      <div
        class="menu"
        :class="$store.state.isUserInfoCardMenuShow ? 'showMenu' : ''"
      >
        <div class="group">
          <div class="menuItem" @click="changeNickName">修改昵称</div>
        </div>
        <div class="group">
          <div class="menuItem" @click="clickAbout">关于</div>
        </div>
        <div class="group">
          <div class="menuItem" @click="logout">退出登录</div>
        </div>
      </div>
      <i class="iconfont icon-setting" @click="showMenu"></i>
    </div>

    <!-- 修改昵称的输入框dialog -->
    <el-dialog
      title="修改昵称"
      :visible.sync="isNickNameDialogShow"
      width="400px"
    >
      <el-input
        v-model="newNickName"
        autocomplete="off"
        class="nickNameInput"
      ></el-input>
      <div slot="footer" class="dialog-footer">
        <el-button @click="isNickNameDialogShow = false" size="small"
          >取 消</el-button
        >
        <el-button type="primary" @click="confirmNickName" size="small"
          >确 定</el-button
        >
      </div>
    </el-dialog>

    <!-- '关于'框dialog -->
    <el-dialog
      title="关于666网盘"
      :visible.sync="isAboutDialogShow"
      width="500px"
      class="aboutDialog"
    >
      使用go-zero + vue + mysql + redis + rabbitmq实现的分布式网盘
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: "UserInfoCard",
  data() {
    return {
      userInfo: {
        name: null,
        email: null,
      },
      // 是否显示menu
      isMenuShow: false,
      // 是否显示修改昵称的dialog
      isNickNameDialogShow: false,
      // 新的昵称
      newNickName: "",
      // 是否显示about的dialog
      isAboutDialogShow: false,
    };
  },
  methods: {
    // 点击退出登录的回调
    logout() {
      // this.showMenu();
      window.localStorage.removeItem("userInfo");
      this.$store.commit("updateUserInfo", {});
      this.$router.push("/login");
    },

    // 显示菜单
    showMenu() {
        // if (this.$store.state.isUserInfoCardMenuShow) {
        //   this.$store.state.isUserInfoCardMenuShow = false;
        // } else {
        //   this.$store.state.isUserInfoCardMenuShow = true;
        // }
        // console.log(!this.$store.state.isUserInfoCardMenuShow);
        setTimeout(() => {
          this.$store.commit("updateIsUserInfoCardMenuShow", true);
        });
    },

    // 请求用户信息
    async getUserInfo() {
      this.$axios.defaults.headers.common["Authorization"] = window.localStorage.getItem("token");
      this.$axios({
        url: "/user/detail",
        method: "post",
      }).then((res) => {
        console.log(res.data);
        this.userInfo.name = res.data.name;
        this.userInfo.email = res.data.email;
        this.$store.commit("updateUserInfo", this.userInfo);
      });
    },

    // 点击修改昵称的回调
    changeNickName() {
      // this.showMenu();
      this.newNickName = this.userInfo.name;
      this.isNickNameDialogShow = true;

      // 聚焦输入框
      // dialog在第一次打开前是不渲染的 所以里面的input需要在下一帧中获取
      this.$nextTick(() => {
        // console.log([document.querySelector(".nickNameInput")]);
        document.querySelector(".nickNameInput").children[0].focus();
      });
    },

    // 确认修改昵称
    async confirmNickName() {
      // 先判断昵称是否合法
      if (this.newNickName == "") {
        this.$message.warning("昵称不能为空哦!");
        return;
      } else if (this.newNickName == this.userInfo.name) {
        this.isNickNameDialogShow = false;
        return;
      }

      // 深拷贝对象 以防出现请求没有成功 但是昵称已经在this.userInfo中被修改的情况
      let data = JSON.parse(JSON.stringify(this.userInfo));
      data.nickname = this.newNickName;
      // 调用此接口以通知后端将上传的文件存入数据库
      this.$axios({
        url: "/user/name/update",
        method: "put",
        data: {
          name: this.newNickName
        }
      }).then((res) => {
        console.log(res.data);
        if (res.data.success) {
          this.$message.success("昵称修改成功!");
          // 更新用户数据
          this.getUserInfo();
          this.isNickNameDialogShow = false;
        } else {
          this.$message.error(res.data.msg);
        }
      });
    },
    clickAbout() {
      // this.showMenu();
      this.isAboutDialogShow = true;
    }
  },
  created() {
    this.getUserInfo();
  },
  mounted() {},
  computed: {
    // 内存进度条
    storageProgress() {
      return ((this.userInfo.neicun / 1048576 / 1024) * 100).toFixed(2) * 1;
    },
    storageSize() {
      return (this.userInfo.neicun / 1048576).toFixed(2);
    },
  },
  watch: {
    "$store.state.userInfo"(current) {
      this.userInfo = current;
    },
    // 上面监听不到内存属性的变化，因为但内存属性发生改变时，userInfo的地址没有发生变化
    "$store.state.userInfo.neicun"(current) {
      this.userInfo.neicun = current;
    },
  },
};
</script>

<style scoped>
.UserInfoCard {
  height: 75px;
  width: 240px;
  border-top: 1px solid #ccc;
  display: flex;
  align-items: center;
  padding: 20px;
  box-sizing: border-box;
  color: #25262b;
  position: relative;
}

.storageProgressContainer {
  position: absolute;
  top: -45px;
  left: 25px;
  width: 100%;
  font-size: 13px;
}

.storageProgress {
  display: flex;
  align-items: center;
  width: 85%;
  align-items: center;
  margin-top: 5px;
}

.storageProgress /deep/ .el-progress__text {
  font-size: 13px !important;
  margin-left: 12px;
}

.avatar img {
  width: 40px;
  height: 40px;
  object-fit: cover;
  border-radius: 50%;
}

.userName {
  font-size: 14px;
  margin-left: 15px;
  width: 110px;
  text-overflow: ellipsis;
  overflow: hidden;
  white-space: nowrap;
}

.iconfont {
  font-size: 18px;
  color: rgb(104, 104, 104);
  cursor: pointer;
}

.menuContainer {
  position: absolute;
  right: 20px;
}

.menu {
  position: absolute;
  background-color: white;
  /* display: none; */
  bottom: 30px;
  width: 150px;
  border: 1px solid #ddd;
  box-shadow: 0px 0px 10px 1px rgba(0, 0, 0, 0.1);
  z-index: 3000;
  padding: 5px 5px;
  border-radius: 7px;
  font-size: 15px;
  display: none;
}

.UserInfoCard /deep/ .el-upload {
  width: 100%;
  cursor: unset;
  text-align: left;
}

.showMenu {
  display: block;
}

.group {
  padding: 4px 0;
  border-bottom: 1px solid #eee;
}

.group:last-child {
  border: none;
}

.group > div {
  padding: 4px 20px;
  font-size: 14px;
  color: rgb(70, 70, 70);
  user-select: none;
}

.group > div:hover {
  background-color: #696bcc;
  color: white;
}

.UserInfoCard /deep/ .el-dialog {
  border-radius: 10px;
}

.UserInfoCard /deep/ .aboutDialog {
  line-height: 20px;
}

.UserInfoCard /deep/ .aboutDialog .el-dialog__body {
  padding-top: 10px;
}
</style>
