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
            {ok:true,msg: uri}, // URI is available 
            {ok:false,msg:""}, // URI is already used
            {ok:false,msg:"internal_error"} // Server error
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
            {ok:true,msg:"Lekjei9"}, // All good, create successful
            {ok:false,msg:"taken"}, // URI is already used
            {ok:false,msg:"missing_params",values:["email","name"]}, // Missing params
            {ok:false,msg:"internal_error"} // Server error
        ];

        var res = a[ Math.floor(Math.random()*(a.length)) ]

        api.running.createIpsum = true;
        setTimeout(function(){
            api.running.createIpsum = false;
            callback(res);
        }, 800);
    }
}