# Service

## Repo Structure

- **/service** - the service source code
- **/utility** - source code of the utility providing possibility to publish 
elements of a JSON array to nats streaming

## Service Endpoints

- **GET /order** - get all orders

- **GET /order/{OrderUid}** - get all order by _OrderUid_

## Utility

```
./utility <nats_url> <cluster_id> <file>
```

Some test data for this utility are located in [data.json](./utility/data.json)

**Example:**

Utility allows to publish elements of a JSON array one by one

```
./utility 127.0.0.1:4222 mycluster ./data.json
```

## Nats Streaming

To start a nats streaming server use [docker-compose.yml](/service/resources/docker-compose.yml)

## Tests

### Unit Tests

Write some unit tests that covers **44.1%** of statements

### Load Testing

Tested service through the **vegeta** tool:

```
cat ./service/resources/vegeta_requests.txt | vegeta attack -format=http -duration=5s --output results.bin
```

Requests that was tested reside in [vegeta_requests.txt](./service/resources/vegeta_requests.txt)

The results I got through `vegeta report results.bin` are listed below:

```
Requests      [total, rate, throughput]  250, 50.20, 50.20
Duration      [total, attack, wait]      4.9796331s, 4.9796331s, 0s
Latencies     [mean, 50, 95, 99, max]    547.049µs, 552.733µs, 1.0795ms, 1.1928ms, 5.903ms
Bytes In      [total, mean]              2376000, 9504.00
Bytes Out     [total, mean]              0, 0.00
Success       [ratio]                    100.00%
Status Codes  [code:count]               200:250  
Error Set:
```