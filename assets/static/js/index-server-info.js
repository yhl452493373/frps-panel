var loadServerInfo = (function ($) {
    var size = filesize.partial({base: 2, standard: "jedec"});

    /**
     * get server info
     * @param lang {{}} language json
     */
    function loadServerInfo(lang) {
        $('#content').empty();
        var loading = layui.layer.load();

        $.ajax({
            url: 'http://127.0.0.1:7500/api/serverinfo',
            dataType: 'jsonp',
            success: function (result) {
                result = {
                    "version": "0.51.3",
                    "bind_port": 7000,
                    "vhost_http_port": 80,
                    "vhost_https_port": 443,
                    "tcpmux_httpconnect_port": 0,
                    "kcp_bind_port": 7000,
                    "quic_bind_port": 0,
                    "subdomain_host": "frp.yanghuanglin.com",
                    "max_pool_count": 100,
                    "max_ports_per_client": 0,
                    "heart_beat_timeout": 90,
                    "total_traffic_in": 1669491,
                    "total_traffic_out": 54422369,
                    "cur_conns": 0,
                    "client_counts": 1,
                    "proxy_type_count": {"http": 9, "https": 8, "tcp": 7}
                };
                renderServerInfo(result);
            },
            complete: function () {
                layui.layer.close(loading);
            }
        });
    }

    function renderServerInfo(data) {
        var html = layui.laytpl($('#serverInfoTemplate').html()).render(data);
        $('#content').html(html);

        renderTrafficChart(data);
        renderCountChart(data);
    }

    function renderTrafficChart(data) {
        var chartData = [
            {value: data.total_traffic_in, name: 'total_traffic_in'},
            {value: data.total_traffic_out, name: 'total_traffic_out'}
        ];
        var chartDom = document.getElementById('trafficChart');
        var chart = echarts.init(chartDom);
        var option = {
            title: {
                text: 'Network Traffic',
                subtext: 'today',
                left: 'center'
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
                data: ['total_traffic_in', 'total_traffic_out'],
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

    function renderCountChart(data) {
        var proxies = data.proxy_type_count;
        var charLegend = [];
        var chartData = [];

        for (var type in proxies) {
            var temp = {
                name: type,
                value: proxies[type]
            };
            charLegend.push(type);
            chartData.push(temp);
        }

        var chartDom = document.getElementById('countChart');
        var chart = echarts.init(chartDom);
        var option = {
            title: {
                text: 'Proxies',
                subtext: 'now',
                left: 'center'
            },
            tooltip: {
                trigger: 'item',
                formatter: function (v) {
                    return v.value;
                },
            },
            legend: {
                orient: 'vertical',
                left: 'left',
                data: charLegend,
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