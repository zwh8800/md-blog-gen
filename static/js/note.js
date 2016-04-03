(function ($) {
    $(function () {
        var script = '<script>(function(){var d=document,s=d.createElement("script");s.src="//lengzzz.disqus.com/embed.js";s.setAttribute("data-timestamp",+new Date());(d.head||d.body).appendChild(s)})();</script>';
        var $window = $(window);
        var $comment = $('#comment');
        var $body = $('body');
        var loaded = false;

        $window.scroll(function () {
            if (loaded) {
                return;
            }
            if ($window.scrollTop() + $window.height() > $comment[0].offsetTop) {
                $body.append(script);

                loaded = true;
            }
        })
    });
})(jQuery);
