// mem_data = ['name you want to show on y axis',a,a,a,a,a,a,a]; number of values you want to show in a time instant
var mem_data=["dummy", 0,0,0,0,0,0,0,0,0,0];
var new_value=5;
var mem_used = ["dummy", 50];
var new_mem_data=60;
/*
var line_chart = c3.generate({
    data: {
        columns: [
            mem_data
        ]
    }
});

// mem data is a list like [23,34,45]

//for loading the chart with new data

while(true){
    // delay(30)  set your delay function
	mem_data.push(new_value); // adds new value
    mem_data.splice(1, 1);
	chart.load({
		columns:[mem_data]
	});
}*/

// for guage-chart

var guage_chart = c3.generate({
    data: {
        columns: [
            mem_used
        ],
        type: 'gauge'
    },
    gauge: {
//        label: {
//            format: function(value, ratio) {
//                return value;
//            },
//            show: false // to turn off the min/max labels.
//        },
//    min: 0, // 0 is default, //can handle negative min e.g. vacuum / voltage / current flow / rate of change
//    max: 100, // 100 is default
//    units: ' %',
//    width: 39 // for adjusting arc thickness
    },
    color: {
        pattern: ['#FF0000', '#F97600', '#F6C600', '#60B044'], // the three color levels for the percentage values.
        threshold: {
//            unit: 'value', // percentage is default
//            max: 200, // 100 is default
              values: [30, 60, 90, 100]
        }
    },
    size: {
        height: 180
    }
});