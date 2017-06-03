$(document).ready(function() {
    $('.js-table').DataTable();
    new Blazy({
        selector: '.js-lazy'
    });
    if (barData) {
        $('.js-bar').each(function (idx, e) {
            new Chart(e, {
                type: 'bar',
                data: barData,
            });
        });
    };
});
