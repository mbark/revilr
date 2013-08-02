function validateForm() {
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
		alert("Invalid username or password");
		password.focus();
	}

	return isValid;
}