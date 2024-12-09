<template>
  <div class="box box-solid box-body">
    <div class="error-main">
      <el-result icon="warning" title="404 not found" sub-title="Sorry~页面飞到火星去啦~">
        <template #extra>
          <el-link type="primary" :underline="false" href="/">回到首页</el-link>
        </template>
      </el-result>
    </div>
  </div>
</template>