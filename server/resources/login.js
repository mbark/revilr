$("#login").submit(function() {
	var name = document.login.username.value;
	var pass = document.login.password.value;

	var isValid = false;

	$.ajax({
		type: "POST",
		url: "/user_valid",
		dataType: "json",
		data: {
			username: name,
			password: pass
		},
		success: function(data) {
			isValid = data.isValid;
		},
		async: false
	});

	if(!isValid) {
		$("#password-group").addClass("error");
		$("#password-error-text").text("Invalid username or password");
		$("#password-error").show();
	} else {
		$("#password-group").removeClass("error");
		$("#password-error").hide();
	}

	return isValid;
});