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
	});
})(jQuery);
