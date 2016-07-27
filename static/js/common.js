(function ($) {
	$(function () {
		$('.search-form').submit(function () {
			$('.search-icon').hide();
			$('.search-loading').show();
		})
	});
})(jQuery);
