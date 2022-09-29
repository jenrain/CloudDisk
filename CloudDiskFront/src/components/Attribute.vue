<template>
  <el-dialog
    title="文件属性"
    :visible.sync="isAttributeShow"
    width="26%"
    :before-close="handleClose"
  >
    <el-form>
      <el-form-item style="margin-top: -20px;" label="文件名: " :label-width="formLabelWidth">
        {{ fileInfo.name }}
      </el-form-item>
      <el-form-item style="margin-top: -20px;" label="文件类型: " :label-width="formLabelWidth">
        {{ fileInfo.filetype }}
      </el-form-item>
      <el-form-item style="margin-top: -20px;" label="文件大小: " :label-width="formLabelWidth">
        {{ (fileInfo.size / 1048576).toFixed(2) + "MB"}}
      </el-form-item>
      <el-form-item style="margin-top: -20px;" label="创建时间: " :label-width="formLabelWidth">
        {{ fileInfo.create_time }}
      </el-form-item>
      <el-form-item style="margin-top: -20px;" label="修改时间: " :label-width="formLabelWidth">
        {{ fileInfo.update_time }}
      </el-form-item>
    </el-form>
  </el-dialog>
</template>

<script>
import { handleLongString } from "../plugins/utils.js";

export default {
  name: "Attribute",
  props: {
    // 当前查看属性的文件的id
    fileId: {
      type: Number,
      default() {
        return 0;
      },
    },
    attributeTop: {
      type: Number,
      default() {
        return 0;
      },
    },
    attributeLeft: {
      type: Number,
      default() {
        return 0;
      },
    },
    isAttributeShow: {
      type: Boolean,
      default() {
        return false;
      },
    },
  },
  data() {
    return {
      isAttributeShow: this.isAttributeShow,
      fileInfo: null,
      fileId: this.fileId,
    };
  },
  methods: {
    //   请求
    // 根据id查看当前文件的属性
    async getFileAttribute() {
      this.$axios.defaults.headers.common["Authorization"] = window.localStorage.getItem("token");
        this.$axios({
          url: "/file/detail/get",
          method: "post",
          data: {
            id: this.fileId,
          }
        }).then((res) => {
          console.log(res.data);
          if (res.data.success) {
            this.fileInfo = res.data
          } else {
            this.$message.error(res.data.msg);
          }
        })
    },
    handleClose() {
      this.isAttributeShow = false;
      this.fileInfo = "";
      this.$emit("closeAttribute");
    }
  },
  created() {},
  watch: {
    // fileId() {
    //   this.getFileAttribute();
    // },
  },
  filters: {
    handleLongString,
  },
};
</script>
