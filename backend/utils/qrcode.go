package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/golang/glog"
	"github.com/skip2/go-qrcode"
	"image"
	"os"
	"path"
)

var fontFile = "/opt/workspace/src/github.com/gbjuno/mpmanager/backend/templates/simhei.ttf"

func GenerateQrcodeImage(url string, comment string, savePath string) error {
	prefix := fmt.Sprintf("%s", "[QRCODE]")
	dir := path.Dir(savePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		errmsg := fmt.Sprintf("cannot create directory for path %s, err %s", savePath, err)
		glog.Errorf("%s %s", prefix, errmsg)
		return errors.New(errmsg)
	}

	var png []byte
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		errmsg := fmt.Sprintf("cannot generate qrcode for url %s, savePath %s, err %s", url, savePath, err)
		glog.Errorf("%s %s", prefix, errmsg)
		return errors.New(errmsg)
	}

	reader := bytes.NewReader(png)
	img, _, err := image.Decode(reader)
	if err != nil {
		errmsg := fmt.Sprintf("cannot generate qrcode for url %s, savePath %s, err %s", url, savePath, err)
		glog.Errorf("%s %s", prefix, errmsg)
		return errors.New(errmsg)
	}

	dc := gg.NewContext(256, 310)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)

	if err := dc.LoadFontFace(fontFile, 15); err != nil {
		errmsg := fmt.Sprintf("cannot generate qrcode for url %s, savePath %s, err %s", url, savePath, err)
		glog.Errorf("%s %s", prefix, errmsg)
		return errors.New(errmsg)
	}

	if len(comment) > 30 {
		comment_part1 := comment[:30]
		comment_part2 := comment[30:]
		dc.DrawStringAnchored(comment_part1, 125, 270, 0.5, 0.5)
		dc.DrawStringAnchored(comment_part2, 125, 290, 0.5, 0.5)
	} else {
		dc.DrawStringAnchored(comment, 125, 270, 0.5, 0.5)
	}

	dc.DrawImage(img, 0, 0)
	dc.Clip()
	dc.SavePNG(savePath)

	msg := fmt.Sprintf("generate qrcode for url %s, savePath %s successfully", url, savePath)
	glog.Infof("%s %s", prefix, msg)
	return nil
}
