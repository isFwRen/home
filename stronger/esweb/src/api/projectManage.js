import service from '@/utils/request'
// @Summary 获取项目配置分页
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysProjectRecordSearch true "获取项目配置分页SysProjectRecordSearch"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /sysProject/getSysProject [post]
export const getSysProjectByPage = (data) => {
  return service({
    url: "/sysProject/getSysProjectByPage",
    method: 'post',
    data: data,
  })
}

// @Summary 获取项目配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /sysProject/getSysProjectList [post]
export const getSysProjectList = (data) => {
  return service({
    url: "/sysProject/getSysProjectList",
    method: 'post',
    data: data,
  })
}

// @Summary 新增项目配置
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body model.SysProject true "新增SysProject"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /sysProject/addSysProject [post]
export const AddSysProject = (data) => {
  return service({
    url: "/sysProject/AddSysProject",
    method: 'post',
    data: data,
  })
}
// @Summary 更新SysProject
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysProjectUpdateRecord true "更新SysProject"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /sysProject/updateSysProjectbyId [post]
export const updateSysProjectbyId = (data) => {
  return service({
    url: "/sysProject/updateSysProjectbyId",
    method: 'post',
    data: data,
  })
}
// @Summary 根据项目名称获取所有模板
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body string true "根据项目名称获取所有模板"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /sysProject/getSysProTempListByProName [get]
export const getSysProTempListByProName = (data) => {
  return service({
    url: "/sysProject/getSysProTempListByProName?name="+data,
    method: 'get',
    // data: data,
  })
}
// @Tags SysProTempBlock
// @Summary 新增项目模板配置
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body model.SysProTemplate true "新增SysProTemplate"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /sysProject/addSysProTemplate [post]
export const addSysProTemplate = (data) => {
  return service({
    url: "/sysProject/addSysProTemplate",
    method: 'post',
    data: data,
  })
}

// @Tags SysProTempBlock
// @Summary 根据id获取项目模板
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body string true "根据id获取项目模板"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /sysProject/getSysProTemplateById [get]
export const getSysProTemplateById = (data) => {
  return service({
    url: "/sysProject/getSysProTemplateById?id="+data,
    method: 'get',
    data: data,
  })
}
// @Tags SysProTempBlock
// @Summary 复制整个模板
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body string true "复制整个模板"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /sysProject/copyTemp [post]
export const copyTemp = (data) => {
  return service({
    url: "/sysProject/copyTemp",
    method: 'post',
    data: data,
  })
}
// @Tags SysProTempBlock
// @Summary 更新模板图片
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce  json
// @Param tempImgs formData file true "tempImgs"
// @Param sysProTempId formData string true "sysProTempId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /sysProject/updateImges [post]
export const updateImges = (data) => {
  return service({
    url: "/sysProject/updateImges",
    method: 'post',
    data: data,
  })
}
// @Tags SysProTempBlock
// @Summary 根据模板id获取项目模板分块
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body model.SysProTempB true "根据模板id获取项目模板分块"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /sysProject/GetSysProTempBlockByTempId [post]
export const getSysProTempBlockByTempId = (data) => {
  return service({
    url: "/sysProject/getSysProTempBlockByTempId",
    method: 'post',
    data: data,
  })
}
// @Tags SysProTempBlock
// @Summary 删除所有然后插入新的
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body request.SysProjectBlocks true "删除所有然后插入新的"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /sysProject/rmAllAndAddSysProTempBlockByTempId [post]
export const rmAllAndAddSysProTempBlockByTempId = (data) => {
  return service({
    url: "/sysProject/rmAllAndAddSysProTempBlockByTempId",
    method: 'post',
    data: data,
  })
}
// @Tags SysProTempBlock
// @Summary 删除所有然后插入新的分块关系
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body request.SysBlocksRelations true "删除所有然后插入新的分块关系"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /sysProject/rmAllAndAddBlockRelations [post]
export const rmAllAndAddBlockRelations = (data) => {
  return service({
    url: "/sysProject/rmAllAndAddBlockRelations",
    method: 'post',
    data: data,
  })
}

/**
 * @Tags SysProTempBlock
 * @Summary 获取所有field
 * @Auth xingqiyi
 * @Date 2020/11/4 3:29 下午
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data body string true "获取所有field"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
 * @Router /sysProject/getSysFieldsList [post]
 */
export const getSysFieldsList = (data) => {
  return service({
    url: "/sysProject/getSysFieldsList?proName="+data,
    method: 'get',
    data: data,
  })
}
/**
 * @Tags SysProTempBlock
 * @Summary 删除所有然后插入新的分块字段关系
 * @Auth xingqiyi
 * @Date 2020/11/5 3:13 下午
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data body request.CasbinInReceive true "删除所有然后插入新的分块字段关系"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
 * @Router /sysProject/rmAllAndAddBFRelations [post]
 */
