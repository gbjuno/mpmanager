package templates

const BIND = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>绑定企业</title>
    <meta content="width=device-width,initial-scale=1.0,maximum-scale=1.0,user-scalable=0" name="viewport"/>
    <link rel="stylesheet" href="/html/weui.css">
    <style>@-moz-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-webkit-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-o-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}embed,object{animation-duration:.001s;-ms-animation-duration:.001s;-moz-animation-duration:.001s;-webkit-animation-duration:.001s;-o-animation-duration:.001s;animation-name:nodeInserted;-ms-animation-name:nodeInserted;-moz-animation-name:nodeInserted;-webkit-animation-name:nodeInserted;-o-animation-name:nodeInserted;}</style>
    <style>@-moz-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-webkit-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-o-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}embed,object{animation-duration:.001s;-ms-animation-duration:.001s;-moz-animation-duration:.001s;-webkit-animation-duration:.001s;-o-animation-duration:.001s;animation-name:nodeInserted;-ms-animation-name:nodeInserted;-moz-animation-name:nodeInserted;-webkit-animation-name:nodeInserted;-o-animation-name:nodeInserted;}</style>
    <style>@-moz-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-webkit-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-o-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}embed,object{animation-duration:.001s;-ms-animation-duration:.001s;-moz-animation-duration:.001s;-webkit-animation-duration:.001s;-o-animation-duration:.001s;animation-name:nodeInserted;-ms-animation-name:nodeInserted;-moz-animation-name:nodeInserted;-webkit-animation-name:nodeInserted;-o-animation-name:nodeInserted;}</style>
    <style>@-moz-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-webkit-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-o-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}embed,object{animation-duration:.001s;-ms-animation-duration:.001s;-moz-animation-duration:.001s;-webkit-animation-duration:.001s;-o-animation-duration:.001s;animation-name:nodeInserted;-ms-animation-name:nodeInserted;-moz-animation-name:nodeInserted;-webkit-animation-name:nodeInserted;-o-animation-name:nodeInserted;}</style>
</head>
<body>
    <div class="weui-cells weui-cells_form">
            <div class="weui-cell">
                <div class="weui-cell__hd"><label class="weui-label">手机号</label></div>
                <div class="weui-cell__bd">
                    <input id="phone" class="weui-input" type="number" pattern="[0-9]*" placeholder="请输入手机号">
                </div>
            </div>
            <div class="weui-cell">
                <div class="weui-cell__hd"><label class="weui-label">密码</label></div>
                <div class="weui-cell__bd">
                    <input id="password" class="weui-input" type="password" placeholder="请输入密码">
                </div>
            </div>
    </div>
    <div class="weui-btn-area">
            <a class="weui-btn weui-btn_primary" href="javascript:" id="submit">确定</a>
    </div>
    <script type="text/javascript" src="https://res.wx.qq.com/open/js/jweixin-1.0.0.js"></script>
    <script src="https://res.wx.qq.com/open/libs/weuijs/1.0.0/weui.min.js"></script>
    <script src="https://www.juntengshoes.cn/html/zepto.min.js"></script>
    <script src="https://www.juntengshoes.cn/html/my.js"></script>
