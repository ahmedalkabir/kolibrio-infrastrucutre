$(function() {

    // Sidebar Toggler
    function sidebarToggle(toogle) {
        var sidebar = $('#sidebar');
        var padder = $('.content-padder');
        if (toogle) {
            $('.notyf').removeAttr('style');
            sidebar.css({ 'display': 'block', 'x': -300 });
            sidebar.transition({ opacity: 1, x: 0 }, 250, 'in-out', function() {
                sidebar.css('display', 'block');
            });
            if ($(window).width() > 960) {
                padder.transition({ marginLeft: sidebar.css('width') }, 250, 'in-out');
            }
        } else {
            $('.notyf').css({ width: '90%', margin: '0 auto', display: 'block', right: 0, left: 0 });
            sidebar.css({ 'display': 'block', 'x': '0px' });
            sidebar.transition({ x: -300, opacity: 0 }, 250, 'in-out', function() {
                sidebar.css('display', 'none');
            });
            padder.transition({ marginLeft: 0 }, 250, 'in-out');
        }
    }

    $('#sidebar_toggle').click(function() {
        var sidebar = $('#sidebar');
        var padder = $('.content-padder');
        if (sidebar.css('x') == '-300px' || sidebar.css('display') == 'none') {
            sidebarToggle(true)
        } else {
            sidebarToggle(false)
        }
    });

    function resize() {
        var sidebar = $('#sidebar');
        var padder = $('.content-padder');
        padder.removeAttr('style');
        if ($(window).width() < 960 && sidebar.css('display') == 'block') {
            sidebarToggle(false);
        } else if ($(window).width() > 960 && sidebar.css('display') == 'none') {
            sidebarToggle(true);
        }
    }

    if ($(window).width() < 960) {
        sidebarToggle(false);
    }

    $(window).resize(function() {
        resize()
    });

    $('.content-padder').click(function() {
        if ($(window).width() < 960) {
            sidebarToggle(false);
        }
    });

})

// it's  going to be used in sending event to the server
function listenSwitch(toggle_switch) {
    console.log('Name: ' + toggle_switch.Name + ' ' + toggle_switch.status());
}

// this fucntion used to establish gauges
function gaugeFactory(element, min, max) {
    let opts = {
        angle: 0.15, // The span of the gauge arc
        lineWidth: 0.44, // The line thickness
        radiusScale: 1, // Relative radius
        pointer: {
            length: 0.6, // // Relative to gauge radius
            strokeWidth: 0.035, // The thickness
            color: '#000000' // Fill color
        },
        limitMax: false, // If false, max value increases automatically if value > maxValue
        limitMin: false, // If true, the min value of the gauge will be fixed
        colorStart: '#6FADCF', // Colors
        colorStop: '#8FC0DA', // just experiment with them
        strokeColor: '#E0E0E0', // to see which ones work best for you
        generateGradient: true,
        highDpiSupport: true, // High resolution support

    };

    let gauge = new Gauge(document.getElementById(element)).setOptions(opts); // create sexy gauge!
    gauge.maxValue = max; // set max gauge value
    gauge.minValue = min; // Prefer setter over gauge.minValue = 0
    gauge.animationSpeed = 32; // set animation speed (32 is default value)
    gauge.setTextField(document.getElementById(element + '_text'));
    gauge.set(1);
    return gauge;
}

// this function used to establish toggle switches
function toggleSwitchFactory(element, func) {
    let toggle_switch = document.getElementById(element);
    let toggle_switch_obj = {
        Name: toggle_switch.id.toString(),
        isOn: (toggle_switch.checked ? true : false),
        isOff: (toggle_switch ? true : false),
        turnOn() {
            toggle_switch.checked = true
        },
        turnOff() {
            toggle_switch.checked = false
        },
        status() {
            return toggle_switch.checked ? 'ON' : 'OFF';
        }
    };

    toggle_switch.addEventListener('click', function() {
        listenSwitch(toggle_switch_obj);
    });

    return toggle_switch_obj;
}


// link gauge with html elements
function get_gauges() {
    let gauges = new Map();
    let gauge_widgets = document.getElementsByClassName('gauge-widget');
    let number_of_gauges = gauge_widgets.length;

    for (let i = 0; i < number_of_gauges; i++) {
        name_of_gauge = gauge_widgets[i].getElementsByClassName('gauge')[0].firstElementChild.id;
        gauges[name_of_gauge] = gaugeFactory(name_of_gauge, 0, 100);
        // gauges.push(gaugeFactory(name_of_gauge, 0, 100));
    }

    return gauges;
}

// to return list of available toggle switch
function get_toggle_switches() {
    let toggle_switches = new Map();
    let toggle_switch = document.getElementsByClassName('toggle-switch-wighet');
    let number_of_toggle_switches = toggle_switch.length;

    for (let i = 0; i < number_of_toggle_switches; i++) {
        name_of_toggle = toggle_switch[i].getElementsByClassName('switch')[0].firstElementChild.id;
        toggle_switches[name_of_toggle] = toggleSwitchFactory(name_of_toggle, listenSwitch);
    }

    return toggle_switches;
}

// this function is used to establish a chart
function chartFactory(element) {
    let options = {
        series: [{
            name: 'Temperature',
            data: [{
                x: '05/06/2014',
                y: 54
            }, {
                x: '05/08/2014',
                y: 17
            }]
        }],

        chart: {
            type: 'area',
            stacked: false,
            height: 350,
            zoom: {
                type: 'x',
                enabled: true,
                autoScaleYaxis: true
            },
            toolbar: {
                autoSelected: 'zoom'
            }
        },
        dataLabels: {
            enabled: false
        },

        markers: {
            size: 0,
        },


        fill: {
            type: 'gradient',
            gradient: {
                shadeIntensity: 1,
                inverseColors: false,
                opacityFrom: 0.5,
                opacityTo: 0,
                stops: [0, 90, 100]
            },
        },
        yaxis: {
            max: 100,
            labels: {
                formatter: function(val) {
                    return (val / 1).toFixed(0);
                },
            },
            title: {
                text: 'Temperature'
            },
        },

        xaxis: {
            type: 'datetime',
        },

        tooltip: {
            shared: false,
            y: {
                formatter: function(val) {
                    return (val / 1).toFixed(0)
                }
            }
        },
    }

    //
    let chart = new ApexCharts(document.getElementById(element), options);

    return chart;
}

function addDashboardToDB() {
    console.log("ADD DASHBOARD");
}

function deleteDashboardDB(dashboard) {
    console.log("DELETE DASHBOARD");
}

function deleteWidget(widget) {
    console.log("DELETE WIDGET");
}