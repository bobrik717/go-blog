$(document).ready(function () {
    $('#content').on('keyup',function () {
        $.post("/gethtml",{md: $(this).val()}, function (res) {
            $("#md_html").html(res.html);
        });
    });
});