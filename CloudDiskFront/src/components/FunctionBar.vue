<template>
  <div>
    <div class="FunctionBar">
      <div class="left">
        <el-upload
          :show-file-list="false"
          multiple
          action=""
          :http-request="uploadFiles"
          class="uploadButton"
        >
          <el-button type="primary" size="small" class="upload">
            <i class="iconfont icon-yunshangchuan"></i> 上传</el-button
          >
        </el-upload>
        <el-button
          size="small"
          class="create"
          @click="createFolder"
          :disabled="!isCreateAble"
        >
          <i class="iconfont icon-add"></i> 新建</el-button
        >
        <el-button
          size="small"
          class="selectAll"
          :class="isSelectAll ? 'select' : ''"
          @click="refresh"
        >
          <i class="el-icon-refresh"></i>
          刷新</el-button
        >
      </div>

      <!-- 右边 -->
      <div class="right">
        <div class="search">
          <el-input
            placeholder="请输入内容"
            suffix-icon="el-icon-search"
            v-model="searchContent"
            @keyup.native.enter="$emit('goSearch', searchContent)"
          >
          </el-input>
        </div>
      </div>
      <div class="goLastFolder" v-if="this.$store.state.currentParentId != 0">
        <el-button
          @click="goLastFolder"
          type="info"
          icon="el-icon-arrow-left"
          size="mini"
          >上一级</el-button
        >
        <el-button @click="goRootFolder" type="primary" size="mini"
          >根目录 &nbsp;<i class="el-icon-s-home"></i
        ></el-button>
      </div>
    </div>
  </div>
</template>

