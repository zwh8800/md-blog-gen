(function ($) {
	$(function () {
		$('.search-form').submit(function () {
			NProgress.start();
			$('.search-icon').hide();
			$('.search-loading').show();
		})
	});
})(jQuery);
