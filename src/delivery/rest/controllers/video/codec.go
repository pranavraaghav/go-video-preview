package controllerVideo

type UploadResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Filename string `json:"filename"`
}
