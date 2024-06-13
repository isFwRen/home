export const getNode = (xml, nodeName) => {

  let node_arr, reg

  node_arr = []

  reg = new RegExp("\\<" + nodeName + "[\\s\\S]*?\\>[\\s\\S]*?\\<\\/" + nodeName + "\\>", "g")

  node_arr = xml.match(reg)

  return node_arr || []

}

export const getNodeValue = (xml, nodeName) => {

  let _data, data, i, len, node_arr, node_data, q, reg

  data = []

  node_arr = []

  reg = new RegExp("\\<" + nodeName + "[\\s\\S]*?\\/\\>");

  if(reg.test(xml)) {
    return ['']
  }

  reg = new RegExp("\\<" + nodeName + "[\\s\\S]*?\\>[\\s\\S]*?\\<\\/" + nodeName + "\\>", 'g')

  node_arr = xml.match(reg)

  if (node_arr !== null) {

    for (i = 0, len = node_arr.length; i < len; i++) {
      q = node_arr[i]

      node_data = q

      reg = new RegExp("(\\<" + nodeName + "[\\s\\S]*?\\>|\\<\\/" + nodeName + "\\>)", 'g')

      _data = node_data.replace(reg, '')

      data.push(_data)
    }
  }

  return data
}