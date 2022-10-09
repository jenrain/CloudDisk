<template>
  <div
    class="iconTypeListContainer"
    :class="[
      banScroll ? 'banScroll' : '',
      $store.state.showType == 'icon' ? '' : 'tableTypeListContainer',
    ]"
  >
    <div class="IconTypeList" v-if="$store.state.showType == 'icon'">
      <!-- 文件夹 -->
      <div
        class="listItem folderItem"
        v-for="(item, index) in currentChildrenFolder"
        @dblclick="openCurrentFolder(item)"
        @contextmenu.prevent="showMenu($event, item, 'folderMenu')"
        @dragstart="(e) => e.preventDefault()"
        @dragover="(e) => e.preventDefault()"
        @drop="dropItem(item)"
        :key="index"
      >
        <div class="imgContainer">
          <img src="../assets/img/folder.png" alt="" :draggable="false" />
        </div>
        <input
          type="text"
          v-model="renameInput"
          class="renameInput"
          @blur="
            rightClickFolderItem.id
              ? renameDone(rightClickFolderItem, index, 'folder')
              : ''
          "
          @keyup.enter="
            rightClickFolderItem.id
              ? renameDone(rightClickFolderItem, index, 'folder')
              : ''
          "
          v-if="isFolderRenameInputShow && rightClickFolderItem.id == item.id"
        />
        <div class="name" v-else>
          {{ item.name == null ? "NoNameFolder" : item.name }}
        </div>
      </div>

      <!-- 新建文件夹的模板 -->
      <div class="listItem createItem" v-if="isCreateFolderShow">
        <img src="../assets/img/folder.png" alt="" />
        <input
          type="text"
          v-model="createdName"
          class="renameInput createNameInput"
          @blur="createNameDone"
          @keyup.enter="createNameDone"
        />
      </div>

      <!-- 文件列表 -->
      <div
        class="listItem"
        :class="selectFiles.find((i) => i.id == item.id) ? 'selectFile' : ''"
        v-for="(item, index) in listData"
        :key="item.id"
        @click="selectCurrentItem(item)"
        @dblclick="openCurrentFile(item)"
        @contextmenu.prevent="showMenu($event, item)"
        :draggable="rightClickMenuType == 'files'"
        @dragend="onDragEndItem"
        slot="reference"
      >
        <div class="imgContainer">
          <img
            src="../assets/img/music.png"
            alt=""
            v-if="item.filetype == 'audio'"
            :draggable="false"
          />
          <img
            :src="item.path"
            alt=""
            v-else-if="item.filetype == 'image'"
            :draggable="false"
          />

          <!-- 视频文件 -->
          <div class="videoContainer" v-else-if="item.ext == '.mp4'">
            <video :src="item.path" class="videoItem" preload="meta"></video>
            <div class="iconContainer">
              <i class="iconfont icon-play_nor"></i>
            </div>
          </div>
          <img
            :src="require(`../assets/img/${computeType(item.filetype)}.png`)"
            alt=""
            v-else
            :draggable="false"
          />
          <img
            src="../assets/img/collect.png"
            alt=""
            class="collectIcon"
            v-if="item.collection"
            :draggable="false"
          />
        </div>
        <input
          type="text"
          v-model="renameInput"
          class="renameInput"
          @blur="rightClickItem.id ? renameDone(rightClickItem, index) : ''"
          @keyup.enter="
            rightClickItem.id ? renameDone(rightClickItem, index) : ''
          "
          v-if="isRenameShow && rightClickItem.id == item.id"
        />
        <div class="name" v-else>
          {{ item.name == null ? "NoNameFile" : item.name }}
        </div>
      </div>
    </div>

    <div class="tableTypeList" v-else>
      <!-- 文件夹 -->
      <div
        class="tableListItem"
        v-for="(item, index) in currentChildrenFolder"
        @dblclick="openCurrentFolder(item)"
        @contextmenu.prevent="showMenu($event, item, 'folderMenu')"
        @dragover="(e) => e.preventDefault()"
        @drop="dropItem(item)"
        :key="index"
      >
        <div class="tableImgContainer">
          <img src="../assets/img/folder.png" alt="" :draggable="false" />
        </div>
        <div class="tableName">
          {{
            item.name == null
              ? "NoNameFolder"
              : item.name.substr(0, item.name.length - 1)
          }}
        </div>
        <div class="tableCollect">-</div>
        <div class="tableItemSize">-</div>
        <div class="tableItemCreateTime">-</div>
      </div>

      <!-- 新建文件夹的模板 -->
      <div class="tableListItem" v-if="isCreateFolderShow">
        <img src="../assets/img/folder.png" alt="" />
        <div class="tableCreateInputContainer">
          <input
            type="text"
            v-model="createdName"
            class="tableCreateInput"
            @blur="createNameDone"
            @keyup.enter="createNameDone"
          />currentFolder
        </div>
      </div>

      <!-- 文件列表 -->
      <div
        class="tableListItem"
        :class="selectFiles.find((i) => i.id == item.id) ? 'selectFile' : ''"
        v-for="(item, index) in listData"
        :key="item.id"
        @click="selectCurrentItem(item)"
        @dblclick="openCurrentFile(item)"
        @contextmenu.prevent="showMenu($event, item)"
        :draggable="rightClickMenuType == 'files'"
        @dragend="onDragEndItem"
        slot="reference"
      >
        <div class="tableImgContainer">
          <img
            src="../assets/img/music.png"
            alt=""
            v-if="item.filetype == 'audio'"
            :draggable="false"
          />
          <el-image
            :src="item.url"
            alt=""
            v-else-if="item.filetype == 'image'"
            :draggable="false"
            fit="cover"
          />

          <div class="tableVideoContainer" v-else-if="item.ext == '.mp4'">
            <video
              :src="item.path"
              class="tableVideoItem"
              preload="meta"
            ></video>
          </div>
          <img
            :src="require(`../assets/img/${computeType(item.type)}.png`)"
            alt=""
            v-else
            :draggable="false"
          />
        </div>
        <div
          class="tableRenameInputContainer"
          v-if="isRenameShow && rightClickItem.id == item.id"
        >
          <input
            type="text"
            v-model="renameInput"
            class="tableRenameInput"
            @blur="rightClickItem.id ? renameDone(rightClickItem, index) : ''"
            @keyup.enter="
              rightClickItem.id ? renameDone(rightClickItem, index) : ''
            "
          />
        </div>

        <div class="tableName" v-else>
          {{ item.name == null ? "NoNameFile" : item.name + "." + item.type }}
        </div>
        <div class="tableCollect">
          <img src="../assets/img/collect.png" alt="" v-if="item.collection" />
        </div>
        <div class="tableItemSize">
          {{
            item.size == null
              ? "未知大小"
              : (item.size / 1048576).toFixed(2) + " MB"
          }}
        </div>
        <div class="tableItemCreateTime">
          {{ item.gmtCreate.substr(0, 16) }}
        </div>
      </div>
    </div>

    <!-- 文件右击菜单框组件 -->
    <right-click-menu
      :menuType="rightClickMenuType == 'collect' ? 'collect' : 'file'"
      :isMenuShow="isMenuShow"
      :menuTop="menuTop"
      :menuLeft="menuLeft"
      :isCurrentFileCollected="isCurrentFileCollected"
      @rename="rename"
      @openCurrentFile="openCurrentFile(rightClickItem)"
      @downloadCurrentFile="downloadCurrentFile('current', rightClickItem)"
      @deleteCurrentFile="deleteCurrentFile('current', rightClickItem)"
      @shareCurrentFile="shareCurrentFile(rightClickItem)"
      @showAttribute="showAttribute"
      @moveFile="moveFile"
    ></right-click-menu>

    <!-- 文件夹右击菜单框组件 -->
    <right-click-menu
      :menuType="'folder'"
      :isFolderMenuShow="isFolderMenuShow"
      :menuTop="menuTop"
      :menuLeft="menuLeft"
      @openCurrentFolder="openCurrentFolder(rightClickFolderItem)"
      @renameCurrentFolder="rename('folder')"
      @deleteCurrentFolder="deleteCurrentFolder(rightClickFolderItem)"
    ></right-click-menu>

    <!-- 图片预览组件 -->
    <image-player
      :currentImg="currentImg"
      :isImgDialogShow="isImgDialogShow"
      @closeDialog="isImgDialogShow = false"
      @shareCurrentFile="shareCurrentFile"
    ></image-player>

    <MusicPlayer> </MusicPlayer>

    <VideoPlayer> </VideoPlayer>

    <!-- 属性框组件 -->
    <Attribute
      :isAttributeShow="isAttributeShow"
      :fileId="rightClickItem.id"
      :attributeTop="attributeTop"
      :attributeLeft="attributeLeft"
      @closeAttribute="isAttributeShow = false"
      ref="attributeDialog"
    >
    </Attribute>

    <!-- 移动到中选择文件夹的dialog -->
    <folder-dialog
      :isFolderDialogShow="isFolderDialogShow"
      :moveType="moveType"
      :currentId="rightClickItem.id"
      @confirmMove="
        (info) => confirmMove(info.selectFolderId, info.currentFileId)
      "
      @closeFolderDialog="isFolderDialogShow = false"
    ></folder-dialog>

    <!-- 搜索为空的提醒 -->
    <div
      class="searchTips"
      v-if="
        $route.params.path &&
        $route.params.path.search('search') != -1 &&
        searchFolder.length == 0 &&
        this.listData.length == 0
      "
    >
      没有找到相应内容哦!
    </div>
    <!-- 返回顶部按钮 -->
    <go-top scrollObj=".iconTypeListContainer"></go-top>
    <!-- 拖动时的预览容器 -->
    <div
      class="dragImgContainer"
      :class="showDragImgContainer ? 'showDragImgContainer' : ''"
      :style="[
        { top: dragImgContainerPosition.y + 'px' },
        { left: dragImgContainerPosition.x + 'px' },
      ]"
    ></div>

    <!-- 分享框 -->
    <share-dialog
      :isShareDialogShow="isShareDialogShow"
      :shareItem="shareItem"
      @closeDialog="isShareDialogShow = false"
    ></share-dialog>
  </div>
