<template>
	<div class="mytips" :type="type" :style="tipstyle">
    <slot></slot>
	</div>
</template>

<script>
  export default {
    props: ["type","tipstyle"],
  }
</script>

<style>
  .mytips {
    font-size: 13px;
    padding-top: 5px;
    line-height: 20px;
  }

  .mytips[type=info] {
    color: #a5a5a5;
  }

  .mytips[type=danger] {
    color: red;
  }
</style>
