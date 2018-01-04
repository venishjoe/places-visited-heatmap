<h1><center>Places Visited Heatmap</center></h1>

## Getting Started
I had a list of cities, addresses, landmarks etc. visited by me over the years and I was looking into ways to visualize those. 
During my search, I came across Google Maps JavaScript API which not only visualizes location data but also has ways to show places which are visited more often than others. But the API needs location data in a specific format with latitue and longitude coordinates in a JavaScript along with weight for each. 

This project is to convert my plain text cities, addresses, landmark etc. data to geographic coordinates and assign weightage to each location. Once I have that, integrate the data with Google Maps to visualize.

End result will look like below. This is generated with a smaller set of test data.

Visit <a href="https://venishjoe.net/post/places-visited-heatmap/"> venishjoe.net/post/places-visited-heatmap </a> for more details. 

<img src="https://venishjoe.net/media/images/00012/places-visited-heatmap.png">

## Installation
You will need Google Geo API Key and Google API Key to run Go code and HTML Google maps render. Visit below site to get a free API key from Google. 

* <a href="https://developers.google.com/maps/documentation/geocoding/get-api-key"> Google Geo API Key </a>

* <a href="https://developers.google.com/maps/documentation/javascript/get-api-key"> Google API Key </a>

Clone the repository and create a directory "data" under source. Create two files

* googlegeoapi.key (contains your Google Geo API Key)
* location.data (contains your location data. sample below)

<b>location.data</b><br>
<i>Las Vegas, NV<br>
Niagara Falls, NY<br>
San Francisco, CA<br>
Empire State Building, NY<br>
Acadia National Park, Bar Harbour, ME<br>
Great Smoky Mountains National Park, TN<br>
The White House, Washington, DC<br>
Old Faithful, Yellowstone National Park, WY<br>
Disneyland, Anaheim, CA<br>
Las Vegas, NV</i><br>


Then run places.go (go run places.go). This will generate locationlatlng.data under data directory and places_data.js under ui directory.

Update Google maps URL in ui/places.html with your Google API key 

<i>src="https://maps.googleapis.com/maps/api/js?key=##YOUR_API_KEY_HERE##&libraries=visualization&callback=initMap"</i>

Load places.html in a browser.


## Contact
[Venish Joe Clarence](http://venishjoe.net)

## License
Places Visited Heatmap source code is available under the MIT [License](/LICENSE).