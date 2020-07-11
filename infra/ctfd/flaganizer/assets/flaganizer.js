if ($('#challenges-board').length) {
    $('<div id=flaganizer></div>').insertBefore('#challenges-board');
    $.get('plugins/flaganizer/assets/submitform.html', function (data) {
        $('#flaganizer').html(data);

        var flaganizer_submit = $('#flaganizer-submit');
        var flaganizer_flag = $('#flaganizer-flag');
        flaganizer_submit.unbind('click');
        flaganizer_submit.click(function (e) {
            e.preventDefault();

            $.post("/flaganizer-submit" , {
                submission: flaganizer_flag.val(),
                nonce: $('#nonce').val()
            }, function (data) {
                console.log(data);
                var result = $.parseJSON(JSON.stringify(data));

                var result_message = $('#flaganizer-message');
                var result_notification = $('#flaganizer-notification');
                result_notification.removeClass();
                result_message.text(result.data.message);

                if (result.data.status == "authentication_required") {
                    window.location = script_root + "/login?next=" + script_root + window.location.pathname + window.location.hash
                    return
                }
                else if (result.data.status == "incorrect") {
                    result_notification.addClass('alert alert-danger alert-dismissable text-center');
                    result_notification.slideDown();

                    flaganizer_flag.removeClass("correct");
                    flaganizer_flag.addClass("wrong");
                    setTimeout(function () {
                        flaganizer_flag.removeClass("wrong");
                    }, 3000);
                }
                else if (result.data.status == "correct") {
                    result_notification.addClass('alert alert-success alert-dismissable text-center');
                    result_notification.slideDown();

                    flaganizer_flag.val("");
                    flaganizer_flag.removeClass("wrong");
                    flaganizer_flag.addClass("correct");
                }
                else if (result.data.status == "already_solved") {
                    result_notification.addClass('alert alert-info alert-dismissable text-center');
                    result_notification.slideDown();

                    flaganizer_flag.addClass("correct");
                }

                setTimeout(function () {
                    $('.alert').slideUp();
                    flaganizer_submit.removeClass("disabled-button");
                    flaganizer_submit.prop('disabled', false);
                }, 4500);
            });
        });

        flaganizer_flag.keyup(function(event){
            if(event.keyCode == 13){
                flaganizer_submit.click();
            }
        });
    });
}
