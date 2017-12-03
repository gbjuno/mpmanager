$(function(){
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
        console.log('data', data)
        console.log('status', status)
        window.location.href="https://www.juntengshoes.cn/html/bindsuccess.html"
	});
    });
});