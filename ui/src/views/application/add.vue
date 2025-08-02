<template>
<el-drawer v-model="show" @opened="afterOpen" direction="rtl" size="750px">
  <template #header>
    <h4 v-if="props.edit">编辑应用</h4>
    <h4 v-else>新增应用</h4>
  </template>
  <template #default>
    <el-form ref="formAdd" :model="form" :rules="rules" label-width="130px">
      <el-divider><span style="color:#b4b4b4">基础配置</span></el-divider>
      <el-form-item label="应用名称" prop="app_name">
        <el-input v-model="form.app_name" size="large" />
      </el-form-item>
      <el-form-item label="部署方式">
        <el-radio-group v-model="form.deploy_type" @change="selectDeployType" :disabled="props.edit">
          <el-radio-button value="虚拟机">虚拟机</el-radio-button>
          <el-radio-button value="容器">容器</el-radio-button>
          <el-radio-button value="Docker">Docker</el-radio-button>
          <el-radio-button value="HTTP">HTTP</el-radio-button>
          <el-radio-button value="GitOps">GitOps</el-radio-button>
          <el-radio-button value="SSH">SSH</el-radio-button>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="代码仓库" prop="git_http_url">
        <el-input v-model="form.git_http_url" size="large" />
        <myTips type="info">仅支持http(s)开头</myTips>
      </el-form-item>
      <el-form-item label="镜像仓库/制品库" prop="repo" v-if="form.deploy_type!=='GitOps'">
        <el-select v-model="form.repo" placeholder="请选择" value-key="id" clearable size="large" style="width:100%">
          <el-option v-for="item in repoList" :key="item.id" :label="item.repo_url" :value="item">
            <span style="float:left">{{item.repo_url}}</span>
            <span style="float:right;color:var(--el-text-color-secondary)">{{item.repo_account}}</span>
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="仓库/项目" prop="repo_name" v-if="form.deploy_type!=='GitOps'">
        <el-input v-model="form.repo_name" placeholder="请填写" size="large" />
        <myTips type="info">Artifactory填写仓库名<br>Harbor填写项目名<br>DockerHub填写namespace<br>S3填写桶名</myTips>
      </el-form-item>
      <el-form-item label="镜像名/制品路径" prop="image_name" v-if="form.deploy_type!=='GitOps'">
        <el-input v-model="form.image_name" placeholder="请填写" size="large" />
      </el-form-item>
      <el-form-item label="所属租户" prop="tenant">
        <el-select 
          v-model="form.tenant" 
          placeholder="请选择租户" 
          clearable 
          filterable 
          :disabled="userInfo.role!=='admin'"
          style="width:100%;" 
          size="large">
          <el-option v-for="(item,index) in props.tenantList" :key="index" :label="item" :value="item" />
        </el-select>
      </el-form-item>
      <el-form-item label="设置标签" prop="tags">
        <el-row v-for="(item,i) in form.tags" :key="i" style="margin-bottom:5px;width:100%">
          <el-button-group>
            <el-input v-model="item.key" size="large" clearable style="width:200px;margin-right:5px" />
            <el-input v-model="item.value" size="large" clearable style="width:200px;" />
            <el-button circle :icon="Delete" size="large" style="float:right" @click="removeTag(i)" />
          </el-button-group>
        </el-row>
        <el-button circle icon="Plus" @click="addTag" />
      </el-form-item>
      <el-form-item label="超时时间(秒)" prop="timeout">
        <el-input type="number" v-model="form.timeout" size="large" style="width:200px" />
      </el-form-item>
      <el-divider v-if="form.deploy_type!=='GitOps'"><span style="color:#b4b4b4">构建配置</span></el-divider>
      <el-form-item v-if="form.deploy_type!=='GitOps'">
        <template #label><el-text>启用自动构建 
          <el-tooltip placement="top">
            <template #content>
              启用自动构建后，Lizardcd将根据配置的默认tekton（必须由管理员在设置中开启持续集成，并设置默认tekton）<br>自动创建相关tekton资源对象，并在代码仓库自动创建webhook触发构建操作。
            </template>
            <el-icon class="text-yellow"><Warning /></el-icon>
          </el-tooltip></el-text>
        </template>
        <el-switch v-model="form.auto_build.enable" />
      </el-form-item>
      <el-form-item label="从构建模板导入" prop="template" v-if="form.deploy_type!=='GitOps'&&form.auto_build.enable===true">
        <el-select 
          v-model="form.template" 
          placeholder="请选择模板" 
          value-key="id" 
          clearable 
          size="large"
          style="width:100%">
          <el-option v-for="item in templateList.filter(n => n.type==='tekton_task')" :key="item.id" :label="item.name" :value="item">
            <span style="float:left">{{item.name}}</span>
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="构建脚本" prop="build_script" v-if="form.deploy_type!=='GitOps'&&form.auto_build.enable===true">
        <v-ace-editor
          v-model:value="form.auto_build.build_script"
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
      <el-form-item label="版本号获取命令" v-if="form.deploy_type!=='GitOps'&&form.auto_build.enable===true">
        <el-input v-model="form.auto_build.version_script" size="large" />
        <myTips type="info">从代码仓库中获取版本号的命令，仅支持单行，如 cat ./VERSION</myTips>
      </el-form-item>
      <el-divider v-if="form.deploy_type==='GitOps'"><span style="color:#b4b4b4">GitOps部署配置</span></el-divider>
      <el-form-item label="Git分支" v-if="form.gitops&&form.deploy_type==='GitOps'">
        <el-input v-model="form.gitops.git_revision" placeholder="请填写" size="large" />
      </el-form-item>
      <el-form-item label="Path" v-if="form.gitops&&form.deploy_type==='GitOps'">
        <el-input v-model="form.gitops.path" placeholder="请填写" size="large" />
      </el-form-item>
      <el-form-item label="Include" v-if="form.gitops&&form.deploy_type==='GitOps'">
        <el-input v-model="form.gitops.include" placeholder="请填写" size="large" />
        <myTips type="info">多个文件请用英文逗号隔开，支持正则表达式</myTips>
      </el-form-item>
      <el-form-item label="Exclude" v-if="form.gitops&&form.deploy_type==='GitOps'">
        <el-input v-model="form.gitops.exclude" placeholder="请填写" size="large" />
        <myTips type="info">多个文件请用英文逗号隔开，支持正则表达式</myTips>
      </el-form-item>
      <el-form-item label="同步方式" v-if="form.gitops&&form.deploy_type==='GitOps'">
        <el-radio-group v-model="form.gitops.sync_type">
          <el-radio value="手动">手动</el-radio>
          <el-radio value="自动">自动</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="同步周期" v-if="form.gitops&&form.deploy_type==='GitOps'&&form.gitops.sync_type==='自动'">
        <el-input v-model="form.gitops.cron" size="large" placeholder="请填写" />
        <myTips type="info">符合Linux Crontab格式： 分 时 日 月 周（五位）<br>
          如选择Dolphinscheduler定时，则格式为：秒 分 时 日 月 周 年（七位）<br>可参考 https://www.matools.com/cron/
        </myTips>
      </el-form-item>
      <el-form-item label="同步参数" v-if="form.gitops&&form.deploy_type==='GitOps'">
        <el-checkbox v-model="form.gitops.prune_on_delete" label="删除应用时同步删除 K8S 资源" />
        <el-checkbox v-model="form.gitops.use_ds" label="使用Dolphinscheduler定时同步" @change="checkDs" v-if="form.gitops.sync_type==='自动'" />
      </el-form-item>
      <el-divider v-if="['容器','GitOps'].includes(form.deploy_type)"><span style="color:#b4b4b4">容器部署配置</span></el-divider>
      <el-form-item v-if="form.deploy_type==='容器'">
        <template #label><el-text>设置为FaaS 
          <el-tooltip placement="top">
            <template #content>
              Function as a Service，用户配置好代码库以后，仅需提交代码，<br>Lizardcd会自动按照您的配置将应用部署于目标集群，<br>并设置好流量分配、自动伸缩、访问域名
            </template>
            <el-icon class="text-yellow"><Warning /></el-icon>
          </el-tooltip></el-text>
        </template>
        <el-switch v-model="form.enable_faas" />
      </el-form-item>
      <el-form-item v-if="form.deploy_type==='容器'&&!form.enable_faas">
        <template #label><el-text>开启流量控制 
          <el-tooltip placement="top">
            <template #content>
              不使用FaaS，仅开启流量控制<br>可用于灰度发布或AB测试
            </template>
            <el-icon class="text-yellow"><Warning /></el-icon>
          </el-tooltip></el-text>
        </template>
        <el-switch v-model="form.enable_traffic_control" />
      </el-form-item>
      <el-form-item label="流量分配策略" v-if="form.deploy_type==='容器'&&form.enable_traffic_control">
        <el-radio-group v-model="form.traffic_policy">
          <el-radio value="weight">基于流量比例</el-radio>
          <el-radio value="header">基于头部字段</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="工作负载" v-if="['容器','GitOps'].includes(form.deploy_type)&&!form.enable_faas">
        <el-card v-for="(m,index) in form.workload" :key="index" style="width:100%">
          <template #header>
            <div class="card-header">
              <span>工作负载 {{ index+1 }}</span>
              <div class="box-tools pull-right">
                <span class="card-header-btn" @click="copyWorkload(index)"><el-icon><CopyDocument /></el-icon></span>
                <span class="card-header-btn" @click="removeWorkload(index)"><el-icon><Close /></el-icon></span>
              </div>
            </div>
          </template>
          <el-form label-width="100px">
            <el-form-item label="容器集群">
              <el-select v-model="m.cluster" clearable placeholder="请选择" size="large" style="width:100%">
                <el-option v-for="(v,k,i) in props.k8scluster" :key="i" :label="k" :value="k" />
              </el-select>
            </el-form-item>
            <el-form-item label="命名空间">
              <el-select v-model="m.namespace" clearable placeholder="请选择" size="large" @change="listResource(m.cluster,m.namespace,m.workload_type)" style="width:100%">
                <el-option v-for="item in props.k8scluster[m.cluster]" :key="item" :label="item" :value="item" />
              </el-select>
            </el-form-item>
            <el-form-item label="工作负载类型" v-if="form.deploy_type==='容器'">
              <el-radio-group v-model="m.workload_type" @change="listResource(m.cluster,m.namespace,m.workload_type)">
                <el-radio-button label="deployments" value="deployments" />
                <el-radio-button label="statefulsets" value="statefulsets" />
                <el-radio-button label="jobs" value="jobs" />
                <el-radio-button label="cronjobs" value="cronjobs" />
                <el-radio-button label="yaml" value="yaml" />
              </el-radio-group>
            </el-form-item>
            <el-form-item label="工作负载名称" v-if="form.deploy_type==='容器'&&m.workload_type!=='yaml'">
              <el-select 
                v-model="m.workload_name" 
                placeholder="请选择工作负载" 
                clearable 
                filterable
                @change="listContainers(m.cluster,m.namespace,m.workload_type,m.workload_name)" 
                style="width:100%" 
                size="large">
                <el-option v-for="item in deploymentList[`${m.cluster}#${m.namespace}`]" :key="item" :label="item" :value="item" />
              </el-select>
            </el-form-item>
            <el-form-item label="容器名称" v-if="form.deploy_type==='容器'&&m.workload_type!=='yaml'">
              <el-select 
                v-model="m.container_name" 
                placeholder="请选择容器" 
                clearable 
                filterable
                style="width:100%" 
                size="large">
                <el-option v-for="item in containerList[`${m.cluster}#${m.namespace}#${m.workload_name}`]" :key="item" :label="item" :value="item" />
              </el-select>
            </el-form-item>
            <el-form-item label="版本号" v-if="form.enable_traffic_control">
              <el-input v-model="m.version" size="large" />
              <myTips type="info">版本号必须和POD的labels.version一致</myTips>
            </el-form-item>
            <el-form-item label="流量比例" v-if="form.enable_traffic_control&&form.traffic_policy==='weight'">
              <el-input-number v-model="m.weight" :max="100" :min="0" size="large" />
            </el-form-item>
            <el-form-item label="匹配头部字段" v-if="form.enable_traffic_control&&form.traffic_policy==='header'">
              <table class="table table-bordered">
                <thead><tr><th width="30%">头部键</th><th width="40%">匹配方式</th><th width="30%">头部值</th></tr></thead>
                <tbody>
                  <tr v-for="(n,i) in m.headers" :key="i">
                    <td><el-input v-model="n.key" /></td>
                    <td>
                      <el-select v-model="n.match_type">
                        <el-option label="exact" value="exact" />
                        <el-option label="prefix" value="prefix" />
                        <el-option label="regex" value="regex" />
                      </el-select>
                    </td>
                    <td><el-input v-model="n.value" /></td>
                  </tr>
                </tbody>
              </table>
              <el-row>
                <el-button icon="Plus" circle @click="addMatchHeader(m.headers)"></el-button>
                <el-button icon="Close" circle @click="removeMatchHeader(m.headers)"></el-button>
              </el-row>
            </el-form-item>
            <el-form-item label="是否启用">
              <el-switch v-model="m.enable" />
            </el-form-item>
          </el-form>
        </el-card>
        <el-row>
          <el-button icon="Plus" circle @click="addWorkload" />
          <el-tooltip content="全部启用">
            <el-button circle @click="setWorkloadEnable(form,true)"><font-awesome-icon icon="toggle-on" /></el-button>
          </el-tooltip>
          <el-tooltip content="全部停用">
            <el-button circle @click="setWorkloadEnable(form,false)"><font-awesome-icon icon="toggle-off" /></el-button>
          </el-tooltip>
        </el-row>
      </el-form-item>
      <el-divider v-if="['虚拟机','SSH'].includes(form.deploy_type)"><span style="color:#b4b4b4">虚拟机/SSH部署配置</span></el-divider>
      <el-form-item label="SSH端口" prop="ssh_port" v-if="form.deploy_type==='SSH'">
        <el-input v-model="form.extra_vars.ssh_port" size="large" />
      </el-form-item>
      <el-form-item label="SSH用户" prop="ssh_user" v-if="form.deploy_type==='SSH'">
        <el-input v-model="form.extra_vars.ssh_user" size="large" />
      </el-form-item>
      <el-form-item label="SSH密码" prop="ssh_pass" v-if="form.deploy_type==='SSH'">
        <el-input v-model="form.extra_vars.ssh_pass" size="large" type="password" />
      </el-form-item>
      <el-form-item label="SSH私钥" prop="ssh_private_key" v-if="form.deploy_type==='SSH'&&edit===false">
        <el-input v-model="form.extra_vars.ssh_private_key" size="large" type="textarea" :autosize="{minRows:5}"/>
      </el-form-item>
      <el-form-item label="部署路径" prop="deploy_path" v-if="['虚拟机','SSH'].includes(form.deploy_type)">
        <el-input v-model="form.extra_vars.deploy_path" size="large" />
      </el-form-item>
      <el-form-item label="部署用户" prop="deploy_user" v-if="form.deploy_type==='虚拟机'">
        <el-input v-model="form.extra_vars.deploy_user" size="large" />
      </el-form-item>
      <el-form-item label="脚本类型" v-if="['虚拟机','SSH'].includes(form.deploy_type)">
        <el-radio-group v-model="form.extra_vars.command_type">
          <el-radio value="shell">shell</el-radio>
          <el-radio value="bat">windows bat</el-radio>
          <el-radio value="ps1">windows powershell</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="部署前命令/脚本" v-if="['虚拟机','SSH'].includes(form.deploy_type)">
        <v-ace-editor
          v-model:value="form.extra_vars.pre_command"
          lang="yaml"
          theme="chrome"
          style="width:100%"
          :options="{
            enableBasicAutocompletion: true,
            enableSnippets: true,
            enableLiveAutocompletion: true,
            tabSize: 2,
            showPrintMargin: false,
            fontSize: 14,
            minLines: 5,
            maxLines: 50,
            wrap: true
          }" />
        <myTips type="info">通常可执行停止服务、清理旧版本等操作</myTips>
      </el-form-item>
      <el-form-item label="启动命令/脚本" prop="start_command" v-if="['虚拟机','SSH'].includes(form.deploy_type)">
        <v-ace-editor
          v-model:value="form.extra_vars.start_command"
          lang="yaml"
          theme="chrome"
          style="width:100%"
          :options="{
            enableBasicAutocompletion: true,
            enableSnippets: true,
            enableLiveAutocompletion: true,
            tabSize: 2,
            showPrintMargin: false,
            fontSize: 14,
            minLines: 5,
            maxLines: 50,
            wrap: true
          }" />
      </el-form-item>
      <el-divider v-if="['虚拟机','SSH'].includes(form.deploy_type)"><span style="color:#b4b4b4">健康检查配置</span></el-divider>
      <el-form-item label="健康检查方式" v-if="['虚拟机','SSH'].includes(form.deploy_type)&&form.extra_vars.health_check">
        <el-radio-group v-model="form.extra_vars.health_check.type">
          <el-radio value="none">无</el-radio>
          <el-radio value="http">HTTP</el-radio>
          <el-radio value="shell">Shell</el-radio>
          <el-radio value="tcp">TCP</el-radio>
        </el-radio-group>
        <myTips type="info">如开启健康检查，则检查通过才表示部署成功</myTips>
      </el-form-item>
      <el-form-item label="健康检查内容" v-if="['虚拟机','SSH'].includes(form.deploy_type)&&form.extra_vars.health_check?.type!=='none'">
        <myTips type="info" v-if="['虚拟机','SSH'].includes(form.deploy_type)&&form.extra_vars.health_check?.type==='http'">返回200OK表示健康检查通过</myTips>
        <table class="table table-bordered" v-if="['虚拟机','SSH'].includes(form.deploy_type)&&form.extra_vars.health_check?.type==='http'">
          <tbody>
          <tr>
            <td width="80">Method</td>
            <td>
              <el-radio-group v-model="form.extra_vars.health_check.method">
                <el-radio value="get">GET</el-radio>
                <el-radio value="post">POST</el-radio>
              </el-radio-group>
            </td>
          </tr>
          <tr>
            <td>Port</td>
            <td><el-input v-model="form.extra_vars.health_check.port" size="large" /></td>
          </tr>
          <tr>
            <td>Uri</td>
            <td><el-input v-model="form.extra_vars.health_check.uri" size="large" /></td>
          </tr>
          </tbody>
        </table>
        <myTips type="info" v-if="['虚拟机','SSH'].includes(form.deploy_type)&&form.extra_vars.health_check?.type==='shell'">脚本返回值为0表示健康检查通过</myTips>
        <v-ace-editor
          v-if="['虚拟机','SSH'].includes(form.deploy_type)&&form.extra_vars.health_check?.type==='shell'"
          v-model:value="form.extra_vars.health_check.shell"
          lang="yaml"
          theme="chrome"
          style="width:100%"
          :options="{
            enableBasicAutocompletion: true,
            enableSnippets: true,
            enableLiveAutocompletion: true,
            tabSize: 2,
            showPrintMargin: false,
            fontSize: 14,
            minLines: 5,
            maxLines: 50,
            wrap: true
          }" />
        <table class="table table-bordered" v-if="['虚拟机','SSH'].includes(form.deploy_type)&&form.extra_vars.health_check?.type==='tcp'">
          <tbody>
          <tr>
            <td>Port</td>
            <td><el-input v-model="form.extra_vars.health_check.port" size="large" /></td>
          </tr>
          </tbody>
        </table>
      </el-form-item>
      <el-divider v-if="form.deploy_type==='Docker'"><span style="color:#b4b4b4">Docker部署配置</span></el-divider>
      <el-form-item label="容器名" prop="container_name" v-if="form.deploy_type==='Docker'">
        <el-input v-model="form.extra_vars.container_name" size="large" />
      </el-form-item>
      <el-form-item label="端口映射" prop="ports" v-if="form.deploy_type==='Docker'">
        <el-input v-model="form.extra_vars.ports" size="large" placeholder="例如：80:8080,90:9090，多个端口用逗号隔开" />
      </el-form-item>
      <el-form-item label="卷绑定" prop="volumes" v-if="form.deploy_type==='Docker'">
        <el-input v-model="form.extra_vars.volumes" size="large" placeholder="例如：/data/volume1:/volume1，多个卷用逗号隔开" />
      </el-form-item>
      <el-form-item label="Network" prop="network" v-if="form.deploy_type==='Docker'">
        <el-input v-model="form.extra_vars.network" size="large" placeholder="网络方式，如使用本地网络，则填写 host" />
      </el-form-item>
      <el-form-item label="DNS" prop="dns" v-if="form.deploy_type==='Docker'">
        <el-input v-model="form.extra_vars.dns" size="large" placeholder="多个dns用逗号隔开" />
      </el-form-item>
      <el-form-item label="WorkingDir" prop="working_dir" v-if="form.deploy_type==='Docker'">
        <el-input v-model="form.extra_vars.working_dir" size="large" />
      </el-form-item>
      <el-form-item label="启动参数" prop="command" v-if="form.deploy_type==='Docker'">
        <el-input v-model="form.extra_vars.command" size="large" placeholder='数组形式，如 ["ls", "/tmp"]' />
      </el-form-item>
      <el-form-item label="目标服务器" v-if="['虚拟机','Docker'].includes(form.deploy_type)">
        <el-select v-model="form.targets" placeholder="请选择" clearable filterable multiple size="large" value-key="ip" style="width:100%">
          <el-option v-for="item in targetList" :key="item" :label="item.ip" :value="item">
            <span style="float:left">{{item.ip}}</span>
            <span style="float:right;color:var(--el-text-color-secondary);font-size:12px" v-if="item.labels">
              <el-tag v-for="label in item.labels" :key="label" size="small">{{ label }}</el-tag>
            </span>
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="目标服务器" v-if="form.deploy_type==='SSH'">
        <el-card v-for="(m,index) in form.workload" :key="index" style="width:100%">
          <template #header>
            <div class="card-header">
              <span>目标服务器 {{ index+1 }}</span>
              <div class="box-tools pull-right">
                <span class="card-header-btn" @click="copyWorkload(index)"><el-icon><CopyDocument /></el-icon></span>
                <span class="card-header-btn" @click="removeWorkload(index)"><el-icon><Close /></el-icon></span>
              </div>
            </div>
          </template>
          <el-form label-width="70px">
            <el-form-item label="服务器IP">
              <el-input v-model="m.workload_name" size="large" />
            </el-form-item>
            <el-form-item label="设置标签" prop="labels">
              <el-row v-for="(item,i) in m.labels" :key="i" style="margin-bottom:5px;width:100%">
                <el-button-group>
                  <el-input v-model="item.key" size="large" clearable style="width:175px;margin-right:5px" />
                  <el-input v-model="item.value" size="large" clearable style="width:175px;" />
                  <el-button circle :icon="Delete" size="large" style="float:right" @click="m.labels.splice(i, 1)" />
                </el-button-group>
              </el-row>
              <el-button circle icon="Plus" @click="m.labels.push({key: '', value: ''})" />
            </el-form-item>
            <el-form-item label="是否启用">
              <el-switch v-model="m.enable" />
            </el-form-item>
          </el-form>
        </el-card>
        <el-row>
          <el-button icon="Plus" circle @click="addWorkload" />
          <el-tooltip content="全部启用">
            <el-button circle @click="setWorkloadEnable(form,true)"><font-awesome-icon icon="toggle-on" /></el-button>
          </el-tooltip>
          <el-tooltip content="全部停用">
            <el-button circle @click="setWorkloadEnable(form,false)"><font-awesome-icon icon="toggle-off" /></el-button>
          </el-tooltip>
        </el-row>
      </el-form-item>
      <el-divider v-if="form.deploy_type==='HTTP'"><span style="color:#b4b4b4">HTTP部署配置</span></el-divider>
      <el-form-item label="Base URL" v-if="form.deploy_type==='HTTP'">
        <el-input v-model="form.extra_vars.http_url" size="large" placeholder="http://<ip or domain>:<port>" />
      </el-form-item>
      <el-form-item label="Http Header" v-if="form.deploy_type==='HTTP'">
        <el-input v-model="form.extra_vars.http_header" size="large" />
        <myTips type="info">例如：Authorization: Bearer xxx (多个header用英文逗号隔开)</myTips>
      </el-form-item>
      <el-form-item label="部署Method" v-if="form.deploy_type==='HTTP'">
        <el-radio-group v-model="form.extra_vars.http_method">
          <el-radio value="post">POST</el-radio>
          <el-radio value="put">PUT</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="部署Path" v-if="form.deploy_type==='HTTP'">
        <el-input v-model="form.extra_vars.http_path" size="large" />
      </el-form-item>
      <el-form-item label="提交Body" v-if="form.deploy_type==='HTTP'">
        <v-ace-editor
          v-model:value="form.extra_vars.http_body"
          lang="json"
          theme="chrome"
          style="width:100%"
          :options="{
            enableBasicAutocompletion: true,
            enableSnippets: true,
            enableLiveAutocompletion: true,
            tabSize: 2,
            showPrintMargin: false,
            fontSize: 14,
            minLines: 5,
            maxLines: 50,
            wrap: true
          }" />
        <myTips type="info"><div v-html="tips1" /></myTips>
      </el-form-item>
      <el-form-item label="Content-Type" v-if="form.deploy_type==='HTTP'">
        <el-radio-group v-model="form.extra_vars.http_content_type">
          <el-radio value="json">json</el-radio>
          <el-radio value="x-www-form-urlencoded">x-www-form-urlencoded</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="部署成功判断" v-if="form.deploy_type==='HTTP'">
        <el-input v-model="form.extra_vars.res_jsonpath" size="large" style="width:200px" placeholder="支持jsonpath" /> =~
        <el-input v-model="form.extra_vars.res_keyword" size="large" style="width:200px" placeholder="支持正则表达式" />
      </el-form-item>
      <el-divider v-if="form.deploy_type==='HTTP'"><span style="color:#b4b4b4">健康检查配置</span></el-divider>
      <el-form-item label="健康检查Method" v-if="form.deploy_type==='HTTP'&&form.extra_vars.health_check">
        <el-radio-group v-model="form.extra_vars.health_check.method">
          <el-radio value="none">无</el-radio>
          <el-radio value="get">GET</el-radio>
          <el-radio value="post">POST</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="健康检查Path" v-if="form.deploy_type==='HTTP'&&form.extra_vars.health_check">
        <el-input v-model="form.extra_vars.health_check.http_path" size="large" />
        <myTips type="info"><div v-html="tips2" /></myTips>
      </el-form-item>
      <el-form-item label="健康检查Body" v-if="form.deploy_type==='HTTP'&&form.extra_vars.health_check">
        <v-ace-editor
          v-model:value="form.extra_vars.health_check.body"
          lang="json"
          theme="chrome"
          style="width:100%"
          :options="{
            enableBasicAutocompletion: true,
            enableSnippets: true,
            enableLiveAutocompletion: true,
            tabSize: 2,
            showPrintMargin: false,
            fontSize: 14,
            minLines: 5,
            maxLines: 50,
            wrap: true
          }" />
      </el-form-item>
      <el-form-item label="检查结束判断" v-if="form.deploy_type==='HTTP'">
        <el-input v-model="form.extra_vars.health_check.finish_jsonpath" size="large" style="width:200px" placeholder="支持jsonpath" /> =~
        <el-input v-model="form.extra_vars.health_check.finish_keyword" size="large" style="width:200px" placeholder="支持正则表达式" />
      </el-form-item>
      <el-form-item label="检查成功判断" v-if="form.deploy_type==='HTTP'&&form.extra_vars.health_check">
        <el-input v-model="form.extra_vars.health_check.success_jsonpath" size="large" style="width:200px" placeholder="支持jsonpath" /> =~
        <el-input v-model="form.extra_vars.health_check.success_keyword" size="large" style="width:200px" placeholder="支持正则表达式" />
      </el-form-item>
      <el-form-item label="提取输出信息" v-if="form.deploy_type==='HTTP'&&form.extra_vars.health_check">
        <el-input v-model="form.extra_vars.health_check.msg_jsonpath" size="large" />
        <myTips type="info"><div v-html="tips3" /></myTips>
      </el-form-item>
    </el-form>
    <newWorkload v-if="form.deploy_type==='容器'&&form.enable_faas===true" ref="refWorkload" :enableFaas="form.enable_faas" :form="form.faas" />
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show=false">取消</el-button>
      <el-button type="primary" @click="confirmClick(formAdd)" :loading="loading.add">提交</el-button>
    </div>
  </template>