</body>
`

const PHOTO = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>拍照上传</title>
    <meta content="width=device-width,initial-scale=1.0,maximum-scale=1.0,user-scalable=0" name="viewport"/>
	<link rel="stylesheet" href="/html/weui.css">
	<style>@-moz-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-webkit-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-o-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}embed,object{animation-duration:.001s;-ms-animation-duration:.001s;-moz-animation-duration:.001s;-webkit-animation-duration:.001s;-o-animation-duration:.001s;animation-name:nodeInserted;-ms-animation-name:nodeInserted;-moz-animation-name:nodeInserted;-webkit-animation-name:nodeInserted;-o-animation-name:nodeInserted;}</style>
	<style>@-moz-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-webkit-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-o-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}embed,object{animation-duration:.001s;-ms-animation-duration:.001s;-moz-animation-duration:.001s;-webkit-animation-duration:.001s;-o-animation-duration:.001s;animation-name:nodeInserted;-ms-animation-name:nodeInserted;-moz-animation-name:nodeInserted;-webkit-animation-name:nodeInserted;-o-animation-name:nodeInserted;}</style>
	<style>@-moz-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-webkit-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-o-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}embed,object{animation-duration:.001s;-ms-animation-duration:.001s;-moz-animation-duration:.001s;-webkit-animation-duration:.001s;-o-animation-duration:.001s;animation-name:nodeInserted;-ms-animation-name:nodeInserted;-moz-animation-name:nodeInserted;-webkit-animation-name:nodeInserted;-o-animation-name:nodeInserted;}</style>
	<style>@-moz-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-webkit-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-o-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}embed,object{animation-duration:.001s;-ms-animation-duration:.001s;-moz-animation-duration:.001s;-webkit-animation-duration:.001s;-o-animation-duration:.001s;animation-name:nodeInserted;-ms-animation-name:nodeInserted;-moz-animation-name:nodeInserted;-webkit-animation-name:nodeInserted;-o-animation-name:nodeInserted;}</style>
</head>
<body>
    <div class="page article js_show">
        <div class="weui-article" style="padding-bottom:0px">
            <p>
                <img id="previewImg" src="/html/pic_article.png" alt="">
            </p>
        </div>
        <div class="weui-btn-area" style="margin-top:0px">
            <a class="weui-btn weui-btn_primary" href="javascript:" id="takephoto">选择图片</a>
            <a class="weui-btn weui-btn_primary" href="javascript:" id="upload">上传图片</a>
        </div>
        <div class="js_dialog" id="confirmDialog" style="opacity: 0;">
            <div class="weui-mask"></div>
            <div class="weui-dialog weui-skin_android">
                <div class="weui-dialog__hd"><strong class="weui-dialog__title">确认上传</strong></div>
                <div class="weui-dialog__bd">
                    已经整改完成，确定上传图片。
                </div>
                <div class="weui-dialog__ft">
                    <a href="javascript:;" class="weui-dialog__btn weui-dialog__btn_default" id="cancel">取消</a>
                    <a href="javascript:;" class="weui-dialog__btn weui-dialog__btn_primary" id="uploadConfirm">确定</a>
                </div>
            </div>
        </div>
    </div>
    <script type="text/javascript" src="https://res.wx.qq.com/open/js/jweixin-1.2.0.js"></script>
    <script src="https://res.wx.qq.com/open/libs/weuijs/1.0.0/weui.min.js"></script>
    <script src="https://www.juntengshoes.cn/html/zepto.min.js"></script>
    <script src="https://www.juntengshoes.cn/html/my.js"></script>
    <script type="text/javascript">
        if(!wx){//验证是否存在微信的js组件
            alert("微信接口调用失败，请检查是否引入微信js！");
        }
        wx.config({
            debug: true,
            appId: '{{ .Wxappid }}',
            timestamp: {{ .Timestamp }},
            nonceStr: '{{ .Noncestr }}',
            signature: '{{ .Signature }}',
            jsApiList: [
                "chooseImage",
                "previewImage",
                "uploadImage"
            ]
        });

        wx.ready(function(){
            initTakePhoto();
            initUploadDialog();
        });

        wx.error(function(res){
            alert("wx init failed")
        });

        $(function(){
            
        });

        function initTakePhoto(){
            var $takephoto = $("#takephoto");
            
            $takephoto.click(function(){
                wx.chooseImage({
                    count: 1,
                    sizeType: ['original', 'compressed'], 
                    sourceType: ['album', 'camera'], 
                    success: function(res) {
                        var localIds = res.localIds; 
                        $("#previewImg").attr('src', localIds[0]);
                        _global_uploadImageId = localIds[0];
                    }
                });
            });

            $("#previewImg").click(function(){
                wx.previewImage({
                    current: this.src, // 当前显示图片的http链接
                    urls: [this.src] // 需要预览的图片http链接列表
                })
            });
        }

        function initUploadDialog(){
            var $confirmDialog = $('#confirmDialog');

            $("#uploadConfirm").click(function(){
                
                wx.uploadImage({
                    localId: _global_uploadImageId, // 需要上传的图片的本地ID，由chooseImage接口获得
                    isShowProgressTips: 1, // 默认为1，显示进度提示
                    success: function(res) {
                        var serverId = res.serverId;
                        alert(serverId);
                        $.post("/backend/download", {
                                serverId: res.serverId
                            },
                            function(data, status) {
                                console.log(status);
                                console.log(data);
                            }
                        );
                        $confirmDialog.fadeOut(200);
                    }
                });
            });

            $('#upload').on('click', function(){
                $confirmDialog.fadeIn(200);
            });

            $('#cancel').on('click', function(){
                $confirmDialog.fadeOut(200);
            });
        }
    </script>
</body>
`

const PICTURE = `
`
