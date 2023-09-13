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
                        document.cookie = 'token=' + result.token + ';path=/'
                        window.location.href = "/"
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