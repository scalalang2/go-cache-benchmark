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
itemSize=500000, workloads=7500000, cacheSize=0.10%, zipf's alpha=0.99, concurrency=1

      CACHE      | HITRATE | MEMORY  |   QPS   |  HITS   | MISSES   
-----------------+---------+---------+---------+---------+----------
  sieve          | 47.42%  | 0.11MiB | 4251701 | 3556217 | 3943783  
  otter          | 47.40%  | 0.45MiB | 3465804 | 3555090 | 3944910  
  tinylfu        | 47.37%  | 0.10MiB | 4242081 | 3553085 | 3946915  
  s3-fifo        | 47.16%  | 0.21MiB | 2683363 | 3537110 | 3962890  
  slru           | 46.48%  | 0.11MiB | 3748126 | 3486209 | 4013791  
  s4lru          | 46.15%  | 0.12MiB | 4934211 | 3461183 | 4038817  
  two-queue      | 45.49%  | 0.16MiB | 2647370 | 3411961 | 4088039  
  clock          | 37.34%  | 0.10MiB | 3890041 | 2800767 | 4699233  
  lru-hashicorp  | 36.59%  | 0.08MiB | 3772636 | 2744181 | 4755819  
  lru-groupcache | 36.59%  | 0.11MiB | 3757515 | 2744181 | 4755819  


results:
itemSize=500000, workloads=7500000, cacheSize=0.10%, zipf's alpha=0.99, concurrency=2

      CACHE      | HITRATE | MEMORY  |   QPS   |  HITS   | MISSES   
-----------------+---------+---------+---------+---------+----------
  sieve          | 47.66%  | 0.12MiB | 2967946 | 3574417 | 3925583  
  otter          | 47.44%  | 0.50MiB | 4263786 | 3558030 | 3941970  
  tinylfu        | 47.36%  | 0.11MiB | 2773669 | 3551639 | 3948361  
  s3-fifo        | 47.16%  | 0.20MiB | 1889645 | 3537211 | 3962789  
  slru           | 46.48%  | 0.11MiB | 2528658 | 3486110 | 4013890  
  s4lru          | 46.15%  | 0.12MiB | 3321523 | 3461103 | 4038897  
  two-queue      | 45.49%  | 0.15MiB | 1823930 | 3411948 | 4088052  
  clock          | 37.34%  | 0.09MiB | 2530364 | 2800765 | 4699235  
  lru-groupcache | 36.59%  | 0.12MiB | 2454992 | 2744121 | 4755879  
  lru-hashicorp  | 36.59%  | 0.08MiB | 2707581 | 2744083 | 4755917  


results:
itemSize=500000, workloads=7500000, cacheSize=0.10%, zipf's alpha=0.99, concurrency=4

      CACHE      | HITRATE | MEMORY  |   QPS   |  HITS   | MISSES   
-----------------+---------+---------+---------+---------+----------
  sieve          | 47.47%  | 0.12MiB | 2711497 | 3560304 | 3939696  
  tinylfu        | 47.37%  | 0.11MiB | 2577320 | 3553093 | 3946907  
  s3-fifo        | 47.16%  | 0.21MiB | 1740139 | 3537120 | 3962880  
  slru           | 46.48%  | 0.12MiB | 2436647 | 3486075 | 4013925  
  s4lru          | 46.15%  | 0.11MiB | 3038898 | 3461004 | 4038996  
  two-queue      | 45.49%  | 0.17MiB | 1851852 | 3411875 | 4088125  
  otter          | 45.21%  | 0.55MiB | 5040323 | 3390975 | 4109025  
  clock          | 37.34%  | 0.09MiB | 2480979 | 2800620 | 4699380  
  lru-groupcache | 36.58%  | 0.11MiB | 2359232 | 2743834 | 4756166  
  lru-hashicorp  | 36.58%  | 0.08MiB | 2515934 | 2743779 | 4756221  


results:
itemSize=500000, workloads=7500000, cacheSize=0.10%, zipf's alpha=0.99, concurrency=8

      CACHE      | HITRATE | MEMORY  |   QPS   |  HITS   | MISSES   
-----------------+---------+---------+---------+---------+----------
  sieve          | 47.65%  | 0.11MiB | 2757353 | 3573923 | 3926077  
  tinylfu        | 47.38%  | 0.11MiB | 2154553 | 3553231 | 3946769  
  s3-fifo        | 47.16%  | 0.22MiB | 1829268 | 3537149 | 3962851  
  slru           | 46.48%  | 0.11MiB | 2103197 | 3486051 | 4013949  
  s4lru          | 46.15%  | 0.11MiB | 2534640 | 3460985 | 4039015  
  two-queue      | 45.49%  | 0.16MiB | 1823930 | 3411853 | 4088147  
  otter          | 44.83%  | 0.54MiB | 4148230 | 3362425 | 4137575  
  clock          | 37.34%  | 0.08MiB | 2103787 | 2800520 | 4699480  
  lru-groupcache | 36.58%  | 0.13MiB | 2046385 | 2743287 | 4756713  
  lru-hashicorp  | 36.57%  | 0.06MiB | 2116850 | 2742948 | 4757052  


results:
itemSize=500000, workloads=7500000, cacheSize=0.10%, zipf's alpha=0.99, concurrency=16

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
```
