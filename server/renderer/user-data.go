package renderer

import (
	"github.com/gin-gonic/gin"
)

type Data interface {
	Set(key string, value interface{}) Data
}

func NewData(c *gin.Context) Data {
	return &UserData{
		data:    map[string]interface{}{},
		context: c,
	}
}

type UserData struct {
	data    map[string]interface{}
	context *gin.Context
}

func (data *UserData) Set(key string, value interface{}) Data {
	data.data[key] = value
	return data
}
