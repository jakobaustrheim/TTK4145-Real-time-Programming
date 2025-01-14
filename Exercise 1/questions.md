Exercise 1 - Theory questions
-----------------------------

### Concepts

What is the difference between *concurrency* and *parallelism*?
> Concurrency is executing several tasks on a singular cpu at once by carefully switching between the tasks. Parallelism is when two or more cpu's execute tasks at the same time.

What is the difference between a *race condition* and a *data race*? 
> A race condition occures when the correctnes of a code is depenendent on the timing of the execution of tasks. A data race occures when a opperation accesses a mutable object at the same time as another operation is already working on it.
 
*Very* roughly - what does a *scheduler* do, and how does it do it?
> A scheduler determines in what order processes should be executed on the cpu. Premtive scheduling is a scheuduler that can stop process during execution in order to run another process if the process exceeds its time-slice, set by the scheduler.


### Engineering

Why would we use multiple threads? What kinds of problems do threads solve?
> Threads ensures responsive programs that handle inputs without delay. This is imortant in real-time programming. It does this by concurrent execution of tasks.

Some languages support "fibers" (sometimes called "green threads") or "coroutines"? What are they, and why would we rather use them over threads?
> Fibers is a lightweight thread-system that utilize cooperative syncronisation, insted of the premtive scheduling used by threads. As a result the thread safety is less of a issue using fibers, meaning less risk of race conditions.

Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?
> Conurrancy allows for more efficiant programs that respons quick to I/O messages. However, the need to write the code "race condition proof" complicates the programming.

What do you think is best - *shared variables* or *message passing*?
> I prefere shared varbiables compared to message passing, due to the fact that i found it to be the easiest method to understand, and therefore implement.


