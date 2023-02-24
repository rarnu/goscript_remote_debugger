package server

type BaseResponse[T any] struct {
    Code int `json:"code"`
    Message string `json:"message"`
    Data T `json:"data"`
}