<script>
import SparkMD5 from "spark-md5";
// 文件分片的大小
const chunkSize = 5;
export default {
  name: "FunctionBar",
  data() {
    return {
      isCreateAble: true,
      // 上传按钮是否是加载状态
      isUploadBtnLoading: false,
      // 是否显示多选操作按钮
      isMultBtnsShow: false,
      // 搜索内容
      searchContent: "",
      // 上传文件的hash值
      fileHash: null,
      // obs返回的key
      key: null,
      // obs返回的uploadId
      uploadId: null,
      // obs返回的etag
      etagArr: [],
      // 分片上传的node
      node: null,
      // 分片上传文件的uuid
      fileIdentity: "",
    };
  },
  props: {
    listData: {
      type: Array,
      default() {
        return [];
      },
    },
    searchFolder: {
      type: Array,
      default() {
        return [];
      },
    },
    // functionbar的类型 file collect
    barType: {
      type: String,
      default() {
        return "file";
      },
    },
  },
  methods: {
    refresh() {
      // this.isSelectAll = !this.isSelectAll;
      // this.$store.commit("updateIsSelectAll", this.isSelectAll);
    },
    // 点击新建的回调
    createFolder() {
      // 先禁用新建按钮，避免重复点击
      this.isCreateAble = false;
      // 更新新建文件夹状态到vuex
      this.$store.commit("updateIsCreateFolder", true);
    },
    // 上传文件
    uploadFiles(params) {
      // 根据文件大小决定上传类型
      // 得到文件大小
      var fileSize = (params.file.size / 1048576).toFixed(2);
      if (fileSize <= 5) {
        // 普通上传
        // this.$message('文件正在上传...');
        this.commonUpload(params.file);
      } else {
        // 分片上传前的准备
        console.log("开始分片上传");
        this.fileToMD5(params.file)
        setTimeout(() => {
          console.log("文件哈希值：", this.fileHash);
          this.chunkUploadPrepare(this.fileHash, params.file.name, this.$store.state.currentParentId, params.file, fileSize);
        },300);
      }
    },
    // 普通上传
    commonUpload(file) {
      console.log("开始普通上传文件");
      const formData = new FormData();
      formData.append("file", file);
      formData.append("parent_id", this.$store.state.currentParentId);
      this.$axios.defaults.headers.common["Authorization"] =
      window.localStorage.getItem("token");
      this.$axios({
        url: "/file/upload",
        method: "post",
        data: formData,
        headers: {
          "Content-Type": "multipart/form-data",
        },
      }).then((res) => {
        console.log(res.data);
        if (res.data.success) {
          this.$axios.defaults.headers.common["Authorization"] = window.localStorage.getItem("token");
          this.$axios({
            url: "/user/file/list",
            method: "post",
            data: {
              id: this.$store.state.currentParentId,
            }
        }).then((res) => {
          console.log(res.data);
          if (res.data.success) {
            this.$emit("update:listData", res.data.files_list);
          } else {
            this.$message.error(res.data.msg);
          }
        });
          this.$message.success("文件上传成功!");
        } else {
          this.$message.error(res.data.msg);
        }
      });
    },
    // chunkUpload(file) {

    // },
    chunkUploadPrepare(md5, name, parent_id, file, fileSize) {
      // MD5      string `json:"md5"`
      // Name     string `json:"name"`
      // ParentId int    `json:"parent_id"`
      this.$axios.defaults.headers.common["Authorization"] =
        window.localStorage.getItem("token");
      this.$axios({
        url: "/file/upload/prepare",
        method: "post",
        data: {
          md5: md5,
          name: name,
          parent_id: parent_id,
        },
      }).then((res) => {
        console.log(res.data);
        if (res.data.success) {
          // 需要分片上传
          console.log("uploadId: ", res.data.upload_id);
          console.log("key: ", res.data.key);
          this.uploadId = res.data.upload_id;
          this.key = res.data.key;
          this.node = res.data.node;
          this.fileIdentity = res.data.file_identity;
          // 开始分片上传
          this.chunkUploadStart(file, fileSize);
        } else {
          // 文件秒传成功
          if (res.data.msg == "文件上传成功") {
            this.$message.success(res.data.msg);
          } else {
            this.$message.error(res.data.msg);
          }
          
        }
      });
    },
    chunkUploadStart(file, fileSize) {
      console.log("分片上传开始");
      this.etagArr = new Array();
      // 分片上传是否失败的标记变量
      var isChunkUploadFail = false;
      // 得到文件分片的数量
      var chunkCount = Math.ceil(fileSize / chunkSize);
      console.log("文件大小：", fileSize);
      console.log("文件分块数量：", chunkCount);
      // console.log("文件：", file.raw);
      for (let i = 0; i < chunkCount; i++) {
        // 分片开始的位置
        let start = i * chunkSize * 1024 * 1024;
        // 分片结束的位置
        let end = Math.min(file.size, start + chunkSize * 1024 * 1024)
        // 截取分片的文件
        let _chunkFile = file.slice(start, end);
        console.log("开始上传第", i + 1, "个分片");
        let formData = new FormData();
        formData.append("file", _chunkFile);
        formData.append("part_number", i + 1)
        formData.append("key", this.key);
        formData.append("upload_id", this.uploadId);
        formData.append("node", this.node);
        this.$axios.defaults.headers.common["Authorization"] =
          window.localStorage.getItem("token");
          this.$axios({
            url: "/file/upload/chunk",
            method: "post",
            data: formData,
            headers: {
              "Content-Type": "multipart/form-data",
            },
          }).then((res) => {
            if (res.data.success) {
              console.log("第几块文件上传成功：", i + 1);
              console.log("对应的etag：", res.data.etag);
              this.etagArr.push({
                part_number: i + 1,
                etag: res.data.etag,
              });
            } else {
              isChunkUploadFail = true;
            }
          });
          if (isChunkUploadFail == true) {
            break;
          }
      }
      setTimeout(() => {
        if (isChunkUploadFail == true) {
          this.$message.error("文件上传失败");
        } else {
          // 分片文件都上传成功，告诉后端
          this.chunkUploadComplete(file);
        }
      }, chunkCount * 2000);
    },
    chunkUploadComplete(file) {
      console.log("分片文件都上传成功，告诉后端");
      console.log("etagArr数组的长度：", this.etagArr.length);
      console.log("待上传的fileIdentity: ", this.fileIdentity);
      for (let i = 0; i < this.etagArr.length; i++) {
        console.log("上传的文件块：", this.etagArr[i].part_number, ", 分块文件的哈希值：", this.etagArr[i].etag);
      }
      // let formData = new FormData();
      //   formData.append("key", this.key);
      //   formData.append("upload_id", this.uploadId)
      //   formData.append("obs_objects", this.etagArr);
      //   formData.append("parent_id", this.$store.state.currentParentId);
      //   formData.append("hash", this.fileHash);
      //   formData.append("name", file.name);
      //   formData.append("size", file.size);
      //   formData.append("node", this.node);
      //   formData.append("file_identity", this.fileIdentity);
      this.$axios.defaults.headers.common["Authorization"] =
        window.localStorage.getItem("token");
      this.$axios({
        url: "/file/upload/chunk/complete",
        method: "post",
        data: {
          key: this.key,
          upload_id: this.uploadId,
          obs_objects: this.etagArr,
          parent_id: this.$store.state.currentParentId,
          hash: this.fileHash,
          name: file.name,
          size: file.size,
          node: this.node,
          file_identity: this.fileIdentity,
        }
      }).then((res) => {
        console.log(res.data);
        if (res.data.success) {
          this.$axios.defaults.headers.common["Authorization"] = window.localStorage.getItem("token");
          this.$axios({
            url: "/user/file/list",
            method: "post",
            data: {
              id: this.$store.state.currentParentId,
            }
        }).then((res) => {
          console.log(res.data);
          if (res.data.success) {
            this.$emit("update:listData", res.data.files_list);
          } else {
            this.$message.error(res.data.msg);
          }
        });
          this.$message.success("文件上传成功");

        } else {
          this.$message.error("文件上传失败");
        }
      });
    },
    

    // 返回上一级文件夹
    goLastFolder() {
      console.log("返回上一级");
      this.$store.commit(
        "updateCurrentParentId",
        this.$store.state.lastParentId
      );
      console.log("currentParentId: ", this.$store.state.currentParentId);
      this.$axios.defaults.headers.common["Authorization"] =
        window.localStorage.getItem("token");
      this.$axios({
        url: "/user/file/list",
        method: "post",
        data: {
          id: this.$store.state.currentParentId,
        },
      }).then((res) => {
        console.log(res.data);
        if (res.data.success) {
          this.$emit("update:listData", res.data.files_list);
          this.$emit("update:searchFolder", res.data.folder_list);
          // 获取当前父目录的ParentId，就是LastParentId
          this.$axios({
            url: "/file/parentid/get",
            method: "post",
            data: {
              id: this.$store.state.currentParentId,
            },
          }).then((res) => {
            console.log(res.data);
            if (res.data.success) {
              this.$store.commit("updateLastParentId", res.data.parent_id);
            } else {
              this.$message.error(res.data.msg);
            }
          });
        } else {
          this.$message.error(res.data.msg);
        }
      });
    },
    // 返回根目录
    goRootFolder() {
      console.log("回到根目录");
      this.$store.commit("updateCurrentParentId", 0);
      this.$store.commit("updateLastParentId", 0);
      this.$axios.defaults.headers.common["Authorization"] =
        window.localStorage.getItem("token");
      this.$axios({
        url: "/user/file/list",
        method: "post",
        data: {
          id: this.$store.state.currentParentId,
        },
      }).then((res) => {
        console.log(res.data);
        if (res.data.success) {
          this.$emit("update:listData", res.data.files_list);
          this.$emit("update:searchFolder", res.data.folder_list);
        } else {
          this.$message.error("文件请求失败");
        }
      });
    },
    // 点击切换展示类型
    changeShowType() {
      if (this.$store.state.showType == "icon") {
        // this.showType = "table";
        // this.$emit("changShowType", "table");
        this.$store.commit("updateShowType", "table");
      } else {
        // this.showType = "icon";
        // this.$emit("changShowType", "icon");
        this.$store.commit("updateShowType", "icon");
      }
    },
    fileToMD5(file) {
        const fileReader = new FileReader()
        fileReader.readAsBinaryString(file);
         fileReader.onload = (e) => {
          this.fileHash = SparkMD5.hashBinary(e.target.result);
       }
    },
  },

  watch: {

    // 监听是否正在创建文件夹
    "$store.state.isCreateFolder"(current) {
      if (current == false) {
        this.isCreateAble = true;
      }
    },
    "$store.state.selectFiles"(current) {
      if (current.length > 0) {
        this.isMultBtnsShow = true;
      } else {
        this.isMultBtnsShow = false;
      }
    },
  },
  created() {},
};
</script>

