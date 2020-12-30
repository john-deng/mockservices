# solar-mock-app

solar-mock-app is an all in one app that build for testing Microservices on BareMetal, Docker, Kubernetes, or Service Mesh, the key features of this app is that you can deploy as many apps as possible with just a single binary file or docker image, then you can get hundreds or thousands microservices for testing. You can specify any to any HTTP request. e.g. starting two apps named order and product, sending HTTP requests from order, you will get a mock response. 

As for the Service Mesh testing, sometimes we want to do fault injection to the app which response the real error code, all you have to do is set specific headers, you will get the result.

| Header   |      Description      |  Example |
|----------|:-------------|:------|
| fi-svc | The app name | product  |
| fi-ver | The app version, it is optional   |   v1 |
| fi-code | Response code |    503 |

Example 

Using httpie to send the request to the mock app

```bash

http http://localhost:8083 fi-svc:product fi-ver:v1 fi-code:503 

```

## Getting started

### From Source Code

Get the source code and build with Go tools

```bash
# Clone the source code
git clone https://github.com/solarmesh-io/solar-mock-app.git

# Build
cd solar-mock-app
go build

```

After the app is built, you can run it directly

### 1. Run the binary file on bare metal.
```bash

samples/platform/baremetal/run.sh

```

### Run with Docker image


```bash

samples/platform/baremetal/run.sh

```


### Test it

using httpie or curl to see the response of these mock apps

```bash
http http://localhost:8083
```

The output will be as follows,

