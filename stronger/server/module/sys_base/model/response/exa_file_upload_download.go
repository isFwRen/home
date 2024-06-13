package response

import "server/module/sys_base/model"

type ExaFileResponse struct {
	File model.ExaFileUploadAndDownload `json:"file"`
}
