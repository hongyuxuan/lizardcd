<template>
  <el-form ref="refNew" :rules="rules" :model="form" label-width="110px">
    <el-form-item label="集群" prop="cluster">
      <el-select v-model="form.cluster" placeholder="请选择集群" clearable filterable size="large">
        <el-option v-for="(v,k) in clusterList" :key="k" :label="k" :value="k" />
      </el-select>
    </el-form-item>
    <el-form-item label="命名空间" prop="namespace">
      <el-select v-model="form.namespace" placeholder="请选择命名空间" clearable filterable size="large">
          <el-option v-for="(item) in clusterList[form.cluster]" :key="item" :label="item" :value="item" />
        </el-select>
    </el-form-item>
    <el-form-item label="从应用模板导入" prop="templates">
      <el-select 
        v-model="form.templates" 
        placeholder="请选择模板" 
        value-key="id" 
        clearable 
        size="large"
        style="width:100%" 
        @change="selectTemplate">
        <el-option v-for="item in templateList" :key="item.id" :label="item.name" :value="item">
          <span style="float:left">{{item.name}}</span>
        </el-option>
      </el-select>
    </el-form-item>
    <el-form-item label="配置YAML" prop="content">
      <v-ace-editor
        v-model:value="form.content"
        lang="yaml"
        theme="chrome"
        style="width:100%;"
        :options="{
          enableBasicAutocompletion: true,
          enableSnippets: true,
          enableLiveAutocompletion: true,
          tabSize: 2,
          showPrintMargin: false,
          fontSize: 14,
          minLines: 10,
          maxLines: 5000
        }" />
    </el-form-item>
    <el-form-item label="模板变量">
      <table class="table table-bordered" style="margin-bottom:0">
        <thead><tr><th>变量名</th><th>默认变量值</th></tr></thead>
        <tbody>
        <tr v-for="(item,index) in form.variables" :key="index" >
          <td><el-input v-model="item.key" size="large" /></td>
          <td><el-input v-model="item.value" size="large" /></td>
          <td width="80">
            <el-button-group>
              <el-button icon="Plus" circle @click="addVar(index)"></el-button>
              <el-button icon="Close" circle @click="removeVar(index)"></el-button>
            </el-button-group>
          </td>
        </tr>
        </tbody>
      </table>
    </el-form-item>
    <el-form-item label="设置版本" v-if="props.enableFaas===true">
      <el-card v-for="(m,index) in form.versions" :key="index" style="width:100%">
        <template #header>
          <div class="card-header">
            <span>版本 {{ index+1 }}</span>
            <div class="box-tools pull-right">
              <span class="card-header-btn" @click="copyVersion(index)"><el-icon><CopyDocument /></el-icon></span>
              <span class="card-header-btn" @click="removeVersion(index)"><el-icon><Close /></el-icon></span>
            </div>
          </div>
        </template>
        <el-form label-width="130px">
          <el-form-item label="版本号">
            <el-input v-model="m.version" size="large" />
          </el-form-item>
          <el-form-item label="流量比例">
            <el-input-number v-model="m.weight" :max="100" :min="0" size="large" />
          </el-form-item>
          <el-form-item>
            <template #label><el-text>启用HPA 
              <el-tooltip placement="top">
                <template #content>
                  Horizontal Pod Autoscaling（Pod 水平自动伸缩）<br>
                  Kubernetes控制平面通过监视Pod的资源利用率指标，动态的扩缩Pod数量
                </template>
                <el-icon class="text-yellow"><Warning /></el-icon>
              </el-tooltip></el-text>
            </template>
            <el-switch v-model="m.enable_hpa" />
          </el-form-item>
          <el-form-item v-if="m.enable_hpa">
            <template #label><el-text>最小实例数量 
              <el-tooltip placement="top">
                <template #content>
                  当Pod根据监控指标进行弹性伸缩时，保留的最小Pod数量
                </template>
                <el-icon class="text-yellow"><Warning /></el-icon>
              </el-tooltip></el-text>
            </template>
            <el-input-number v-model="m.min_pod" :max="100" :min="1" size="large" />
          </el-form-item>
          <el-form-item v-if="m.enable_hpa">
            <template #label><el-text>最大实例数量 
              <el-tooltip placement="top">
                <template #content>
                  当Pod根据监控指标进行弹性伸缩时，扩充的最大Pod数量
                </template>
                <el-icon class="text-yellow"><Warning /></el-icon>
              </el-tooltip></el-text>
            </template>
            <el-input-number v-model="m.max_pod" :max="100" :min="0" size="large" />
          </el-form-item>
          <el-form-item v-if="m.enable_hpa">
            <template #label><el-text>目标CPU利用率 
              <el-tooltip placement="top">
                <template #content>
                  当工作负载的所有Pod的平均CPU利用率低于此目标时，<br>将减少Pod数量，直至最小实例数量；<br>
                  当工作负载的所有Pod的平均CPU利用率高于此目标时，<br>将增加Pod数量，直至最大实例数量
                </template>
                <el-icon class="text-yellow"><Warning /></el-icon>
              </el-tooltip></el-text>
            </template>
            <el-input-number v-model="m.cpu_usage" :max="100" :min="0" size="large" />
          </el-form-item>
          <el-form-item v-if="m.enable_hpa">
            <template #label><el-text>目标内存利用率 
              <el-tooltip placement="top">
                <template #content>
                  当工作负载的所有Pod的平均内存利用率低于此目标时，<br>将减少Pod数量，直至最小实例数量；<br>
                  当工作负载的所有Pod的平均内存利用率高于此目标时，<br>将增加Pod数量，直至最大实例数量
                </template>
                <el-icon class="text-yellow"><Warning /></el-icon>
              </el-tooltip></el-text>
            </template>
            <el-input-number v-model="m.mem_usage" :max="100" :min="0" size="large" />
          </el-form-item>          
        </el-form>
      </el-card>
      <el-row>
        <el-button icon="Plus" circle @click="addVersion()" />
      </el-row>
    </el-form-item>
  </el-form>
