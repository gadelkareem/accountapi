# Account API Library
[![Build Status](https://travis-ci.org/gadelkareem/accountapi.svg)](https://travis-ci.org/gadelkareem/accountapi)

Account API library for Form3 in Golang

```bash
> go test -v -bench=. -benchmem
go test -v -bench=. -benchmem
=== RUN   TestAccountClient_Fetch
=== PAUSE TestAccountClient_Fetch
=== RUN   TestAccountClient_List
=== PAUSE TestAccountClient_List
=== RUN   TestAccountClient_Create
=== PAUSE TestAccountClient_Create
=== RUN   TestAccountClient_Delete
=== PAUSE TestAccountClient_Delete
=== CONT  TestAccountClient_Fetch
=== CONT  TestAccountClient_Delete
=== CONT  TestAccountClient_Create
=== CONT  TestAccountClient_List
--- PASS: TestAccountClient_Delete (0.00s)
--- PASS: TestAccountClient_List (0.00s)
--- PASS: TestAccountClient_Fetch (0.00s)
--- PASS: TestAccountClient_Create (0.00s)
goos: darwin
goarch: amd64
pkg: github.com/gadelkareem/accountapi/v1
BenchmarkAccountClient_Fetch-8    	    8138	    132862 ns/op	   13345 B/op	     218 allocs/op
BenchmarkAccountClient_List-8     	    8766	    147909 ns/op	   12285 B/op	     212 allocs/op
BenchmarkAccountClient_Create-8   	    5265	    237498 ns/op	   20133 B/op	     382 allocs/op
BenchmarkAccountClient_Delete-8   	   12759	    111243 ns/op	    5333 B/op	      63 allocs/op
PASS
ok  	github.com/gadelkareem/accountapi/v1	7.072s
```

Some BDD tests are also available.


 