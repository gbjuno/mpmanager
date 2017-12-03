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
        window.location.href="./bindsuccess.html"
	});
    });
});