// @ts-check

var selectedChat = 'general';

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
  var newMessage = document.getElementById('message');
  if (newMessage != null) {
    console.log('newMessage --> ', newMessage);
  }
  return false;
}

/**
 * Initializing the app
 */
window.onload = function () {
  console.log('::: INITIALIZING :::');
  var chatroomSelection = document.getElementById('chatroom-selection');
  var chatroomMessage = document.getElementById('chatroom-message');

  if (chatroomSelection) {
    chatroomSelection.onsubmit = changeChatRoom;
  }

  if (chatroomMessage) {
    chatroomMessage.onsubmit = sendMessage;
  }

  if (window['WebSocket']) {
    console.log('::: Connecting to Websockets :::');
    var conn = new WebSocket('ws://' + document.location.host + '/ws');
  } else {
    alert('Browser does not support Websocket');
  }
};
