package util_test

import (
	"testing"

	"github.com/ecator/gofile/util"
)

func TestMd5FromStr(t *testing.T) {
	t.Log(util.Md5FromStr("hello"))
}
