(function ($) {
    $(function () {
        var script = '<script>(function(){var d=document,s=d.createElement("script");s.src="//lengzzz.disqus.com/embed.js";s.setAttribute("data-timestamp",+new Date());(d.head||d.body).appendChild(s)})();</script>';
        var $window = $(window);
        var $comment = $('#comment');
        var $body = $('body');
        var loaded = false;
        function checkAndLoad() {
            if (loaded) {
                return;
            }
            if ($window.scrollTop() + $window.height() > $comment[0].offsetTop) {
                $body.append(script);

                loaded = true;
            }
        }
        checkAndLoad();

        $window.scroll(checkAndLoad);

        var $qrcode = $('.qrcode');
        var $qrcodeBtn = $('.qrcode-button');
        $qrcodeBtn.click(function () {
            $qrcode.toggleClass('hidden');
        });
    });
    $(window).load(function () {
        setTimeout(function () {
            NProgress.configure({
                showSpinner: false,
                trickle: false,
                minimum: 0
            });

            var $window = $(window);
            var $document = $(document);
            $window.scroll(function () {
                var progress = $window.scrollTop() / ($document.height() - $window.height());
                NProgress.set(progress);
            });
        }, 100);
    })
})(jQuery);
