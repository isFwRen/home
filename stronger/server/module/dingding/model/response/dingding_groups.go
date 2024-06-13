/**
 * @Author: 阿先
 * @Description:
 * @Date: 2020/11/2 11:47
 */

package response

import (
	"server/module/dingding/model"
)

type DingdingGroupResponse struct {
	DingdingGroup model.DingdingGroups `json:"dingdingGroup"`
}
