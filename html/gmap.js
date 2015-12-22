// Markers and other controls are stored in the variables

// variable for map
var map;

// marker, info window variables
var anchor_marker;
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
var anchor_icon = 'http://maps.google.com/mapfiles/ms/icons/blue-dot.png'

var mrkr;

// Initialize the main map
function initializeMap() {
	console.log("Initializing the map")
	
	// Create the SF map with the center of city as initial value
	var sfcity = new google.maps.LatLng(37.78, -122.454150)
	map = new google.maps.Map(document.getElementById('gmap-canvas'), {
		zoom: 14,
		center: sfcity,
		mapTypeId: google.maps.MapTypeId.ROADMAP
	});

	index = 0;
	curRadius = 800;
	curCount = 10;

	// create all the google map services
	geocoder = new google.maps.Geocoder;
	directionsService = new google.maps.DirectionsService;
	directionsDisplay = new google.maps.DirectionsRenderer;

	directionsDisplay.setMap(map);
	directionsDisplay.setOptions( { suppressMarkers: true } );

	infowindow = new google.maps.InfoWindow;
	infowindow.setContent('San Francisco');

	m_infowindow = new google.maps.InfoWindow;

	// Anchor marker ... user's selection, or typed in the address bar

	anchor_marker = new google.maps.Marker({
		map: map,
		position: sfcity,
		icon: anchor_icon,
		draggable: true
	});

	// Create the search box and link it to the UI element.
	var defaultBounds = new google.maps.LatLngBounds(
		new google.maps.LatLng(37.5902, -122.1759),
		new google.maps.LatLng(37.8474, -122.5631)
	);

	input = document.getElementById('pac-input');
	//var autocomplete  = new google.maps.places.Autocomplete(input);
	searchBox = new google.maps.places.SearchBox(input, {
		bounds: defaultBounds
	});
	map.controls[google.maps.ControlPosition.TOP_LEFT].push(input);

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
		strokeOpacity: 0.9,
		strokeWeight: 4,
		fillColor: '#000000',
		fillOpacity: 0.35,
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
	map.addListener('click', function(e) {
		placeMarkerAndPanTo(e.latLng, map);
	});

	// Place marker for the initial location
	placeMarkerAndPanTo(sfcity, map)
}


// Add marker at truck location

function addTruckLocationMarker(location, id, name, address, fooditems) {

	var marker = new google.maps.Marker({
		position: location,
		map: map,
		clickable: true
	});


	google.maps.event.addListener(marker, 'click', function() {
		m_infowindow.setContent("<b>" + name + "</b>" + "<br>" + address + "<br>" + fooditems)
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
	clearMarkers();
	markers = [];
}

// Get the json for a given URL
function GetJson(yourUrl) {
	var Httpreq = new XMLHttpRequest(); // a new request
	Httpreq.open("GET", yourUrl, false);
	Httpreq.send(null);
	return Httpreq.responseText;
}

// Place the anchor_marker and pan to the location
function placeMarkerAndPanTo(latLng, map) {
	index = index + 1;
	map.panTo(latLng);

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

	curLat = latLng.lat();
	curLng = latLng.lng();
	doSearchAndUpdate()
}

function doSearchAndUpdate() {
	var searchurl = "/search?lat=" + curLat + "&lng=" + curLng + "&radius=" + curRadius + "&count=" + curCount;
	var json_obj = JSON.parse(GetJson(searchurl));

	deleteMarkers();

	// Create table entries for the results from the query
	var innerhtml = "";
	if (json_obj) {
		for (var i = 0; i < json_obj.length; i++) {
			var obj = json_obj[i];
			var loc = {
				lat: obj.Lat,
				lng: obj.Lng
			};
			var id = "";
			id = id + i;
			addTruckLocationMarker(loc, id, obj.Name, obj.Address, obj.Fooditems)

			innerhtml = innerhtml + "<tr style=color:black onclick= " + 'mouseover(this,' + i + ') '
			innerhtml = innerhtml + " onmouseout= " + 'mouseout(this,' + i + ')> '
			innerhtml = innerhtml + "<td>" + "<b>" + obj.Name + "</b>" + "<br>"
			innerhtml = innerhtml + obj.Address + '<br><font color="green">' + obj.Distance.toFixed(2) + " miles </font><br>"
			innerhtml = innerhtml + "<a style=\"color:blue\;\" onclick= " + 'foo(' + curLat + ',' + curLng + ',' + obj.Lat + ',' + obj.Lng + ',\"DRIVING\")>' + 'Drive</a>'
			innerhtml = innerhtml + "<a style=\"color:blue\;\" onclick= " + 'foo(' + curLat + ',' + curLng + ',' + obj.Lat + ',' + obj.Lng + ',\"WALKING\")>&nbsp;&nbsp;&nbsp;&nbsp' + 'Walk</a>'
			innerhtml = innerhtml + "</td>" + "</tr>" ;
		}
	} else {
		innerhtml = "<p>" + "No results found !! " + "</p>";
	}

	document.getElementById('results').innerHTML = innerhtml;
}

/*
directionsService.route(request, function(response, status) {
    if (status == google.maps.DirectionsStatus.OK) {
      directionsDisplay.setDirections(response);
      // add start and end markers
      startMarker = new google.maps.Marker({
        position: response.mc.origin,
        icon: 'http://maps.google.com/mapfiles/ms/icons/green-dot.png',
        map: map
      });
      endMarker = new google.maps.Marker({
        position: response.mc.destination,
        map: map
      });
    }
  });
*/
function foo(clat, clng, dlat, dlng, mode) {
	var origin_latlng = new google.maps.LatLng(clat, clng)
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
function mouseover(row, id) {
	google.maps.event.trigger(markers[id], 'click');
	directionsDisplay.setDirections({routes: []});
}

function mouseout(row, id) {

}

function truckCountUpdate(val) {
	document.querySelector('#trucks').value = val;
	curCount = val;
	doSearchAndUpdate()
}

function radiusChangeUpdate(val) {
	document.querySelector('#radius').value = val;
	// Converting to meters
	curRadius = val * 1600;
	circle_marker.setRadius(curRadius)
	doSearchAndUpdate()
}

// Initialize the google map on load
google.maps.event.addDomListener(window, 'load', initializeMap);