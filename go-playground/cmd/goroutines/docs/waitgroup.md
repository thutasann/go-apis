# sync.WaitGroup

```
+------------------+
|  main goroutine  |
+------------------+
          |
          | calls wg.Add(3)
          |
          |-------------------+------------------+------------------+
          |                   |                  |
     goroutine A         goroutine B        goroutine C
          |                   |                  |
     does work           does work          does work
          |                   |                  |
     wg.Done()           wg.Done()          wg.Done()
          |                   |                  |
          +-------------------+------------------+
                              |
                          wg.Wait() unblocks

```
