// @ts-check

const net = require('net');
const { encodeFrame, decodeFrames } = require('./protocol');

const socket = net.createConnection({ port: 9000 });

let buffer = Buffer.alloc(0);

socket.on('connect', () => {
  console.log('Connected to server');
  socket.write(encodeFrame(Buffer.from('hello tcp')));
});

socket.on('data', (/** @type any */ chunk) => {
  buffer = Buffer.concat([buffer, chunk]);
  const { frames, remaining } = decodeFrames(buffer);

  // @ts-ignore
  buffer = remaining;

  for (const frame of frames) {
    console.log('Server response:', frame.payload.toString());
  }
});
