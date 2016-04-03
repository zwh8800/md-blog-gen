(function ($) {
    $(function () {
        var script = '<script charset="utf-8" type="text/javascript" src="https://changyan.sohu.com/upload/changyan.js" ></script><script type="text/javascript">window.changyan.api.config({appid: "cyskgEYXj", conf: "prod_7cfa6ac36d43f0c53cd1cd5eb8ba5701"});</script>';
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
