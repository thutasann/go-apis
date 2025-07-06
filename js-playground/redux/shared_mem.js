function shared_between_threads() {
  const shared = new SharedArrayBuffer(4);
  const view = new Int32Array(shared);

  view[0] = 100;
  console.log(view[0]);
}
shared_between_threads();
