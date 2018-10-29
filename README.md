# IGCparaglinding


This api is an online service that will allow users to browse information about IGC files. 
IGC is an international file format for soaring track files that are used by paragliders and gliders. 
The program will store IGC files metadata in a NoSQL Database (persistent storage). The system will generate events and it will monitor for new events happening from the outside services.
These files can be browsed later by entering the specific URL paths defined later in this README.

For the development of the IGC processing, I have been using the open source IGC library for Go: <a href="https://github.com/marni/goigc">goigc</a>



## Getting Started

The project's name is **paragliding**. The root of the Igc Api is /paragliding/ which if you enter that path you get redirected to the /paragliding/api path. From that root you can write the other paths to get results.


### GET     /paragliding/api

For this get request you get meta information about the API

    {
      "uptime": <uptime>
      "info": "Service for IGC tracks."
      "version": "v1"
    }

where: <uptime> is the current uptime of the service formatted according to Duration format as specified by ISO 8601. 


### POST    /paragliding/api/track

You can do track registration with this request by putting in the request body the URL of the Igc File

    {
      "url": "<url>"
    }


As a response you get the ID for that track:

    {
      "id": "<id>"
    }



### GET     /paragliding/api/track

This request returns the array of all tracks ids as an array:
    
    [<id1>, <id2>, ...]
        


### GET     /paragliding/api/track/\<id\>

Returns the meta information about a given track with the provided <id>, or NOT FOUND response code with an empty body.
    
    {
        "H_date": <date from File Header, H-record>,
        "pilot": <pilot>,
        "glider": <glider>,
        "glider_id": <glider_id>,
        "track_length": <calculated total track length>
        "track_src_url": <the original URL used to upload the track, ie. the URL used with POST>
    }


### GET     /paragliding/api/track/\<id\>/\<field\>

Returns the single detailed meta information about a given track with the provided <id>, or NOT FOUND response code with an empty body. The response should always be a string, with the exception of the calculated track length, that should be a number.

Response will be: 

    <pilot> for pilot
    
    <glider> for glider
    
    <glider_id> for glider_id
    
    <calculated total track length> for track_length
    
    <H_date> for H_date
    
    <track_src_url>` for `track_src_url
   
  
  
## Ticker API

The system allows multiple people to upload tracks. To facilitate sharing track information, the API allows people to search through the entire collection of tracks, by track ID, to obtain the details about a given track. To notify dependent systems about new tracks uploaded to the system, there is a ticker API. The purpose of the ticker API is to notify dependant applications (such as complex IGC visualisation webapps) about the new tracks being available. For example: imagine that another webapp needs to know what new tracks have been uploaded from the last sync that the app done with your system. They will keep track of the last timestamp of the last track that they know about. If they ask your API about the last timestamp, and it is different to the one they have, that means that your system has more tracks now that they know about. I.e. new tracks have been uploaded since they have queried the track information. So that they can request to get IDs of all the new tracks that the system now has. The ticker API provide this facility, of obtaining information about updates. It provides simple paging functionality. 


### GET /api/ticker/latest

* What: returns the timestamp of the latest added track
* Response: `<timestamp>` for the latest added track


### GET /api/ticker/

* Returns the JSON struct representing the ticker for the IGC tracks. The first track returned should be the oldest. The array of track ids returned should be capped at 5, to emulate "paging" of the responses. The cap (5) should be a configuration parameter of the application (ie. easy to change by the administrator).
* Response

```
{
"t_latest": <latest added timestamp>,
"t_start": <the first timestamp of the added track>, this will be the oldest track recorded
"t_stop": <the last timestamp of the added track>, this might equal to t_latest if there are no more tracks left
"tracks": [<id1>, <id2>, ...],
"processing": <time in ms of how long it took to process the request>
}
```

### GET /api/ticker/`<timestamp>`

* Returns the JSON struct representing the ticker for the IGC tracks. The first returned track should have the timestamp HIGHER than the one provided in the query. The array of track IDs returned should be capped at 5, to emulate "paging" of the responses. The cap (5) should be a configuration parameter of the application (ie. easy to change by the administrator).
* Response:

```
{
   "t_latest": <latest added timestamp of the entire collection>,
   "t_start": <the first timestamp of the added track>, this must be higher than the parameter provided in the query
   "t_stop": <the last timestamp of the added track>, this might equal to t_latest if there are no more tracks left
   "tracks": [<id1>, <id2>, ...],
   "processing": <time in ms of how long it took to process the request>
}
```



## Webhooks API

The system will allow subscribing a webhook such that it can notify subscribers about the events in your system. The event that we have is adding (registering) new track via the `POST /api/track`. Thus, your system will notify all interested subscribers to the webhook API, with the notification about new track being added. 

### POST /api/webhook/new_track/

* Registration of new webhook for notifications about tracks being added to the system. Returns the details about the registration. The `webhookURL` is required parameter of the request. The `minTriggerValue` is optional integer, that defaults to 1 if ommited. It indicated the frequency of updates - after how many new tracks the webhook should be called. 

* **Request**

```
{
    "webhookURL": {
      "type": "string"
    },
    "minTriggerValue": {
      "type": "number"
    }
}
```


* **Response**

The response body should contain the id of the created resource (aka webhook registration), as string. Note, the response body will contain only the created id, as string, not the entire path; no json encoding. Response code upon success should be 200 or 201.


### After invoking a registered webhook

Everytime a track is added you get in Discord a message containing the information like this:
```
{
   "t_latest": <latest added timestamp of the entire collection>,
   "tracks": [<id1>, <id2>, ...],
   "processing": <time in ms of how long it took to process the request>
}
```


### GET /api/webhook/new_track/`<webhook_id>`

* With this you can access the registered webhooks. Registered webhooks should be accessible using the GET method and the webhook id generated during registration.
* **Response body**

```
{
    "webhookURL": {
      "type": "string"
    },
    "minTriggerValue": {
      "type": "number"
    }
}
```

### DELETE /api/webhook/new_track/`<webhook_id>`

* Deletes a registered webhooks. Registered webhooks can further be deleted using the DELETE method and the webhook id.
* Response body:

```
{
    "webhookURL": {
      "type": "string"
    },
    "minTriggerValue": {
      "type": "number"
    }
}
```


## Clock trigger

The idea behind the clock is to have a task that happens on regular basis without user interventions. In this system is implemented a task, that checks every 10min if the number of tracks differs from the previous check, and if it does, it will notify a predefined Discord webhook. The clock trigger has been deployed in OpenStack.



## Admin API 

### GET /admin/api/tracks_count

* Returns the current count of all tracks in the DB
* Response: current count of the DB records


### DELETE /admin/api/tracks

* Deletes all tracks in the DB
* Response: count of the DB records removed from DB




## Resources

* [Go IGC library](https://github.com/marni/goigc)
* [official MongoDB Go driver](https://github.com/mongodb/mongo-go-driver)


## Deployment

Deployed in Heroku: https://fierce-oasis-92488.herokuapp.com/
Deployed in OpenStack and it has these floating IPs: 


## Built With

Go Language
