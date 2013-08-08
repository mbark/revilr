var username = document.register.username;
var email = document.register.email;
var password = document.register.password;
var verification = document.register.verification;

$("#username").blur(isUsernameValid);

$("#email").blur(isEmailValid);

$("#password").blur(isPasswordValid);

$("#verification").blur(isVerificationValid);

$("#register").submit(function() {
	isValid = true;

	isValid = isValid && isUsernameValid();
	isValid = isValid && isEmailValid();
	isValid = isValid && isPasswordValid();
	isValid = isValid && isVerificationValid();

	return isValid;
});

function isUsernameValid() {
	var minLength = 5;
	var maxLength = 20;

	if(isIncorrectLength(username, 5, 20)) {
		$("#username-group").addClass("error");
		$("#username-error-text").text("Username have a length of between " + minLength + " and " + maxLength + " characters");
		$("#username-error").show();
		return false;
	} else if(isUsernameTaken(username)) {
		$("#username-group").addClass("error");
		$("#username-error-text").text("Username is already taken");
		$("#username-error").show();
		return false;
	} else {
		$("#username-group").removeClass("error");
		$("#username-error").hide();
		return true;
	}
}

function isEmailValid() {
	var regex = new RegExp("([a-z]*.)+@[a-z]+.[a-z]+");
	var value = email.value
	if (value == undefined) {
		$("#email-group").addClass("error");
		$("#email-error-text").text("Can't leave email field blank");
		$("#email-error").show();
		return false;
	} else if (!regex.test(email.value)) {
		$("#email-group").addClass("error");
		$("#email-error-text").text("Invalid email provided");
		$("#email-error").show();
		return false;
	} else if (isEmailTaken(email)) {
		$("#email-group").addClass("error");
		$("#email-error-text").text("An account already exists with the provded email");
		$("#email-error").show();
		return false;
	} else {
		$("#email-group").removeClass("error");
		$("#email-error").hide();
		return true;
	}
}

function isPasswordValid() {
	var minLength = 8;
	if(isIncorrectLength(password, 8, -1)) {
		$("#password-group").addClass("error");
		$("#password-error-text").text("Password must be " + minLength + " characters or longer");
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
		$("#verification-error-text").text("Passwords do not match");
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

function isEmailTaken(email) {
	var value = email.value;
	var isTaken = false;

	$.ajax({
		type: 'POST',
		url: "/email_taken",
		dataType: "json",
		data: {
			email: value
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