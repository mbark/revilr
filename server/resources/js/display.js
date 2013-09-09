$(document).ready(function() {
	$('.revil-img').fancybox({
		helpers: {
			title : {
				type : 'float'
			}
		}
	});
});

$(".delete-revil").submit(function() {
	var form = $(this);
	form.ajaxSubmit({
		url: "/delete",
		type: "POST",
		success: function() {
			$("#submit-alert").addClass("alert-success");
			$("#submit-alert-text").text("The revil was successfuly removed!");
			$("#submit-alert").show();
			removeRevil(form);
		},
		error: function() {
			$("#submit-alert").addClass("alert-danger");
			$("#submit-alert-text").text("An internal error occurred when removing the revil");
			$("#submit-alert").show();
		}
	});

	return false;
});

function removeRevil(form) {
	id = form.find("input[name=id]").val();
	$("#" + id).remove();
}