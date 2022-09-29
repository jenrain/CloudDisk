<template>
    <div class="Files">
      <function-bar
        v-bind:listData.sync="listData"
        v-bind:searchFolder.sync="searchFolder"
        @multDelete="$refs.iconTypeList.deleteCurrentFile('mult')"
        @multDownload="$refs.iconTypeList.downloadCurrentFile('mult')"
        @multMove="$refs.iconTypeList.moveFile('mult')"
        @changeSortType="(type) => (sortType = type)"
        @changShowType="(type) => (showType = type)"
        @multCollect="
          (flag) => {
            $refs.iconTypeList.collectCurrentFile(flag, 'mult');
          }
        "
      ></function-bar>
      <icon-type-list
      v-bind:listData.sync="listData"
      v-bind:searchFolder.sync="searchFolder"
      :folderList="folderList"
      :sortType="sortType"
      :showType="showType"
      ref="iconTypeList"
      >
      </icon-type-list>
    </div>
  </template>
  
  <script>
  import FunctionBar from "../components/FunctionBar.vue";
  import IconTypeList from '../components/IconTypeList.vue';
  export default {
    name: "Files",
    data() {
      return {
        listData: [],
        folderList: {},
        // 是否已经被创建
        // isCreated: false,
        // 排序方式
        sortType: "time",
        // 展示方式 icon table
        showType: "icon",
        // 搜索的文件夹
        searchFolder: [],
      };
    },
    components: {
      FunctionBar,
      IconTypeList,
    },
    methods: {
    },
    async created() {
        this.$store.commit("updateCurrentParentId", 0);
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
            this.listData = res.data.files_list;
            this.searchFolder = res.data.folder_list;
          } else {
            this.$message.error("文件请求失败");
          }
        });
    },
    watch: {},
  };
  </script>
  
  <style scoped>
  .Files {
    width: 100%;
  }
  </style>
  