package templates

const BIND = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>绑定企业</title>
    <meta content="width=device-width,initial-scale=1.0,maximum-scale=1.0,user-scalable=0" name="viewport"/>
    <link rel="stylesheet" href="/html/weui.css">
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
    <script type="text/javascript" src="https://res.wx.qq.com/open/js/jweixin-1.2.0.js"></script>
    <script src="https://res.wx.qq.com/open/libs/weuijs/1.0.0/weui.min.js"></script>
    <script src="https://{{ .Domain }}/html/zepto.min.js"></script>
    <script type="text/javascript">
    
        $(function(){
            /**
            * 初始化函数
            */
            initBind();
        });
        
        function initBind(){
            $("#submit").click(function(){
                var phoneVal = $("#phone").val();
                var passwordVal = $("#password").val();
                console.log(phoneVal);
                console.log(passwordVal);
                $.post("/backend/confirm",{
                    phone:phoneVal,
                    password:passwordVal
                    },
                    function(data,status) {
                        var jdata = JSON.parse(data);
                        if (jdata.status) {
                            alert(jdata.message)
                            window.location.href="https://{{ .Domain }}/html/bindsuccess.html"
                        } else {
                            alert(jdata.message)
                        }
                    }
                );
            });
        }
    </script>
</body>
`

const PHOTO = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>拍照上传</title>
    <meta content="width=device-width,initial-scale=1.0,maximum-scale=1.0,user-scalable=0" name="viewport"/>
	<link rel="stylesheet" href="/html/weui.css">
    <link rel="stylesheet" href="/html/my.css">
</head>
<body>
    <div class="page article js_show">
        <div class="weui-form-preview__bd">
            <div class="weui-form-preview__item">
                <label class="weui-form-preview__label">当前拍照地点</label>
                <span class="weui-form-preview__value">{{ .PlaceName }}</span>
            </div>
        </div>
        <div class="weui-article" style="padding-bottom:0px">
            <p>
                <img id="previewImg" src="/html/pic_article.png" alt="">
            </p>
        </div>
        <div class="weui-btn-area" style="margin-top:0px">
            <a class="weui-btn weui-btn_primary" href="javascript:" id="takephoto">选择图片</a>
            <a class="weui-btn weui-btn_primary" href="javascript:" id="upload">上传图片</a>
        </div>
        <div class="js_dialog" id="confirmDialog" style="display:none; opacity: 0;">
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
        <div class="js_dialog" id="forwardDialog" style="display:none; opacity: 0;">
            <div class="weui-mask"></div>
            <div class="weui-dialog weui-skin_android">
                <div class="weui-dialog__hd"><strong class="weui-dialog__title">操作超时</strong></div>
                <div class="weui-dialog__bd">
                    未在2分钟内完成拍照，请重新扫描二维码。
                </div>
                <div class="weui-dialog__ft">
                    <a href="javascript:;" class="weui-dialog__btn weui-dialog__btn_primary" id="forwardConfirm">确定</a>
                </div>
            </div>
        </div>
    </div>
    <script type="text/javascript" src="https://res.wx.qq.com/open/js/jweixin-1.2.0.js"></script>
    <script src="https://res.wx.qq.com/open/libs/weuijs/1.0.0/weui.min.js"></script>
    <script src="https://{{ .Domain }}/html/zepto.min.js"></script>
    <script src="https://{{ .Domain }}/html/my.js"></script>
    <script type="text/javascript">
        if(!wx){//验证是否存在微信的js组件
            alert("微信接口调用失败，请检查是否引入微信js！");
        }
        wx.config({
            debug: false,
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

        });

        wx.error(function(res){
            alert("wx init failed")
        });

        var _uploadImageId;

        $(function(){
            setTimeout(function(){
                $("#forwardDialog").fadeIn(200);
            }, 120000);
            initClickForward();
            initTakePhoto();
            initUploadDialog();
        });


        function initClickForward(){
            $("#forwardConfirm").click(function(){
                window.location.href="https://open.weixin.qq.com/connect/oauth2/authorize?appid={{ .Wxappid }}&redirect_uri=https%3A%2F%2F{{ .Domain }}%2Fbackend%2Fscanqrcode&response_type=code&scope=snsapi_base&state=scanqrcode#wechat_redirect"
            })
        }

        function initTakePhoto(){
            $("#takephoto").click(function(){
                wx.chooseImage({
                    count: 1,
                    sizeType: ['compressed'], 
                    sourceType: ['camera'], 
                    success: function(res) {
                        var localIds = res.localIds; 
                        $("#previewImg").attr('src', localIds[0]);
                        _uploadImageId = localIds[0];
                    }
                });
            });

            $("#previewImg").click(function(){
                wx.previewImage({
                    current: this.src, // 当前显示图片的http链接
                    urls: [this.src] // 需要预览的图片http链接列表
                });
            });
        }

        function initUploadDialog(){
            var $confirmDialog = $('#confirmDialog');

            $("#uploadConfirm").click(function(){
                
                wx.uploadImage({
                    localId: _uploadImageId, // 需要上传的图片的本地ID，由chooseImage接口获得
                    isShowProgressTips: 1, // 默认为1，显示进度提示
                    success: function(res) {
                        var serverId = res.serverId;
                        $.post("/backend/download", {
                                corrective: {{ .Corrective }},
                                userId: {{ .Userid }},
                                placeId: {{ .Placeid }},
                                serverId: res.serverId
                            },
                            function(data, status) {
                                var jdata = JSON.parse(data)
                                if(jdata.status) {
                                   alert(jdata.message)
                                   window.location.href="https://open.weixin.qq.com/connect/oauth2/authorize?appid={{ .Wxappid }}&redirect_uri=https%3A%2F%2F{{ .Domain }}%2Fbackend%2Fscanqrcode&response_type=code&scope=snsapi_base&state=scanqrcode#wechat_redirect" 
                                } else {
                                   alert(jdata.message) 
                               }
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

const SCANQRCODE = `
<html>
<head>
    <meta charset="UTF-8">
    <title>扫描二维码</title>
    <meta content="width=device-width,initial-scale=1.0,maximum-scale=1.0,user-scalable=0" name="viewport"/>
    <link rel="stylesheet" href="/html/weui.css">
