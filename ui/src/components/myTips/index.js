import myTipsComponent from './myTips.vue'

const myTips = {
	install: (Vue) => {
		Vue.component('myTips', myTipsComponent)
	}
}

export default myTips