<style scoped>
.FunctionBar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: calc(100vw - 280px);
  min-width: 980px;
  /* background-color: pink; */
  height: 80px;
  padding: 25px 20px;
  position: relative;
  top: 0;
  left: 0px;
  z-index: 50;
}

.el-input /deep/ input {
  width: 14vw;
  border-radius: 30px;
  height: 35px;
  line-height: 35px;
  background-color: #f1f2f3;
  border: 1px solid #f1f2f3;
}

.el-input /deep/ i {
  line-height: 35px !important;
  color: #687176;
}

.right {
  display: flex;
  align-items: center;
}

.right i {
  font-size: 20px;
  color: #a0a0a0;
  font-weight: bold;
}

.right > div {
  margin-right: 15px;
}

.select {
  background-color: #696bcc !important;
  color: white !important;
  border: 1px solid #696bcc !important;
}

.selectAll:hover {
  background-color: #595bb3;
  color: white;
}

i {
  cursor: pointer;
}

.sortTypeItem {
  font-size: 13px;
  padding: 10px 0 10px 40px;
  position: relative;
  cursor: pointer;
  color: #595bb3;
}

.sortTypeItem i {
  position: absolute;
  left: 13px;
  color: #595bb3;
}

.sortTypeItem:hover {
  background-color: #efeff5;
}

