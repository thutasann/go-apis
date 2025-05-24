# Go Cache

Building own simple cache in GoLang

## LRU

Maintains a specific length of items in cache
LRU - Only the most recently used items are stored

for a true LRU cache -

1. if an item already exists, we need to remove it and add it to the beginning
2. an order the items is maintained
3. deletion happens at the tail and addition happens at the end.