</head>
<body>
    <div class="page preview js_show">
    <div class="page__bd">
        <div class="weui-form-preview">
            <div class="weui-form-preview__hd">
                <div class="weui-form-preview__item">
                    <label class="weui-form-preview__label">拍照用户</label>
                    <em class="weui-form-preview__value">{{ .User }}</em>
                </div>
            </div>
            <div class="weui-form-preview__bd">
                <div class="weui-form-preview__item">
                    <label class="weui-form-preview__label">手机号码</label>
                    <span class="weui-form-preview__value">{{ .Phone }}</span>
                </div>
                <div class="weui-form-preview__item">
                    <label class="weui-form-preview__label">所属企业</label>
                    <span class="weui-form-preview__value">{{ .Company }}</span>
                </div>
            </div>
        </div>
        <div class="weui-btn-area" >
            <a class="weui-btn weui-btn_primary" href="javascript:" id="scanqrcode">扫描地点二维码</a>
        </div>
    </div>
    </div>

    <script type="text/javascript" src="https://res.wx.qq.com/open/js/jweixin-1.2.0.js"></script>
    <script src="https://res.wx.qq.com/open/libs/weuijs/1.0.0/weui.min.js"></script>
    <script src="https://{{ .Domain }}/html/zepto.min.js"></script>


    <script type="text/javascript">
    if(!wx){//验证是否存在微信的js组件
        alert("微信接口调用失败，请检查是否引入微信js！");
    }
    wx.config({
        debug: false,
        appId: '{{ .Wxappid }}',
        timestamp: {{ .Timestamp }},
        nonceStr: '{{ .Noncestr }}',
        signature: '{{ .Signature }}',
        jsApiList: [
            "scanQRCode",
        ]
    });

    wx.ready(function(){
        initQrcode();
    });

    wx.error(function(res){
        alert("wx init failed")
    });

    $(function(){
        
    });

    function initQrcode(){
        $("#scanqrcode").click(function(){
            wx.scanQRCode({
                needResult: 1,
                scanType: ["qrCode"],
                success: function (res) {
                    var result = res.resultStr;
                    console.log(result);
                    window.location.href=result;
                }
            });
        });
    }
    </script>
