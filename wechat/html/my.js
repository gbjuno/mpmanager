/**
 * 全局变量
 */
var _global_uploadImageId;

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
        $.post("/backend/confirm",{
		phone:phoneVal,
		password:passwordVal
	},
	function(data,status) {
        window.location.href="https://www.juntengshoes.cn/html/bindsuccess.html"
	});
    });
}
