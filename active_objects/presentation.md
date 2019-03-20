class: center, middle

# Active objects
### How to remove mutexes from your codebase



March 25, 2019, Alexander Stavonin, alex@sysdev.me
---
# Problem description

- Go gives us very chip concurrency;
- Some applications has _state_;
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
- But is it possible to have something like...?

```go
type SomeStruct struct {
	
    someMap map[int]string
    someArray []string
    
    ???
}
```
---

# Why this transformation is needed? (1/2)

### Races and deadlocks:

- Developer should carefully choose big enough code block protected by lock;

    - How big protected area should be? Whole function? `if` block?

--

- Developer should choose right lock type (`RLock` or `Lock`);

    - should I call `someMapLock.Lock()` or `someMapLock.RLock()` and why?

--

- Developer must not forget to unlock mutex.

    - `defer` helps a lot, but it works only in function scope, but not `if` or `for` scopes.

---
# Why this transformation is needed? (2/2)

### Amdahl's law

$$S_\text{latency}(s) = \frac 1 {(1 - p) + \frac p s}$$

where:

- S<sub>latency</sub> is the theoretical speedup of the execution of the whole task;
- `s` is the speedup of the part of the task that benefits from improved system resources (parallel execution count);
- `p` is the proportion of execution time that the part benefiting from improved resources originally occupied (in %).

Impact: if you can run in  parallel 25% of your code on 10 executors (threads) then **theoretically** achievable speedup is 3.077 and only 3.883 for 100 executors!

---
# Concurrent interactions and communication

There are 2 main types of interaction in concurrent applications:

- Shared memory;

- Message passing:
    
    - Communicating sequential processes (CSP);
    
    - Actors.

---
# Concurrent interactions and communication
### Shared memory

- Most common way for cross process/thread interaction in concurrent world;

- Supported by almost all languages (at least for cross thread communication);

- Fastest way across all other techniques;

- Error prone: deadlocks and races are very common for shared memory access.

---
# Concurrent interactions and communication
### Communicating sequential processes (CSP)


---
# Concurrent interactions and communication
### Actors

---
# Questions and comments are welcome! :)