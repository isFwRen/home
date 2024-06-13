import { localStorage, sessionStorage } from 'vue-rocket'
import PouchDB from 'pouchdb'
import store from '@/store'
import { tools as lpTools } from '@/libs/util'

const { baseURL: constBaseURL } = lpTools.constBaseURL()

// 获取lp2下的常量表信息
export const LP2 = 'lp2'
const LP2_localDB = new PouchDB(LP2)
const LP2_remoteDB = new PouchDB(`${constBaseURL}${LP2}`)

const LP2_CONSTANTS = 'lp2_constants'
const projectConstants = {}

// 记录每个常量表的信息
let [total, count] = [0, 0]

// 获取每个项目对应的常量表
const getProInfo = (docs) => {
  const proInfo = {
    total: 0,
    codes: []
  }
  const EXCLUDE_KEYS = ['total', 'codes']
  const proItems = localStorage.get('auth')?.proItems
  const proCode = sessionStorage.get('proCode')

  console.log(proItems)
  console.log(proCode)

  // 初始化
  proItems.map(item => {
    proInfo['codes'].push(item.value)

    proInfo[item.value] = {
      total: 0,
      docs: []
    }
  })

  console.log(proInfo)

  // 分类不同项目
  docs.map(doc => {
    if (proInfo['codes'].includes(doc.proCode)) {
      proInfo[doc.proCode].docs.push(doc)
    }
  })

  // 统计所有项目数量
  let num = 0

  for (let key in proInfo) {
    if (!EXCLUDE_KEYS.includes(key)) {
      const item = proInfo[key]
      item.total = item.docs.length
      num += item.total
    }
  }

  proInfo.total = num

  return proInfo
}

// 同步常量表lp2  
export const syncLP2 = async () => {
  return new Promise(resolve => {
    PouchDB.replicate(LP2_remoteDB, LP2_localDB, {
      style: 'main_only'
    })
      .on('change', (info) => {
        console.warn(`changed, ${LP2}`)

        if (info.ok) {
          localStorage.set(LP2_CONSTANTS, info.docs)
        }
        else {
          this.toasted.dynamic('同步常量表失败', 400)
        }
      })
      .on('complete', async (info) => {
        console.log(`%c completed, ${LP2}`, 'color: #4caf50')

        const result = await LP2_localDB.allDocs({
          attachments: true,
          include_docs: true,
          style: 'main_only'
        })

        const docs = []

        result.rows.map(row => {
          docs.push(row.doc)
        })

        const proInfo = getProInfo(docs)

        if (proInfo.total === 0) {
          resolve({
            code: 200,
            data: null
          })
        }

        localStorage.set(LP2_CONSTANTS, docs)
        total = proInfo.total
        store.commit('UPDATE_CONSTANTS', { total })

        count = 0

        if (info.ok) {
          resolve(syncLP2Constants())
        }
      })
      .on('error', (error) => {
        console.error(`同步常量表 ${LP2} 出错.`, error)
      })
  })
}

// 同步lp2常量表下的每个常量表
export const syncLP2Constants = async () => {
  return new Promise(resolve => {
    const mapPro = localStorage.get('auth')?.mapPro
    const constants = localStorage.get(LP2_CONSTANTS)

    for (let constant of constants) {
      if (!mapPro[constant.proCode]) {
        continue
      }

      const { chineseName, dbName } = constant

      projectConstants[constant.proCode] = {}

      const localDB = new PouchDB(chineseName)
      const remoteDB = new PouchDB(`${constBaseURL}${dbName}`)

      PouchDB.replicate(remoteDB, localDB, {
        style: 'main_only'
      })
        .on('change', () => {
          console.warn({ change: `${chineseName}，updating...` })
        })
        .on('complete', async () => {
          projectConstants[constant.proCode][chineseName] = {
            headers: [],
            desserts: []
          }

          const result = await remoteDB.allDocs({ attachments: true, include_docs: true })

          result?.rows.map(row => {
            if (/^TableTop::/.test(row.id)) {
              projectConstants[constant.proCode][chineseName].headers = row.doc?.tabletop
            }

            if (/^const::/.test(row.id)) {
              projectConstants[constant.proCode][chineseName].desserts.push(row.doc?.arr)
            }
          })

          ++count
          store.commit('UPDATE_CONSTANTS', { count })

          console.log({ total, count, complete: `${chineseName}，已同步到 IndexedDB.` })

          if (count === total) {
            resolve({
              code: 200,
              data: projectConstants
            })
          }
        })
        .on('error', (error) => {
          console.error(`${chineseName}，同步错误.`, error)
        })
    }
  })
}