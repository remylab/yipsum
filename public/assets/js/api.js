// Server methods
function getKeyFromURL() {
    var x = window.location.href.split("?")[0];
    var y = x.split("/");
    return y[5];
}
function getIpsumFromURL() {
    var x = window.location.href.split("?")[0];
    var y = x.split("/");
    return y[3];
}
var api = {
    ipsumKey:getKeyFromURL(),
    ispumUri:getIpsumFromURL(),
    
    running:{
        checkName:false,
        createIpsum:false,
        addQuote:false,
        editQuote:false,
        deleteQuote:false
    },
    deleteQuote:function($e, callback) {
        if (api.running.deleteQuote)  { return; }

        var res = {ok:false,msg:"internal_error"}

        $.post( "/api/s/"+api.ispumUri+"/deletetext", {
            id:$e.attr('data-id'),
            csrf:$('#csrf').val(),
            key:api.ipsumKey
        })
        .done(function(data) {
            callback($e, data);
        })
        .fail(function(data, statusText, xhr) {
            if ( xhr == "Forbidden" ) {
                res.msg = "forbidden";
            }
            callback($e, res);
        })
        .always(function() {
            api.running.deleteQuote = false;
        });
    },
    editQuote:function($e, t1, t2, callback) {
        if (api.running.editQuote)  { return; }

        var res = {ok:false,msg:"internal_error"}

        $.post( "/api/s/"+api.ispumUri+"/updatetext", {
            id:$e.attr('data-id'),
            text:t2,
            csrf:$('#csrf').val(),
            key:api.ipsumKey
        })
        .done(function(data) {
            callback($e, t1, t2, data);
        })
        .fail(function(data, statusText, xhr) {
            if ( xhr == "Forbidden" ) {
                res.msg = "forbidden";
            }
            callback($e, t1, t2, res);
        })
        .always(function() {
            api.running.editQuote = false;
        });

    },
    addQuote:function($e, text, callback) {
        if (api.running.addQuote)  { return; }
        api.running.addQuote = true;

        var res = {ok:false, msg:"internal_error"}

        $.post( "/api/s/"+api.ispumUri+"/addtext", {
            text:text,
            csrf:$('#csrf').val(),
            key:api.ipsumKey
        })
        .done(function(data) {
            callback($e, data);
        })
        .fail(function(data,statusText, xhr) {
            if ( xhr == "Forbidden" ) {
                res.msg = "forbidden"
            }
            callback($e, res);
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