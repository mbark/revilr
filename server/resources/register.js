function validateForm() {
	var username = document.register.username;
	var password1 = document.register.password;
	var password2 = document.register.password2;

	if(!isUsernameValid(username, 5, 12)) {
		return false;
	}
	if(!isUsernameFree(username)) {
		return false;
	}
	if(!isPasswordValid(password1, 8)) {
		return false;
	}
	if(!doPasswordsMatch(password1, password2)) {
		return false;
	}

	return true;
}

function isUsernameValid(username, min, max) {
	var length = 0;
	if(username.value != undefined) {
		length = username.value.length;
	}
	if (length == 0) {
		alert("Username can not be empty!");
	} else if(length > max || length < min) {
		alert("Username must between " + min + " and " + max + " in length");
	} else {
		return true;
	}
	username.focus();
	return false;
}

function isUsernameFree(username) {
	var name = username.value;
	var isTaken = false;

	$.ajax({
		type: 'POST',
		url: "/user_taken",
		dataType: "json",
		data: {
			username: name
		},
		success: function(data) {
			isTaken = data.isTaken;
		},
		async: false
	});

	if(isTaken) {
		alert("Username is taken");
		username.focus();
		return false;
	}

	return true;
}

function isPasswordValid(password, min) {
	var length = 0;
	if(password.value != undefined) {
		length = password.value.length;
	}

	if(length == 0) {
		alert("Password can not be blank!");
	} else if(length < min) {
		alert("Password must be longer than or equal to " + min)	
	} else {
		return true;
	}

	password.focus();
	return false;
}

function doPasswordsMatch(password1, password2) {
	if(password1.value != password2.value) {
		alert("Verification does not match first password!");
	} else {
		return true;
	}

	password2.focus();
	return false;
}