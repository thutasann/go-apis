// @ts-check

/** selecte chat @type { string }  */
let selectedChat = 'general';

/** ws connection @type { WebSocket | null} */
let conn = null;

/**
 * Socket Events
 */
class SocketEvent {
  /**
   * Socket Event
   * @param {string} type - Socket Event Type
   * @param {*} payload - Socket Event payload
   */
  constructor(type = '', payload = null) {
    this.type = type;
    this.payload = payload;
  }
}

/**
 * Route Event
 * @param {{type: SocketEventType | undefined}} event - socket event
 */
function routeEvent(event) {
  if (event.type === undefined) {
    alert('no type field in the socket event');
  }

  switch (event.type) {
    case 'new_message':
      console.log('🚀 new message...');
      break;
    default:
      alert('unsupported message type');
      break;
  }
}

/**
 * send socket event
 * @param {SendEvent} eventName - socket event name
 * @param {*} payload - socket payload
 */
function sendEvent(eventName, payload) {
  if (!conn) {
    alert('socket connection failed...');
    return;
  }
  const event = new SocketEvent(eventName, payload);
  conn.send(JSON.stringify(event));
}

/**
 * Change Chat Room
 */
function changeChatRoom() {
  const newChat = /** @type {HTMLInputElement | null} */ (document.getElementById('chatroom'));
  if (newChat != null && newChat.value != selectedChat) {
    console.log('newChat --> ', newChat);
  }
  return false;
}

/**
 * Send Mesasge
 */
function sendMessage() {
  const newMessage = /** @type {HTMLInputElement | null} */ (document.getElementById('message'));
  if (newMessage != null && conn) {
    sendEvent('send_message', newMessage.value);
  }
  return false;
}

/**
 * Initializing the app
 */
window.onload = function () {
  console.log('::: INITIALIZING :::');
  const chatroomSelection = document.getElementById('chatroom-selection');
  const chatroomMessage = document.getElementById('chatroom-message');

  if (chatroomSelection) {
    chatroomSelection.onsubmit = changeChatRoom;
  }

  if (chatroomMessage) {
    chatroomMessage.onsubmit = sendMessage;
  }

  if (window['WebSocket']) {
    console.log('::: Connecting to Websockets :::');
    conn = new WebSocket('ws://' + document.location.host + '/ws');

    conn.onmessage = function (evt) {
      const eventData = JSON.parse(evt.data);
      const event = Object.assign(new SocketEvent(), eventData);
      routeEvent(event);
    };
  } else {
    alert('Browser does not support Websocket');
  }
};