</template>

<script setup>
import { onBeforeMount, ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { axios } from '/src/assets/util/axios'
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/ext-language_tools'
/* 变量定义 */
const props = defineProps({
  enableFaas: { type: Boolean },
  form: { type: Object }, 
})
const clusterList = ref({})
const refNew = ref(null)
const rules = reactive({
  cluster: [{required: true, message: '请选择集群'}],
  namespace: [{required: true, message: '请选择命名空间'}],
})
const form = ref({content:'',variables:[],versions:[]})
const templateList = ref([])
/* 生命周期函数 */
onBeforeMount(async () => {
  await getClusterList()
  getTemlates()
  if(props.form?.id) {
    form.value = Object.assign({}, props.form)
    form.value.variables = []
    for(let [k,v] of Object.entries(props.form.variables)) {
      if(k == "Versions") continue
      form.value.variables.push({
        key: k,
        value: v
      })
    }
  }
})
/* methods */
const getClusterList = async () => {
  clusterList.value = await axios.get(`/lizardcd/server/clusters`)
}
const getTemlates = async () => {
  let response = await axios.get(`/lizardcd/db/yaml_template?filter=type==application&page=1&size=100&sort=update_at desc`)
  templateList.value = response.results.map(x => {
    x.variables = JSON.parse(x.variables)
    return x
  })
}
const selectTemplate = (val) => {
  if(val) {
    form.value.content = val.content
    form.value.variables = val.variables
  }
  else {
    form.value.content = ""
  }
}
const getFormData = async (callback) => {
  await refNew.value.validate(async (valid) => {
    if(valid) {
      let params = Object.assign({}, form.value)
      delete params.templates
      let vars = {
        Versions: params.versions
      }
      for(let x of params.variables) {
        vars[x.key] = x.value
      }
      params.variables = vars
      callback(params)
    }
    else {
      ElMessage.warning('必填项未填完')
    }
  })
}
const addVar = (index) => {
  form.value.variables.splice(index+1, 0, {key:"", value:""})
}
const removeVar = (index) => {
  form.value.variables.splice(index, 1)
}
const addVersion = () => {
  form.value.versions.push({version:"", weight:50, min_pod: 1, max_pod: 10, cpu_usage: 50, mem_usage: 50, enable_hpa: false})
}
const removeVersion = (index) => {
  form.value.versions.splice(index, 1)
}
const copyVersion = (index) => {
  let version = Object.assign({}, form.value.versions[index])
  form.value.versions.splice(index, 0, Object.assign({}, version))
}
defineExpose({ getFormData })
</script>