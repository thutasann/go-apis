## Memory Visualized (Heap Diagram)

```
Memory Heap:
┌────────────┬──────────────────────────────────────┐
│ Address    │ Contents                             │
├────────────┼──────────────────────────────────────┤
│ 0xABC123   │ { id: '123', name: 'Thuta' }         │  ← original object
└────────────┴──────────────────────────────────────┘

Variables (Stack):
┌────────────────────┬────────────┐
│ Variable            │ Pointer    │
├────────────────────┼────────────┤
│ user                │ 0xABC123   │
│ store.all['abc']    │ 0xABC123   │
│ store.byId['123']   │ 0xABC123   │
└────────────────────┴────────────┘
```

No matter how many times you assign the same object reference — only one memory block exists.
