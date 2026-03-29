package pack

type CreateRequest struct {
	Size int `json:"size"`
}

type CreateBatchRequest struct {
	Sizes []int `json:"sizes"`
}
