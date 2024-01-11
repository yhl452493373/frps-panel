var httpPort, httpsPort, pageOptions;
(function ($) {
    $(function () {
        function init() {
            var langLoading = layui.layer.load()
            $.getJSON('/lang.json').done(function (lang) {
                pageOptions = {
                    limitTemplet: function (item) {
                        return item + lang['PerPage'];
                    },
                    skipText: [lang['Goto'], '', lang['Confirm']],
                    countText: [lang['Total'], lang['Items']]
                };

                $.ajaxSetup({
                    error: function (xhr,) {
                        if (xhr.status === 401) {
                            layui.layer.msg(lang['TokenInvalid'], function () {
                                window.location.reload();
                            });
                        }
                    },
                });

                layui.element.on('nav(leftNav)', function (elem) {
                    var id = elem.attr('id');
                    var title = elem.text();
                    if (id === 'serverInfo') {
                        loadServerInfo(lang, title.trim());
                    } else if (id === 'userList') {
                        loadUserList(lang, title.trim());
                    } else if (elem.closest('.layui-nav-item').attr('id') === 'proxyList') {
                        if (id != null && id.trim() !== '') {
                            var suffix = elem.closest('.layui-nav-item').children('a').text().trim();
                            loadProxyInfo(lang, title + " " + suffix, id);
                        }
                    }
                });

                $('#leftNav .layui-this > a').click();
            }).always(function () {
                layui.layer.close(langLoading);
            });
        }

        function logout() {
            $.get("/logout", function (result) {
                window.location.reload();
            });
        }

        $(document).on('click.logout', '#logout', function () {
            logout();
        });

        init();
    });
})(layui.$);