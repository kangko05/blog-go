# Local Performance Test Results

## Test Environment

- OS: Arch Linux (EndeavourOS)
- Go Version: go1.25.0
- Database: SQLite (local file)
- Test Tool: Apache Bench 2.3
- Date: 2025.9.16

## Summary

| Endpoint     | RPS    | Avg Response Time | Description              |
| ------------ | ------ | ----------------- | ------------------------ |
| /checkhealth | 50,728 | 0.197ms           | Health check (no DB)     |
| /posts/      | 37,469 | 0.267ms           | List all posts (DB scan) |
| /posts/1     | 39,082 | 0.256ms           | Single post (DB index)   |

## Detailed Results

### Health Check Endpoint

```plain
This is ApacheBench, Version 2.3 <$Revision: 1923142 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)


Server Software:
Server Hostname:        localhost
Server Port:            8000

Document Path:          /checkhealth
Document Length:        0 bytes

Concurrency Level:      10
Time taken for tests:   0.020 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      75000 bytes
HTML transferred:       0 bytes
Requests per second:    50727.95 [#/sec] (mean)
Time per request:       0.197 [ms] (mean)
Time per request:       0.020 [ms] (mean, across all concurrent requests)
Transfer rate:          3715.43 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.0      0       0
Processing:     0    0   0.0      0       0
Waiting:        0    0   0.0      0       0
Total:          0    0   0.0      0       0

Percentage of the requests served within a certain time (ms)
  50%      0
  66%      0
  75%      0
  80%      0
  90%      0
  95%      0
  98%      0
  99%      0
 100%      0 (longest request)

```

### Posts List Endpoint

```plain
This is ApacheBench, Version 2.3 <$Revision: 1923142 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)


Server Software:
Server Hostname:        localhost
Server Port:            8000

Document Path:          /posts/
Document Length:        417 bytes

Concurrency Level:      10
Time taken for tests:   0.027 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      541000 bytes
HTML transferred:       417000 bytes
Requests per second:    37468.62 [#/sec] (mean)
Time per request:       0.267 [ms] (mean)
Time per request:       0.027 [ms] (mean, across all concurrent requests)
Transfer rate:          19795.43 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.0      0       0
Processing:     0    0   0.1      0       1
Waiting:        0    0   0.1      0       1
Total:          0    0   0.1      0       1

Percentage of the requests served within a certain time (ms)
  50%      0
  66%      0
  75%      0
  80%      0
  90%      0
  95%      0
  98%      1
  99%      1
 100%      1 (longest request)
```

### Single Post Endpoint

```plain
This is ApacheBench, Version 2.3 <$Revision: 1923142 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)


Server Software:
Server Hostname:        localhost
Server Port:            8000

Document Path:          /posts/1
Document Length:        204 bytes

Concurrency Level:      10
Time taken for tests:   0.026 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      328000 bytes
HTML transferred:       204000 bytes
Requests per second:    39082.35 [#/sec] (mean)
Time per request:       0.256 [ms] (mean)
Time per request:       0.026 [ms] (mean, across all concurrent requests)
Transfer rate:          12518.56 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.0      0       0
Processing:     0    0   0.1      0       1
Waiting:        0    0   0.1      0       1
Total:          0    0   0.1      0       1

Percentage of the requests served within a certain time (ms)
  50%      0
  66%      0
  75%      0
  80%      0
  90%      0
  95%      0
  98%      1
  99%      1
 100%      1 (longest request)
```

### Stress Test

```plain
This is ApacheBench, Version 2.3 <$Revision: 1923142 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)


Server Software:
Server Hostname:        localhost
Server Port:            8000

Document Path:          /posts/
Document Length:        417 bytes

Concurrency Level:      500
Time taken for tests:   0.028 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      541000 bytes
HTML transferred:       417000 bytes
Requests per second:    35696.44 [#/sec] (mean)
Time per request:       14.007 [ms] (mean)
Time per request:       0.028 [ms] (mean, across all concurrent requests)
Transfer rate:          18859.15 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    3   1.5      4       5
Processing:     2    8   3.5      7      18
Waiting:        0    7   3.2      6      16
Total:          4   11   3.3     11      23

Percentage of the requests served within a certain time (ms)
  50%     11
  66%     13
  75%     14
  80%     14
  90%     16
  95%     17
  98%     21
  99%     22
 100%     23 (longest request)
```
