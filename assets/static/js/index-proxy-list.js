var loadProxyInfo = (function ($) {
    var size = filesize.partial({base: 2, standard: "jedec"});
    var i18n = {};

    /**
     * get proxy info
     * @param lang {{}} language json
     * @param title page title
     * @param proxyType proxy type
     */
    function loadProxyInfo(lang, title, proxyType) {
        i18n = lang;
        $("#title").text(title);
        $('#content').empty();
        var loading = layui.layer.load();

        $.getJSON('/proxy/api/proxy/' + proxyType).done(function (result) {
            if (result.success) {
                $('#content').html($('#proxyListTableTemplate').html());
                renderProxyListTable(JSON.parse(result.data)['proxies'], proxyType);
            } else {
                layui.layer.msg(result.message);
            }
        }).always(function () {
            layui.layer.close(loading);
        });
    }

    /**
     * render proxy list table
     * @param data proxy data
     * @param proxyType proxy type
     */
    function renderProxyListTable(data, proxyType) {
        data = data.concat(data)
        data = data.concat(data)
        var $section = $('#content > section');
        var cols = [
            {field: 'id', title: '', width: 30, templet: '#toggleProxyInfoArrowTemplate'},
            {field: 'name', title: 'Name', sort: true},
            {
                field: 'port', title: 'Port', sort: true, templet: '<span>{{ d.conf.remote_port }}</span>'
            },
            {field: 'cur_conns', title: 'Connections', sort: true},
            {
                field: 'today_traffic_in', title: 'Traffic In', sort: true, templet: function (d) {
                    return size(d.today_traffic_in);
                }
            },
            {
                field: 'today_traffic_out', title: 'Traffic Out', sort: true, templet: function (d) {
                    return size(d.today_traffic_out);
                }
            },
            {field: 'client_version', title: 'ClientVersion', sort: true},
            {field: 'status', title: 'Status', sort: true}
        ];
        var proxyListTable = layui.table.render({
            elem: '#proxyListTable',
            height: $section.height(),
            text: {none: i18n['EmptyData']},
            cols: [cols],
            page: navigator.language.indexOf("zh") !== -1,
            data: data,
            done: function (res, curr, count, origin) {
                //向每一行tr后面追加显示子table的tr
                var $tr = $('.layui-table-view[lay-id=' + this.id + '] tbody tr');
                var expandTrTemplateHtml = $('#expandTrTemplate').html();
                for (var i = 0; i < $tr.length; i++) {
                    var html = layui.laytpl(expandTrTemplateHtml).render({
                        index: i,
                        colspan: cols.length - 1,
                        proxyType: proxyType,
                        data: res.data[i]
                    });
                    $($tr[i]).after(html);
                }
            }
        });

        layui.table.on('tool(proxyListTable)', function (obj) {
            var index = obj.index;
            $(this).toggleClass('open');
            var open = $(this).hasClass('open');
            $('#childTr_' + index).toggleClass('layui-hide', !open);
            proxyListTable.resize();
        });

        window.onresize = function () {
            proxyListTable.resize();
        }
    }

    return loadProxyInfo;
})(layui.$);