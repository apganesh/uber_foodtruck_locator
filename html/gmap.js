// Markers and other controls are stored in the variables

// variable for map
var map;

// marker, info window variables
var anchor_marker;
var pop_marker;
var infowindow;

var markers = [];
var m_infowindow;

// variable for covering the radius
var circle_marker;

var index;
var geocoder;

// variable for directions
var directionsService;
var directionsDisplay;

var input;
var searchBox;

// cache variables for the current values

var curLat;
var curLng;
var curRadius;
var curCount;
var curTypes;
var fitToBounds;
var curBounds;

// locations
var defaultLocation
var curLocation

// copied the icon from http://maps.google.com/mapfiles/ms/icons/blue-dot.png'
var anchor_icon = 'blue-dot.png'
var truck_icon = 'truck16_black.png'
var pop_icon = 'red-dot.png'
var mrkr;

// Initialize the main map
function initializeMap() {
	console.log("Initializing the map")
	
	// Create the SF map with the center of city as initial value
	var sfcity = new google.maps.LatLng(37.78, -122.454150)
	
	defaultLocation = new google.maps.LatLng(37.78, -122.454150)
	curLocation = new google.maps.LatLng(37.78, -122.454150)

	map = new google.maps.Map(document.getElementById('gmap-canvas'), {
		zoom: 1,
		center: defaultLocation,
		mapTypeId: google.maps.MapTypeId.ROADMAP
	});

	index = 0;
	curRadius = 800;
	curCount = 10;
	curTypes = [];
	fitToBounds = true;

	// create all the google map services
	geocoder = new google.maps.Geocoder;
	directionsService = new google.maps.DirectionsService;
	directionsDisplay = new google.maps.DirectionsRenderer;

	directionsDisplay.setMap(map);
	directionsDisplay.setOptions( { suppressMarkers: true, preserveViewport: true } );

	infowindow = new google.maps.InfoWindow;
	infowindow.setContent('San Francisco');

	m_infowindow = new google.maps.InfoWindow;

	// Anchor marker ... user's selection, or typed in the address bar

	anchor_marker = new google.maps.Marker({
		map: map,
		position: defaultLocation,
		icon: anchor_icon,
		draggable: true
	});

	pop_marker = new google.maps.Marker({
		map:map,
		animation: google.maps.Animation.DROP,
		icon: pop_icon,
		clickable: false,
		zIndex: 10,
		optimized: false
	});

	// Create the search box and link it to the UI element.
	var defaultBounds = new google.maps.LatLngBounds(
		new google.maps.LatLng(37.5902, -122.1759),
		new google.maps.LatLng(37.8474, -122.5631)
	);

	input = document.getElementById('pac-input');

	searchBox = new google.maps.places.SearchBox(input, {
		bounds: defaultBounds
	});
	//map.controls[google.maps.ControlPosition.TOP_LEFT].push(input);

	// Bias the SearchBox results towards current map's viewport.
	map.addListener('bounds_changed', function() {
		searchBox.setBounds(map.getBounds());
	});

	searchBox.addListener('places_changed', function() {
		var places = searchBox.getPlaces();

		if (places.length == 0) {
			return;
		}

		var curloc = places[0].geometry.location
		var curlatlng = new google.maps.LatLng(curloc.lat(), curloc.lng());

		placeMarkerAndPanTo(curlatlng, map)

	});

	// Circle marker shows the desired area of search
	circle_marker = new google.maps.Circle({
		strokeColor: '#FFA700',
		strokeOpacity: 0.8,
		strokeWeight: 2,
		fillColor: '#000040',
		fillOpacity: 0.2,
		map: map,
		radius: curRadius,
		clickable: false
	});
	circle_marker.bindTo('center', anchor_marker, 'position')


	// Callbacks for anchor_marker
	google.maps.event.addListener(anchor_marker, 'click', function() {
		infowindow.open(anchor_marker.get('map'), anchor_marker);
	});

	google.maps.event.addListener(anchor_marker, 'dragend', function(event) {
		placeMarkerAndPanTo(event.latLng, map);
	});

	// Callback for selecting a location
	// Disbling click on the map as this causes too much confusion to the user 
	// if he accidentaly clicks on map which clicking on a marker
	//map.addListener('click', function(e) {
	//	placeMarkerAndPanTo(e.latLng, map);
	//});

	// Place marker for the initial location
	placeMarkerAndPanTo(sfcity, map)
}


function setDefaultLocation() {
	map.setCenter(defaultLocation)
	setCurLocation(defaultLocation.lat, defaultLocation.lng)
	placeMarkerAndPanTo(defaultLocation, map)
}

