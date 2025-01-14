Task 3:

The resulting value of i is not zero. This is due to swapping. The result is dependent on when the operations is executed in time. In this case the work of the increment function can be overwritten by the decrement function.

Task 4:

C - In this case we should use mutex to achive the correct behaviour. Mutex will block another thread from alternating a shared variable, ensuring the correct result. Semaphores could be used, but are more versitile, and will not simplest solution in this case.



