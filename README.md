# ethapi

# EthAPI
A simple HTTP API to query block and transaction information from the Ethereum Mainnet via Infura.

## Table Of Contents

* [Intro](#intro)
* [API documentation](#api-documentation)
* [Requirements](#requirements)
* [How To Build](#how-to-build)
* [Configuration](#configuration)
* [Load test](#load-test)
* [Improvements And Limitations](#improvements-and-limitations)

## Intro
The HTTP API exposes transaction and block data via REST endpoints.

## API Documentation
### Endpoints:
#### Block  
`GET api/v1/get-block`  
  
One of the following parameters:  
| Parameter | Type | Description |
| :--- | :--- | :--- |
| `latest` | `bool` | Whether to query for the latest block |
| `hash` | `string` | Hash of the block to be queried |
| `number` | `int` | Number of the block to be queried |

Sample response:
```
curl "http://localhost:9090/api/v1/get-block?number=1223333"
{"blockNumber":1223333,"hash":"0x793874f11b63154cbbaca5a5a4dd6d8a7fc83e2b01e153fe5e86b5dffe58f5d0","difficulty":22367179797762,"timestamp":1459032995,"transactionsCount":3,"transactions":[{"hash":"0xb4083356d8d959c18d70bb78a24df377bd5225176579f8bb2b23db861f541235","to":"0x120A270bbC009644e35F0bB6ab13f95b8199c4ad","value":"999580000000000000","nonce":0,"pending":false,"gas":21000,"gasPrice":20000000000},{"hash":"0x4697bfec59462de684ce2dd2be03595ef18d22eaa3b85c47854870f8d9c164cb","to":"0xE58622BC31Aaad2B71a42A54DBa85d6fa44E0e56","value":"1019981460000000000","nonce":353338,"pending":false,"gas":90000,"gasPrice":20000000000},{"hash":"0xde477d2269c9fe9f469c92c47e0fe5038aed44e59568d3351c4e29e57cc140e6","to":"0x34a5f2C9d68C3fF0E52D4D2F8C77e0466f4072c2","value":"1019762110000000000","nonce":353339,"pending":false,"gas":90000,"gasPrice":20000000000}]}
```
#### Transaction  
`GET api/v1/get-tx`  
  
One of the following parameters:  
| Parameter | Type | Description |
| :--- | :--- | :--- |
| `hash` | `string` | Hash of the transaction to be queried |

Sample response:
```
curl "http://localhost:9090/api/v1/get-tx?hash=0x63a169366501a8e5aeaa1a131615e854a5a62e067a68dedfe54dd4e013954a63"
{"hash":"0x63a169366501a8e5aeaa1a131615e854a5a62e067a68dedfe54dd4e013954a63","to":"0x36d593440B0c2c91232DB96C315cd10765388E54","value":"5000000000000000","nonce":1348014,"pending":false,"gas":50000,"gasPrice":31000000000}
```
#### Status Codes
| Status Code | Description |
| :--- | :--- |
| 200 | `OK` |
| 400 | `BAD REQUEST` |
| 404 | `NOT FOUND` |
| 500 | `INTERNAL SERVER ERROR` |

## Requirements
* Golang (~>1.15 for local build)
* Docker

## How To Build
Build locally, binary written to bin/ethapi:
```
make build
```

Build with Docker:
```
make build-docker
```

Tag and push to the Docker registry:
```
make distribute
```

## How To Run Locally
Without Docker:
```
make build
LOGLEVEL=DEBUG INFURA_ADDR=https://mainnet.infura.io/v3/<PROJECT_ID> ./bin/ethapi
```

With Docker:
```
make build-docker
docker run -it --rm -p 9090:9090 -e "INFURA_ADDR=https://mainnet.infura.io/v3/<PROJECT_ID>" docker.io/gpaggi/ethapi:latest
```

## Configuration
The following environment variables can be set to configure the application:
* `INFURA_ADDR` - Address of the Infura endpoint, including project ID [*required*]
* `LISTEN_ADDR` - Bind address for the webserver (IP:PORT) [*optional*]
* `LOG_LEVEL` - Log level (INFO, WARN, DEBUG, TRACE) [*optional*]

## Load Test
The load test is executed using [Vegeta](https://github.com/tsenart/vegeta) and a [list of static targets](test/target.list).  
The webserver (Gorilla Mux) used for this project is not the most performant but it can generally easily handle thousands of requests per second with a simple setup and on a recent Macbook pro (2,4 GHz 8-Core Intel Core i9). Given the nature of this API we are bound to possible rate limiting by upstream and also we can't outperform upstream in terms of latencies.  
The test can be run using `make load-test` and it will output a report:
```
Requests      [total, rate, throughput]         125259, 2080.26, 2058.73
Duration      [total, attack, wait]             1m1s, 1m0s, 594.307ms
Latencies     [min, mean, 50, 90, 95, 99, max]  869.358µs, 240.877ms, 137.284ms, 417.045ms, 686.073ms, 1.685s, 4.958s
Bytes In      [total, mean]                     13798881, 110.16
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100%
Status Codes  [code:count]                      200:125259
Error Set:  
```
The API has been tested up to 2k rps with no errors and it is expected to scale linearly.

## Improvements And Limitations
* Gorilla Mux can be replaced with [httprouter](https://github.com/julienschmidt/httprouter) if performances are of concern.
* There currently is no support for authentication / client secret.
* Only access logs are being logged to Stdout. Service logs should be implemented for better monitoring.
* No metrics are exposed. The API should be instrumented to expose them in [Prometheus](github.com/prometheus/client_golang/prometheus) format. 