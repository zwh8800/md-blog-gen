(function ($) {
    $(function () {
        $('.toushi-price').click(function (e) {
            $('.toushi-price').removeClass('selected');
            var $target = $(e.target);
            $target.addClass('selected');
            $('#price').val($target.data('price'));
        });
    });
})(jQuery);
