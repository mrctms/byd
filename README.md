### BYD

Byd is yet another YouTube downloader wrapper around youtube-dl or yt-dlp.

Does not provide any kind of user interface, only the back-end is available

## How it works
Everything revolves around hubs, where each hub contains a list of download request        
Each hub creates a zip with all the files of download requests


## Environment variables

Environment variables that must be provided:

- YTD_PATH - youtube-dl or yt-dlp path
- FFMPEG_PATH - ffmpeg path
- BYD_OUTPUT_PATH - output path where all the zips are saved
- BYD_BUS_CS - RabbitMQ connection string
- BYD_DB_CS - PostgreSQL connection string
- BYD_API_URL - URL to use for the REST API

## REST API

Available endpoints:

`/hubs/downloads` - POST            
Accept on body:
```json
{
    "data":[
        {
            "url":"",
            "dir":""
        }  
    ],
    "format":"",
    "type":""
}
```
Format must be `audio` or `video`           
Returns:
```json
{
    "Message":"",
    "Details":""
}
```
`/hubs/{hubId}/downloads` - GET         
Returns:
```json
{
    "hubId": 0,
    "downloads":[
        {
            "downloadId":0,
            "eta":"",
            "speed":"",
            "total":"",
            "progress":"",
            "url":"",
            "type":"",
            "format":"",
            "dir":"",
            "details":""
        }
    ]
}
```
`/hubs/{hubId}` - GET           
Returns:
```json
{
    "hubId":0,
    "result":"",
    "details":"",
    "zipName":""

}
```
The result could be:
- ok - everything ok
- errors - some download has failed
- fatal-error - entire hub downloads has failed

`/hubs/{hubId}/downloads/{downloadId}` - GET            
Returns:
```json
 {
    "downloadId":0,
    "eta":"",
    "speed":"",
    "total":"",
    "progress":"",
    "url":"",
    "type":"",
    "format":"",
    "dir":"",
    "details":""
}
```