</body>
`

const NOTICEPAGE = `
<html>
<head>
    <meta charset="UTF-8">
    <title>{{ .Title }}</title>
    <meta content="width=device-width,initial-scale=1.0,maximum-scale=1.0,user-scalable=0" name="viewport"/>
    <link rel="stylesheet" href="/html/weui.css">
</head>
<body>
    <div class="weui-msg">
        <div class="weui-msg__icon-area"><i class="weui-icon-{{ .Type }} weui-icon_msg"></i></div>
        <div class="weui-msg__text-area">
            <h2 class="weui-msg__title">{{ .Msgtitle }}</h2>
            <p class="weui-msg__desc">{{ .Msgbody }}</p>
        </div>
        <div class="weui-msg__extra-area">
            <div class="weui-footer">
                <p class="weui-footer__text">Copyright © 2016-2017</p>
            </div>
        </div>
    </div>
    <script type="text/javascript" src="https://res.wx.qq.com/open/js/jweixin-1.2.0.js"></script>
    <script src="https://res.wx.qq.com/open/libs/weuijs/1.0.0/weui.min.js"></script>
    <script src="https://{{ .Domain }}/html/zepto.min.js"></script>
</body>

`

const AHREF = `<a href="{{ .URL }}">{{ .Content }}</a>`

const COMPANYSTAT = `
<html>
<head>
    <meta charset="UTF-8">
    <title>本日拍照进度</title>
    <meta content="width=device-width,initial-scale=1.0,maximum-scale=1.0,user-scalable=0" name="viewport"/>
    <link rel="stylesheet" href="/html/weui.css">
    <style>@-moz-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-webkit-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-o-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}embed,object{animation-duration:.001s;-ms-animation-duration:.001s;-moz-animation-duration:.001s;-webkit-animation-duration:.001s;-o-animation-duration:.001s;animation-name:nodeInserted;-ms-animation-name:nodeInserted;-moz-animation-name:nodeInserted;-webkit-animation-name:nodeInserted;-o-animation-name:nodeInserted;}</style>
    <style>@-moz-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-webkit-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-o-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}embed,object{animation-duration:.001s;-ms-animation-duration:.001s;-moz-animation-duration:.001s;-webkit-animation-duration:.001s;-o-animation-duration:.001s;animation-name:nodeInserted;-ms-animation-name:nodeInserted;-moz-animation-name:nodeInserted;-webkit-animation-name:nodeInserted;-o-animation-name:nodeInserted;}</style>
    <style>@-moz-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-webkit-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-o-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}embed,object{animation-duration:.001s;-ms-animation-duration:.001s;-moz-animation-duration:.001s;-webkit-animation-duration:.001s;-o-animation-duration:.001s;animation-name:nodeInserted;-ms-animation-name:nodeInserted;-moz-animation-name:nodeInserted;-webkit-animation-name:nodeInserted;-o-animation-name:nodeInserted;}</style>
    <style>@-moz-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-webkit-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-o-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}embed,object{animation-duration:.001s;-ms-animation-duration:.001s;-moz-animation-duration:.001s;-webkit-animation-duration:.001s;-o-animation-duration:.001s;animation-name:nodeInserted;-ms-animation-name:nodeInserted;-moz-animation-name:nodeInserted;-webkit-animation-name:nodeInserted;-o-animation-name:nodeInserted;}</style>
