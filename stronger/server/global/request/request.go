package request

type RmById struct {
	Ids []string `json:"ids" form:"ids"`
}