export const rmAllAndAddBFRelations = (data) => {
  return service({
    url: "/sysProject/rmAllAndAddBFRelations",
    method: 'post',
    data: data,
  })
}

// @Tags SysProTempBlock
// @Summary 根据分块id获取分块关系
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body string true "根据分块id获取分块关系"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /sysProject/getBlockRelationsByBId [get]
export const getBlockRelationsByBId = (data) => {
  return service({
    url: "/sysProject/getBlockRelationsByBId?id="+data,
    method: 'get',
    data: data,
  })
}

/**
 * @Tags SysProTempBlock
 * @Summary 根据项目名称获取分页
 * @Auth xingqiyi
 * @Date 2020/12/22 下午3:25
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data body request.PageResult true "根据项目名称获取分页"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
 * @Router /sysProject/getExportAndNodesByProId [post]
 */
export const getExportAndNodesByProId = (data) => {
  return service({
    url: "/sysProject/getExportAndNodesByProId",
    method: 'post',
    data: data,
  })
}
/**
 * @Tags SysProTempBlock
 * @Summary 新增导出配置
 * @Auth xingqiyi
 * @Date 2020/12/22 下午3:05
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data body request.SysExport true "新增导出配置"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
 * @Router /sysProject/addExport [post]
 */
export const addExport = (data) => {
  return service({
    url: "/sysProject/addExport",
    method: 'post',
    data: data,
  })
}
/**
 * @Tags SysProTempBlock
 * @Summary 增加新的节点
 * @Auth xingqiyi
 * @Date 2020/11/12 10:08 上午
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data body model.AddExportNode true "增加新的节点"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"新增成功"}"
 * @Router /sysProject/addExportNode [post]
 */
export const addExportNode = (data) => {
  return service({
    url: "/sysProject/addExportNode",
    method: 'post',
    data: data,
  })
}
/**
 * @Tags SysProTempBlock
 * @Summary 根据id更新节点
 * @Auth xingqiyi
 * @Date 2020/12/22 下午3:05
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data body model.SysExport true "根据id更新节点"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
 * @Router /sysProject/updateExport [post]
 */
export const updateExport = (data) => {
  return service({
    url: "/sysProject/updateExport",
    method: 'post',
    data: data,
  })
}
/**
 * @Tags SysProTempBlock
 * @Summary 根据id更新节点
 * @Auth xingqiyi
 * @Date 2020/11/12 10:20 上午
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data body model.SysExportNode true "根据id更新节点"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
 * @Router /sysProject/updateExportNode [post]
 */
export const updateExportNode = (data) => {
  return service({
    url: "/sysProject/updateExportNode",
    method: 'post',
    data: data,
  })
}
/**
 * @Tags SysProTempBlock
 * @Summary 根据id删除
 * @Auth xingqiyi
 * @Date 2020/12/23 上午11:54
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data body request.IdsIntReq true "根据id删除"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
 * @Router /sysProject/deleByIds [post]
 */
export const deleByIds = (data) => {
  return service({
    url: "/sysProject/deleByIds",
    method: 'post',
    data: data,
  })
}
/**
 * @Tags SysProject
 * @Summary 将序号插到序号前
 * @Auth xingqiyi
 * @Date 2020/12/23 下午2:00
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data body request.InsertTo true "将序号插到序号前"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"插入成功"}"
 * @Router /sysProject/insertTo [post]
 */
export const insertTo = (data) => {
  return service({
    url: "/sysProject/insertTo",
    method: 'post',
    data: data,
  })
}
/**
 * @Tags SysProject
 * @Summary 根据exportId导出到Excel
 * @Auth xingqiyi
 * @Date 2020/12/28 上午11:20
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data body request.ExportId true "根据exportId导出到Excel"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
 * @Router /sysProject/ExportNodeByExportId [post]
 */
export const exportNodeByExportId = (data) => {
  return service({
    url: "/sysProject/exportNodeByExportId?exportId="+data,
    method: 'get',
    data: data,
  })
}

/**
 * @Tags SysProTempBlock
 * @Summary 增加审核配置
 * @Auth xingqiyi
 * @Date 2020/11/23 下午5:29
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data body request.SysInspection true "增加审核配置"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
 * @Router /sysProject/addInspection [post]
 */
export const addInspection = (data) => {
  return service({
    url: "/sysProject/addInspection",
    method: 'post',
    data: data,
  })
}
/**
 * @Tags SysProTempBlock
 * @Summary 根据id更新审核配置
 * @Auth xingqiyi
 * @Date 2020/11/23 下午5:49
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data body model.SysInspection true "根据id更新审核配置"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
 * @Router /sysProject/updateInspection [post]
 */