</head>
<body>
    <div class="weui-cell">
        <div class="weui-cell__bd">
            <p>公司</p>
        </div>
        <div class="weui-cell__ft">{{ .CompanyName }}</div>
    </div>

    <div class="weui-grids">
    {{ range .MonitorPlaceSummaryList }}
    {{ if eq .EverUpload "T" }}
        <div class="weui-gallery" id="gallery__{{ .MonitorPlaceID }}" style="opacity: 0; display: none;">
            <span class="weui-gallery__img" id="galleryImg__{{ .MonitorPlaceID }}" style="background-image:url(./backend/photolist?id={{ .MonitorPlaceID }})"></span>
        </div>
    {{ end }}
        <a href="" class="weui-grid">
            <div class="weui-grid__icon">
                {{ if eq .IsUpload "T" }}<i id="icon__{{ .MonitorPlaceID }}" class="weui-icon-success-no-circle"></i>{{ else }}<i class="weui-icon-warn"></i>{{ end }}
            </div>
            <p class="weui-grid__label">{{ .MonitorPlaceName }}</p>
        </a>
    {{ end }}</div> 
    <script type="text/javascript" src="https://res.wx.qq.com/open/js/jweixin-1.2.0.js"></script>
    <script src="https://res.wx.qq.com/open/libs/weuijs/1.0.0/weui.min.js"></script>
    <script src="https://{{ .Domain }}/html/zepto.min.js"></script>

    <script type="text/javascript">
        
    $(function(){
        var $everUploadClass = $(".weui-icon-success-no-circle"), $currentGallery;

        $everUploadClass.on("click", function(){
            var iconId = $(this).attr("id");
            var placeId = iconId.split("__")[1];
            $currentGallery = $("#gallery__" + placeId);
            $currentGallery.fadeIn(100);
        });
        $currentGallery.on("click", function(){
            $currentGallery.fadeOut(100);
        });
    });
    </script>
</body>
`

const PHOTOLIST = `
<html>
<head>
    <meta charset="UTF-8">
    <title>本日拍照具体地点</title>
    <meta content="width=device-width,initial-scale=1.0,maximum-scale=1.0,user-scalable=0" name="viewport"/>
    <link rel="stylesheet" href="/html/weui.css">
    <style>@-moz-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-webkit-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-o-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}embed,object{animation-duration:.001s;-ms-animation-duration:.001s;-moz-animation-duration:.001s;-webkit-animation-duration:.001s;-o-animation-duration:.001s;animation-name:nodeInserted;-ms-animation-name:nodeInserted;-moz-animation-name:nodeInserted;-webkit-animation-name:nodeInserted;-o-animation-name:nodeInserted;}</style>
    <style>@-moz-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-webkit-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-o-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}embed,object{animation-duration:.001s;-ms-animation-duration:.001s;-moz-animation-duration:.001s;-webkit-animation-duration:.001s;-o-animation-duration:.001s;animation-name:nodeInserted;-ms-animation-name:nodeInserted;-moz-animation-name:nodeInserted;-webkit-animation-name:nodeInserted;-o-animation-name:nodeInserted;}</style>
    <style>@-moz-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-webkit-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-o-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}embed,object{animation-duration:.001s;-ms-animation-duration:.001s;-moz-animation-duration:.001s;-webkit-animation-duration:.001s;-o-animation-duration:.001s;animation-name:nodeInserted;-ms-animation-name:nodeInserted;-moz-animation-name:nodeInserted;-webkit-animation-name:nodeInserted;-o-animation-name:nodeInserted;}</style>
    <style>@-moz-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-webkit-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-o-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}embed,object{animation-duration:.001s;-ms-animation-duration:.001s;-moz-animation-duration:.001s;-webkit-animation-duration:.001s;-o-animation-duration:.001s;animation-name:nodeInserted;-ms-animation-name:nodeInserted;-moz-animation-name:nodeInserted;-webkit-animation-name:nodeInserted;-o-animation-name:nodeInserted;}</style>
</head>
<body>
    <div class="weui-cell">
        <div class="weui-cell__bd">
            <p>{{ .MonitorPlaceName }}</p>
        </div>
        <div class="weui-cell__ft">{{ .CompanyName }}</div>
    </div>
    <div class="weui-gallery" style="display: block">
        <span class="weui-gallery__img">
             <img src="{{ .PictureURL }}" />
        </span>
        <div class="weui-gallery__opr">
             <a href="javascript:;" class="weui-btn weui-btn_default" id="showToast">{{ .JudgeComment }}</a>
        </div>
    </div>
    <script type="text/javascript" src="https://res.wx.qq.com/open/js/jweixin-1.2.0.js"></script>
    <script src="https://res.wx.qq.com/open/libs/weuijs/1.0.0/weui.min.js"></script>
    <script src="https://{{ .Domain }}/html/zepto.min.js"></script>
