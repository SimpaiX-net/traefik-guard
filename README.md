<!-- header -->

<div align="center">   
    <div>
        <img src="img/guard.png" width=300 style="border: 2px solid grey;"><br>
    </div>
</div>

## Intro
--- **Guard** is an elegant IPQS plugin for Traefik. Acting as a middleware or microservice between your web server.


--- **Questions?** feel free to ask by [contacting me](https://t.me/z3ntl3)! 

### Install
[Instructions](https://plugins.traefik.io/install)

### Configuration
```yaml
proxy: # optional, good to set a rotating proxy to prevent rate limits
ttl: 24h # optional, otherwise defaults to 1 week of TTL
timeout: 300ms # must be set
ipHeaders: # must be set
  - "cf-connecting-ip" # if your backend hides behind CF's reverse proxies
  - "X-Forwarded-For"
```

### Additional notes
Guard uses **InternetDB** to determine the reputation of an ioT device. It's completely free, and allows high traffic throughput. You can always use ``proxy``to allow a limitless quota when needed. 

To be fast and not halter or negatively impact your avg response times while sitting as an intermediary between your backend, Guard is effectively using an in memory-cache.

Here's the performance benchmark (for the in memory cache):
```
Running tool: C:\Program Files\Go\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkClient$ github.com/SimpaiX-net/ipqs/tests

goos: windows
goarch: amd64
pkg: github.com/SimpaiX-net/ipqs/tests
cpu: AMD Ryzen 7 4800H with Radeon Graphics         
BenchmarkClient-16    	 8923340	       135.7 ns/op	     256 B/op	       4 allocs/op
PASS
ok  	github.com/SimpaiX-net/ipqs/tests	2.911s

```


### Credits
--- Programmed by [z3ntl3](https://z3ntl3.com)
