var loadServerInfo = (function ($) {
    var size = filesize.partial({base: 2, standard: "jedec"});

    /**
     * get server info
     * @param lang {{}} language json
     */
    function loadServerInfo(lang, title) {
        console.log(title)
        $("#title").text(title);
        $('#content').empty();
        var loading = layui.layer.load();

        $.getJSON('/proxy/api/serverinfo').done(function (result) {
            if (result.success) {
                renderServerInfo(JSON.parse(result.data));
            } else {
                layui.layer.msg(result.message);
            }
        }).always(function () {
            layui.layer.close(loading);
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