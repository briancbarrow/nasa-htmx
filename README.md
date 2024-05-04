# Project nasa-htmx

This is a small application meant to provide data from the NASA API. Specifically, it provides data about Solar Flares and Asteroids near Earth. There are two ways to access the data, either through a web interface or through a REST API. The web interface is built using htmx and the REST API is built using Go and Chi.

## REST API

### Solar Flares

#### GET /solar?startDate=YYYY-MM-DD&endDate=YYYY-MM-DD

**Parameters**:

- startDate: The start date of the query (Defaults to 30 days ago)
- endDate: The end date of the query (Defaults to the current date)

**Example Response**:

```json
{
  "count": 1,
  "solarFlares": [
    {
      "beginTime": "2021-09-07T20:48:00",
      "peakTime": "2021-09-07T20:48:00",
      "endTime": "2021-09-07T20:48:00",
      "classType": "C1.1",
      "link": "https://www.swpc.noaa.gov/products/goes-x-ray-flux"
    }
  ],
  "startDate": "2021-09-07",
  "endDate": "2021-09-07"
}
```

This endpoint returns an object with the following fields:

- count: The number of solar flares in the response
- solarFlares: An array of objects with information about each solar flare
- startDate: The start date of the query
- endDate: The end date of the query

Each solar flare object has the following fields:

- beginTime: The start time of the solar flare
- peakTime: The peak time of the solar flare
- endTime: The end time of the solar flare
- classType: The class of the solar flare
- link: A link to more information about the solar flare

### Asteroids

#### GET /asteroids

**Example Response**:

```json
{
  "count": 1,
  "asteroids": [
    {
      "name": "2021 QL2",
      "nasa_jpl_url": "https://ssd.jpl.nasa.gov/sbdb.cgi?sstr=2021%20QL2;old=0;orb=1;cov=0;log=0;cad=0#orb",
      "is_potentially_hazardous_asteroid": false,
      "close_approach_date": "2021-09-07",
      "estimated_diameter": {
        "feet": {
          "min": 134.5,
          "max": 301.8
        }
      }
    }
  ]
}
```

This endpoint returns an object with the following fields:

- count: The number of asteroids in the response
- asteroids: An array of objects with information about each asteroid

Each asteroid object has the following fields:

- name: The name of the asteroid
- nasa_jpl_url: A link to the NASA JPL page for the asteroid
- is_potentially_hazardous_asteroid: Whether the asteroid is potentially hazardous
- close_approach_date: The date of the closest approach to Earth
- estimated_diameter: An object with the estimated diameter of the asteroid in feet

## Web Interface

The web interface provides a simple way to access the data from the NASA API. You can view solar flare data and asteroid data by visiting the respective URLs. The data is loaded using htmx and updated dynamically without refreshing the page.

### Solar Flares

#### [/web/solar](/web/solar)

This page initially displays two input fields for the start and end dates of the query. You can enter the dates in the format YYYY-MM-DD and click the "Submit" button to load the data. The solar flare data will be displayed below the input fields.

You can submit the form without entering any dates to use the default values (30 days ago to the current date).

### Asteroids

#### [/web/asteroid](/web/asteroid)

This page displays the asteroid data directly when you visit the URL. The asteroid data will be displayed below the header.
