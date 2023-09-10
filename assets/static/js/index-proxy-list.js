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
        var $section = $('#content > section');
        var cols = [
            {field: 'id', type: 'space', width: 60, align: 'center', templet: '#toggleProxyInfoArrowTemplate'},
            {field: 'name', title: 'Name', sort: true},
            {
                field: 'port', title: 'Port', width: '12%', sort: true, templet: '<span>{{ d.conf.remote_port }}</span>'
            },
            {field: 'cur_conns', title: 'Connections', minWidth: 140, width: '12%', sort: true},
            {
                field: 'today_traffic_in',
                title: 'Traffic In',
                minWidth: 140,
                width: '12%',
                sort: true,
                templet: function (d) {
                    return size(d.today_traffic_in);
                }
            },
            {
                field: 'today_traffic_out',
                title: 'Traffic Out',
                minWidth: 140,
                width: '12%',
                sort: true,
                templet: function (d) {
                    return size(d.today_traffic_out);
                }
            },
            {field: 'client_version', title: 'ClientVersion', minWidth: 140, width: '12%', sort: true},
            {field: 'status', title: 'Status', width: '12%', sort: true}
        ];

        proxyType = proxyType.toLowerCase();
        if (proxyType === 'http' || proxyType === 'https') {
            data.forEach(function (temp) {
                if (proxyType === 'http') {
                    temp.conf.remote_port = http_port;
                } else if (proxyType === 'https') {
                    temp.conf.remote_port = https_port;
                }
            });
        }

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

    /**
     * load traffic statistics data
     */
    function loadTrafficStatistics() {
        var proxyName = $(this).closest('.layui-row').find('input').val();
        var loading = layui.layer.load();
        $.getJSON('/proxy/api/traffic/' + proxyName).done(function (result) {
            if (result.success) {
                renderTrafficChart(JSON.parse(result.data));
            } else {
                layui.layer.msg(result.message);
            }
        }).always(function () {
            layui.layer.close(loading);
        });
    }

    /**
     * render traffic statistics chart
     * @param data traffic data
     */
    function renderTrafficChart(data) {
        var html = layui.laytpl($('#trafficStaticTemplate').html()).render();
        var dates = [];
        var now = new Date();
        for (var i = 0; i < data.traffic_in.length; i++) {
            dates.push(now.getFullYear() + "/" + (now.getMonth() + 1) + "/" + now.getDate());
            now.setDate(now.getDate() - 1);
        }
        layui.layer.open({
            title: 'Traffic Statistics',
            type: 1,
            content: html,
            area: ['800px', '400px'],
            success: function () {
                var chartDom = document.getElementById('trafficBarChart');
                var chart = echarts.init(chartDom);
                var option = {
                    tooltip: {
                        trigger: 'axis',
                        axisPointer: {
                            type: 'shadow',
                        },
                        formatter: function (data) {
                            var html = ''
                            if (data.length > 0) {
                                html += data[0].name + '<br/>'
                            }
                            for (var v of data) {
                                var colorEl = '<span style="display:inline-block;margin-right:5px;border-radius:10px;width:9px;height:9px;background-color:' + v.color + '"></span>'
                                html += colorEl + v.seriesName + ': ' + size(v.value) + '<br/>'
                            }
                            return html
                        },
                    },
                    legend: {
                        data: ['Traffic In', 'Traffic Out'],
                    },
                    grid: {
                        left: '3%',
                        right: '4%',
                        bottom: '3%',
                        containLabel: true,
                    },
                    xAxis: [
                        {
                            type: 'category',
                            data: dates.reverse(),
                        },
                    ],
                    yAxis: [
                        {
                            type: 'value',
                            axisLabel: {
                                formatter: function (value) {
                                    return size(value)
                                },
                            },
                        },
                    ],
                    series: [
                        {
                            name: 'Traffic In',
                            type: 'bar',
                            data: data.traffic_in.reverse(),
                        },
                        {
                            name: 'Traffic Out',
                            type: 'bar',
                            data: data.traffic_out.reverse(),
                        },
                    ],
                };

                option && chart.setOption(option);
            }
        });
    }

    /**
     * document event
     */
    (function bindDocumentEvent() {
        $(document).on('click.trafficStatistics', '.traffic-statistics', function () {
            loadTrafficStatistics.call(this);
        });
    })();

    return loadProxyInfo;
})(layui.$);