</el-drawer>
</template>
<script setup>
import { Warning,CopyDocument,Close,Delete } from '@element-plus/icons-vue'
import { onBeforeMount, ref, reactive, computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { ElMessage } from 'element-plus'
import newWorkload from '/src/views/kubernetes/new.vue'
import MyTips from '/src/components/myTips/myTips.vue'
import { axios } from '/src/assets/util/axios'
import { getImageBase, parseLabels } from '@/assets/util/common'
import _ from 'lodash'
import moment from 'moment'
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/mode-json'
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/ext-language_tools'
/* 变量定义 */
const props = defineProps({
  form: { type: Object }, 
  k8scluster: { type: Object },
  repoList: { type: Array },
  tenantList: { type: Array },
  edit: { type: Boolean },
})
const store = useStore()
const userInfo = computed(() => {
  return store.state.userInfo
})
const tenant = localStorage.tenant.split(",")[0]
const show = ref(false)
const form = ref({auto_build:{}})
const rules = reactive({
  app_name: [{required: true, message: '请填写应用名称'}],
  git_http_url: [{ validator: (rule, value, callback) => {
    if(value && !value.startsWith('http')) {
      callback(new Error("代码仓库必须以http(s)开头"))
    } else {
      callback()
    }
  }, trigger: 'blur' }],
})
const deploymentList = ref({})
const containerList = ref({})
const templateList = ref([])
const targetList = ref([])
const formAdd = ref(null)
const refWorkload = ref(null)
const loading = ref({
  add: false
})
// 添加标签的三个变量
const tips1 = ref("部署制品可用变量 {{artifact_url}} 进行占位")
const tips2 = ref("支持使用变量 {{response$.data.id}} 接收部署接口的返回json，提取jsonpath字段")
const tips3 = ref("支持使用变量 {{$.data.msg}} 提取jsonpath字段")
const defaultTekton = ref({})
/* 生命周期函数 */
onBeforeMount(async () => {
  getSettings()
})
/* methods */
const getTemlates = async () => {
  let response = await axios.get(`/lizardcd/db/yaml_template?search=type==tekton_&page=1&size=100&sort=update_at desc`)
  templateList.value = response.results.map(x => {
    x.variables = JSON.parse(x.variables)
    return x
  })
}
const getSettings = async () => {
  let response = await axios.get(`/lizardcd/db/settings?filter=setting_key==default_tekton`)
  if(response.total > 0) {
    defaultTekton.value = JSON.parse(response.results[0].setting_value)
    if(defaultTekton.value.cluster&&defaultTekton.value.namespace) {
      response = await axios.get(`/lizardcd/tekton/cluster/${defaultTekton.value.cluster}/namespace/${defaultTekton.value.namespace}/eventlisteners`)
      let el = response.results.find(n => n.metadata.name === 'default-eventlistener')
      defaultTekton.value.trigger_endpoint = el?.metadata?.annotations?.endpoint
    }
  }
}
const listResource = async (cluster, namespace, workload_type) => {
  if(cluster && namespace && workload_type !== 'yaml') {
    let response = await axios.get(`/lizardcd/kubernetes/cluster/${cluster}/namespace/${namespace}/${workload_type}`)
    deploymentList.value[`${cluster}#${namespace}`] = response.results.map(x => x.metadata.name)
  }
}
const listContainers = async (cluster, namespace, workload_type, workload_name) => {
  if(cluster && namespace && workload_name && workload_type !== 'yaml') {
    let response = await axios.get(`/lizardcd/kubernetes/cluster/${cluster}/namespace/${namespace}/${workload_type}/${workload_name}`)
    if(workload_type === 'cronjobs') {
      containerList.value[`${cluster}#${namespace}#${workload_name}`] = response.spec.jobTemplate.spec.template.spec.containers.map(x => x.name).concat(response.spec.jobTemplate.spec.template.spec.initContainers?.map(x => x.name)||[])
    } else {
      containerList.value[`${cluster}#${namespace}#${workload_name}`] = response.spec.template.spec.containers.map(x => x.name).concat(response.spec.template.spec.initContainers?.map(x => x.name)||[])
    }
  }
}
const getTargets = async () => {
  targetList.value = await axios.get(`/lizardcd/server/targets`)
}
const selectDeployType = (val) => {
  if(val === 'HTTP') {
    form.value.extra_vars.http_method = 'post'
    form.value.extra_vars.health_check.method = 'none'
    form.value.extra_vars.http_content_type = 'json'
  }
  if(val === 'GitOps') {
    form.value.gitops = {
      sync_type: '手动',
    }
  }
}
const afterOpen = async () => {
  await getTemlates()
  await getTargets()
  if(props.form?.id) {
    form.value = _.cloneDeep(props.form)
    form.value.auto_build ||= {enable: false,build_script:'',version_script:''}
    form.value.timeout ||= 300
    form.value.template = templateList.value.find(n => n.id === form.value.auto_build?.template_id)
  } else {
    form.value = {
      workload: [],
      traffic_policy: 'weight',
      enable_traffic_control: false,
      tenant: tenant,
      tags: [],
      deploy_type: '容器',
      targets: [],
      extra_vars: {health_check:{type:'none'},command_type:'shell',pre_command:'',start_command:''},
      auto_build: {enable: false,build_script:'',version_script:''},
      timeout: 300
    }
  }
  form.value.tenant ||= localStorage.tenant
  form.value.repo = props.repoList.find(n => n.id === form.value.repo_id)
  if(form.value.template?.id === 0) form.value.template = undefined
  let tags = []
  for (let tag of form.value.tags||[]) {
    let arr = tag.replace(":", "=").split("=")
    tags.push({key: arr[0].trim(), value: arr[1].trim()})
  }
  form.value.tags = tags
  for(let w of form.value.workload) {
    w.headers ||= []
    if(!deploymentList.value.hasOwnProperty(`${w.cluster}#${w.namespace}`)) {
      listResource(w.cluster, w.namespace, w.workload_type)
    }
    if(!containerList.value.hasOwnProperty(`${w.cluster}#${w.namespace}${w.workload_name}`)) {
      listContainers(w.cluster, w.namespace, w.workload_type, w.workload_name)
    }
  }
  if(['虚拟机','Docker'].includes(form.value.deploy_type)) {
    form.value.targets = form.value.workload.map(x => { 
      return { ip: x.workload_name, labels: x.labels}
    })
    try {
      form.value.extra_vars = await JSON.parse(form.value.extra_vars)
      form.value.extra_vars.health_check ||= {type:'none',method:'get',shell:''}
    }
    catch(e) {
      form.value.extra_vars = {}
      form.value.extra_vars.health_check = {}
    }
  }
  if(form.value.deploy_type === 'SSH') {
    for(let x of form.value.workload) {
      x.labels = parseLabels(x.labels)
    }
    try {
      form.value.extra_vars = await JSON.parse(form.value.extra_vars)
      form.value.extra_vars.health_check ||= {type:'none',method:'get',shell:''}
    }
    catch(e) {
      form.value.extra_vars = {}
      form.value.extra_vars.health_check = {}
    }
  }
  if(form.value.deploy_type === 'HTTP') {
    form.value.extra_vars = JSON.parse(form.value.extra_vars)
    let headers = []
    for(let [k,v] of Object.entries(form.value.extra_vars.http_header)) {
      headers.push(`${k}: ${v}`)
    }
    form.value.extra_vars.http_header = headers.join(",")
  }
  // 获取faas
  if(form.value.id) {
    let response = await axios.get(`/lizardcd/db/application_faas?filter=application_id==${form.value.id}`)
    if(response.total > 0) {
      form.value.enable_faas = true
      form.value.faas = response.results[0]
    }
  }
}
const addTag = () => {
  form.value.tags.push({key: '', value: ''})
}
const removeTag = (index) => {
  form.value.tags.splice(index, 1)
}
const addWorkload = () => {
  form.value.workload.push({
    cluster: '',
    namespace: '',
    workload_type: 'deployments',
    workload_name: '',
    container_name: '',
    weight: 50,
    headers: [],
    labels: [],
    enable: true
  })
}
const removeWorkload = (index) => {
  form.value.workload.splice(index, 1)
}
const copyWorkload = (index) => {
  let workload = Object.assign({}, form.value.workload[index])
  workload.headers = Object.assign([], workload.headers)
  form.value.workload.splice(index, 0, Object.assign({}, workload))
}
const setWorkloadEnable = (form, enable) => {
  for(let x of form.workload||form.app_name.workload) {
    x.enable = enable
  }
}
const addMatchHeader = (headers) => {
  headers.push({
    key: "",
    match_type: "exact",
    value: ""
  })
}
const removeMatchHeader = (headers) => {
  headers.pop()
}
const confirmClick = async (f) => {
  if(!f) return
  await f.validate(async (valid) => {
    if(valid) {
      let params = _.cloneDeep(form.value)
      params.update_at = moment()
      params.timeout = parseInt(params.timeout)
      params.repo_id = params.repo?.id
      params.auto_build.template_id = params.template?.id
      delete params.repo
      delete params.template
      params.tags = params.tags.map(x => `${x.key}=${x.value}`)
      if(['虚拟机','Docker'].includes(params.deploy_type)) {
        params.workload = params.targets.map(x => {
          return {
            workload_type: "vm",
            workload_name: x.ip,
            labels: x.labels,
            enable: true
          }
        })
        params.extra_vars = JSON.stringify(params.extra_vars)
      } else if (params.deploy_type === 'SSH') {
        params.workload = params.workload.map(x => {
          return {
            workload_type: "ssh",
            workload_name: x.workload_name,
            labels: x.labels.map(y => `${y.key}=${y.value}`),
            enable: true
          }
        })
        params.extra_vars = JSON.stringify(params.extra_vars)
      } else if (params.deploy_type === 'HTTP') {
        delete params.extra_vars.pre_command
        delete params.extra_vars.start_command
        delete params.extra_vars.health_check.type
        let headers = {}
        for(let x of params.extra_vars.http_header.split(",")) {
          let h = x.split(":")
          headers[h[0].trim()] = h[1].trim()
        }
        params.extra_vars.http_header = headers
        params.extra_vars = JSON.stringify(params.extra_vars)
      } else {
        delete params.extra_vars
      }
      delete params.targets
      if(params.enable_traffic_control === true && params.traffic_policy === 'weight') {
        let weightTotal = 0
        for(let x of form.value.workload) {
          weightTotal += x.weight
        }
        if(weightTotal != 100) {
          ElMessage.warning('所有工作负载权重之和必须等于100')
          return
        }
      }
      if(params.deploy_type !== 'GitOps') {
        delete params.gitops
      }
      let enable_faas = params.enable_faas
      delete params.enable_faas
      if(props.edit === false) {
        delete params.id
        let response = await axios.post(`/lizardcd/db/application`, {body:params})
        form.value.id = response.id
        if(enable_faas) {
          refWorkload.value.getFormData(async (w) => {
            w.application_id = response.id
            if(!w.variables.hasOwnProperty("Appname"))
              ElMessage.warning('模板变量必须设置变量名Appname')
            else {
              try { 
                await axios.post(`/lizardcd/db/application_faas`, {body:w}) 
                for(let v of w.versions) {
                  await axios.put(`/lizardcd/kubernetes/cluster/${w.cluster}/namespace/${w.namespace}/deployments/${w.variables.Appname}-${v.version}/hpa`, {
                    "max": v.max_pod,
                    "min": v.min_pod,
                    "cpu": v.cpu_usage,
                    "memory": v.mem_usage
                  }) 
                }
              } catch(e){}
            }
          })
        }
      }
      else {  // 编辑应用
        delete params.faas
        await axios.put(`/lizardcd/db/application/${params.id}`, {body:params})
        if(enable_faas === true) {
          refWorkload.value.getFormData(async (w) => {
            w.application_id = params.id
            if(w.id)
              await axios.put(`/lizardcd/db/application_faas/${w.id}`, {body:w})
            else
              await axios.post(`/lizardcd/db/application_faas`, {body:w})
            for(let v of w.versions) {
              await axios.put(`/lizardcd/kubernetes/cluster/${w.cluster}/namespace/${w.namespace}/deployments/${w.variables.Appname}-${v.version}/hpa`, {
                "max": v.max_pod,
                "min": v.min_pod,
                "cpu": v.cpu_usage,
                "memory": v.mem_usage
              }) 
            }
          })
        } else if(form.value.faas&&form.value.faas.id !== 0) {
          await axios.delete(`/lizardcd/db/application_faas/${form.value.faas.id}`)
        }
      }
      if(form.value.auto_build.enable) {
        await applyTektonCRD(params, form.value.template.content)
      }
      show.value = false
      loading.value.add = false
      emit('submit')
    }
  })
}
const applyTektonCRD = async (params, template_content) => {
  // 创建 build task
  let scripts = []
  if(params.auto_build.build_script) {
    scripts = params.auto_build.build_script.split('\n')
    for(let i=0; i < scripts.length; i++){
      if(i != 0) {
        scripts[i] = '      ' + scripts[i]
      }
    }
  }
  try {
    await axios.post(`/lizardcd/tekton/cluster/${defaultTekton.value.cluster}/namespace/${defaultTekton.value.namespace}/apply?kind=Task`, {
      content: template_content,
      variables: {
        Appname: params.app_name,
        BuildScript: scripts.join('\n'),
        VersionScript: params.auto_build.version_script || undefined,
        Username: localStorage.username
      }
    })
    ElMessage.success({message: `创建/更新自动构建Task成功`})
  } catch(e) {
    ElMessage.error({message: `创建/更新自动构建Task失败: ${e}`})
    return Promise.reject()
  }
  // 创建 pipeline
  try {
    let template = templateList.value.find(n => n.name === 'tekton_template_pipeline')
    if(!template) {
      ElMessage.warning({message: `未找到模板: tekton_template_pipeline`})
      return Promise.reject()
    }
    await axios.post(`/lizardcd/tekton/cluster/${defaultTekton.value.cluster}/namespace/${defaultTekton.value.namespace}/apply?kind=Pipeline`, {
      content: template.content,
      variables: {
        Appname: params.app_name,
        ImageUrl: getImageBase(form.value.repo, form.value.repo_name, form.value.image_name),
        Context: './',
        Username: localStorage.username
      }
    })
    ElMessage.success({message: `创建/更新自动构建Pipeline成功`})
  } catch(e) {
    ElMessage.error({message: `创建/更新自动构建Pipeline失败: ${e}`})
    return Promise.reject()
  }
  // 创建 trigger
  let response = await axios.get(`/lizardcd/db/ci_trigger?filter=app_id==${params.id}`)
  if(response.total === 0) {
    await axios.post(`/lizardcd/db/ci_trigger`, { body: {
      app_id: params.id,
      trigger_name: params.app_name,
      git_http_url: params.git_http_url,
      trigger_path: [""],
      trigger_type: "pipelinerun",
      ref_pattern: ".*",
      match_labels: params.tags,
      trigger_event: "push,merge_request",
      tenant,
      update_at: moment()
    }})
  }
}
const checkDs = async (val) => {
  if(val === true) {
    let response = await axios.get(`/lizardcd/db/settings?filter=setting_key==dolphinscheduler,tenant==${tenant}`)
    if(response.total === 0) {
      ElMessage.warning({message: `未找到Dolphinscheduler设置`})
      return
    }
    let ds = JSON.parse(response.results[0].setting_value)
    if(!ds.base_url) {
      ElMessage.warning({message: `未配置Dolphinscheduler的base_url`})
      return
    }
    if(!ds.token) {
      ElMessage.warning({message: `未配置Dolphinscheduler的token`})
      return
    }
  }
}
const open = () => {
  show.value = true
}
defineExpose({ open })
const emit = defineEmits(['submit'])
</script>