.left {
  display: flex;
}

.uploadButton {
  margin-right: 10px;
  /* padding-bottom: 20px; */
  /* margin-bottom: 20px; */
}

.uploadProgress {
  width: 230px;
}

.multButtons {
  display: flex;
  align-items: center;
  font-size: 14px;
  transform: scale(0.9);
}

.tips {
  color: rgb(177, 177, 177);
}

.multButtonsContainer {
  display: flex;
  align-items: center;
  border: 1px solid #ccc;
  border-radius: 8px;
  color: rgb(97, 97, 97);
  margin-left: 10px;
  overflow: hidden;
}

.multButtonsContainer div {
  padding: 8px 15px;
  border-right: 1px solid #ccc;
  cursor: pointer;
}

.multButtonsContainer div:nth-last-child(1) {
  border: none;
}

.multButtonsContainer div:hover {
  background-color: #595bb3;
  color: white;
}

.goLastFolder {
  cursor: pointer;
  color: #595bb3;
  font-size: 10px;
  margin-left: 20px;
  position: absolute;
  bottom: 0;
  left: 0;
  margin-top: 20px;
}

.goLastFolder a {
  /* margin: 0 5px; */
  /* padding-top: 20px; */
}

.tableHead {
  display: flex;
  line-height: 50px;
  height: 50px;
  width: calc(100vw - 260px);
  padding: 0 25px;
  box-sizing: border-box;
  font-size: 15px;
}

.tableName {
  line-height: 43px;
  width: 50%;
}
.tableHeadName {
  width: calc(50% + 43px);
  padding: 0;
}

.tableItemSize {
  width: 20%;
  padding-left: 80px;
  box-sizing: border-box;
}

.tableItemCreateTime {
  padding-left: 60px;
  box-sizing: border-box;
  width: 25%;
}

.tableCollect {
  width: 10%;
  text-align: center;
}
</style>
