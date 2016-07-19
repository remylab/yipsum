// CREATE IPUSM FORM

$( "#name" ).keyup(function() {
	var v = $(this).val();
	$("#uri").val( getSlug(v) );
});

$( "#uri" ).focusout(function() {

})

$( "#createipsum" ).submit(function( event ) {
  
  event.preventDefault();
});