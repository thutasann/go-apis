function obj_assign() {
  const original = { name: 'Alice' };

  const copy = original;

  copy.name = 'Bob';

  console.log(original.name);
}

function assign_to_another_object_key() {
  const user = { name: 'Alice' };
  const container = {};

  // Assigning object as a proprty value
  container.data = user;

  container.data.name = 'bob';

  console.log(user.name); // "Bob" — still the same reference!
}

function shallow_copy() {
  const original = { name: 'Alice' };
  const copy = { ...original };
  copy.name = 'Bob';
  console.log('copy', copy.name);
  console.log('original', original.name);
}

function redux_basic() {
  const user = { id: '123', name: 'thuta' };

  const store = {
    all: {},
    byId: {},
  };

  store.all['abc'] = user;
  store.byId['123'] = user;

  console.log(store);

  store.all['abc'].name = 'updated';

  console.log('after store', store);
}

function array_ref() {
  const user = [{ id: '123', name: 'thuta' }];

  const store = {
    all: {},
    byId: {},
  };

  store.all['abc'] = user;
  store.byId['123'] = user;

  store.all['abc'][0].name = 'updated';

  console.log('after store', JSON.stringify(store, null, 2));
}

function mutable_data() {
  const user = { name: 'Thuta' };
  user.name = 'Sann'; // ✅ mutated

  const arr = [1, 2, 3];
  arr.push(4); // ✅ mutated
}

function immutable_data() {
  let name = 'Thuta';
  name = 'Sann'; // ❌ not mutation — creates a new string

  let x = 5;
  x = x + 1; // ❌ not mutation — creates a new number
}

function gc_sample_one() {
  let obj = { value: 10 };
  let ref = obj;

  obj = null;
  console.log('obj --> ', obj);
  console.log('ref --> ', ref.value); // 10, object still alive via `ref`

  ref = null;
  console.log('after ref --> ', ref);
}

function weak_map() {
  let obj = {};
  const wm = new WeakMap();

  wm.set(obj, 'some value');
  console.log('wm --> ', wm.get(obj));

  obj = null; // obj can be garbage collected, WeakMap entry will be removed automatically

  console.log('after wm --> ', wm.get(obj));
}
