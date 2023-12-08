# Cache benchmark for Go

This benchmark compares cache eviction algorithms using Zipfian distribution.

## Cache libraries used
| Name           | Ref                                           |
|----------------|-----------------------------------------------|
| s3-fifo        | https://github.com/scalalang2/golang-fifo     |
| clock          | https://github.com/Code-Hex/go-generics-cache |
| lru-hashicorp  | https://github.com/hashicorp/golang-lru       |
| lru-groupcache | https://github.com/golang/groupcache/lru      |
| two-queue      | https://github.com/hashicorp/golang-lru       |
| s4-lru         | https://github.com/dgryski/go-s4lru           |
| tinylfu        | https://github.com/dgryski/go-tinylfu         |

## Benchmark result
- The **golang-fifo** algorithm shows the best cache efficiency (= low miss ratio) relative to other LRU-based algorithm
when cache size is much smaller than the total size of item set.
  - golang-fifo is a implementation of modern cache eviction algorithm introduced in following papers.
  - [FIFO queues are all you need for cache eviction | ACM SOSP'23](https://dl.acm.org/doi/10.1145/3600006.3613147)
  - [SIEVE is Simpler than LRU: an Efficient Turn-Key Eviction Algorithm for Web Caches | USENIX NSDI'24](https://junchengyang.com/publication/nsdi24-SIEVE.pdf)
- As the cache size becomes closer to the overall item size, the efficiency difference between cache algorithms decreases.

```shell
$ go run *.go

results:
itemSize=500000, workloads=2500000, cacheSize=0.10%, zipf's alpha=0.99  

      CACHE      | HITRATE | MEMORY  |  DURATION  |  HITS   | MISSES   
-----------------+---------+---------+------------+---------+----------
  s3-fifo        | 48.85%  | 0.20MiB | 857.0013ms | 1221268 | 1278732  
  tinylfu        | 47.36%  | 0.10MiB | 830.4299ms | 1183944 | 1316056  
  slru           | 46.39%  | 0.11MiB | 930.5827ms | 1159650 | 1340350  
  s4lru          | 46.09%  | 0.11MiB | 703.051ms  | 1152255 | 1347745  
  two-queue      | 45.48%  | 0.16MiB | 1.2742908s | 1136954 | 1363046  
  clock          | 37.35%  | 0.10MiB | 846.2553ms |  933784 | 1566216  
  lru-hashicorp  | 36.60%  | 0.13MiB | 927.4414ms |  915083 | 1584917  
  lru-groupcache | 36.60%  | 0.11MiB | 973.254ms  |  915083 | 1584917  


results:
itemSize=500000, workloads=2500000, cacheSize=1.00%, zipf's alpha=0.99  

      CACHE      | HITRATE | MEMORY  |  DURATION  |  HITS   | MISSES   
-----------------+---------+---------+------------+---------+----------
  s3-fifo        | 64.28%  | 1.86MiB | 839.5197ms | 1606974 |  893026  
  tinylfu        | 63.77%  | 0.97MiB | 878.5349ms | 1594331 |  905669  
  slru           | 62.31%  | 1.05MiB | 952.0577ms | 1557846 |  942154
  s4lru          | 62.26%  | 1.05MiB | 753.4121ms | 1556430 |  943570
  two-queue      | 61.88%  | 1.45MiB | 1.1751269s | 1547014 |  952986
  clock          | 56.09%  | 0.88MiB | 846.6917ms | 1402170 | 1097830
  lru-hashicorp  | 55.36%  | 1.02MiB | 868.35ms   | 1383915 | 1116085
  lru-groupcache | 55.36%  | 1.06MiB | 869.5358ms | 1383915 | 1116085


results:
itemSize=500000, workloads=2500000, cacheSize=10.00%, zipf's alpha=0.99

      CACHE      | HITRATE |  MEMORY  |  DURATION  |  HITS   | MISSES
-----------------+---------+----------+------------+---------+---------
  tinylfu        | 79.05%  | 9.99MiB  | 1.0189901s | 1976206 | 523794
  s3-fifo        | 78.78%  | 12.81MiB | 927.0793ms | 1969553 | 530447
  two-queue      | 78.09%  | 14.43MiB | 1.246231s  | 1952361 | 547639
  s4lru          | 76.53%  | 10.04MiB | 866.5103ms | 1913318 | 586682
  clock          | 76.04%  | 8.65MiB  | 803.2425ms | 1901104 | 598896
  slru           | 75.58%  | 9.67MiB  | 1.0344438s | 1889583 | 610417
  lru-hashicorp  | 75.52%  | 10.22MiB | 930.0334ms | 1888047 | 611953
  lru-groupcache | 75.52%  | 10.27MiB | 985.3425ms | 1888047 | 611953
```
