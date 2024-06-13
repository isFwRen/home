import { getNode, getNodeValue } from './tools'

export const codesList46 = [
  ['fc138', 'fc139', 'fc140', 'fc141', 'fc247', 'fc215'],
  ['fc142', 'fc143', 'fc144', 'fc145', 'fc248', 'fc216'],
  ['fc146', 'fc147', 'fc148', 'fc149', 'fc249', 'fc217'],
  ['fc150', 'fc151', 'fc152', 'fc153', 'fc250', 'fc218'],
  ['fc154', 'fc155', 'fc156', 'fc157', 'fc251', 'fc219'],
  ['fc158', 'fc159', 'fc160', 'fc161', 'fc252', 'fc220'],
  ['fc162', 'fc163', 'fc164', 'fc165', 'fc253', 'fc221'],
  ['fc166', 'fc167', 'fc168', 'fc169', 'fc254', 'fc222']
]

export const codesList48 = [
  ['fc138', 'fc139', 'fc140', 'fc141'],
  ['fc142', 'fc143', 'fc144', 'fc145'],
  ['fc146', 'fc147', 'fc148', 'fc149'],
  ['fc150', 'fc151', 'fc152', 'fc153'],
  ['fc154', 'fc155', 'fc156', 'fc157'],
  ['fc158', 'fc159', 'fc160', 'fc161'],
  ['fc162', 'fc163', 'fc164', 'fc165'],
  ['fc166', 'fc167', 'fc168', 'fc169']
]

export const validate50CommonFn = ({ bill, fieldsObject }) => {
  const otherInfo = bill.otherInfo
  const xmls = getNode(otherInfo, 'applyCauses')
  let values = []

  for(let xml of xmls) {
    const vals = getNodeValue(xml, 'causeCode')
    values = [...values, ...vals]
  }

  const fc179Values = []

  for(let key in fieldsObject) {
    const sessionStorage = fieldsObject[key].sessionStorage
    const fieldsList = fieldsObject[key].fieldsList

    if(sessionStorage) {
      for(let fields of fieldsList) {
        for(let field of fields) { 
          const { code, resultValue } = field

          if(code === 'fc179') {
            resultValue && fc179Values.push(+resultValue)
          }
        }
      }
    }
  }

  return { values, fc179Values }
}

