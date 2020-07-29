$(document).ready(function () {
    // formatTable('#example')

    // setInterval(loadTable, 1000)

    summaryTable = $('#summary').DataTable({
        // "processing": true,
        // "serverSide": true,
        paging: false,
        ajax: {
            url: "/api/data-table",
            dataSrc: function (rsp) {
                for (var i = 0, iLen = rsp.aaData.length; i < iLen; i++) {
                    for (var j = 0, jLen = rsp.aaData[i].length; j < jLen; j++) {
                        rsp.aaData[i][j] = (parseInt(rsp.aaData[i][j]) || rsp.aaData[i][j]).toLocaleString()
                    }
                }
                return rsp.aaData
            }
        }
    })

    // $("select.duration").on('change', function () {
    //     loadChart("request-rate", "http://45.119.83.111:9081/api/histories")
    // });
    loadProfilerInfo()
    loadChartWrapper("request-rate", window.location.origin + "/api/highchart/request-rate", "Apis' Request Rate (Req/s)")
    loadChartWrapper("process-rate", window.location.origin + "/api/highchart/process-rate", "Apis' Process Rate Per Routine (Req/s)")

    setInterval(function () { summaryTable.ajax.reload() }, 1000)
})

maxTickInterval = 30 //have the same value as server

loadChart = function (id, csvUrl, title) {
    wrapperId = `${id}-wrapper`
    Highcharts.chart(id, {

        chart: {
            scrollablePlotArea: {
                minWidth: 700
            }
        },

        data: {
            csvURL: `${csvUrl}?duration=${$(`#${wrapperId} select.duration`).val()}`,
            beforeParse: function (csv) {
                return csv.replace(/\n\n/g, '\n');
            },
            enablePolling: true,
            dataRefreshRate: $(`#${wrapperId} select.duration`).val() * 60 / maxTickInterval,
            switchRowsAndColumns: true
        },

        time: {
            timezoneOffset: -7 * 60
        },

        title: {
            text: title
        },

        xAxis: {
            type: 'datetime',
            label: {
                format: '%H:%M:%S.%L'
            }
        },

        tooltip: {
            shared: true,
            crosshairs: true
        }
    });
}

loadChartWrapper = function (id, csvUrl, title) {
    wrapperId = `${id}-wrapper`
    $(`#${wrapperId} select.duration`).on('change', function () {
        loadChart(id, csvUrl, title)
    })

    loadChart(id, csvUrl, title)
}

loadProfilerInfo = function () {
    $.ajax({
        url: "/api/profiler-info",
        success: function (result) {
            $("#profiler-name").html(`Profiler Name: ${result.name}`);
            $("#start-time").html(`Start Time: ${result.startTime}`);
        }
    });
}
