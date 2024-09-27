package models

type GenericResponseMessage struct {
	Message string `json:"message"`
	Result  bool   `json:"result"`
}
