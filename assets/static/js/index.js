var httpPort, httpsPort;
(function ($) {
    $(function () {
        function init() {
            var langLoading = layui.layer.load()
            $.getJSON('/lang.json').done(function (lang) {
                $.ajaxSetup({
                    error: function (xhr,) {
                        if (xhr.status === 401) {
                            layui.layer.msg(lang['TokenInvalid'], function () {
                                window.location.reload();
                            });
                        }
                    },
                })

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

var pageOptions = {
    layout: navigator.language.indexOf("zh") === -1 ? ['prev', 'page', 'next', 'skip', 'limit'] : ['prev', 'page', 'next', 'skip', 'count', 'limit'],
    limitTemplet: function (item) {
        if (navigator.language.indexOf("zh") === -1) {
            return item + ' / Page';
        }
        return item + ' 条/页';
    },
    skipText: navigator.language.indexOf("zh") === -1 ? ['Go to', '', 'Confirm'] : ['&#x5230;&#x7B2C;', '&#x9875;', '&#x786e;&#x5b9a;']
};