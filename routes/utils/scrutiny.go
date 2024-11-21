package utils

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"regexp"
	"stew/constants"
	"stew/types"
	"strconv"
)

func InputInvalidResponse(c *gin.Context) {
	c.AbortWithStatus(http.StatusBadRequest)
}

func ValidateIPv4(ipString string, allowEmpty bool, ctx *gin.Context) bool {
	if ipString != "" {
		ip := net.ParseIP(ipString)
		if ip == nil {
			return false
		}

		ip4 := ip.To4()
		if ip4 == nil {
			return false
		}

		if ip4[0] >= 224 || ip4[0] <= 0 {
			return false
		}

		if ip4.IsUnspecified() ||
			ip4.IsLinkLocalMulticast() ||
			ip4.IsInterfaceLocalMulticast() ||
			ip4.IsLinkLocalUnicast() ||
			ip4.IsMulticast() {
			return false
		}

		return true
	} else if allowEmpty {
		return true
	}
	return false
}

func ValidateIgn(name string, allowEmpty bool, ctx *gin.Context) bool {
	if name != "" {
		ignPattern := "^[a-zA-Z0-9_]{3,16}$"
		ignRe := regexp.MustCompile(ignPattern)
		return ignRe.MatchString(name)
	} else if allowEmpty {
		return true
	}
	return false
}

func ValidateID(id string, allowEmpty bool, ctx *gin.Context) bool {
	if id != "" {
		idNum, err1 := strconv.Atoi(id)
		if err1 != nil {
			return false
		} else if idNum <= 0 {
			return false
		} else {
			return true
		}
	} else if allowEmpty {
		return true
	}
	return false
}

func ValidateVersion(v string, allowEmpty bool, ctx *gin.Context) bool {
	if v != "" {
		verNum, err0 := strconv.Atoi(v)
		if err0 != nil {
			return false
		} else if !constants.IsKnownProtocolNumber(verNum) {
			return false
		} else {
			return true
		}
	} else if allowEmpty {
		return true
	}
	return false
}

func ValidateUUID(uuid string, allowEmpty bool, ctx *gin.Context) bool {
	if uuid != "" {
		uuidPattern := "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$"
		uuidRe := regexp.MustCompile(uuidPattern)
		return uuidRe.MatchString(uuid)
	} else if allowEmpty {
		return true
	}
	return false
}

func GetQueryData(field string, ctx *gin.Context) string {
	return ctx.Query(field)
}

func GetFormData(field string, ctx *gin.Context) string {
	return ctx.PostForm(field)
}

func ValidateAllData(fields []types.UnvalidatedField, ctx *gin.Context, allowAllEmpty bool) bool {
	allAllowEmpty := true
	allEmpty := true
	for _, field := range fields {
		if !field.Required {
			continue
		}
		if !field.AllowEmpty {
			allAllowEmpty = false
			break
		}
	}
	if allAllowEmpty && !allowAllEmpty {
		for _, field := range fields {
			if !field.Required {
				continue
			}
			f := field.Getter(field.Name, ctx)
			if f != "" {
				allEmpty = false
				break
			}
		}
		if allEmpty {
			InputInvalidResponse(ctx)
			return false
		}
	}

	for _, field := range fields {
		if !field.Required {
			continue
		}
		postData := field.Getter(field.Name, ctx)
		if field.Validator != nil {
			res := field.Validator(postData, field.AllowEmpty, ctx)
			if !res {
				InputInvalidResponse(ctx)
				return false
			}
		} else if field.AllowEmpty {
			continue
		}
	}
	return true
}

func ValidateData(field types.UnvalidatedField, ctx *gin.Context) bool {
	return ValidateAllData([]types.UnvalidatedField{field}, ctx, field.AllowEmpty)
}
