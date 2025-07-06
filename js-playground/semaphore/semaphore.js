class Semaphore {
  /**
   * Creates a semaphore that limits the number of concurrent Promises being handled
   * @param {number} maxConcurrentRequests max number of concurrent promises being handled at any time
   */
  constructor(maxConcurrentRequests = 1) {
    this.currentRequests = [];
    this.runningRequests = 0;
    this.maxConcurrentRequests = maxConcurrentRequests;
  }

  /**
   * Returns a Promise that will eventually return the result of the function passed in
   * Use this to limit the number of concurrent function executions
   * @param {*} fnToCall function that has a cap on the number of concurrent executions
   * @param  {...any} args any arguments to be passed to fnToCall
   * @returns Promise that will resolve with the resolved value as if the function passed in was directly called
   */
  callFunction(fnToCall, ...args) {
    return new Promise((resolve, reject) => {
      this.currentRequests.push({
        resolve,
        reject,
        fnToCall,
        args,
      });
      this.tryNext();
    });
  }

  tryNext() {
    if (!this.currentRequests.length) {
      return;
    } else if (this.runningRequests < this.maxConcurrentRequests) {
      let { resolve, reject, fnToCall, args } = this.currentRequests.shift();
      this.runningRequests++;
      let req = fnToCall(...args);
      req
        .then((res) => resolve(res))
        .catch((err) => reject(err))
        .finally(() => {
          this.runningRequests--;
          this.tryNext();
        });
    }
  }
}

const throttler = new Semaphore(2);
throttler
  .callFunction(fetch, 'https://www.facebook.com')
  .then((res) => res.text())
  .then((body) => {
    console.log('Fetched Facebook:', body.slice(0, 100)); // log first 100 chars
  })
  .catch((err) => console.error('Error:', err));

throttler
  .callFunction(fetch, 'https://github.com')
  .then((res) => res.text())
  .then((body) => {
    console.log('Fetched Github:', body.slice(0, 100)); // log first 100 chars
  })
  .catch((err) => console.error('Error:', err));

throttler
  .callFunction(fetch, 'https://instagram.com')
  .then((res) => res.text())
  .then((body) => {
    console.log('Fetched Instagram:', body.slice(0, 100)); // log first 100 chars
  })
  .catch((err) => console.error('Error:', err));
