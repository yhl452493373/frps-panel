var loadServerInfo = (function ($) {
    var size = filesize.partial({base: 2, standard: "jedec"});
    var i18n = {};

    /**
     * get server info
     * @param lang {{}} language json
     * @param title page title
     */
    function loadServerInfo(lang, title) {
        i18n = lang;
        $("#title").text(title);
        $('#content').empty();
        var loading = layui.layer.load();

        $.getJSON('/proxy/api/serverinfo').done(function (result) {
            if (result.success) {
                var data = JSON.parse(result.data);
                data.proxy_counts = 0;
                http_port = data.vhost_http_port;
                https_port = data.vhost_https_port;
                for (var proxy in data.proxy_type_count) {
                    data.proxy_counts = data.proxy_counts + data.proxy_type_count[proxy];
                }
                data.bind_port = data.bind_port || i18n['Disable'];
                data.kcp_bind_port = data.kcp_bind_port || i18n['Disable'];
                data.quic_bind_port = data.quic_bind_port || i18n['Disable'];
                data.vhost_http_port = data.vhost_http_port || i18n['Disable'];
                data.vhost_https_port = data.vhost_https_port || i18n['Disable'];
                data.tcpmux_httpconnect_port = data.tcpmux_httpconnect_port || i18n['Disable'];
                data.subdomain_host = data.subdomain_host || i18n['NotSet'];
                data.max_pool_count = data.max_pool_count || i18n['NotSet'];
                data.max_ports_per_client = data.max_ports_per_client || i18n['NotLimit'];
                data.heart_beat_timeout = data.heart_beat_timeout || i18n['NotSet'];
                data.allow_ports_str = data.allow_ports_str || i18n['NotLimit'];
                data.tls_only = i18n[data.tls_only || false];
                renderServerInfo(data);
            } else {
                layui.layer.msg(result.message);
            }
        }).always(function () {
            layui.layer.close(loading);
        });
    }

    /**
     * render server info page
     * @param data server info data
     */
    function renderServerInfo(data) {
        var html = layui.laytpl($('#serverInfoTemplate').html()).render(data);
        $('#content').html(html);
        $('#frpVersion').text(data.version);

        renderTrafficChart(data);
        renderCountChart(data);
    }

    /**
     * render traffic chart with echarts
     * @param data traffic data
     */
    function renderTrafficChart(data) {
        var chartLegend = [i18n['TrafficIn'], i18n['TrafficOut']];
        var chartData = [
            {value: data.total_traffic_in, name: i18n['TrafficIn']},
            {value: data.total_traffic_out, name: i18n['TrafficOut']}
        ];
        var chartDom = document.getElementById('trafficPieChart');
        var chart = echarts.init(chartDom);
        var option = {
            title: {
                text: i18n['NetworkTraffic'],
                subtext: i18n['today'],
                left: 'center',
                textStyle: {
                    textBorderColor: '#fff',
                    textBorderWidth: 2
                },
                subtextStyle: {
                    textBorderColor: '#fff',
                    textBorderWidth: 2
                }
            },
            tooltip: {
                trigger: 'item',
                formatter: function (v) {
                    return size(v.value) + ' (' + v.percent + '%)';
                },
            },
            legend: {
                orient: 'vertical',
                left: 'left',
                data: chartLegend,
                textStyle: {
                    textBorderColor: '#fff',
                    textBorderWidth: 2
                }
            },
            series: [
                {
                    type: 'pie',
                    radius: '55%',
                    center: ['50%', '60%'],
                    data: chartData,
                    emphasis: {
                        itemStyle: {
                            shadowBlur: 10,
                            shadowOffsetX: 0,
                            shadowColor: 'rgba(0, 0, 0, 0.5)'
                        }
                    }
                }
            ]
        };

        option && chart.setOption(option);
    }

    /**
     * render proxy count chat with echarts
     * @param data proxy count data
     */
    function renderCountChart(data) {
        var proxies = data.proxy_type_count;
        var chartLegend = [];
        var chartData = [];

        for (var type in proxies) {
            var temp = {
                name: type.toUpperCase(),
                value: proxies[type]
            };
            chartLegend.push(type.toUpperCase());
            chartData.push(temp);
        }

        var chartDom = document.getElementById('countPieChart');
        var chart = echarts.init(chartDom);
        var option = {
            title: {
                text: i18n['Proxy'],
                subtext: i18n['now'],
                left: 'center',
                textStyle: {
                    textBorderColor: '#fff',
                    textBorderWidth: 2
                },
                subtextStyle: {
                    textBorderColor: '#fff',
                    textBorderWidth: 2
                }
            },
            tooltip: {
                trigger: 'item',
                formatter: function (v) {
                    return v.value + ' (' + v.percent + '%)';
                }
            },
            legend: {
                orient: 'vertical',
                left: 'left',
                data: chartLegend,
                textStyle: {
                    textBorderColor: '#fff',
                    textBorderWidth: 2
                }
            },
            series: [
                {
                    type: 'pie',
                    radius: '55%',
                    center: ['50%', '60%'],
                    data: chartData,
                    emphasis: {
                        itemStyle: {
                            shadowBlur: 10,
                            shadowOffsetX: 0,
                            shadowColor: 'rgba(0, 0, 0, 0.5)'
                        }
                    }
                }
            ]
        };

        option && chart.setOption(option);
    }

    return loadServerInfo;
})(layui.$);