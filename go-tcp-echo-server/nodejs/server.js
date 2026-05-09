// @ts-check

const net = require('net');
const { handleConnection } = require('./connection');

/**
 * @typedef {Object} ServerConfig
 * @property {number} port
 * @property {string} host
 */

/** @type { ServerConfig} */
const config = {
  port: 9000,
  host: '0.0.0.0',
};

const server = net.createServer(handleConnection);

server.listen(config.port, config.host, () => {
  console.log(`🚀 TCP server listening on ${config.host}:${config.port}`);
});

process.on('SIGINT', shutdown);
process.on('SIGTERM', shutdown);

function shutdown() {
  console.log('🛑 Shutting down TCP server...');
  server.close(() => {
    console.log('✅ Server closed');
    process.exit(0);
  });
}
