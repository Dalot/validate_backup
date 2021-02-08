# Deep validation of two json files.

## Assumptions:
 1. Each object must have an ID field.
 2. Both files, equal or not, must be an array JSON (see below)
 3. Nested objects or slices are not supported.
 4. The solution has to be memory efficient to support a number of GBs of files (It uses stream parsing and concurrency for such purpose).
 
 Example of an array JSON
 ```json
 [
    {
        "id": "1",
        "name": "John Doe"
    },
    {
        "id": "2",
        "name": "John Doe"
    }
 ]
 ```

> Although it does not support nested objects or slices, it can have an arbitrary number of fields in the object, additionally to those provided on the example above.


## Building the project
```
$ git clone https://github.com/Dalot/validate_backup
$ cd validate_backup
$ go build
$ go test -timeout 30s ./stream
$ ./validate_backup.exe -f1=before.json -f2=after.json`
```
