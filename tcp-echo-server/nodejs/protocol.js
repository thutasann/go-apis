// @ts-check

/**
 * @typedef {Object} Frame
 * @property {number} length
 * @property {Buffer} payload
 */

/**
 * Encode a payload into length-prefixed frame
 * @param {Buffer} payload
 * @returns {Buffer}
 */
function encodeFrame(payload) {
  const header = Buffer.allocUnsafe(4);
  header.writeUint32BE(payload.length, 0);
  return Buffer.concat([header, payload]);
}

/**
 * Attempt to decode frames from a buffer
 * @param {Buffer} buffer
 * @returns {{ frames: Frame[], remaining: Buffer }}
 */
function decodeFrames(buffer) {
  /** @type {Frame[]} */
  const frames = [];

  let offset = 0;

  while (buffer.length - offset >= 4) {
    const length = buffer.readUint32BE(offset);
    if (buffer.length - offset - 4 < length) break;

    const payload = buffer.subarray(offset + 4, offset + 4 + length);
    frames.push({ length, payload });
    offset += 4 + length;
  }

  return {
    frames,
    remaining: buffer.subarray(offset),
  };
}

module.exports = {
  encodeFrame,
  decodeFrames,
};
