package kvuser

import (
	"github.com/32leaves/ruruku/pkg/types"
	"strings"
)

const pathSeparator = "/"

func pathUser(name string) []byte {
	return []byte(strings.Join([]string{"u", name}, pathSeparator))
}

func pathUserToken(token string) []byte {
	return []byte(strings.Join([]string{"t", token}, pathSeparator))
}

func pathUserPermissions(name string) []byte {
    return []byte(strings.Join([]string{"r", name}, pathSeparator))
}
func pathUserPermission(name string, permission types.Permission) []byte {
	return []byte(strings.Join([]string{"r", name, string(permission)}, pathSeparator))
}
