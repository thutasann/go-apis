let memory = []; // [{id, name}]
let store = {}; // key: reference

function randomName() {
  const names = ['thuta', 'sann', 'alice', 'bob', 'eve'];
  return names[Math.floor(Math.random() * names.length)];
}

function addUser() {
  const user = {
    id: Math.random().toString(36).substring(2, 6),
    name: randomName(),
  };
  memory.push(user);
  updateMemoryView();
  updateMemorySelect();
}

function assignReference() {
  const key = document.getElementById('storeKey').value;
  const selectedIndex = document.getElementById('memorySelect').value;
  if (!key || selectedIndex === '') return;
  store[key] = memory[selectedIndex];
  updateStoreView();
}

function assignCopy() {
  const key = document.getElementById('storeKey').value;
  const selectedIndex = document.getElementById('memorySelect').value;
  if (!key || selectedIndex === '') return;
  store[key] = { ...memory[selectedIndex] };
  updateStoreView();
}

function clearAll() {
  memory = [];
  store = {};
  updateMemoryView();
  updateStoreView();
  updateMemorySelect();
}

function updateMemoryView() {
  const container = document.getElementById('memory');
  container.innerHTML = '';
  memory.forEach((obj, i) => {
    const box = document.createElement('div');
    box.className = 'box';
    box.innerHTML = `<b>Memory[${i}]</b><br>ID: ${obj.id}<br>Name: ${obj.name}`;
    container.appendChild(box);
  });
}

function updateMemorySelect() {
  const select = document.getElementById('memorySelect');
  select.innerHTML = '';
  memory.forEach((_, i) => {
    const option = document.createElement('option');
    option.value = i;
    option.text = `Memory[${i}]`;
    select.appendChild(option);
  });
}

function updateStoreView() {
  const container = document.getElementById('store');
  container.innerHTML = '';
  for (const key in store) {
    const obj = store[key];
    const indexInMemory = memory.findIndex((mem) => mem === obj);
    const box = document.createElement('div');
    box.className = 'box';
    box.innerHTML = `<b>${key}</b><br>ID: ${obj.id}<br>Name: ${obj.name}<br>`;
    if (indexInMemory !== -1) {
      box.innerHTML += `<span class="pointer">→ Memory[${indexInMemory}]</span>`;
    } else {
      box.innerHTML += `<span class="copied">🧾 Copied Object</span>`;
    }
    container.appendChild(box);
  }
}

console.log('store ---> ', store);
