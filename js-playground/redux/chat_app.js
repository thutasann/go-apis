/**
 * All Data are normalized
 * - Messages are stored once.
 * - Channels reference message IDs.
 * - Messages reference user IDs.
 */
const initialState = {
  users: {
    u1: { id: 'u1', name: 'Thuta' },
    u2: { id: 'u2', name: 'Sann' },
  },
  messages: {
    m1: { id: 'm1', text: 'Hello!', sender: 'u1' },
    m2: { id: 'm2', text: 'Hey there', sender: 'u2' },
  },
  channels: {
    general: { id: 'general', name: 'General', messages: ['m1', 'm2'] },
  },
  activeChannelId: 'general',
};

function selectMessagesForActiveChannel(state) {
  const channel = state.channels[state.activeChannelId];
  return channel.messages.map((id) => state.messages[id]);
}
