# POST /v1/scan 

## Example Request

```
curl -s -F "name=eicar" -F "file=@test/eicar.txt" http://clamav-http/v1/scan
```

## Response

### File is infected

```
HTTP/1.1 200 OK

{
    "errorCode": "SCAN_OK",
    "listFiles": [
        {
            "name": "eicar.txt",
            "isPassed": false,
            "found": "stream: Win.Test.EICAR_HDB-1 FOUND"
        }
    ],
    "errorDescription": ""
}
```

### File is clean

```
HTTP/1.1 200 OK

{
    "errorCode": "SCAN_OK",
    "listFiles": [
        {
            "name": "test.txt",
            "isPassed": true,
            "found": ""
        }
    ],
    "errorDescription": ""
}
```

### No File sent

```
HTTP/1.1 200 OK

{
    "errorCode": "NO_FILE_SENT",
    "listFiles": [],
    "errorDescription": ""
}
```

## Response Model

### errorCode: possible values
```
"SCAN_OK" | "MULTIPART_PARSE_ERROR" | "NO_FILE_SENT" | "FILE_ERROR" | "SCAN_FILE_ERROR" | "SCAN_RESPONSE_ERROR" | "SCAN_PARSE_ERROR"
```
### errorDescription: contains possible errors on multipart parsing, opening file, processing file
### listFiles: Array of scanned listFiles

## Scanned file Model
### name: Name of the file scanned
### isPassed: Boolean (True if file passed correctly the scan / False if the scan found something in the file)
### found: Description of the what have been found during the scan (Example: stream: Win.Test.EICAR_HDB-1 FOUND)