// Mock server 
var api = {
    running:{
        checkName:false,
        createIpsum:false,
        addQuote:false,
        editQuote:false
    },
    editQuote:function($e, text, callback) {
        if (api.running.editQuote)  { return; }
        console.log("mockserver : /api/s/:ipsum/updatetext");

        // ajax call will populate the res variable
        var a = [
            {ok:true,msg: ""}, // URI is available 
            {ok:false,msg:"forbidden"}, // Unknown user
            {ok:false,msg:"internal_error"}, // Server error
            {ok:false,msg:"missing_params"}, // frontend is broken
        ];        

        var res = a[ Math.floor(Math.random()*(a.length)) ]

        api.running.editQuote = true;
        setTimeout(function(){
            console.log('res : ' + res.ok + ", " + res.msg);
            api.running.editQuote = false;
            callback($e, text, res);
        }, 600);

    },
    addQuote:function(ipsum, text, $e, callback) {
        if (api.running.addQuote)  { return; }
        console.log("mockserver : /api/s/:ipsum/addtext");

        // ajax call will populate the res variable
        var a = [
            {ok:true,msg: ""}, // URI is available 
            {ok:false,msg:"forbidden"}, // Unknown user
            {ok:false,msg:"internal_error"}, // Server error
            {ok:false,msg:"missing_params"}, // frontend is broken
        ];        

        var res = a[ Math.floor(Math.random()*(a.length)) ]

        api.running.addQuote = true;
        setTimeout(function(){
            console.log('res : ' + res.ok + ", " + res.msg);
            api.running.addQuote = false;
            callback($e, res);
        }, 600);

    },
    checkName:function(uri, callback){
        if (api.running.checkName || api.running.createIpsum)  { return; }
        console.log("mockserver : /api/checkname");

        // ajax call will populate the res variable
        var a = [
            {ok:true,msg: uri}, // URI is available 
            {ok:false,msg:""}, // URI is already used
            {ok:false,msg:"internal_error"} // Server error
        ];

        var res = a[ Math.floor(Math.random()*(a.length)) ]

        api.running.checkName = true;
        setTimeout(function(){
            console.log('res : ' + res.ok + ", " + res.msg);
            api.running.checkName = false;
            callback(res);
        }, 600);
    },
    createIpsum:function($form, callback){
        if (api.running.createIpsum)  {  return; }
        console.log("mockserver : /api/createipsum");

        // ajax call will populate the res variable
        var a = [
            {ok:true,msg:"Lekjei9"}, // All good, create successful
            {ok:false,msg:"taken"}, // URI is already used
            {ok:false,msg:"missing_params",values:["email","name"]}, // Missing params
            {ok:false,msg:"internal_error"} // Server error
        ];

        var res = a[ Math.floor(Math.random()*(a.length)) ]

        api.running.createIpsum = true;
        setTimeout(function(){
            console.log('res : ' + res.ok + ", " + res.msg);
            api.running.createIpsum = false;
            callback(res);
        }, 600);
    }
}