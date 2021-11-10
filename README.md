dtom
====

`dtom` is a tool for importing DynamoDB JSONs into MongoDB.

Background
----------
DynamoDB is great for managing data without worrying about instances, servers, etc. To analyze the data more deeply, sometimes you could not always rely on the DynamoDB's query or scan functions, and you may need other tools like Athena.

To analyze the data locally, one can easily export a table dump gzipped into S3, and download it. 

Since DynamoDB can export all data into JSONs, we can in turn import them into MongoDB and use its powerful query engine.

### DynamoDB JSON
This is an example of a DynamoDB JSON exported to S3:
```js
{
    "Item":
    {
        "key1":
        {
            "S": "val1"
        },
        "key2":
        {
            "N": "42"
        },
        "key3":
        {
            "BOOL": false
        },
        "key4":
        {
            "M":
            {
                "key4.1":
                {
                    "S": "val4.1"
                },
                "key4.2":
                {
                    "M":
                    {
                        "key4.2.1":
                        {
                            "S": "val4.2.1"
                        }
                    }
                }
            }
        },
        "key5":
        {
            "L":
            [
                {
                    "M":
                    {
                        "key5.0":
                        {
                            "S": "val5.0"
                        },
                        "key5.1":
                        {
                            "N": "42"
                        }
                    }
                }
            ]
        }
    }
}
```

Usually, every item is stored as a single unformatted JSON on a line in a [JSON lines format](https://jsonlines.org/) in a text file. Then all files are distributed and compressed in various gzip archives.

Instructions
------------
- Make sure to have a running MongoDB instance.

- Once compiled, use `./dtom --help` to check the help:
```
Usage of ./dtom:
  -blockSize int
        How many items to insert on each request to Mongo (default 100)
  -collection string
        MongoDB collection
  -db string
        MongoDB database
  -key string
        DynamoDB partition key that will be used as _id
  -mongoURI string
        MongoDB URI
```

- The program reads input from the standard input, so you can use it in a pipeline:

```bash
zcat *.gz | ./dtom -db testDB -collection testColl -key testId --mongoURI 'mongodb://localhost' --blockSize 100
```

Developing
------------
To init the environment and install all dependencies, enter

``` 
go mod tidy
```

To build the executable, enter

``` 
go build
```

To run the tests, enter
```
go test
```

Contributing
------------
If you find any bugs or want to suggest an improvement, please open an issue.

You're more than welcome to open pull requests :smile:
