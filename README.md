# Cache benchmark for Go

This benchmark compares cache eviction algorithms using Zipfian distribution.

This referred to [this code](https://github.com/vmihailenco/go-cache-benchmark) a lot.

## Cache libraries used
| Name            | Ref                                           |
|-----------------|-----------------------------------------------|
| sieve           | https://github.com/scalalang2/golang-fifo     |
| s3-fifo         | https://github.com/scalalang2/golang-fifo     |
| s3-fifo (otter) | https://github.com/maypok86/otter             |
| clock           | https://github.com/Code-Hex/go-generics-cache |
| lru-hashicorp   | https://github.com/hashicorp/golang-lru       |
| lru-groupcache  | https://github.com/golang/groupcache/lru      |
| two-queue       | https://github.com/hashicorp/golang-lru       |
| s4-lru          | https://github.com/dgryski/go-s4lru           |
| tinylfu         | https://github.com/dgryski/go-tinylfu         |

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

      CACHE      | HITRATE | MEMORY  |   QPS   |  HITS   | MISSES   
-----------------+---------+---------+---------+---------+----------
  sieve          | 47.41%  | 0.11MiB | 2756340 | 3555962 | 3944038  
  tinylfu        | 47.37%  | 0.11MiB | 2003205 | 3552824 | 3947176  
  s3-fifo        | 47.17%  | 0.21MiB | 1835985 | 3537405 | 3962595  
  slru           | 46.49%  | 0.11MiB | 1877817 | 3486475 | 4013525  
  s4lru          | 46.13%  | 0.12MiB | 2351834 | 3459613 | 4040387  
  two-queue      | 45.49%  | 0.18MiB | 1821715 | 3411796 | 4088204  
  otter          | 45.02%  | 0.52MiB | 3810976 | 3376497 | 4123503  
  clock          | 37.34%  | 0.10MiB | 1936483 | 2800279 | 4699721  
  lru-groupcache | 36.57%  | 0.12MiB | 1843658 | 2742635 | 4757365  
  lru-hashicorp  | 36.56%  | 0.08MiB | 2006957 | 2741779 | 4758221


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
