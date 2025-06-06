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
 * Login the user
 */
function login() {
  const username = /** @type {HTMLInputElement| null} */ (document.getElementById('username'));
  const password = /** @type {HTMLInputElement | null} */ (document.getElementById('password'));

  if (!username?.value || !password?.value) {
    alert('invalid request');
    return;
  }

  /** @type { LoginForm } */
  let formData = {
    username: username.value,
    password: password.value,
  };

  fetch('login', {
    method: 'post',
    body: JSON.stringify(formData),
    mode: 'cors',
  })
    .then((response) => {
      if (response.ok) {
        return response.json();
      } else {
        throw 'unauthorized';
      }
    })
    .then((/** @type { LoginRes } */ data) => {
      // connect websocket
      connectWebsocket(data.otp);
    })
    .catch((err) => {
      console.error('cannot login: ', err);
    });

  return false;
}

/**
 * connect to websocket
 * @param {string} otp - login otp string
 */
function connectWebsocket(otp) {
  if (window['WebSocket']) {
    console.log('::: Connecting to Websockets :::', otp);
    conn = new WebSocket('ws://' + document.location.host + '/ws?otp=', otp);

    conn.onmessage = function (evt) {
      const eventData = JSON.parse(evt.data);
      const event = Object.assign(new SocketEvent(), eventData);
      routeEvent(event);
    };
  } else {
    alert('Browser does not support Websocket');
  }
}

/**
 * Initializing the app
 */
window.onload = function () {
  console.log('::: INITIALIZING :::');
  const chatroomSelection = document.getElementById('chatroom-selection');
  const chatroomMessage = document.getElementById('chatroom-message');
  const loginForm = document.getElementById('login-form');

  if (chatroomSelection) {
    chatroomSelection.onsubmit = changeChatRoom;
  }

  if (chatroomMessage) {
    chatroomMessage.onsubmit = sendMessage;
  }

  if (loginForm) {
    loginForm.onsubmit = login;
  }
};
