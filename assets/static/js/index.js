(function ($) {
    $(function () {
        var langLoading = layui.layer.load()
        $.getJSON('/lang').done(function (lang) {
            layui.element.on('nav(leftNav)', function (elem) {
                if (elem.attr('id') === 'serverInfo') {
                    loadServerInfo(lang, elem.text().trim());
                } else if (elem.attr('id') === 'userList') {
                    loadUserList(lang, elem.text().trim());
                }
            });

            $('#leftNav .layui-this > a').click();
        }).always(function () {
            layui.layer.close(langLoading);
        });
    });
})(layui.$);
