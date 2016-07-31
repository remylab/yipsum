// Server methods
var api = {
    running:{
        checkName:false,
        createIpsum:false,
        addQuote:false
    },
    addQuote:function(ipsum, text, $e, callback) {
        if (api.running.addQuote)  { return; }
        api.running.addQuote = true;

        var res = {ok:false,msg:"internal_error"}

        $.post( "/api/s/"+ipsum+"/addtext", {text:text,csrf:$('#csrf').val()} )
        .done(function(data) {
            callback($e, data);
        })
        .fail(function(data,statusText, xhr) {
            if ( xhr.status == 503 ) {
                res.msg = "forbidden"
            }
            callback($e, data);
        })
        .always(function() {
            api.running.addQuote = false;
        });

    },
    checkName:function(uri, callback){
        if (api.running.checkName || api.running.createIpsum)  {  return; }
        api.running.checkName = true;

        $.get( "/api/checkname", {uri:uri})
        .done(function(data) {
            callback(data);
        })
        .fail(function(data) {
            callback({ok:false,msg:"internal_error"});
        })
        .always(function() {
            api.running.checkName = false;
        });
    },
    createIpsum:function($form, callback){
        if (api.running.createIpsum)  {  return; }
        api.running.createIpsum = true;

        $.post( "/api/createipsum", $form.serialize() )
        .done(function(data) {
            callback(data);
        })
        .fail(function(data) {
            callback({ok:false,msg:"internal_error"});
        })
        .always(function() {
            api.running.createIpsum = false;
        });
    }
}