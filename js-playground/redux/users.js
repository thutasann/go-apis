/**
 * You're not creating 2 users in memory.
 * Both `store.all['abc']` and `store.byId['123']` point to the same object
 * 
 * ```.md
 *         ┌────────────┐
store ─▶│ all        │
        │ byId       │
        └────┬───────┘
             │
             ▼
        { id: '123', name: 'Thuta' }
 * ```
 */
const user_store = {
  all: {},
  byId: {},
};

const user = { id: '123', name: 'thuta' };

user_store.all['abc'] = user;
user_store.byId['123'] = user;

// ---- Ref
const a = { name: 'thuta' };
const b = a; // b points to the same object
b.name = 'Sann';
console.log(a.name);

// ---- Memory Usage and Growth

const users = [{ id: '1' }, { id: '2' }];

users.forEach((user) => {
  user_store.all[user.id] = user;
  user_store.byId[user.id] = user;
});

user_store.byId['1'].id = 'new id';
user_store.byId['123'].name = { ...user_store.byId['123'], name: 'new new new name' };

console.log(user_store);
