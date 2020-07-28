$(document).ready(function () {
    // formatTable('#example')

    // setInterval(loadTable, 1000)

    summaryTable = $('#summary').DataTable({
        // "processing": true,
        // "serverSide": true,
        paging: false,
        // "sAjaxSource": "http://45.119.83.111:9081/api/data-table",
        ajax: {
            "url": "http://45.119.83.111:9081/api/data-table",
            dataSrc: function(rsp){
                for (var i = 0, iLen = rsp.aaData.length; i < iLen; i++){
                    for (var j = 0, jLen = rsp.aaData[i].length; j < jLen; j++){
                        rsp.aaData[i][j] = (parseInt(rsp.aaData[i][j]) || rsp.aaData[i][j]).toLocaleString() 
                    }
                }
                return rsp.aaData
            }
        }
    })

    $("select#duration").on('change', function () {
        loadChart()
    });

    $("select#interval").on('change', function () {
        loadChart()
    });

    loadChart()

    setInterval(function () { summaryTable.ajax.reload() }, 1000)
})

maxTickInterval = 30 //have to the same with server

loadChart = function () {
    Highcharts.chart('container', {

        chart: {
            scrollablePlotArea: {
                minWidth: 700
            }
        },

        data: {
            csvURL: `http://45.119.83.111:9081/api/histories?duration=${$("select#duration").val()}`,
            beforeParse: function (csv) {
                return csv.replace(/\n\n/g, '\n');
            },
            enablePolling: true,
            dataRefreshRate: $("select#duration").val()*60/maxTickInterval,
            switchRowsAndColumns: true
        },

        time: {
            timezoneOffset: -7 * 60
        },

        title: {
            text: "Apis' Interval Request Count"
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
        },

        // series: [{
        //     name: 'GetProfile',
        //     lineWidth: 4,
        //     marker: {
        //         radius: 4
        //     }
        // }, {
        //     name: 'LastIntervalUpdate'
        // }]
    });
}
