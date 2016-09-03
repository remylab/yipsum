// Mock server 
var api = {
    running:{
        checkName:false,
        createIpsum:false,
        addQuote:false,
        editQuote:false,
        deleteQuote:false,
        generateIpsum:false,
        settingsAction:false
    },
    settingsAction:function(type, captchaResp, callback) {
        if (api.running.settingsAction)  { return; }
        console.log("mockserver : /api/s/:ipsum/resetkey");

        // ajax call will populate the res variable
        var a = [
            {ok:true,msg: ""}, // all good
            {ok:false,msg:"wrong_captcha"}, // wrong captcha !
            {ok:false,msg:"internal_error"} // Server error
        ];        

        var res = a[ Math.floor(Math.random()*(a.length)) ]

        api.running.settingsAction = true;

        res = {ok:true,msg: ""};
        
        setTimeout(function(){
            console.log('res : ' + res.ok + ", " + res.msg);
            api.running.settingsAction = false;
            callback(res);
        }, 600);

    },
    generateIpsum:function(callback) {
        if (api.running.generateIpsum)  { return; }
        console.log("mockserver : /api/:ipsum/generate");

        var res = [
        "all lies and jest, still, a man hears what he wants to hear and disregards the rest",
        "all of us get lost in the darkness, dreamers learn to steer by the stars",
        "all you need is love, love. Love is all you need",
        "an honest man's pillow is his peace of mind",
        "and in the end, the love you take is equal to the love you make",
        "before you accuse me take a look at yourself",
        "bent out of shape from society's pliers, cares not to come up any higher, but rather get you down in the hole that he's in",
        "different strokes for different folks, and so on and so on and scooby dooby dooby",
        "don't ask me what I think of you, I might not give the answer that you want me to",
        "don't you draw the Queen of Diamonds, boy, she'll beat you if she's able. You know, the Queen of Hearts is always your best bet",
        "every new beginning comes from some other beginning's end",
        "even the genius asks questions",
        "evil feeds off a source of apathy, weak in the mind, and of course you have to be. Less than a man, more like a thing, no knowledge you're nothing, knowledge is king",
        "fathers be good to your daughters, daughters will love like you do. Girls become lovers, who turn into mothers, so mothers be good to your daughters, too",
        "fear is the lock and laughter the key to your heart",
        "for what is a man, what has he got? If not himself, then he has naught. To say the things he truly feels, and not the words of one who kneels",
        "freedom, well, that's just some people talking. Your prison is walking through this world all alone",
        "freedom's just another word for nothing left to lose. Nothing ain't nothing, but it's free",
        "get your head out of the mud, baby. Put flowers in the mud, baby",
        "heard ten thousand whispering and nobody listening. Heard one person starve, I heard many people laughing. Heard the song of a poet who died in the gutter",
        "hero not the handsome actor who plays a hero's role. Hero not the glamour girl who'd love to sell her soul",
        "how many ears must one man have before he can hear people cry? Yes, and how many deaths will it take 'til he knows that too many people have died?",
        "i can dig it, he can dig it, she can dig it, we can dig it, they can dig it, you can dig it, oh let's dig it. Can you dig it, baby?",
        "i don't need no money, fortune, or fame. I got all the riches baby, one man can claim"
        ];

        api.running.generateIpsum = true;
        setTimeout(function(){
            api.running.generateIpsum = false;
            callback(res);
        }, 600);
    },
    deleteQuote:function($e, callback) {
        if (api.running.deleteQuote)  { return; }
        console.log("mockserver : /api/s/:ipsum/deletetext");

        // ajax call will populate the res variable
        var a = [
            {ok:true,msg: ""}, // URI is available 
            {ok:false,msg:"forbidden"}, // Unknown user
            {ok:false,msg:"internal_error"}, // Server error
            {ok:false,msg:"missing_params"}, // frontend is broken
        ];        

        var res = a[ Math.floor(Math.random()*(a.length)) ]

        api.running.deleteQuote = true;
        setTimeout(function(){
            console.log('res : ' + res.ok + ", " + res.msg);
            api.running.deleteQuote = false;
            callback($e, res);
        }, 600);

    },
    editQuote:function($e, t1, t2, callback) {
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
            callback($e, t1, t2, res);
        }, 600);

    },
    addQuote:function($e, text, callback) {
        if (api.running.addQuote)  { return; }
        console.log("mockserver : /api/s/:ipsum/addtext");

        // ajax call will populate the res variable
        var a = [
            {ok:true,msg: "562"}, // URI is available 
            {ok:false,msg:"too_many"}, // max 1000 quotes by ipsum
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