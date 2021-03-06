class: center, middle

# Active objects
### How to remove mutexes from your codebase



March 25, 2019, Alexander Stavonin, alex@sysdev.me
---
# Problem description

- Go gives us very chip concurrency.

- Some applications have _state_.

- Chip concurrency + states often leads us to...

--

```go
type SomeStruct struct {
	
    someMap map[int]string
    someMapLock sync.RWMutex
    
    someArray []string
    someArrayLock sync.RWMutex
}
```

--
- But is it possible to have something like this?

```go
type SomeStruct struct {
	
    someMap map[int]string
    someArray []string
}
```
---

# Motivation (1/2)

### Races and deadlocks:

- Developer should choose carefully code block which is big enough and protect it by lock.

    - How big protected area should be? Whole function? `if` block?

--

- Developer should choose right lock type (`RLock` or `Lock`).

    - should I call `someMapLock.Lock()` or `someMapLock.RLock()` and why?

--

- Developer must not forget to unlock mutex.

    - `defer` helps a lot, but it works only in function scope, but not `if` or `for` scopes.

---
# Motivation (2/2)

### Amdahl's law

$$S_\text{latency}(s) = \frac 1 {(1 - p) + \frac p s}$$

where:

- S<sub>latency</sub> is the theoretical speedup of the execution of the whole task.
- `s` is the speedup of the part of the task that benefits from improved system resources (parallel execution count).
- `p` is the proportion of execution time that the part benefiting from improved resources originally occupied (in %).

Impact: if you can run 25% of your code in parallel on 10 executors (threads) then **theoretically** achievable speedup is 3.077 and only 3.883 for 100 executors!

---
# Concurrent interactions

There are 2 main types of interactions in concurrent applications:

- Shared memory;

- Message passing:
    
    - Communicating sequential processes (CSP).
    
    - Actors model.

---
# Concurrent interactions
### Shared memory

- Most common way for cross process/thread interaction in concurrent world.

- Supported by almost all languages (at least for cross thread communication).

- Fastest way across all other techniques.

- Error prone: deadlocks and races are very common for shared memory access.

---
# Concurrent interactions
### Communicating sequential processes (CSP)

- CSP was first described in a 1978 paper by Tony Hoare and significantly improved in 1985.

- Go is one of multiple languages that support CSP.
    - CSP is also implemented in Rust, Clojure, OCaml, etc.
    
- `Channel`s as a cross process/thread communication way is the core component of CSP.

- All modern CSP implementations are much more secure then Shared memory.

- Deadlocks are still possible as any Goroutine may have multiple Channels.

---
# Concurrent interactions
### Actors

- Very similar to CSP, but...

    - CSP processes are anonymous, while actors have identities.
    
    - CSP uses explicit channels for message passing, whereas actor systems transmit messages to named destination actors.

    - *Invalid for Go CSP implementation statement*: CSP message-passing fundamentally involves a rendezvous between the processes involved in sending and receiving the message, i.e. the sender cannot transmit a message until the receiver is ready to accept it.

- Some languages (Erlang, for instance) have Actors model implementation out of the box

- or have very famous implementation like AKKA from Scala.

---
# What is Active Object?

The idea of Active Object is based on Actors model and was initially described in Pattern-Oriented Software Architecture (POSA book, 1996).

![](./ao.svg)

---
#Communication interface

Each request should have unique type as we need to implement behavior similar to function call.

```go
type requestCommand int

const(
	genNum requestCommand = iota
	runCmd
)
```

--
Request is not only command type, but also command data (if any) and the way to push response back.
```go
type request struct {
	cmd requestCommand
	data interface {}
	out chan response
}
```

--
Processor is able to generate any type of response. In lack of metaprogramming `interface {}` is a good option.
```go
type response interface {}
```

---
# Requests processing 1/3
### Main processor

```go
func processRequest(req request)  {
	switch req.cmd {
	case genNum:
		req.out <- rand.Intn(100)
	case runCmd:
		go doExec(req)
	}
}
```
where:
- `req.cmd` is command type.
- `req.out` is outcoming channel (the way to tell caller result).
- fast command `genNum` should be executed "in place".
- we need special handling for slow command `runCmd`.

---
# Requests processing 2/3
### Slow commands

```go
func doExec(req request) {
	time.Sleep(1*time.Second)

	cmd := req.data.(string)
	out, err := exec.Command(cmd).Output()
	if err != nil {
		req.out <- err.Error()
	} else {
		req.out <- string(out)
	}
}
```
where:
- `Sleep` is just for illustration proposes of delay.
- `exec.Command` is external command execution.
- `req.out` is outcoming channel (same as earlier).

---
# Requests processing 3/3
## Event loop

```go
func eventLoop(done chan struct{}, inCh chan request) {
	for {
		select {
		case req := <-inCh:
			processRequest(req)
		case <-done:
			break
		}
	}
}
```
where:
- `inCh` is an input channel. Input channel is the only way to interact with Processor.
- `done` is signal for shutting processor down.

---
# Request sending 1/2
### Event loop initialization

```go
	done := make(chan struct{})
	inCh := make(chan request, 10)

	go eventLoop(done, inCh)
```
where:
- `inCh` shouldn't block callers, which means buffer size is reasonably big.
- `eventLoop` has own thread and `inCh` is the only way to communicate with it.

---
# Request sending 2/2
### Requests and response

```go
	resp := make(chan response, 1)

	inCh <- request{
		runCmd,
		"ls",
		resp,
	}
	fmt.Println(<-resp1)
```
where:
- `resp` is response channel with at least one element for non-blocking response delivery.
---
# Reusing asynchronously calculated results

```go
func doExec(req request, inCh chan request) {
	cmd := req.data.(string)
	out, _ := exec.Command(cmd).Output()
    inCh <- request{
        runCmdResult,
        out,
        nil,
    }
...
```
where:
- `runCmdResult` is a new command with asynchronously calculated results.
- `processRequest` call should process this new command and, for example, to cache calculated result.

---
# Questions and comments are welcome! :)

```




Demo and slides:
    
    https://github.com/astavonin/presentations/
```