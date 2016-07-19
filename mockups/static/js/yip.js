String.prototype.isEmpty = function() {
    return (this.length === 0 || !this.trim());
};

var util = {
    displayFieldHelp:function(el, action = "reset", msg = ""){

        var $g = el.closest('.form-group')
        $i = $g.find('.glyphicon')
        $h = $g.find('.help-block');

        $g.removeClass("has-error has-success");
        $i.removeClass("glyphicon-remove glyphicon-ok");
        $h.html("");

        switch( action ) {
            case "error":
                $g.addClass("has-error");
                $i.addClass("glyphicon-remove");
                $h.html(msg);
                break;

            case "success":
                $g.addClass("has-success");
                $i.addClass("glyphicon-ok");
                break;
        } 
    },
    validateForm:function($form){
        var isValid = true;

        $('input[required]',$form).each(function(){
            $e = $(this);

            if ( $e.val().isEmpty()) {
                util.displayFieldHelp($e,"error","This field is required");
                isValid = false
            } else {

                // ignore fields validated server side
                attr = $e.attr('serverval')
                if (typeof attr !== typeof undefined && attr !== false) { return; }

                util.displayFieldHelp($e);
            }
        }) 
        return isValid;
    }
}

$(function() {
    // Form validation
    var $reqFields = $('input[required]','form.validate');

    $reqFields.focusin(function(){
        util.displayFieldHelp($(this));
    })

    $reqFields.focusout(function(){
        $e = $(this);

        // ignore fields validated server side
        attr = $e.attr('serverval')
        if (typeof attr !== typeof undefined && attr !== false) { return; }

        if ( $e.val().isEmpty() ) {
            util.displayFieldHelp($e,"error","This field is required");
        }
    })
});

var CreateIpsum = (function() {
    "use strict";
    var 
    $form = $( "#createipsum" ),
    $name = $( "#name", $form ),
    $uri = $( "#uri", $form ),
    uri,
    uriLength = $uri.attr('maxlength'),
    // methods
    init,
    bind,
    bindUIActions,
    updateUri,
    onCheckNameResult,
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
            util.displayFieldHelp($uri,false)
            $("#btn-yipurl").html(""); 
        });

        $uri.focusout(function() {
            updateUri($(this).val());

            if ( !uri.isEmpty() ) {
                api.checkName(uri, onCheckNameResult);
            } else {
                util.displayFieldHelp($uri,"error","This field is required");
            }
        });

        $form.submit(function( event ) {

            event.preventDefault();

            if ( util.validateForm($form) ) {
                api.createIpsum($form, onCreateIpsumResult);
            }

        });

    };

    onCheckNameResult = function(res) {

        if (res.ok) {
            uri = res.msg
            util.displayFieldHelp($uri,"success");
            $("#btn-yipurl").html("yipsum.com/"+uri);
        } else {
            if ( res.msg == "internal_error") {
               util.displayFieldHelp($uri,"error","Unexpected error, URL validation failed");
            } else {
                util.displayFieldHelp($uri,"error","Sorry this URL is already taken, please choose a new one");
            }
            $("#btn-yipurl").html("");
        }

    }

    onCreateIpsumResult = function(res){
        var $msg = $( "#messages" );
        if (res.ok) {

            $('.helloyip').hide();
            $('#yipurladm').html("yipsum.com/"+uri+"/adm/"+res.msg).attr('href',"/"+uri+"/adm/"+res.msg);
            $('#yipurl').html("yipsum.com/"+uri).attr('href',"/"+uri);
            $('#createsuccess').fadeIn();

        } else {

            $("#btn-yipurl").html("");
            switch(res.msg) {
                case "taken":
                    util.displayFieldHelp($uri,"error","Sorry this URL is already taken, please choose a new one")
                    break;

                case "missing_params":
                    $('#messages').html("Please check the required fields")
                    break;
                default:  
                    $('#messages').html("Sorry something unexpected happened, please come back later...")
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

// Mock server 
var api = {
    running:{
        checkName:false,
        createIpsum:false,
    },
    checkName:function(uri, callback){
        if (api.running.checkName || api.running.createIpsum)  {  return; }
        console.log("mockserver : /api/checkname");

        // ajax call will populate the res variable
        var a = [
            {ok:true,msg: uri},
            {ok:false,msg:"internal_error"}
        ];

        var res = a[ Math.floor(Math.random()*(a.length)) ]

        api.running.checkName = true;
        setTimeout(function(){
            api.running.checkName = false;
            callback(res);
        }, 800);
    },
    createIpsum:function($form, callback){
        if (api.running.createIpsum)  {  return; }
        console.log("mockserver : /api/createipsum");

        // ajax call will populate the res variable
        var a = [
            {ok:true,msg:"Lekjei9"},
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