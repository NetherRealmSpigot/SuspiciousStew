package types

import (
	"github.com/gin-gonic/gin"
)

type UnvalidatedField struct {
	Name       string
	Getter     UnvalidatedDataGetterFunction
	Validator  ValidatorFunction
	Required   bool
	AllowEmpty bool
}

type ValidatorFunction func(postData string, allowEmpty bool, ctx *gin.Context) bool

type UnvalidatedDataGetterFunction func(field string, ctx *gin.Context) string
