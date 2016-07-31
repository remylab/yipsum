String.prototype.isEmpty = function() {
    return (this.length === 0 || !this.trim());
};

var util = {
    displayHelp:function(el, action = "reset", msg = ""){

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
                util.displayHelp($e,"error","This field is required");
                isValid = false
            } else {

                // ignore fields validated server side
                attr = $e.attr('serverval')
                if (typeof attr !== typeof undefined && attr !== false) { return; }

                util.displayHelp($e);
            }
        }) 
        return isValid;
    }
}

$(function() {
    // Form validation
    var $reqFields = $('input[required]','form.validate');

    $reqFields.focusin(function(){
        util.displayHelp($(this));
    })

    $reqFields.focusout(function(){
        $e = $(this);

        // ignore fields validated server side
        attr = $e.attr('serverval')
        if (typeof attr !== typeof undefined && attr !== false) { return; }

        if ( $e.val().isEmpty() ) {
            util.displayHelp($e,"error","This field is required");
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

        // Pre poluate URL field from name value
        $name.keyup(function() {
            updateUri($(this).val())
        });

        // Validate URL from DB
        $uri.focusin(function() {
            util.displayHelp($uri,false)
            $("#btn-yipurl").html(""); 
        });
        $uri.focusout(function() {
            updateUri($(this).val());

            if ( !uri.isEmpty() ) {
                api.checkName(uri, onCheckNameResult);
            } else {
                util.displayHelp($uri,"error","This field is required");
            }
        });

        // Create new Yipsum
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
            util.displayHelp($uri,"success");
            $("#btn-yipurl").html("yipsum.com/"+uri);
        } else {
            if ( res.msg == "internal_error") {
                util.displayHelp($uri,"error","Unexpected error, URL validation failed");
            } else {
                util.displayHelp($uri,"error","Sorry this URL is already taken, please choose a new one");
            }
            $("#btn-yipurl").html("");
        }
    }

    onCreateIpsumResult = function(res){
        var $msg = $( "#messages" );
        if (res.ok) {

            $('#yipurladm').html("yipsum.com/"+uri+"/adm/"+res.msg).attr('href',"/"+uri+"/adm/"+res.msg);
            $('#yipurl').html("yipsum.com/"+uri).attr('href',"/"+uri);

            var $yipin = $('#yipurladm-in');
            $yipin.val("http://yipsum.com/"+uri+"/adm/"+res.msg);
            $yipin.attr('size', $yipin.val().length);

            $('.helloyip').hide();
            $('#createsuccess').fadeIn();

        } else {

            $("#btn-yipurl").html("");
            switch(res.msg) {
                case "taken":
                    util.displayHelp($uri,"error","Sorry this URL is already taken, please choose a new one")
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


var Admin = (function() {
    "use strict";
    var 
    init,
    bindUIActions,
    onClickEdit,
    onClickAdd,
    onClickAddResult
    ;

    init = function(){
        bindUIActions();
    };

    bindUIActions = function() {

        $('.btn-edit').click(onClickEdit);

        $('.btn-add').click(function() {

            var $row = $(this).closest('.row-yiptext-add');
            var text = $('.yiptext-add textarea',$row).val().trim();

            if ( text.length == 0 ) { return; }

            var $e = $('<div class="row row-yiptext">'+
                '<div class="col-xs-10 col-yiptext">'+
                    '<div class="yiptext" style="display:none;">'+text+'</div>'+
                    '<div class="yiptext-edit">'+
                        '<textarea wrap="soft" maxlength="136">'+text+'</textarea>'+
                    '</div>'+
                '</div>'+
                '<div class="col-xs-2 col-edit">'+
                    '<button type="button" class="btn btn-default btn-edit glyphicon glyphicon-ok"></button>'+
                    '<span class="btn btn-default btn-saved glyphicon glyphicon-ok" style="display:none;"></span>'+
                    '<button type="button" class="btn btn-default glyphicon glyphicon-remove btn-delete" ></button>'+
                '</div>'+
            '</div>');

            onClickAdd($e,text);
        });
    };

    onClickAdd = function($e,text) {
        if (api.running.addQuote)  { return; }

        $('.yiptest-list').prepend($e);

        // register Edit event
        $('.btn-edit', $e).click(onClickEdit);

        api.addQuote( $('#ipsumId').val(), text, $e, onClickAddResult );

    };

    onClickAddResult = function($e, res) {
        if (res.ok) { return ; }

        if ( res.msg == "forbidden") {
            $('.col-edit', $e).hide();
            $('.yiptext-edit', $e).hide().after('<div>Your session has expired, please reload this page to continue </div>');
            $e.delay(5000).fadeOut();
        } else {
            $('.col-edit', $e).hide();
            $('.yiptext-edit', $e).hide().after('<div>Sorry, server error. Please try again later...</div>');
            $e.delay(2000).fadeOut();
        }

    };

    onClickEdit = function($e) {

        var $row = $(this).closest('.row-yiptext') 
        var t1 = $('.yiptext',$row).html().trim(), t2 = $('.yiptext-edit textarea',$row).val().trim();

        if ( t1 != t2 ) { 

            $('.yiptext',$row).html(t2); 

            $('.btn-edit',$row).hide();
            $('.btn-saved',$row).fadeIn(800).delay(1200).fadeOut(600).delay(800,function() {
                $('.btn-edit',$row).show();
            });
        }  
    };

    return {
        init: init
    };
}());

CreateIpsum.init();
Admin.init()