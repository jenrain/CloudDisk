<template>
  <el-dialog
    custom-class="shareDialog"
    :title="'分享文件: ' + shareItem.name"
    :visible.sync="isShow"
    @close="
      () => {
        $emit('closeDialog');
        showItem = {};
      }
    "
  >
    <input class="urlInput" v-model="showItem" readonly />
    <el-button
      size="small"
      type="primary"
      data-clipboard-target=".urlInput"
      class="btn"
      >复制链接</el-button
    >
    <div slot="footer" class="dialog-footer">
      tips: 请复制上述链接以分享文件
    </div>
  </el-dialog>
</template>

<script>
// 剪切板插件
import ClipboardJS from "clipboard";
new ClipboardJS(".btn");

export default {
  name: "shareDialog",
  data() {
    return {
      // 是否显示分享框
      isShow: this.isShareDialogShow,
      // 用于展示的数据
      showItem: "",
    };
  },
  props: {
    // 是否显示分享框
    isShareDialogShow: {
      type: Boolean,
      default() {
        return false;
      },
    },
    // 分享的文件对象
    shareItem: {
      type: Object,
      default() {
        return {};
      },
    },
  },
  methods: {
  },
  watch: {
    async isShareDialogShow(current) {
      this.isShow = current;
      this.showItem = this.shareItem.path;
    },
  },
  computed: {},
};
</script>

<style scoped>
.el-button {
  padding: 10px 60px;
  border-radius: 20px;
  margin-top: 20px;
}

.urlInput {
  display: block;
  width: 90%;
  height: 35px;
  border-radius: 7px;
  outline: none;
  border: 1px solid #ccc;
  color: rgb(87, 87, 87);
  padding: 0 10px;
}
</style>
