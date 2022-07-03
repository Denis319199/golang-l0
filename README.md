# Service

## Repo Structure

- **/service** - the service source code
- **/utility** - source code of the utility providing possibility to publish 
files content located in some directory

## Service Endpoints

- **GET /order** | **Req params: page int; size int** - get all orders

- **GET /order/{OrderUid}** - get all order by _OrderUid_

## Utility

The utility opens the given directory and publishes each file to NATS Streaming 

```
./utility <nats_url> <cluster_id> <dir>
```

Some test data for this utility are located in [data directory](./utility/data)

**Example:**

Utility allows to publish elements of a JSON array one by one

```
./utility 127.0.0.1:4222 mycluster ./data
```

## Nats Streaming

To start a nats streaming server use [docker-compose.yml](./service/resources/docker-compose.yml)

## Tests

### Unit Tests

Write some unit tests that covers **44.6%** of statements

### Load Testing

Tested service through the **vegeta** tool:

```
cat ./service/resources/vegeta_requests.txt | vegeta attack -format=http -duration=5s --output results.bin
```

Requests that was tested reside in [vegeta_requests.txt](./service/resources/vegeta_requests.txt)

The results I got through `vegeta report results.bin` are listed below:

```
Requests      [total, rate, throughput]  250, 50.23, 50.23
Duration      [total, attack, wait]      4.9768656s, 4.9768656s, 0s
Latencies     [mean, 50, 95, 99, max]    589.917µs, 562.966µs, 1.2505ms, 1.8133ms, 6.4169ms
Bytes In      [total, mean]              2678500, 10714.00
Bytes Out     [total, mean]              0, 0.00
Success       [ratio]                    100.00%
Status Codes  [code:count]               200:250  
Error Set:
```