</template>

<script>
let isClickSelectAll = true;

import { getTypeIcon } from "../plugins/utils.js";

import ImagePlayer from "./ImagePlayer.vue";
import RightClickMenu from "./RightClickMenu.vue";
import FolderDialog from "./FolderDialog.vue";
import GoTop from "./GoTop.vue";
import ShareDialog from "./ShareDialog.vue";
import MusicPlayer from "./MusicPlayer.vue";
import VideoPlayer from "./VideoPlayer.vue";
import App from "@/App.vue";
import Attribute from "./Attribute.vue";

export default {
  components: {
    ImagePlayer,
    RightClickMenu,
    FolderDialog,
    GoTop,
    ShareDialog,
    MusicPlayer,
    VideoPlayer,
    App,
    Attribute,
  },
  name: "IconTypeList",
  data() {
    return {
      selectFiles: [],
      isMenuShow: false,
      menuTop: 0,
      menuLeft: 0,
      cardoffsetTop: 0,
      cardoffsetLeft: 0,
      // 重命名输入框
      renameInput: "",
      // 是否显示重命名输入框
      isRenameShow: false,
      // 右击的item
      rightClickItem: {},
      // 新建文件夹名称
      createdName: "",
      isCreateFolderShow: false,
      // 是否显示图片diaload
      isImgDialogShow: false,
      //   当前要打开diaload的图片
      currentImg: {},
      // 当前目录的子目录
      currentChildrenFolder: [],
      //   是否显示文件右键菜单
      isFolderMenuShow: false,
      // 下载的文件信息
      downloadFileInfo: {
        url: "",
        name: "",
      },
      // 视频列表
      videoList: [],
      // 是否显示属性框
      isAttributeShow: false,
      // 页面可见区域宽
      pageWidth: document.body.clientWidth,
      // 页面可见区域高
      pageHeight: document.body.clientHeight,
      // 属性dialog的位置
      attributeTop: 0,
      attributeLeft: 0,
      // 是否禁止页面滚动
      banScroll: false,
      // 是否显示文件选择框
      isFolderDialogShow: false,
      // 当前文件是否已收藏
      isCurrentFileCollected: false,
      // 文件移动类型  'mult'和'current'
      moveType: "current",
      // 右键选中的文件夹item
      rightClickFolderItem: {},
      // 当前文件夹的id
      currentFolderId: 0,
      // 是否显示文件夹重命名输入框
      isFolderRenameInputShow: false,
      // 被拖动的itemList
      dragItemList: [],
      // 是否显示dragImgContainer
      showDragImgContainer: false,
      // dragImgContainer的位置
      dragImgContainerPosition: {
        x: 0,
        y: 0,
      },
      // 是否显示分享框
      isShareDialogShow: false,
      // 分享的文件
      shareItem: {},
    };
  },
  props: {
    // type name
    listData: {
      type: Array,
      default() {
        return [];
      },
    },
    // 右键菜单类型
    rightClickMenuType: {
      type: String,
      default() {
        return "files";
      },
    },
    // item展示方式
    showType: {
      type: String,
      default() {
        return "icon";
      },
    },
    searchFolder: {
      type: Array,
      default() {
        return [];
      },
    },
  },

  methods: {
    // 单击item的回调
    selectCurrentItem(item) {
      // 操作dom  直接操作dom可以减少循环，提高性能
      // let listItem = document.querySelectorAll(".listItem");
      // 先判断该选项是否已经被选中
      let i = this.selectFiles.findIndex((item1) => item1.id == item.id);
      if (i == -1) {
        this.selectFiles.push(item);
      } else {
        this.selectFiles.splice(i, 1);
      }
    },

    // 右键单击item的事件
    showMenu(e, item, type) {
      // 如果属性框打开，则先关闭属性框
      if (this.isAttributeShow == true) {
        this.isAttributeShow = false;
      }
      this.banScroll = true;
      console.log(e, item);
      // console.log(this.cardoffsetLeft);
      // 获取菜单的高度
      let menu = document.querySelector(".RightClickMenu");
      let menuHeight = menu.offsetTop;
      console.log([menu], menuHeight);
      // 计算菜单dialog 的位置
      // files的菜单高度和collect不一样
      if (this.rightClickMenuType == "files") {
        this.menuTop =
          e.pageY + 250 > this.pageHeight ? this.pageHeight - 250 : e.pageY;
      } else {
        this.menuTop =
          e.pageY + 220 > this.pageHeight ? this.pageHeight - 220 : e.pageY;
      }
      this.menuLeft =
        e.pageX + 140 > this.pageWidth ? this.pageWidth - 140 : e.pageX;

      // 计算属性dialog的位置
      this.attributeTop =
        e.pageY + 230 > this.pageHeight ? this.pageHeight - 230 : e.pageY;
      this.attributeLeft =
        e.pageX + 340 > this.pageWidth ? this.pageWidth - 340 : e.pageX;

      if (!type || type == "menu") {
        // 判断右击文件是否已收藏
        this.isCurrentFileCollected = item.collection == 1 ? true : false;
        this.rightClickItem = item;
        // this.$store.commit("updateRightClickItem", item);
        this.isMenuShow = true;
        this.isFolderMenuShow = false;
      } else if (type == "folderMenu") {
        this.isFolderMenuShow = true;
        this.isMenuShow = false;
        // this.$store.commit("updateRightClickFolderItem", item);
        this.rightClickFolderItem = item;
      }
    },

    // 点击重命名的回调
    rename(type) {
      // 获取点击重命名的索引
      // this.rightClickItem = this.$store.state.rightClickItem;
      if (!type) {
        // 文件
        this.isRenameShow = true;
        this.renameInput = this.rightClickItem.name;
      } else if (type == "folder") {
        // 文件夹
        this.isFolderRenameInputShow = true;
        this.renameInput = this.rightClickFolderItem.name;
      }
      //在input的属性中添加autofocus只能触发一次 这里改用操作DOM
      this.$nextTick(() => {
        if (this.$store.state.showType == "icon") {
          document.querySelector(".renameInput").focus();
        } else if (this.$store.state.showType == "table") {
          document.querySelector(".tableRenameInput").focus();
        }
      });
      //   console.log(this.rightClickIndex, 123456);
    },

    // 重命名完成后的回调  失去焦点或者回车
    async renameDone(item, index, type) {
      console.log(item);
      // 判断输入内容是否为空
      if (this.renameInput.trim().length == 0) {
        if (!type) {
          this.$message.warning("文件名不能为空哦!");
          this.isRenameShow = false;
        } else {
          this.$message.warning("文件夹名称不能为空哦!");
          this.isFolderRenameInputShow = false;
        }
      }
      // console.log("文件类型：", type);
      if (!type) {
        if (this.rightClickItem.name != this.renameInput.trim()) {
          this.$axios.defaults.headers.common["Authorization"] =
            window.localStorage.getItem("token");
          this.$axios({
            url: "/user/file/update",
            method: "post",
            data: {
              identity: item.identity,
              name: this.renameInput.trim(),
            },
          }).then((res) => {
            console.log(res.data);
            if (res.data.success) {
              this.listData[index].name = res.data.name;
              this.$emit("update:listData", this.listData);
              this.$message.success("重命名成功!");
            } else {
              this.$message.error("重命名失败,请稍后重试!");
            }
          });
        }
        this.rightClickItem = {};
        this.isRenameShow = false;
      } else {
        // 文件夹
        console.log(this.renameInput);
        if (this.rightClickFolderItem.name != this.renameInput.trim()) {
          this.$axios.defaults.headers.common["Authorization"] =
            window.localStorage.getItem("token");
          this.$axios({
            url: "/user/file/update",
            method: "post",
            data: {
              identity: item.identity,
              name: this.renameInput.trim(),
            },
          }).then((res) => {
            if (res.data.success) {
              this.searchFolder[index].name = res.data.name;
              this.$emit("update:searchFolder", this.searchFolder);
              console.log("修改后的值：", this.renameInput.trim());
              this.$message.success("重命名成功!");
            } else {
              this.$message.error("重命名失败,请稍后重试!");
            }
          });
        }
        this.rightClickFolderItem = {};
        this.isFolderRenameInputShow = false;
      }
      // this.$store.commit("updateRightClickItem", {});
      // 清空renameInput
      this.renameInput = "";
      // }
    },

    // 新建文件夹命名完成的回调
    createNameDone() {
      // 避免 enter和 blur冲突调用两次此函数
      this.$nextTick(async () => {
        let input;
        if (this.$store.state.showType == "icon") {
          input = document.querySelector(".createNameInput");
        } else if (this.$store.state.showType == "table") {
          input = document.querySelector(".tableCreateInput");
        }
        if (input) {
          if (this.createdName == "") {
            this.$message.warning("文件夹名称不能为空哦!");
            this.$store.commit("updateIsCreateFolder", false);
          } else {
            this.$axios.defaults.headers.common["Authorization"] =
              window.localStorage.getItem("token");
            this.$axios({
              url: "/user/folder/create",
              method: "post",
              data: {
                parent_id: this.$store.state.currentParentId,
                name: this.createdName.trim(),
              },
            }).then((res) => {
              console.log(res.data);
              if (res.data.success) {
                this.$store.commit("updateIsCreateFolder", false);
                this.isCreateFolderShow = false;
                this.$message.success("创建成功!");
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
                    this.$emit("update:searchFolder", res.data.folder_list);
                  } else {
                    this.$message.error(res.data.msg);
                  }
                });
              } else {
                this.$message.error(res.data.msg);
                this.$store.commit("updateIsCreateFolder", false);
                this.isCreateFolderShow = false;
                this.$message.success("创建失败!");
              }
            });
            this.createdName = "";
          }
        }
      });
    },

    // 双击打开文件
    openCurrentFile(item) {
      // 先判断打开的文件类型
      console.log("isVideoPlayerShow: ", this.$store.state.isVideoPlayerShow);
      let filetype = item.filetype;
      console.log("filetype: ", filetype);
      console.log("item: ", item);
      // 打开的是video文件
      if (filetype == "video") {
        this.$store.commit("updateIsVideoPlayerShow", true);
        this.$store.commit("updateCurrentVideoInfo", item);
      }
      // 打开的是音频文件
      else if (filetype == "audio") {
        this.$store.commit("updateIsMusicPlayerShow", true);
        this.$store.commit("updateCurrentMusicInfo", item);
      }
      // 打开的是图片文件
      else if (filetype == "image") {
        this.currentImg = item;
        this.isImgDialogShow = true;
      }
      // 其它文件暂时无法直接打开 提醒暂时无法直接打开，可以下载后打开
      else {
        this.$message.warning("该文件暂时不能直接打开哦,可以下载后在本地打开!");
      }
    },

    // 删除文件
    async deleteCurrentFile(type, item) {
      let res;
      let arr = [];
      if (type == "current") {
        arr.push(item.id);
      } else if (type == "mult") {
        this.selectFiles.forEach((item) => {
          arr.push(item.id);
        });
      }
      console.log("identity: ", item.identity);
      this.$axios.defaults.headers.common["Authorization"] =
        window.localStorage.getItem("token");
      this.$axios({
        url: "/user/file/delete",
        method: "delete",
        data: {
          identity: item.identity,
        },
      }).then((res) => {
        console.log(res.data);
        if (res.data.success) {
          let index = this.listData.findIndex(
            (item) => item.id == this.rightClickItem.id
          );
          this.listData.splice(index, 1);
          this.$message.success("删除成功!");
        } else {
          this.$message.error("删除失败,请稍后重试!");
        }
      });
    },

    // 打开当前双击的文件夹
    openCurrentFolder(item) {
      this.$store.commit(
        "updateLastParentId",
        this.$store.state.currentParentId
      );
      this.$store.commit("updateCurrentParentId", item.id);
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
          // this.listData = res.data.files_list
          // this.searchFolder = res.data.folder_list
          this.$emit("update:listData", res.data.files_list);
          this.$emit("update:searchFolder", res.data.folder_list);
        } else {
          this.$message.error("文件请求失败");
        }
      });
    },

    // 点击下载文件的回调
    downloadCurrentFile(type, item) {
      window.open(item.path + "?attname=" + item.name);
    },

    // 移动文件
    async confirmMove(selectFolderId, currentFileId, list) {
      // 把parentid传过来
      console.log("selectFolderId: ", selectFolderId);
      console.log("currentFileId: ", currentFileId);
      // 先判断是否是当前文件夹 是的话直接return
      if (selectFolderId == this.$store.state.currentParentId) {
        return;
      }
      this.$axios.defaults.headers.common["Authorization"] =
        window.localStorage.getItem("token");
      this.$axios({
        url: "/user/file/move",
        method: "put",
        data: {
          id: currentFileId,
          parent_id: selectFolderId,
        },
      }).then((res) => {
        console.log(res.data);
        if (res.data.success) {
          this.$message.success("移动成功！");
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
            } else {
              this.$message.error(res.data.msg);
            }
          });
        } else {
          this.$message.error(res.data.msg);
        }
      });
    },

    // 分享当前文件
    shareCurrentFile(item) {
      this.shareItem = item;
      this.isShareDialogShow = true;
    },

    // 点击移动文件的回调
    moveFile(type) {
      this.isFolderDialogShow = true;
      if (!type || type == "current") {
        this.moveType = "current";
      } else {
        this.moveType = "mult";
      }
    },

    // 删除当前文件夹
    async deleteCurrentFolder(item) {
      this.$axios.defaults.headers.common["Authorization"] =
        window.localStorage.getItem("token");
      this.$axios({
        url: "/user/file/delete",
        method: "delete",
        data: {
          identity: item.identity,
        },
      }).then((res) => {
        console.log(res.data);
        if (res.data.success) {
          let index = this.listData.findIndex(
            (item) => item.id == this.rightClickItem.id
          );
          this.searchFolder.splice(index, 1);
          this.$message.success("删除成功!");
          // this.$emit("getFolderList");
        } else {
          this.$message.error("删除失败,请稍后重试!");
        }
      });
    },

    // 重命名当前文件夹
    renameCurrentFolder(item) {
      this.isFolderRenameShow = true;
    },
    showAttribute() {
      this.isAttributeShow = true;
      this.$nextTick(() => {
        this.$refs.attributeDialog.getFileAttribute();
      });
    },
    onDragEndItem() {
      document.ondragover = null;
      this.showDragImgContainer = false;
      let imgContainer = document.querySelector(".dragImgContainer");
      imgContainer.innerHTML = "";
    },
  },
  computed: {
    computeType() {
      return getTypeIcon;
    },
  },
  watch: {
    selectFiles(current) {
      let listItem = document.querySelectorAll(".listItem");
      if (current.length == listItem.length) {
        this.$store.commit("updateIsSelectAll", true);
      } else {
        if (this.$store.state.isSelectAll == true) {
          isClickSelectAll = false;
          this.$store.commit("updateIsSelectAll", false);
        }
      }
    },

    // 监听是否正在新建文件夹的状态
    "$store.state.isCreateFolder"(current) {
      this.isCreateFolderShow = current;
      if (current == true) {
        if (this.$store.state.showType == "icon") {
          this.$nextTick(() => {
            document.querySelector(".createNameInput").focus();
          });
        } else if (this.$store.state.showType == "table") {
          this.$nextTick(() => {
            document.querySelector(".tableCreateInput").focus();
          });
        }
      }
    },

    // 监听选中文件的变化
    selectFiles(current) {
      this.$store.commit("updateSelectFiles", current);
    },

    // 监听是否正在获取文件夹
    "$store.state.isGetingFolder"(current) {
      if (current == false) {
        this.createdName = "";
        this.$store.commit("updateIsCreateFolder", false);
      }
    },

    // 监听排序方式
    "$store.state.sortType"(current) {
      if (current == "time") {
        this.listData.sort((a, b) => {
          return Date.parse(a.gmtCreate) - Date.parse(b.gmtCreate);
        });
      } else if (current == "size") {
        this.listData.sort((a, b) => {
          return a.size - b.size;
        });
      }
    },
    searchFolder(current) {
      this.currentChildrenFolder = current;
    },
  },
  created() {},
  mounted() {
    // 获取组件的offset
    let Card = document.querySelector(".iconTypeListContainer");
    // console.log(Card);
    this.cardoffsetTop = Card.offsetTop;
    this.cardoffsetLeft = Card.offsetLeft;

    // 监听页面窗口大小变化
    window.addEventListener("resize", (e) => {
      // console.log(e);
      this.pageWidth = document.body.clientWidth;
      this.pageHeight = document.body.clientHeight;
    });

    // 监听所有点击事件
    document.onclick = (e) => {
      // console.log(e);
      if (this.isMenuShow == true || this.isFolderMenuShow == true) {
        this.isMenuShow = false;
        this.isFolderMenuShow = false;
        this.banScroll = false;
      }
      if (this.$store.state.isUserInfoCardMenuShow) {
        this.$store.commit("updateIsUserInfoCardMenuShow", false);
      }
    };
  },
};
</script>

