var username = $("#username");
var email = $("#email");
var password = $("#password");
var verification = $("#verification");

username.blur(isUsernameValid);

email.blur(isEmailValid);

password.blur(isPasswordValid);

verification.blur(isVerificationValid);

$("#register").submit(function() {
	isValid = true;

	isValid = isValid && isUsernameValid();
	isValid = isValid && isEmailValid();
	isValid = isValid && isPasswordValid();
	isValid = isValid && isVerificationValid();

	if(isValid) {
		$(this).ajaxSubmit({
			url: "/register",
			type: "POST",
			success: function() {
				$("#submit-alert").addClass("alert-success");
				$("#submit-alert-text").text("A mail has been sent your email, follow instructions there to complete registration.");
				$("#submit-alert").show();
			},
			error: function() {
				$("#submit-alert").addClass("alert-danger");
				$("#submit-alert-text").text("Unable to register user.");
				$("#submit-alert").show();
			}
		});
	}

	return false;
});

function isUsernameValid() {
	var minLength = 5;
	var maxLength = 20;

	if(isIncorrectLength(username, 5, 20)) {
		$("#username-group").addClass("has-error");
		$("#username-error").text("Username have a length of between " + minLength + " and " + maxLength + " characters");
		return false;
	} else if(isUsernameTaken(username)) {
		$("#username-group").addClass("has-error");
		$("#username-error").text("Username is already taken");
		return false;
	} else {
		$("#username-group").removeClass("has-error");
		$("username-error").text("");
		return true;
	}
}

function isEmailValid() {
	var regex = new RegExp("([a-z]*.)+@[a-z]+.[a-z]+");
	var value = email.val();
	if (value == undefined) {
		$("#email-group").addClass("has-error");
		$("#email-error").text("Can't leave email field blank");
		return false;
	} else if (!regex.test(value)) {
		$("#email-group").addClass("has-error");
		$("#email-error").text("Invalid email provided");
		return false;
	} else if (isEmailTaken(email)) {
		$("#email-group").addClass("has-error");
		$("#email-error").text("An account already exists with that email");
		return false;
	} else {
		$("#email-group").removeClass("has-error");
		$("#email-error").text("");
		return true;
	}
}

function isPasswordValid() {
	var minLength = 8;
	if(isIncorrectLength(password, 8, -1)) {
		$("#password-group").addClass("has-error");
		$("#password-error").text("Password must be " + minLength + " characters or longer");
		return false;
	} else {
		$("#password-group").removeClass("has-error");
		$("#password-error").text("");
		return true;
	}
}

function isVerificationValid() {
	var minLength = 8;
	if(password.val() != verification.val()) {
		$("#verification-group").addClass("has-error");
		$("#verification-error").text("Passwords do not match");
		return false;
	} else {
		$("#verification-group").removeClass("has-error");
		$("#verification-error").text("");
		return true;
	}
}

function isIncorrectLength(username, min, max) {
	var length = 0;
	val = username.val()
	if(val != undefined) {
		length = val.length
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
	var name = username.val();
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
	var value = email.val();
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