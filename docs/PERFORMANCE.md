# Performance

Tested on AMD EPYC 9274F, 256 GB RAM, 2×2 TB NVMe, 10 GbE.

|
 Endpoint 
|
 P99 latency 
|
 RPS/node 
|
|
----------
|
-------------
|
----------
|
|
 getSignaturesForAddress (2 T rows) 
|
 82 ms 
|
 12 k 
|
|
 getBlock 
|
 3 ms 
|
 50 k 
|
|
 getTransaction 
|
 4 ms 
|
 45 k 
|

Scaling: add nodes + shards; fractal hash keeps data local.