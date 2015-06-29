package main

// contents of _templates/index.html
const index_html_str = `<html>
	<head>
		<style>
			body {
				font: 15px sans-serif;
			}
			#graph {
				font: 10px sans-serif;
			}

			.axis path,
			.axis line {
			  fill: none;
			  stroke: #000;
			  shape-rendering: crispEdges;
			}

			.area {
			  fill: steelblue;
			}
		</style>
		<script src="http://d3js.org/d3.v3.js"></script>
		<script language="JavaScript">
			
			var showGraph = function(serverData){
				var margin = {top: 20, right: 20, bottom: 30, left: 50},
				    width = 960 - margin.left - margin.right,
				    height = 500 - margin.top - margin.bottom;

				var parseDate = d3.time.format("%Y-%m-%dT%H:%M:%S%Z").parse;
				var parseISODate = d3.time.format.utc("%Y-%m-%dT%H:%M:%S.%LZ").parse;

				var x = d3.time.scale()
				    .range([0, width]);

				var y = d3.scale.linear()
				    .range([height, 0]);

				var xAxis = d3.svg.axis()
				    .scale(x)
				    .orient("bottom");

				var yAxis = d3.svg.axis()
				    .scale(y)
				    .orient("left");

				var area = d3.svg.area()
				    .x(function(d) { return x(d.date); })
				    .y0(height)
				    .y1(function(d) { return y(d.level); });

				var svg = d3.select("#graph").append("svg")
				    .attr("width", width + margin.left + margin.right)
				    .attr("height", height + margin.top + margin.bottom)
				  .append("g")
				    .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

				var data = serverData.WaterLevels.map(function(d) {
					// remove second fraction "2015-02-21T23:21:22.847943775-05:00" -> "2015-02-21T23:21:22-05:00"
					var timeStr = d.Time;
					var myDate;
					if( timeStr.charAt(timeStr.length-1) == 'Z' ) {
						myDate = parseISODate(timeStr.substring(0,23) + 'Z');
					} else {
						var dotPos = timeStr.lastIndexOf('.');
						if(dotPos > -1) {
							var hyphenPos = timeStr.lastIndexOf('-');
							if(hyphenPos > -1) {
								timeStr = timeStr.substring(0, dotPos) + timeStr.substring(hyphenPos);
							}
						}

						// remove hyphen from time zone "2015-02-21T23:21:22-05:00" -> "2015-02-21T23:21:22-0500"
						var lastIndex = timeStr.lastIndexOf(":")
						timeStr = timeStr.slice(0, lastIndex) + timeStr.slice(lastIndex+1);
						myDate = parseDate(timeStr);
					}
					// return our data point
					return {
						date: myDate,
						level: d.Level
					};
				});
				console.log(data);

			  x.domain(d3.extent(data, function(d) { return d.date; }));
			  y.domain([0, d3.max(data, function(d) { return d.level; })]);

			  svg.append("path")
			      .datum(data)
			      .attr("class", "area")
			      .attr("d", area);

			  svg.append("g")
			      .attr("class", "x axis")
			      .attr("transform", "translate(0," + height + ")")
			      .call(xAxis);

			  svg.append("g")
			      .attr("class", "y axis")
			      .call(yAxis)
			    .append("text")
			      .attr("transform", "rotate(-90)")
			      .attr("y", 6)
			      .attr("dy", ".71em")
			      .style("text-anchor", "end")
			      .text("Water Level (inches)");
			};

			var loadFunc = function(){
				var xmlhttp = new XMLHttpRequest();
				xmlhttp.onreadystatechange=function()
				{
					if (xmlhttp.readyState==4 && xmlhttp.status==200)
					{
						var jsonData = JSON.parse(xmlhttp.responseText);
						var waterLevels = jsonData.WaterLevels;
						if(waterLevels.length > 0){
							document.getElementById("currentDepth").innerHTML="Current depth: " + waterLevels[waterLevels.length -1].Level + " inches."
						}
						showGraph(jsonData);
					}
				}
				xmlhttp.open("GET","/info", true);
				xmlhttp.send();
			};
			window.onload = loadFunc;
			
		</script>
	</head>
	<body>
		<h1>Blake's Raspberry Pi Sump Water Level Monitor</h1>
		<div id="currentDepth" style="font: 15px sans-serif;"></div>
		<div id="graph"></div>
		<br>
		<div>View Memory Usage: <a href="/profiler/info.html">here</a></div>
		<br>
		<h3>Source Code</h3>
		<div>Sump: Raspberry Pi-based Sump Monitor - <a href="https://github.com/wblakecaldwell/sump">here</a></div>
		<div>Memory profiler - <a href="https://github.com/wblakecaldwell/profiler">here</a></div>
	</body>
</html>
`
