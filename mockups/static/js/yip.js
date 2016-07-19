// Mock server api
var api = {
    running:{
        checkName:false,
        createIpsum:false,
    },
    checkName:function(callback){
        if (api.running.checkName)  {  return; }
        // ajax call will populate the res variable
        var a = [
            {ok:true,msg:""},
            {ok:false,msg:"internal_error"}
        ];

        var res = a[ Math.floor(Math.random()*(a.length)) ]

        api.running.checkName = true;
        setTimeout(function(){
            api.running.checkName = false;
            callback(res);
        }, 800);
    },
    createIpsum:function(callback){
        if (api.running.createIpsum)  {  return; }

        // ajax call will populate the res variable
        var a = [
            {ok:true,msg:""},
            {ok:false,msg:"taken"},
            {ok:false,msg:"internal_error"},
            {ok:false,msg:"missing_params",values:["email","name"]}
        ];

        var res = a[ Math.floor(Math.random()*(a.length)) ]

        api.running.createIpsum = true;
        setTimeout(function(){
            api.running.createIpsum = false;
            callback(res);
        }, 800);
    }
}


var CreateIpsum = (function() {
    "use strict";
    var 
    $form = $( "#createipsum" ),
    $name = $( "#name", $form ),
    $uri = $( "#uri", $form ),
    uri,
    uriLength = 32,
    // methods
    init,
    bind,
    bindUIActions,
    updateUri,
    onCheckNameResult,
    resetUriMessage,
    onCreateIpsumResult
    ;


    init = function(){
        bindUIActions();
    };

    bind = function() {
        if (!$CreateIpsum) {
            $CreateIpsum = $(CreateIpsum);
        }
        $CreateIpsum.bind.apply($CreateIpsum, arguments);
    };

    bindUIActions = function() {

        $name.keyup(function() {
            updateUri($(this).val())
        });

        $uri.focusin(function() {
            resetUriMessage()
        });

        $uri.focusout(function() {
            updateUri($(this).val());
            api.checkName(onCheckNameResult);
        });

        $form.submit(function( event ) {
            resetUriMessage()
            api.createIpsum(onCreateIpsumResult);
            event.preventDefault();
        });

    };

    resetUriMessage = function() {
        $uri.closest('.form-group').attr("class","form-group has-feedback");
        $( "#uri-feedback").attr("class","glyphicon form-control-feedback");
        $("#uri-help").html("");
        $("#yipurl").html(""); 
    }

    onCheckNameResult = function(res) {

        if (res.ok) {

            $uri.closest('.form-group').attr("class","form-group has-success has-feedback");
            $( "#uri-feedback").attr("class","glyphicon glyphicon-ok form-control-feedback");
            $("#yipurl").html("yipsum.com/"+uri);

        } else {
            $uri.closest('.form-group').attr("class","form-group has-error has-feedback");
            $( "#uri-feedback").attr("class","glyphicon glyphicon-remove form-control-feedback");
            $("#uri-help").html("Sorry this uri is already taken, please choose a new one");
            $("#yipurl").html("");
        }

    }

    onCreateIpsumResult = function(res){
        var $msg = $( "#messages" );

        if (res.ok) {


        } else {

            switch(res.msg) {
                case "taken":

                    $uri.closest('.form-group').attr("class","form-group has-error has-feedback");
                    $( "#uri-feedback").attr("class","glyphicon glyphicon-remove form-control-feedback");
                    $("#uri-help").html("Sorry this uri is already taken, please choose a new one");
                    $("#yipurl").html("");

                    break;

                case "missing_params":
                    break;
                default:  
                   
            }
        }
    }

    updateUri =  function(s) {
        uri = getSlug(s);

        if (uri.length > uriLength) {
            uri = uri.substring(0, uriLength);
        }
        $uri.val( uri );

    };

    return {
        init: init,
        bind: bind
    };

}());

CreateIpsum.init();


