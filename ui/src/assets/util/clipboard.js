import Clipboard from 'clipboard'
import { ElMessage } from 'element-plus'

const copyDocument = () => {
  let clipboard = new Clipboard('.copy')
  clipboard.on('success', () => {
    ElMessage.success({message: '复制成功'})
    clipboard.destroy()
  })
  clipboard.on('error', () => {
    ElMessage.error({message: '该浏览器不支持自动复制'})
    clipboard.destroy()
  })
}

export {copyDocument}