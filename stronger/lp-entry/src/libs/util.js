import Vue from 'vue'
import tools from './util.tools'

const bus = new Vue({})

const util = {
  tools,
  bus
}

export {
  tools,
  bus
}

export default util