export const updateInspection = (data) => {
  return service({
    url: "/sysProject/updateInspection",
    method: 'post',
    data: data,
  })
}
/**
 * @Tags SysProTempBlock
 * @Summary 批量根据id删除
 * @Auth xingqiyi
 * @Date 2020/11/23 下午8:23
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data body request.IdsIntReq true "批量根据id删除"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
 * @Router /sysProject/rmSysInspectionByIds [post]
 */
export const rmSysInspectionByIds = (data) => {
  return service({
    url: "/sysProject/rmSysInspectionByIds",
    method: 'post',
    data: data,
  })
}
/**
 * @Tags SysProTempBlock
 * @Summary 分页获取审核配置
 * @Auth xingqiyi
 * @Date 2021/1/5 上午11:17
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data body request.SearchSysInspection true "分页获取审核配置"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
 * @Router /sysProject/getSysInspectionByPage [post]
 */
export const getSysInspectionByPage = (data) => {
  return service({
    url: "/sysProject/getSysInspectionByPage",
    method: 'post',
    data: data,
  })
}
/**
 * @Tags SysProject
 * @Summary 增加质检配置
 * @Auth xingqiyi
 * @Date 2021/1/4 下午5:19
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data body request.SysQuality true "增加质检配置"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"增加成功"}"
 * @Router /sysProject/addSysQuality [post]
 */
export const addSysQuality = (data) => {
  return service({
    url: "/sysProject/addSysQuality",
    method: 'post',
    data: data,
  })
}
/**
 * @Tags SysProTempBlock
 * @Summary 更新质检配置
 * @Auth xingqiyi
 * @Date 2021/1/5 上午11:00
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data body model.SysQuality true "更新质检配置"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
 * @Router /sysProject/updateSysQuality [post]
 */
export const updateSysQuality = (data) => {
  return service({
    url: "/sysProject/updateSysQuality",
    method: 'post',
    data: data,
  })
}
/**
 * @Tags SysProTempBlock
 * @Summary 批量根据id删除质检配置
 * @Auth xingqiyi
 * @Date 2021/1/5 上午9:35
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data body request.IdsIntReq true "批量根据id删除"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
 * @Router /sysProject/rmAddSysQualityByIds [post]
 */
export const rmAddSysQualityByIds = (data) => {
  return service({
    url: "/sysProject/rmAddSysQualityByIds",
    method: 'post',
    data: data,
  })
}
/**
 * @Tags SysProTempBlock
 * @Summary 分页获取质检配置
 * @Auth xingqiyi
 * @Date 2021/1/5 上午11:17
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data body request.SearchSysQuality true "分页获取质检配置"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
 * @Router /sysProject/getSysQualityByPage [post]
 */
export const getSysQualityByPage = (data) => {
  return service({
    url: "/sysProject/getSysQualityByPage",
    method: 'post',
    data: data,
  })
}
/**
 * @Tags SysProTempBlock
 * @Summary 根据类型获取质检配置
 * @Auth xingqiyi
 * @Date 2021/1/7 上午10:28
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data body request.ProId true "根据类型获取质检配置"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
 * @Router /sysProject/getSysQualityByType [post]
 */
export const getSysQualityByType = (data) => {
  return service({
    url: "/sysProject/getSysQualityByType",
    method: 'post',
    data: data,
  })
}
// @Tags SysProTempBlock
// @Summary 删除字段
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "删除字段"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysProject/rmSysProjectbyIds [post]
export const rmSysFieldsById = (data) => {
  return service({
    url: "/sysProject/rmSysFieldsById",
    method: 'post',
    data: data,
  })
}
// @Tags SysProTempBlock
// @Summary 获取项目字段分页
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysProFieldsSearch true "获取项目字段分页"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sysProject/getSysFieldsByPage [post]
export const getSysFieldsByPage = (data) => {
  return service({
    url: "/sysProject/getSysFieldsByPage",
    method: 'post',
    data: data,
  })
}
// @Tags SysProTempBlock
// @Summary 更新字段
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysProField true "更新字段"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sysProject/updateSysFieldsById [post]
export const updateSysFieldsById = (data) => {
  return service({
    url: "/sysProject/updateSysFieldsById",
    method: 'post',
    data: data,
  })
}
/**
 * @Tags casbin
 * @Summary 新增字段
 * @Auth xingqiyi
 * @Date 2020/11/3 4:11 下午
 * @Security ApiKeyAuth
 * @accept application/json
 * @Produce application/json
 * @Param data body request.CasbinInReceive true "新增字段"
 * @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
 * @Router /sysProject/addFields [post]
 */
export const addFields = (data) => {
  return service({
    url: "/sysProject/addFields",
    method: 'post',
    data: data,
  })
}
