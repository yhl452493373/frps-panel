(function ($) {
    $(function () {
        function login() {
            if (!layui.form.validate('#loginForm')) {
                return;
            }

            $.ajax({
                url: "/login",
                type: 'post',
                data: {
                    username: $('#username').val(),
                    password: $('#password').val()
                },
                success: function (result) {
                    if (result.success) {
                        window.location.reload();
                    } else {
                        layui.layer.msg(result.message);
                    }
                }
            });
        }

        $(document).on('click.login', '#login', function () {
            login();
        }).on('keydown', function (e) {
            if (e.keyCode === 13) {
                login();
            }
        });
    })
})(layui.$)