// 忽略问号
export const ignoreQestionMark = function(value) {
  return value.includes('?')
}

// 根据字段配置的可通过字符忽略当前校验
export const ignoreFreeValue = function({ effectValidations, value }) {
  const freeRule = effectValidations?.freeValue?.rule

  if(!freeRule?.includes(value)) {
    return false
  }

  return true
}

export const ignore = function({ effectValidations, value }) {
  if(ignoreQestionMark(value)) return true
  
  return ignoreFreeValue({ effectValidations, value })
}