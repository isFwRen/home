import service from '@/utils/request'
// @Summary 添加tabs
// @Produce  application/json
// @Param data body
// @Router /tabs/addTabs [post]
export const addTabs = (data) => {
  return service({
    url: "/tabs/addTabs",
    method: 'post',
    data: data
  })
}
// @Summary 根据name更新tabs
// @Produce  application/json
// @Param data body
// @Router /tabs/updateTabs [post]
export const updateTabs = (data) => {
  return service({
    url: "/tabs/updateTabs",
    method: 'post',
    data: data
  })
}

// @Summary 获取所有Tabs
// @Produce  application/json
// @Param data body
// @Router /tabs/getTabsList [post]
export const getTabsList = (data) => {
  return service({
    url: "/tabs/getTabsList",
    method: 'get',
    data: data
  })
}
// @Summary 获取最后一条Tabs
// @Produce  application/json
// @Param data body
// @Router /tabs/getTabsLast [post]
export const getTabsLast = (data) => {
  return service({
    url: "/tabs/getTabsLast",
    method: 'get',
    data: data
  })
}
// @Summary 根据ID删除Tabs
// @Produce  application/json
// @Param data body
// @Router /tabs/removeTabs [post]
export const removeTabs = (data) => {
  return service({
    url: "/tabs/removeTabs",
    method: 'post',
    data: data
  })
}

