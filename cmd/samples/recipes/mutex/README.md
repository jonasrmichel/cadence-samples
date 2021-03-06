This mutex workflow demos an ability to lock/unlock a particular resource within a particular cadence domain
so that other workflows within the same domain would wait until a resource lock is released. This is useful 
when we want to avoid race conditions or parallel mutually exclusive operations on the same resource.

One way of coordinating parallel processing is to use cadence signals with SignalWithStartWorkflow and
make sure signals are getting processed sequentially, however the logic might become too complex if we
need to lock two or more resources at the same time. Mutex workflow pattern can simplify that.

This example enqueues two long running SampleWorkflowWithMutex workflows in parallel. And each of the workflows has a mutex section. 
When SampleWorkflowWithMutex reaches Mutex section, it starts a mutex workflow via local activity, and blocks until
"acquire-lock-event" is received. Once "acquire-lock-event" is received, it enters critical section,
and finally releases the lock once processing is over by sending "releaseLock" a signal to the MutexWorkflow.


### Steps to run this sample:
1) You need a cadence service running. See details in cmd/samples/README.md
2) Run the following command to start the worker

    ```
    ./bin/mutex -m worker
    ```

3) Run the following command to start the example

    ```
    ./bin/mutex -m trigger
    ```
  
You should see that second workflow critical section is executed when first workflow
critical operation is finished.
