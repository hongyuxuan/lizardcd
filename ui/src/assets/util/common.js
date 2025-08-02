import moment from 'moment'
import axios from 'axios'
import { customAlphabet } from 'nanoid'

const getImageBase = (repo, repo_name, image_name) => {
  if(repo.repo_type === 'Artifactory' || repo.repo_type === 'Harbor') {
    return `${repo.repo_url.split('://')[1]}/${repo_name}/${image_name}`
  }
  return ""
}

const getDuration = (timeTill, timeFrom) => {
  timeTill ||= moment()
  let duration = moment.duration(moment(timeTill).diff(moment(timeFrom)))
  let hours = duration.hours()
  let minutes = duration.minutes()
  let seconds = duration.seconds()
  return `${hours > 0 ? hours + '小时' : ''}${minutes > 0 ? minutes + '分钟':''}${seconds + '秒'}`
}

const getSettings = async () => {
  let tenants = localStorage.tenant?.split(",") || []
  let response = await axios.get(`/lizardcd/db/settings?size=1000&filter=tenant==${tenants[0]}`)
  let settings = {}
  for(let x of response.results) {
    if(x.setting_value === 'true' || x.setting_value === 'false')
      x.setting_value = JSON.parse(x.setting_value)
    settings[x.setting_key] = x.setting_value
  }
  return settings
}

function generateShortId(length) {
  // 定义只包含k8s允许字符的字母表
  const k8sAlphabet = 'abcdefghijklmnopqrstuvwxyz0123456789'
  const generateK8sId = customAlphabet(k8sAlphabet, length) // 生成长度为length的ID
  return generateK8sId()
}

const getColor = (reason) => {
  switch(reason) {
    case 'Succeeded': return 'success'
    case 'Pending': 
    case 'Started': return 'warning'
    case 'Running': return 'primary'
    default: return 'danger'
  }
}
const getColor2 = (status) => {
  switch(status.status) {
    case "pending":
    case "initialize":
    case "waiting": return "warning"
    case "running": return "primary"
    case "terminated": return "danger"
    case "finished": {
      if(status.success) return "success"
      else return "danger"
    }
    case "completed": {
      if(status.results.decision === "true") return "success"
      else return "danger"
    }
  }
}
const getPodClass = (state, reason, ready) => {
  switch(state) {
    case 'running': {
      if(ready === 'True' || ready === true) 
        return 'text-green'
      return 'twinkling text-yellow'
    }
    case 'waiting':
    case 'deleting': return 'twinkling text-yellow'
    case 'terminated': {
      if(reason === 'Error')
        return 'text-red'
      if(reason === 'Completed')
        return 'text-gray'
    }
  }
}

const forPodList = (results) => {
  return results.map(x => {
    let conditionReady = x.status.conditions.find(n => {
      return n.type == 'Ready'
    })
    let m = {
      node_name: x.spec.nodeName,
      hostip: x.status.hostIP,
      podip: x.status.podIP,
      pod_name: x.metadata.name,
      creationTimestamp: x.metadata.creationTimestamp,
      status: {
        ready: conditionReady?.status||'False',
      },
    }
    let containerStatuses = x.status.containerStatuses?.map(y => {
      y.state_message = y.image
      y.status = Object.keys(y.state)[0]
      if(y.status!=='running') {
        y.state_message = y.state[y.status].reason
        y.reason = y.state[y.status].reason
      } else if(y.ready === false) { // 即使为running，也可能ready=False，需要显示ready为False的reason
        y.state_message = conditionReady.reason
      } 
      return y
    }) || x.spec.containers.map(y => {  // 无法调度的pod没有containerStatuses字段
      return {
        name: y.name,
        state_message: x.status.conditions[0].reason,
        reason: x.status.conditions[0].reason,
        status: 'waiting'
      }
    })
    let initContainerStatuses = x.status.initContainerStatuses?.map(y => {
      y.initContainer = true
      y.state_message = y.image
      y.status = Object.keys(y.state)[0]
      if(y.status!=='running') {
        y.state_message = y.state[y.status].reason
        y.reason = y.state[y.status].reason
      } else if(y.ready === false) { // 即使为running，也可能ready=False，需要显示ready为False的reason
        y.state_message = conditionReady.reason
      }
      return y
    }) || x.spec.initContainers?.map(y => {  // 无法调度的pod没有containerStatuses字段
      return {
        name: y.name,
        state_message: x.status.conditions[0].reason,
        reason: x.status.conditions[0].reason,
        status: 'waiting'
      }
    }) || []
    m.status.containerStatuses = containerStatuses.concat(initContainerStatuses)
    m.state = m.status.containerStatuses[0].status // running/waiting/terminated
    if(x.metadata.deletionTimestamp) m.state = 'deleting'
    m.reason = m.status.containerStatuses[0].reason
    if(m.state === 'running' && conditionReady.status === 'True')
      m.state_message = `Created ${moment.duration(moment(m.status.containerStatuses[0].state.running.startedAt)-moment()).humanize(true)}`
    else
      m.state_message = m.status.containerStatuses[0].state_message
    return m
  })
}

const parseLabels = (labels) => {
  labels ||= []
  return labels.map(x => {
    let arr = x.split('=')
    return {
      key: arr[0],
      value: arr[1]
    }
  })
}

export {
  getImageBase,
  getDuration,
  getSettings,
  generateShortId,
  getColor,
  getColor2,
  forPodList,
  getPodClass,
  parseLabels
}