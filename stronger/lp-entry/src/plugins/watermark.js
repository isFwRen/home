/**  水印添加方法  */

const setWatermark = (name) => {
  const id = 'watermark_lp2.0'

  if (document.getElementById(id) !== null) {
    document.body.removeChild(document.getElementById(id))
  }

  const can = document.createElement('canvas')
  // 设置canvas画布大小
  can.width = 190
  can.height = 170

  const cans = can.getContext('2d')
  cans.rotate(-20 * Math.PI / 180) // 水印旋转角度
  cans.font = '14px Vedana'
  cans.fillStyle = '#000000'
  cans.textAlign = 'center'
  cans.textBaseline = 'Middle'
  cans.fillText(name, can.width / 2, can.height) // 水印在画布的位置x，y轴

  const div = document.createElement('div')
  div.id = id
  div.style.pointerEvents = 'none'
  div.style.top = '0px'
  div.style.left = '0px'
  div.style.opacity = '0.15'
  div.style.position = 'fixed'
  div.style.zIndex = '100000'
  div.style.width = document.documentElement.clientWidth + 'px'
  div.style.height = document.documentElement.clientHeight  + 'px'
  div.style.background = 'url(' + can.toDataURL('image/png') + ') left top repeat'
  document.body.appendChild(div)
  return id
}

// 添加水印
export const newWaterMark = (name) => {
  let id = setWatermark(name)
  if (document.getElementById(id) === null) {
    id = setWatermark(name)
  }
}

// 移除水印
export const delWatermark = () => {
  const id = 'watermark_lp2.0'

  if (document.getElementById(id) !== null) {
    document.body.removeChild(document.getElementById(id))
  }
}
