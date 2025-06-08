// @ts-check

/** selecte chat @type { string }  */
let selectedChat = 'general';

/** ws connection @type { WebSocket | null} */
let conn = null;

/**
 * Socket Events
 */
class SocketEvent {
  type = '';
  payload = null;

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
 * Send Message Event
 */
class SendMesasgeEvent {
  message = '';
  from = '';

  /**
   * Send Message Event
   * @param {string} message - Send Message
   * @param {string} from - From Information
   */
  constructor(message, from) {
    this.message = message;
    this.from = from;
  }
}

/**
 * New Message Event
 */
class NewMessageEvent {
  message = '';
  from = '';
  sent = new Date();

  /**
   * Send Message Event
   * @param {string} message - Send Message
   * @param {string} from - From Information
   * @param {Date} sent - sent date
   */
  constructor(message = '', from = '', sent = new Date()) {
    this.message = message;
    this.from = from;
    this.sent = sent;
  }
}

/**
 * Change ChatRoom Event
 */
class ChangeChatRoomEvent {
  name = '';
  /**
   * Change Chat Room Event
   * @param {string} name - chat room name
   */
  constructor(name) {
    this.name = name;
  }
}

/**
 * Route Event
 * @param {{type: SocketEventType | undefined, payload: any}} event - socket event
 */
function routeEvent(event) {
  if (event.type === undefined) {
    alert('no type field in the socket event');
  }

  switch (event.type) {
    case 'new_message':
      console.log('🚀 new message...');
      const messageEvent = Object.assign(new NewMessageEvent(), event.payload);
      appendChatMesage(messageEvent);
      break;
    default:
      alert('unsupported message type');
      break;
  }
}

/**
 * Append Chat Message
 * @param {NewMessageEvent} messageEvent
 */
function appendChatMesage(messageEvent) {
  let date = new Date(messageEvent.sent);
  const formattedMsg = `${date.toLocaleDateString()}: ${messageEvent.message}`;

  const textarea = document.getElementById('chatmessages');
  if (textarea) {
    textarea.innerHTML = textarea.innerHTML + '\n' + formattedMsg;
    textarea.scrollTop = textarea.scrollHeight;
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
    selectedChat = newChat.value;
    const header = document.getElementById('chat-header');
    if (header) {
      header.innerHTML = 'Currently in chatroom: ' + selectedChat;

      let changeEvent = new ChangeChatRoomEvent(selectedChat);
      sendEvent('change_room', changeEvent);

      const textarea = document.getElementById('chatmessages');
      if (textarea) {
        textarea.innerHTML = `You changed room into: ${selectedChat}`;
      }
    }
  }
  return false;
}

/**
 * Send Mesasge
 */
function sendMessage() {
  const newMessage = /** @type {HTMLInputElement | null} */ (document.getElementById('message'));
  if (newMessage != null && conn) {
    let outgoingEvent = new SendMesasgeEvent(newMessage.value, 'test');
    sendEvent('send_message', outgoingEvent);
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
    conn = new WebSocket('wss://' + document.location.host + '/ws?otp=' + otp);

    conn.onopen = function (evt) {
      const connection_header = document.getElementById('connection-header');
      if (connection_header) {
        connection_header.innerHTML = 'Connected to Websocket: true';
      }
    };

    conn.onclose = function (evt) {
      const connection_header = document.getElementById('connection-header');
      if (connection_header) {
        connection_header.innerHTML = 'Connected to Websocket: false';
      }
    };

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
