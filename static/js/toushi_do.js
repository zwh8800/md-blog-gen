(function ($) {
    $(function () {
        var cancel = setInterval(function () {
            $.get('/alipay/query', function (data) {
                console.log(data);
                handleResult(data);
            });
        }, 5500);
        var ws = new WebSocket("wss://lengzzz.com/alipay/ws");

        ws.onopen = function(evt) {
            console.log("Connection open ...");
        };

        ws.onmessage = function(evt) {
            console.log( "Received Message: " + evt.data);
            handleResult(JSON.parse(evt.data));
        };

        ws.onclose = function(evt) {
            console.log("Connection closed.");
        };

        function handleResult(data) {
            if (!data) {
                return
            }
            if (data.status === "SUCCESS" || data.status === "FAILED") {
                window.location.href = "/alipay/status";
                clearInterval(cancel);
            } else if (data.status === "ORDER_CREATED") {
                $('.toushi-scan-ok').show();
            }
        }
    });
})(jQuery);