function setCurLocation(lat, lng) {
	curLocation = new google.maps.LatLng(lat,lng);
}

function updateCurrentLocation() {
	
	if (navigator.geolocation) {
		navigator.geolocation.getCurrentPosition(function(position) {
			var pos  = new google.maps.LatLng(position.coords.latitude,position.coords.longitude);
			map.setCenter(pos);
			placeMarkerAndPanTo(pos,map)
			setCurLocation(position.coords.latitude,position.coords.longitude);
		}, function() {
			handleLocationError(true, infoWindow, map.getCenter());
		});
	} else {
    	// Browser doesn't support Geolocation
    	handleLocationError(false, infoWindow, map.getCenter());
    }

}

function handleLocationError(browserHasGeolocation, infoWindow, pos) {
  infoWindow.setPosition(pos);
  infoWindow.setContent(browserHasGeolocation ?
                        'Error: The Geolocation service failed.' :
                        'Error: Your browser doesn\'t support geolocation.');
}

function initializeFoodCategories() {
// READ the food categories file
  	$.getJSON('./food_categories.json', function(data) {
    	initializeFoodtypes(data)
	});
}


function initializeFoodtypes(data) {
	var innerhtml = ""
	
	for (var i = 0; i < data.length ; i++) {
		innerhtml = innerhtml + "<li><input class=\"ftype_cbox\" type=\"checkbox\" value=" + data[i].name + " name=\"foodtype\"/> <label>"+data[i].name+"</label></li>"
	}
	document.getElementById('foodtypes').innerHTML = innerhtml;
}

function initializeApp() {
	initializeMap()
	initializeFoodCategories()
}

$('#foodtypes').on('change', 'input[type=checkbox]', function(e) {
	var types = [];
	$.each($("input[name='foodtype']:checked"), function(){            
        types.push($(this).val());
    });
    foodTypeUpdate(types)
 });

$('#fitbounds').click(
	function() {
		 if ($('#fitbounds').is(':checked')) {
		 	fitToBounds = true;
        }  else {
        	fitToBounds = false;
        }
        fitToBoundsUpdate()
	}
);



function popMarker(lat, lng) {
	var loc = new google.maps.LatLng(lat, lng) 
	pop_marker.setPosition(loc);
	pop_marker.setVisible(true);
	m_infowindow.close();
	directionsDisplay.setDirections({routes: []});
}

function addTruckLocationMarker(id,obj) { 
	var loc = {lat: obj.LL.Lat, lng: obj.LL.Lng};

	var marker = new google.maps.Marker({
		position: loc,
		map: map,
		icon: truck_icon,
		clickable: true
	});

	google.maps.event.addListener(marker, 'click', function() {
		var htmlcontent = "<b>" + obj.Name + "</b>" + "<br>" + obj.Address + "<br><p>" 
		htmlcontent = htmlcontent + obj.Fooditems + "</p><br>" 
		htmlcontent = htmlcontent + "<b>Hours: </b>" + obj.Dayhours + "<br/>"
		htmlcontent = htmlcontent + "<b>Food Types: </b>" + obj.Foodtypes

		m_infowindow.setContent(htmlcontent)
		m_infowindow.open(map, marker)
		directionsDisplay.setDirections({routes: []});
	});

	markers.push(marker);
	return marker;
}



// Sets the map on all markers in the array.
function setMapOnAll(map) {
	for (var i = 0; i < markers.length; i++) {
		markers[i].setMap(map);
	}
}

// Removes the markers from the map, but keeps them in the array.
function clearMarkers() {
	setMapOnAll(null);
}

// Shows any markers currently in the array.
function showMarkers() {
	setMapOnAll(map);
}

// Deletes all markers in the array by removing references to them.
function deleteMarkers() {
	directionsDisplay.setDirections({routes: []});
	pop_marker.setVisible(false)
	clearMarkers();
	markers = [];
}

// Place the anchor_marker and pan to the location
function placeMarkerAndPanTo(latLng, map) {
	index = index + 1;
	//map.panTo(latLng);

	geocoder.geocode({
		'location': latLng
	}, function(results, status) {
		if (status === google.maps.GeocoderStatus.OK) {
			if (results[1]) {
				infowindow.setContent(results[1].formatted_address);
			} else {
				window.alert('No results found');
			}
		} else {
			window.alert('Geocoder failed due to: ' + status);
		}
	});

	anchor_marker.setPosition(latLng);
	setCurLocation(latLng.lat(), latLng.lng())
	map.fitBounds(circle_marker.getBounds());

	curLat = latLng.lat();
	curLng = latLng.lng();
	doSearchAndUpdate()
}

