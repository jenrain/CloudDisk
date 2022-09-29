<template>
  <div class="FolderDialog">
    <el-dialog
      title="移动到"
      width="300px"
      :visible.sync="isFolderDialogShow"
      @close="$emit('closeFolderDialog')"
    >
      <el-tree
        highlight-current
        :data="folderList"
        :props="defaultProps"
        accordion
        @node-click="handleNodeClick"
      ></el-tree>
      <span slot="footer" class="dialog-footer">
        <el-button
          @click="$emit('closeFolderDialog')"
          size="mini"
          class="cancel"
          >取 消</el-button
        >
        <el-button @click="confirmMove" size="mini" class="confirm"
          >确 定</el-button
        >
      </span>
    </el-dialog>
  </div>
</template>

<script>
import IconTypeList from "./IconTypeList.vue";

export default {
  name: "FolderDialog",
  components: {
    IconTypeList,
  },
  props: {
    isFolderDialogShow: {
      type: Boolean,
      default() {
        return false;
      },
    },
    moveType: {
      type: String,
      default() {
        return "current";
      },
    },
    currentId: {
      type: Number,
      default() {
        return 0;
      },
    }
  },
  data() {
    return {
      currentId: this.currentId,
      folderList: [],
      // 显示规则
      showProps: {
        children: "childrenList",
        label: (data) => {
          return data.name.slice(0, data.name.length - 1);
        },
      },

      // 选中的文件夹的名称
      selectFolder: 0,
      
    };
  },
  methods: {
    handleFolderList() {
        console.log("currentId: ", this.currentId);
        this.$axios.defaults.headers.common["Authorization"] = window.localStorage.getItem("token");
        this.$axios({
          url: "/user/folder/list",
          method: "post",
          data: {
          }
        }).then((res) => {
          console.log(res.data);
          if (res.data.success) {
            // 文件树的最大深度为三层
            this.folderList = new Array();
            for (let i = 0; i < res.data.folder_list.length; i++) {
              if (res.data.folder_list[i].parent_id == 0) {
                this.folderList.push({
                  id: res.data.folder_list[i].id,
                  parent_id: res.data.folder_list[i].parent_id,
                  label: res.data.folder_list[i].name + " ",
                  children: [],
                })
                for (let j = i + 1; j < res.data.folder_list.length; j++) {
                  if (res.data.folder_list[j].parent_id == res.data.folder_list[i].id) {
                    this.folderList[i].children.push({
                        id: res.data.folder_list[j].id,
                        parent_id: res.data.folder_list[j].parent_id,
                        label: res.data.folder_list[j].name + " ",
                        children: [],
                    })
                    for (let k = j + 1; k < res.data.folder_list.length; k++) {
                      if (res.data.folder_list[k].parent_id == res.data.folder_list[j].id) {
                            this.folderList[i].children[j].push({
                            id: res.data.folder_list[k].id,
                            parent_id: res.data.folder_list[k].parent_id,
                            label: res.data.folder_list[k].name + " ",
                            children: [],
                        })
                      }
                    }
                  }
                }
              }
            }
          } else {
            this.$message.error("文件请求失败");
          }
        });
    },

    getFullPath(node, path) {
      if (node != null) {
        path = node.parent.data.name + path;
        return this.getFullPath(node.parent, path);
      } else {
        return path;
      }
    },

    // js递归遍历树形json数据，根据关键字查找节点
    //@leafId  查找的id，
    //@nodes   原始Json数据
    //@path    供递归使用
    findPathByLeafId(leafId, nodes, path) {
      if (path === undefined) {
        path = [];
      }
      for (var i = 0; i < nodes.length; i++) {
        var tmpPath = path.concat();
        tmpPath.push(nodes[i].name);
        if (leafId == nodes[i].name) {
          return tmpPath;
        }
        if (nodes[i].childrenList) {
          var findResult = this.findPathByLeafId(
            leafId,
            nodes[i].childrenList,
            tmpPath
          );
          if (findResult) {
            return findResult;
          }
        }
      }
    },

    handleNodeClick(e, node) {
      // 选中文件的id
      this.selectFolder = e.id;
    },

    // 点击确定的回调
    confirmMove() {
      // 关闭选择框
      this.$emit("closeFolderDialog");
      // 将选中的文件夹名称发给父组件
      this.$emit("confirmMove", {
        selectFolderId: this.selectFolder,
        currentFileId: this.currentId,
      });
    },
  },
  watch: {
    isFolderDialogShow(current) {
      if (current) {
        this.handleFolderList();
      }
    },
  },
  created() {
        
  },
};
</script>

<style scoped>
.FolderDialog /deep/ .el-dialog {
  border-radius: 10px;
  height: 300px;
}

.FolderDialog /deep/ .el-dialog__header {
  border-bottom: 1px solid #eee;
  padding: 10px 20px;
}

.FolderDialog /deep/ .el-dialog__title {
  font-size: 15px;
  color: rgb(95, 95, 95);
}

.FolderDialog /deep/ .el-dialog__body {
  height: 160px;
  /* overflow-y: scroll; */
  padding: 15px;
}

.FolderDialog /deep/ .el-tree-node__content {
  height: 30px;
}

.FolderDialog /deep/ .el-tree-node__content:hover {
  background-color: #f0f0fc;
}

.FolderDialog /deep/ .is-current > .el-tree-node__content {
  background-color: #e6e6f5;
}

.cancel:hover {
  background-color: white;
  color: #696bcc;
  border-color: #a8a9d4;
}

.confirm {
  background-color: #7c7dd1;
  border-color: #7c7dd1;
  color: white;
}

.confirm:hover {
  background-color: #6d6fce;
}

.FolderDialog /deep/ .el-button:focus {
  background-color: #7c7dd1 !important;
  color: white !important;
  border-color: #7c7dd1 !important;
}
</style>
