var http_port, https_port;
(function ($) {
    $(function () {
        var langLoading = layui.layer.load()
        $.getJSON('/lang.json').done(function (lang) {
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
    });
})(layui.$);