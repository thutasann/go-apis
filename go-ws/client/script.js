// @ts-check

/** @type { string } */
let selectedChat = 'general';

/** @type { WebSocket | null} */
let conn = null;

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
    conn.send(newMessage.value);
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
      console.log('ℹ️ Event --> ', evt.data);
    };
  } else {
    alert('Browser does not support Websocket');
  }
};
