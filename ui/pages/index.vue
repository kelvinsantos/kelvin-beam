<template>
  <section class="section">
    <div id="googleMap" style="width:100%;height:400px;"></div>
    <div class="columns">
      <div class="column">
        <b-field label="Latitude">
          <b-input v-model="latitude"></b-input>
        </b-field>
      </div>
      <div class="column">
        <b-field label="Longitude">
          <b-input v-model="longitude"></b-input>
        </b-field>
      </div>
      <div class="column">
        <b-field label="Distance">
          <b-input v-model="distance"></b-input>
        </b-field>
      </div>
      <div class="column">
        <b-field label="Unit">
            <b-select placeholder="Select a unit" v-model="unit" expanded>
              <option
                v-for="unit in units"
                :value="unit.value"
                :key="unit.value"
              >
                {{ unit.name }}
              </option>
            </b-select>
        </b-field>
      </div>
      <div class="column find-scooters-container">
        <div>
          <b-button type="is-primary" class="find-scooters-btn" @click="onFindScootersClicked">Find Scooters</b-button>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import axios from 'axios'

import { scooters } from '~/assets/scooters.json'
import gmapsUtil from '../utils/gmaps';

function calculateDistance(lat1, lon1, lat2, lon2, unit) {
  var radlat1 = Math.PI * lat1/180
  var radlat2 = Math.PI * lat2/180
  var theta = lon1-lon2
  var radtheta = Math.PI * theta/180
  var dist = Math.sin(radlat1) * Math.sin(radlat2) + Math.cos(radlat1) * Math.cos(radlat2) * Math.cos(radtheta);
  
  if (dist > 1) {
      dist = 1;
  }
  
  dist = Math.acos(dist)
  dist = dist * 180/Math.PI
  dist = dist * 60 * 1.1515
  
  if (unit == "M") { dist = dist * 1609.34 } // 1 mile
  if (unit == "K") { dist = dist * 1.609344 } // 1 mile
  
  return dist
}

// Removes the markers from the map, but keeps them in the array.
function clearMarkers(component) {
  setMapOnAll(component, null);
}

// Sets the map on all markers in the array.
function setMapOnAll(component, map) {
  if (component.scootersMarker) {
    for (var i = 0; i < component.scootersMarker.length; i++) {
      component.scootersMarker[i].setMap(map);
    }
  }

  if (component.ownMarker) {
    component.ownMarker.setMap(map);
  }
}

// Plot all the available scooters in the map
function findScootersOnTheMap(component, scootersCopy) {
  var image = {
    url: '../images/scooter.png',
    size: new google.maps.Size(32, 32),
    origin: new google.maps.Point(0, 0),
    anchor: new google.maps.Point(0, 0)
  };

  return scootersCopy.map(scooter => new google.maps.Marker({
    position: {
      lat: parseFloat(scooter.lat),
      lng: parseFloat(scooter.lng)
    },
    map: component.map,
    icon: image, 
    title: scooter.name
  }));
}

// Plot own marker on the map
function findOwnMarkerOnTheMap(component) {
  return new google.maps.Marker({ 
    position: {
      lat: parseFloat(component.latitude),
      lng: parseFloat(component.longitude)
    },
    map: component.map
  });
}

export default {
  name: 'HomePage',
  data() {
    return {
      scooters: [],
      latitude: "1.285066",
      longitude: "103.823453",
      map: null,
      scootersMarker: null,
      ownMarker: null,
      distance: 3000,
      unit: "M",
      units: [
        { "name": "Meters", "value": "M" },
        { "name": "Kilometers", "value": "KM" },
      ],
      currentZoomLevel: 14
    }
  },
  methods: {
    getScooterList: async function() {
      return await axios.get(`http://localhost:9090/api/get-scooters`)
      .then((res) => {
        if (res.status === 200) {
          return res.data;
        }
      })
    },
    onFindScootersClicked: async function(event) {
      this.renderMap();
    },
    renderMap: async function() {
      try {
        const google = await gmapsUtil();

        var myLatlng = { lat: parseFloat(this.latitude), lng: parseFloat(this.longitude) };

        let map = this.map;
        if (!map) {
          map = new google.maps.Map(document.getElementById("googleMap"), {
            zoom: this.currentZoomLevel,
            center: myLatlng
          });
          this.map = map;
        }

        let scootersCopy = Object.assign([], this.scooters);
        scootersCopy = scootersCopy.filter(scooter => {
          return calculateDistance(this.latitude, this.longitude, scooter.lat, scooter.lng, this.unit) <= this.distance
        })

        let scootersMarker = this.scootersMarker;
        if (!scootersMarker) {
          scootersMarker = findScootersOnTheMap(this, scootersCopy);
          this.scootersMarker = scootersMarker;
        }

        let ownMarker = this.ownMarker;
        if (!ownMarker) {
          ownMarker = findOwnMarkerOnTheMap(this);
          this.ownMarker = ownMarker;
        }

        if (scootersMarker && ownMarker) {
          clearMarkers(this);

          scootersMarker = findScootersOnTheMap(this, scootersCopy);
          this.scootersMarker = scootersMarker;

          ownMarker = findOwnMarkerOnTheMap(this);
          this.ownMarker = ownMarker;
        }

        map.setCenter(ownMarker.getPosition());

        map.addListener('zoom_changed', function() {
          this.currentZoomLevel = map.getZoom();
        });
      } catch (error) {
        console.error(error);
      }
    }
  },
  async mounted() {
    this.scooters = await this.getScooterList();
    this.renderMap();
  }
}
</script>

<style>
html,
body {
  margin: 0;
  padding: 0;
}

.App {
  width: 100vw;
  height: 100vh;
}

.find-scooters-container {
  margin-top: auto;
}

.find-scooters-btn {
  width: 100%;
}
</style>