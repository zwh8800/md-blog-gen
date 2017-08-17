(function ($) {
	$(function () {
		$('.search-form').submit(function () {
			NProgress.start();
			$('.search-icon').hide();
			$('.search-loading').show();
		});
		$(window).on('beforeunload', function(){
			NProgress.configure({
				showSpinner: true,
				trickle: true,
				minimum: 0.08
			});
			NProgress.start();
			NProgress.set(0);
		});
		$("#keyword").autocomplete({
			source: function(request, response) {
				$.get("/api/search/" + encodeURIComponent(request.term), function(data) {
					response(data);
				});
			},
			select: function(event, ui) {
				event.preventDefault();
                $("#keyword").val(ui.item.value.replace('<em>', '').replace('</em>', ''));
                window.location.href = "/note/" + ui.item.id;
			},
            focus: function(event, ui) {
                event.preventDefault();
                $("#keyword").val(ui.item.value.replace('<em>', '').replace('</em>', ''));
            },
            open: function(event, ui) {
                $(this).autocomplete("widget").css({
                    "width": 195
                });
            }

		}).data('ui-autocomplete')._renderItem = function( ul, item ) {
            return $('<li class="ui-menu-item"></li>')
                .data("ui-autocomplete-item", item)
				.append(
					$('<div class="ui-menu-item-wrapper"></div>').
					append(item.label)
				)
                .appendTo(ul);
        };
	});
})(jQuery);