function doSearchAndUpdate() {
	var searchurl = "/Trucks?lat=" + curLat + "&lng=" + curLng + "&radius=" + curRadius + "&count=" + curCount +"&types=" + curTypes;
	$.ajax({
		type: 'GET',
		url: searchurl,
		data: {},
		dataType: 'json',
		success: function(data) 
		{ 
			populateTruckLocations(data) 
		},
		error: function() { alert('Error occured while getting food truck locations !!!'); }
	});
}

function populateTruckLocations(data) {
	deleteMarkers();

	// Create table entries for the results from the query
	var innerhtml = "<tr><td></td></tr>";


	curBounds = new google.maps.LatLngBounds()
	curBounds.extend(curLocation)

	if (data) {
		for (var i = 0; i < data.length; i++) {
			var obj = data[i];
			var loc = {lat: obj.LL.Lat, lng: obj.LL.Lng};

			addTruckLocationMarker(i,obj)
			var gLL = new google.maps.LatLng(loc.lat,loc.lng);
			curBounds.extend(gLL);

			innerhtml = innerhtml + "<tr" +  ' onmouseover= popMarker(' + loc.lat + ',' + loc.lng+ ')>' 
			innerhtml = innerhtml + "<td style=\"line-height:1.75\">" 
			innerhtml = innerhtml + "<a onclick="+ 'showInfoWindow(' + i + ") </a>"
			innerhtml = innerhtml + "<span class=\"leftalign\" style=\"color:navyblue;font-weight:bold\"><u>" + obj.Name + "</u></span></a>" 
			innerhtml = innerhtml + "<span class=\"rightalign\" style=\"color:darkolivegreen\">" + obj.Distance.toFixed(2) + " miles </span><br/>"
			innerhtml = innerhtml + "<span class=\"leftalign\">" + obj.Address + "</span>"
			innerhtml = innerhtml + "<a class=\"blah\" onclick=" + 'findDirection(' + loc.lat + ',' + loc.lng + ',\"WALKING\")>'  
			innerhtml = innerhtml + "<span class=\"rightalign\" style=\"color:darkolivegreen;font-weight:bold\">" + "<u>Directions</u>" + "</span></a><br/>"

			innerhtml = innerhtml + '<b>Serving: </b>' + obj.Foodtypes[0]
			for(var j = 1; j < obj.Foodtypes.length; j++) {
				innerhtml = innerhtml  + ', ' + obj.Foodtypes[j] 
			}

			innerhtml = innerhtml + "</td></tr>" ;
		}
	} else {
		innerhtml = "<p>" + "No results found !! " + "</p>";
	}

	if(fitToBounds) {
		map.fitBounds(curBounds)
	} else {
		map.fitBounds(circle_marker.getBounds())
	}
	document.getElementById('results').innerHTML = innerhtml;
}

function showInfoWindow(id) {
	google.maps.event.trigger(markers[id], 'click');
}

function findDirection(dlat, dlng, mode) {
	popMarker(dlat, dlng);
	var origin_latlng = new google.maps.LatLng(curLat, curLng)
	var dest_latlng = new google.maps.LatLng(dlat, dlng);

	directionsService.route({
		origin: origin_latlng,
		destination: dest_latlng,
		travelMode: google.maps.TravelMode[mode]
	}, function(response, status) {
		if (status === google.maps.DirectionsStatus.OK) {
			directionsDisplay.setDirections(response);
		} else {
			window.alert('Directions request failed due to ' + status);
		}
	});
}

// Function for mouse click on a row of the table element
// function mouseover(row, id) {
// 	google.maps.event.trigger(markers[id], 'click');
// 	directionsDisplay.setDirections({routes: []});
// }

// function mouseout(row, id) {

// }
function fitToBoundsUpdate() {
	if(fitToBounds)
		map.fitBounds(curBounds)
	else
		map.fitBounds(circle_marker.getBounds())
}

function foodTypeUpdate(val) {
	//document.querySelector('#fcat').val = val;
	curTypes = val;
	doSearchAndUpdate()
}

function truckCountUpdate(val) {
	document.querySelector('#trucks').value = val;
	curCount = val;
	doSearchAndUpdate()
}

function radiusChangeUpdate(val) {
	document.querySelector('#radius').value = val;
	// Converting to meters
	curRadius = val * 1609;
	doSearchAndUpdate()
	circle_marker.setRadius(curRadius)
	map.fitBounds(circle_marker.getBounds())
}

// Initialize the google map on load
google.maps.event.addDomListener(window, 'load', initializeApp);