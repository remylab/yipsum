// Server methods
var api = {
    running:{
        checkName:false,
        createIpsum:false,
    },
    checkName:function(uri, callback){
        if (api.running.checkName || api.running.createIpsum)  {  return; }

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