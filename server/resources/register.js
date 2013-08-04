var username = document.register.username;
var password = document.register.password;
var verification = document.register.verification;

$("#username").blur(isUsernameValid);

$("#password").blur(isPasswordValid);

$("#verification").blur(isVerificationValid);

$("#register").submit(function() {
	isValid = true;

	isValid = isValid && isUsernameValid();
	isValid = isValid && isPasswordValid();
	isValid = isValid && isVerificationValid();

	return isValid;
});

function isUsernameValid() {
	var minLength = 5;
	var maxLength = 20;

	if(isIncorrectLength(username, 5, 20)) {
		$("#username-group").addClass("error");
		$("#username-error").text("Username have a length of between " + minLength + " and " + maxLength + " characters");
		$("#username-error").show();
		return false;
	} else if(isUsernameTaken(username)) {
		$("#username-group").addClass("error");
		$("#username-error").text("Username is already taken");
		$("#username-error").show();
		return false;
	} else {
		$("#username-group").removeClass("error");
		$("#username-error").hide();
		return true;
	}
}

function isPasswordValid() {
	var minLength = 8;
	if(isIncorrectLength(password, 8, -1)) {
		$("#password-group").addClass("error");
		$("#password-error").text("Password must be " + minLength + " characters or longer");
		$("#password-error").show();
		return false;
	} else {
		$("#password-group").removeClass("error");
		$("#password-error").hide();
		return true;
	}
}

function isVerificationValid() {
	var minLength = 8;
	if(isNotSame(password, verification)) {
		$("#verification-group").addClass("error");
		$("#verification-error").text("Passwords do not match");
		$("#verification-error").show();
		return false;
	} else {
		$("#verification-group").removeClass("error");
		$("#verification-error").hide();
		return true;
	}
}

function isIncorrectLength(username, min, max) {
	var length = 0;
	if(username.value != undefined) {
		length = username.value.length;
	}
	var isCorrectLength = true;
	if(min >= 0) {
		isCorrectLength = isCorrectLength && length >= min;
	}
	if(max >= 0) {
		isCorrectLength = isCorrectLength && length <= max;
	}

	return !isCorrectLength;
}

function isUsernameTaken(username) {
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

	return isTaken;
}

function isNotSame(password1, password2) {
	return password1.value != password2.value;
}