<style scoped>
.iconTypeListContainer {
  height: calc(100vh - 100px);
  width: calc(100vw - 240px);
  overflow-y: scroll;
  position: relative;
}

.tableTypeListContainer {
  height: calc(100vh - 150px);
}

.banScroll {
  overflow-y: hidden;
}

.IconTypeList {
  display: flex;
  flex-wrap: wrap;
  position: relative;
  padding: 10px 20px 0;
}

.listItem {
  margin: 0 5px 10px;
  padding: 15px 5px 10px;
  text-align: center;
  width: 100px;
  font-size: 14px;
  border-radius: 10px;
  cursor: pointer;
  position: relative;
  height: 118px;
  display: flex;
  flex-direction: column;
  align-items: center;
  /* box-sizing: border-box; */
}

.collectIcon {
  position: absolute;
  width: 15px;
  height: 15px;
  top: 0px;
  right: 7px;
}

.folderItem:hover {
  background-color: unset !important;
}

.name {
  word-break: break-all;
  text-overflow: ellipsis;
  /* overflow: hidden; * 隐藏超出的内容 * */
  white-space: nowrap;
  padding-top: 10px;
  width: 100px;
  height: 18px;
}

.imgContainer {
  height: 90px;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.listItem img {
  max-height: 90px;
  max-width: 90%;
  user-select: none;
}

.listItem:hover {
  background-color: rgba(105, 107, 204, 0.2);
}

.listItem:hover i {
  background-color: rgba(255, 255, 255, 0.5);
}

.selectFile {
  background-color: rgba(105, 107, 204, 0.4) !important;
}

.renameInput {
  width: 92px;
  position: absolute;
  bottom: 10px;
  left: 50%;
  transform: translateX(-50%);
}

.videoContainer {
  height: 90px;
  width: 90px;
  margin-bottom: 2px;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #000;
}

.videoItem {
  height: 90px;
  width: 90px;
  border-radius: 5px;
  overflow: hidden;
}

.iconContainer {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.iconContainer i {
  color: #fcfcfc;
  font-size: 32px;
  padding: 5px;
  background-color: rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.1s;
}

.tableListItem {
  display: flex;
  /* background-color: pink; */
  width: calc(100vw - 260px);
  height: 55px;
  padding: 6px 25px;
  box-sizing: border-box;
  /* align-items: center; */
  border-bottom: 1px solid #eee;
  cursor: pointer;
  font-size: 14px;
  color: rgb(66, 66, 66);
}

.tableListItem > div {
  height: 43px;
  line-height: 43px;
  box-sizing: border-box;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tableListItem:hover {
  background-color: #e1e1f5;
}

.tableImgContainer {
  width: 48px;
}

.tableImgContainer img,
.tableImgContainer .el-image {
  height: 43px;
  width: 43px;
  border-radius: 5px;
}

.tableVideoContainer {
  height: 43px;
  width: 43px;
  background-color: black;
  border-radius: 5px;
}

.tableVideoItem {
  height: 43px;
  width: 43px;
}

.tableName {
  padding-left: 13px;
  line-height: 43px;
  width: 50%;
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

.tableCollect img {
  height: 15px;
}

.tableCreateInput {
  margin-left: 15px;
  height: 23px;
}

.tableRenameInput {
  height: 23px;
  margin-left: 10px;
}

.searchTips {
  position: absolute;
  top: 30%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 22px;
  color: rgb(63, 63, 63);
}
</style>