```json
{
    "Code": 200,
    "Data": {
        "App": "payment",
        "Cluster": "cluster02",
        "Header": {
            "Accept": [
                "*/*"
            ],
            "Accept-Encoding": [
                "gzip, deflate"
            ],
            "Connection": [
                "keep-alive"
            ],
            "Uber-Trace-Id": [
                "60dd508736c414a:684c72590f7e454d:60dd508736c414a:1",
                "60dd508736c414a:4e913b9d36f7ed3a:60dd508736c414a:1"
            ],
            "User-Agent": [
                "HTTPie/2.3.0"
            ]
        },
        "MetaData": " -> payment -> ",
        "SourceApp": "",
        "SourceAppVersion": "",
        "Upstream": [
            {
                "Code": 200,
                "Data": {
                    "App": "inventory",
                    "Cluster": "cluster02",
                    "Header": {
                        "Accept": [
                            "*/*"
                        ],
                        "Accept-Encoding": [
                            "gzip, deflate"
                        ],
                        "Connection": [
                            "close",
                            "keep-alive"
                        ],
                        "Uber-Trace-Id": [
                            "60dd508736c414a:684c72590f7e454d:60dd508736c414a:1",
                            "60dd508736c414a:e7b99b1d734b0ad:7db2e51b4dee559:1"
                        ],
                        "User-Agent": [
                            "HTTPie/2.3.0"
                        ]
                    },
                    "MetaData": " -> inventory -> ",
                    "SourceApp": "payment",
                    "SourceAppVersion": "v1",
                    "Upstream": [
                        {
                            "Code": 200,
                            "Data": {
                                "App": "product",
                                "Cluster": "cluster01",
                                "Header": {
                                    "Accept": [
                                        "*/*"
                                    ],
                                    "Accept-Encoding": [
                                        "gzip, deflate"
                                    ],
                                    "Connection": [
                                        "close",
                                        "keep-alive"
                                    ],
                                    "Uber-Trace-Id": [
                                        "60dd508736c414a:684c72590f7e454d:60dd508736c414a:1",
                                        "60dd508736c414a:e7b99b1d734b0ad:7db2e51b4dee559:1"
                                    ],
                                    "User-Agent": [
                                        "HTTPie/2.3.0"
                                    ]
                                },
                                "MetaData": " -> product",
                                "SourceApp": "inventory",
                                "SourceAppVersion": "v1",
                                "Upstream": null,
                                "Url": "localhost:8080/",
                                "UserData": "baremetal",
                                "Version": "v1"
                            },
                            "Message": "Success"
                        }
                    ],
                    "Url": "localhost:8081/",
                    "UserData": "demo",
                    "Version": "v1"
                },
                "Message": "Success"
            },
            {
                "Code": 200,
                "Data": {
                    "App": "order",
                    "Cluster": "cluster02",
                    "Header": {
                        "Accept": [
                            "*/*"
                        ],
                        "Accept-Encoding": [
                            "gzip, deflate"
                        ],
                        "Connection": [
                            "close",
                            "keep-alive"
                        ],
                        "Uber-Trace-Id": [
                            "60dd508736c414a:684c72590f7e454d:60dd508736c414a:1",
                            "60dd508736c414a:4e913b9d36f7ed3a:60dd508736c414a:1",
                            "60dd508736c414a:659dbe59cfab6b68:1e61d36379e11215:1",
                            "60dd508736c414a:7afb61ab158b358b:1e61d36379e11215:1"
                        ],
                        "User-Agent": [
                            "HTTPie/2.3.0"
                        ]
                    },
                    "MetaData": " -> order -> ",
                    "SourceApp": "payment",
                    "SourceAppVersion": "v1",
                    "Upstream": [
                        {
                            "Code": 200,
                            "Data": {
                                "App": "product",
                                "Cluster": "cluster01",
                                "Header": {
                                    "Accept": [
                                        "*/*"
                                    ],
                                    "Accept-Encoding": [
                                        "gzip, deflate"
                                    ],
                                    "Connection": [
                                        "close",
                                        "keep-alive"
                                    ],
                                    "Uber-Trace-Id": [
                                        "60dd508736c414a:684c72590f7e454d:60dd508736c414a:1",
                                        "60dd508736c414a:4e913b9d36f7ed3a:60dd508736c414a:1",
                                        "60dd508736c414a:659dbe59cfab6b68:1e61d36379e11215:1"
                                    ],
                                    "User-Agent": [
                                        "HTTPie/2.3.0"
                                    ]
                                },
                                "MetaData": " -> product",
                                "SourceApp": "order",
                                "SourceAppVersion": "v2",
                                "Upstream": null,
                                "Url": "localhost:8080/",
                                "UserData": "baremetal",
                                "Version": "v1"
                            },
                            "Message": "Success"
                        },
                        {
                            "Code": 200,
                            "Data": {
                                "App": "inventory",
                                "Cluster": "cluster02",
                                "Header": {
                                    "Accept": [
                                        "*/*"
                                    ],
                                    "Accept-Encoding": [
                                        "gzip, deflate"
                                    ],
                                    "Connection": [
                                        "close",
                                        "keep-alive"
                                    ],
                                    "Uber-Trace-Id": [
                                        "60dd508736c414a:684c72590f7e454d:60dd508736c414a:1",
                                        "60dd508736c414a:4e913b9d36f7ed3a:60dd508736c414a:1",
                                        "60dd508736c414a:659dbe59cfab6b68:1e61d36379e11215:1",
                                        "60dd508736c414a:7afb61ab158b358b:1e61d36379e11215:1",
                                        "60dd508736c414a:7115aa258de5e54:efaca992213cfd5:1"
                                    ],
                                    "User-Agent": [
                                        "HTTPie/2.3.0"
                                    ]
                                },
                                "MetaData": " -> inventory -> ",
                                "SourceApp": "order",
                                "SourceAppVersion": "v2",
                                "Upstream": [
                                    {
                                        "Code": 200,
                                        "Data": {
                                            "App": "product",
                                            "Cluster": "cluster01",
                                            "Header": {
                                                "Accept": [
                                                    "*/*"
                                                ],
                                                "Accept-Encoding": [
                                                    "gzip, deflate"
                                                ],
                                                "Connection": [
                                                    "close",
                                                    "keep-alive"
                                                ],
                                                "Uber-Trace-Id": [
                                                    "60dd508736c414a:684c72590f7e454d:60dd508736c414a:1",
                                                    "60dd508736c414a:4e913b9d36f7ed3a:60dd508736c414a:1",
                                                    "60dd508736c414a:659dbe59cfab6b68:1e61d36379e11215:1",
                                                    "60dd508736c414a:7afb61ab158b358b:1e61d36379e11215:1",
                                                    "60dd508736c414a:7115aa258de5e54:efaca992213cfd5:1"
                                                ],
                                                "User-Agent": [
                                                    "HTTPie/2.3.0"
                                                ]
                                            },
                                            "MetaData": " -> product",
                                            "SourceApp": "inventory",
                                            "SourceAppVersion": "v1",
                                            "Upstream": null,
                                            "Url": "localhost:8080/",
                                            "UserData": "baremetal",
                                            "Version": "v1"
                                        },
                                        "Message": "Success"
                                    }
                                ],
                                "Url": "localhost:8081/",
                                "UserData": "demo",
                                "Version": "v1"
                            },
                            "Message": "Success"
                        }
                    ],
                    "Url": "localhost:8082/",
                    "UserData": "demo",
                    "Version": "v2"
                },
                "Message": "Success"
            }
        ],
        "Url": "localhost:8083/",
        "UserData": "demo",
        "Version": "v1"
    },
    "Message": "Success"
}


```