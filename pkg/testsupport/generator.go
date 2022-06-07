package testsupport

import (
	"fmt"
	mrand "math/rand"
	"strings"
)

func RandomAccount() string {
	return fmt.Sprintf("40903810%s", strings.Trim(strings.ReplaceAll(fmt.Sprint(mrand.Perm(11)), " ", ""), "[]"))
}
