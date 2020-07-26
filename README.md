[![Run on Repl.it](https://repl.it/badge/github/matt4biz/go-class-profile)](https://repl.it/github/matt4biz/go-class-profile)

# Go class: Profile example
This example creates server which gets TODOs from Typicode's test server, but leaks goroutines and sockets in the process.

It incorporates both pprof and Prometheus to show different ways to see the leak.

Run the server with

`go run .`

and access it at `http://localhost:8080`

You can exercise it as

```shell
$ curl http://localhost:8080/1
[ ] 1 - delectus aut autem

$ curl http://localhost:8080/2
[ ] 2 - quis ut nam facilis et officia qui

$ curl http://localhost:8080/3
[ ] 3 - fugiat veniam minus

$ curl http://localhost:8080/4
[x] 4 - et porro tempora

$ curl http://localhost:8080/5
[ ] 5 - laboriosam mollitia et enim quasi adipisci quia provident illum

$ curl http://localhost:8080/6
[ ] 6 - qui ullam ratione quibusdam voluptatem quia omnis

$ curl http://localhost:8080/7
[ ] 7 - illo expedita consequatur quia in

$ curl http://localhost:8080/20
[x] 20 - ullam nobis libero sapiente ad optio sint
```

and then check the number of queries, e.g.

```shell
$ curl -s http://localhost:8080/metrics | head -3
# HELP all_queries How many queries we've received.
# TYPE all_queries counter
all_queries 8
```

or the goroutines with 

```shell
$ curl -s http://localhost:8080/metrics | grep goroutine
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 15
```

or just open the metrics page [http://localhost:8080/metrics](http://localhost:8080/metrics).

To see pprof, you'll need to open your browser to [http://localhost:8080/debug/pprof](http://localhost:8080/debug/pprof) where you can see goroutines hung on the net poller:

	goroutine profile: total 13
	9 @ 0x1038770 0x1031e0a 0x1031375 0x10c7aa5 0x10c89a1 0x10c8983 0x11b721f 0x11c8e0e 0x1220810 0x10ee2b1 0x1220a5c 0x121f015 0x122313b 0x1223146 0x1171def 0x1075307 0x12a1ce7 0x12a1c9a 0x12a2521 0x12c2a0d 0x12c21ef 0x1067c51
	#	0x1031374	internal/poll.runtime_pollWait+0x54        /usr/local/Cellar/go/1.14.4/libexec/src/runtime/netpoll.go:203
	#	0x10c7aa4	internal/poll.(*pollDesc).wait+0x44        /usr/local/Cellar/go/1.14.4/libexec/src/internal/poll/fd_poll_runtime.go:87
	#	0x10c89a0	internal/poll.(*pollDesc).waitRead+0x200   /usr/local/Cellar/go/1.14.4/libexec/src/internal/poll/fd_poll_runtime.go:92
	. . .
