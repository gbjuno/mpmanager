package utils

import (
	"testing"
)

func Test_GenerateQrcode(t *testing.T) {
	if err := GenerateQrcodeImage("http://www.juntengshoes.cn", "数据库查询一下就知道原因了吧数据库查询一下就知道原因了吧", "/tmp/out3.png"); err != nil {
		t.Fatalf("Test_GenerateQrcode failed, err %s", err)
	}
}
