(function ($) {
    $(function () {
        $('.search-form').submit(function () {
            NProgress.start();
            $('.search-icon').hide();
            $('.search-loading').show();
        });
        $(window).on('beforeunload', function () {
            NProgress.configure({
                showSpinner: true,
                trickle: true,
                minimum: 0.08
            });
            NProgress.start();
            NProgress.set(0);
        });
        $("#keyword").autocomplete({
            source: function (request, response) {
                $.get("/api/search/" + encodeURIComponent(request.term), function (data) {
                    data.push({
                        id: 0,
                        value: '查看更多...'
                    });
                    response(data);
                });
            },
            select: function (event, ui) {
                event.preventDefault();
                var $keyword = $("#keyword");
                if (ui.item.id === 0) {
                    window.location.href = "/search/" + encodeURIComponent($keyword.val());
                    return;
                }
                $keyword.val(ui.item.value.replace(/<em>/g, '').replace(/<\/em>/g, ''));
                window.location.href = "/note/" + ui.item.id;
            },
            focus: function (event, ui) {
                event.preventDefault();
            },
            open: function (event, ui) {
                if ($(window).width() > 700) {
                    $(this).autocomplete("widget").css({
                        "width": 195
                    });
                }
            }

        }).data('ui-autocomplete')._renderItem = function (ul, item) {
            return $('<li class="ui-menu-item"></li>')
                .data("ui-autocomplete-item", item)
                .append(
                    $('<div class="ui-menu-item-wrapper"></div>').append(item.label)
                )
                .appendTo(ul);
        };
        if ($(window).width() <= 700) {
        	var duration = 1;
        	if (sessionStorage.getItem('scrolled') !== 'true') {
                duration = 500
            }
			$("html, body").animate({
				scrollTop: $('.search-box').outerHeight() + "px"
			}, {
				duration: duration,
				easing: "swing"
			});
            sessionStorage.setItem('scrolled', 'true');
        }
    });
})(jQuery);
