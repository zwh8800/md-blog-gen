(function ($) {
    $(function () {
        var cancel = setInterval(function () {
            $.get('/alipay/query', function (data) {
                console.log(data);
                handleResult(data);
            });
        }, 5500);
        var ws = new WebSocket("ws://localhost:3336/alipay/ws");

        ws.onopen = function(evt) {
            console.log("Connection open ...");
        };

        ws.onmessage = function(evt) {
            console.log( "Received Message: " + evt.data);
            handleResult(JSON.parse(evt.data));
            ws.close();
        };

        ws.onclose = function(evt) {
            console.log("Connection closed.");
        };

        function handleResult(data) {
            if (!data) {
                return
            }
            if (data.status === "SUCCESS") {
                $('.alipay').append('<p>投食成功</p>');
                clearInterval(cancel);
            } else if (data.status === "FAILED") {
                $('.alipay').append('<p>投食失败</p>');
                clearInterval(cancel);
            } else {

            }
        }
    });
})(jQuery);
