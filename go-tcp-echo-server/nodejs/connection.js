// @ts-check

const { decodeFrames, encodeFrame } = require('./protocol');
const net = require('net');

/**
 * @typedef {Object} ConnectionContext
 * @property {net.Socket} socket
 * @property {Buffer} buffer
 * @property {string} id
 */

/**
 * Handle a TCP connection
 * @param {net.Socket} socket
 */
function handleConnection(socket) {
  /** @type { ConnectionContext} */
  const ctx = {
    socket,
    buffer: Buffer.alloc(0),
    id: `${socket.remoteAddress}:${socket.remotePort}`,
  };

  console.log(`[CONNECTED] ${ctx.id}`);

  socket.on('data', (/** @type any */ chunk) => onData(ctx, chunk));
  socket.on('close', () => onClose(ctx));
  socket.on('error', (err) => onError(ctx, err));
}

/**
 * @param {ConnectionContext} ctx
 * @param {Buffer} chunk
 */
function onData(ctx, chunk) {
  ctx.buffer = Buffer.concat([ctx.buffer, chunk]);

  const { frames, remaining } = decodeFrames(ctx.buffer);
  ctx.buffer = remaining;

  for (const frame of frames) {
    handleMessage(ctx, frame.payload);
  }
}

/**
 * @param {ConnectionContext} ctx
 * @param {Buffer} payload
 */
function handleMessage(ctx, payload) {
  const message = payload.toString('utf8');
  console.log(`[RECV] ${ctx.id}:`, message);

  // Example: echo + transform
  const response = Buffer.from(
    JSON.stringify({
      ok: true,
      echo: message,
      timestamp: Date.now(),
    }),
    'utf8',
  );

  ctx.socket.write(encodeFrame(response));
}

/**
 * @param {ConnectionContext} ctx
 */
function onClose(ctx) {
  console.log(`[DISCONNECTED] ${ctx.id}`);
}

/**
 * @param {ConnectionContext} ctx
 * @param {Error} err
 */
function onError(ctx, err) {
  console.error(`[ERROR] ${ctx.id}`, err);
  ctx.socket.destroy();
}

module.exports = {
  handleConnection,
};