</body>
`

const SUBSCRIBE = `
<html>
<head>
    <meta charset="UTF-8">
    <title>本日拍照进度</title>
    <meta content="width=device-width,initial-scale=1.0,maximum-scale=1.0,user-scalable=0" name="viewport"/>
    <link rel="stylesheet" href="/html/weui.css">
    <style>@-moz-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-webkit-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-o-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}embed,object{animation-duration:.001s;-ms-animation-duration:.001s;-moz-animation-duration:.001s;-webkit-animation-duration:.001s;-o-animation-duration:.001s;animation-name:nodeInserted;-ms-animation-name:nodeInserted;-moz-animation-name:nodeInserted;-webkit-animation-name:nodeInserted;-o-animation-name:nodeInserted;}</style>
    <style>@-moz-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-webkit-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-o-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}embed,object{animation-duration:.001s;-ms-animation-duration:.001s;-moz-animation-duration:.001s;-webkit-animation-duration:.001s;-o-animation-duration:.001s;animation-name:nodeInserted;-ms-animation-name:nodeInserted;-moz-animation-name:nodeInserted;-webkit-animation-name:nodeInserted;-o-animation-name:nodeInserted;}</style>
    <style>@-moz-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-webkit-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-o-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}embed,object{animation-duration:.001s;-ms-animation-duration:.001s;-moz-animation-duration:.001s;-webkit-animation-duration:.001s;-o-animation-duration:.001s;animation-name:nodeInserted;-ms-animation-name:nodeInserted;-moz-animation-name:nodeInserted;-webkit-animation-name:nodeInserted;-o-animation-name:nodeInserted;}</style>
    <style>@-moz-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-webkit-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@-o-keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}@keyframes nodeInserted{from{opacity:0.99;}to{opacity:1;}}embed,object{animation-duration:.001s;-ms-animation-duration:.001s;-moz-animation-duration:.001s;-webkit-animation-duration:.001s;-o-animation-duration:.001s;animation-name:nodeInserted;-ms-animation-name:nodeInserted;-moz-animation-name:nodeInserted;-webkit-animation-name:nodeInserted;-o-animation-name:nodeInserted;}</style>
</head>
<body>
<div class="page article js_show">
    <div class="page__hd">
        <h1 class="page__title">Article</h1>
        <p class="page__desc"></p>
    </div>
    <div class="page__bd">
        <article class="weui-article">
            <h1>大标题</h1>
            <section>
                <h2 class="title">章标题</h2>
                <section>
                    <h3>1.1 节标题</h3>
                    <p>
                        Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod
                        tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam,
                        quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo
                        consequat.
                    </p>
                    <p>
                        <img src="./images/pic_article.png" alt="">
                        <img src="./images/pic_article.png" alt="">
                    </p>
                </section>
                <section>
                    <h3>1.2 节标题</h3>
                    <p>
                        Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod
                        tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam,
                        cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non
                        proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
                    </p>
                </section>
            </section>
        </article>
    </div>
    <div class="page__ft j_bottom">
        <a href="javascript:home()"><img src="./images/icon_footer_link.png"></a>
    </div>
</div>
    <script type="text/javascript" src="https://res.wx.qq.com/open/js/jweixin-1.2.0.js"></script>
    <script src="https://res.wx.qq.com/open/libs/weuijs/1.0.0/weui.min.js"></script>
    <script src="https://{{ .Domain }}/html/zepto.min.js"></script>
</body>
</html>
`

const TEMPLATEPAGE = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>{{ .Name }}</title>
    <meta content="width=device-width,initial-scale=1.0,maximum-scale=1.0,user-scalable=0" name="viewport"/>
    <link rel="stylesheet" href="/html/weui.css">
</head>
<body>
<div class="weui-panel weui-panel_access">
{{ range .HtmlChapters }}
<div class="weui-panel__bd">
    <a href="{{ .HtmlUrl }}" class="weui-media-box weui-media-box_appmsg">
        <div class="weui-media-box__hd">
            <img class="weui-media-box__thumb" src="{{ .PictureUrl }}" alt="">
        </div>
        <div class="weui-media-box__bd">
            <h4 class="weui-media-box__title">{{ .Title }}</h4>
            <p class="weui-media-box__desc">{{ .Digest }}</p>
        </div>
    </a>
</div>
{{ end }}
</div>
</body>
`
