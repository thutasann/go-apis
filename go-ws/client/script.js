var selectedChat = 'general';

/**
 * Change Chat Room
 */
function changeChatRoom() {
  var newChat = document.getElementById('chatroom');
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

window.onload = function () {
  console.log('::: INITIALIZING... :::');

  var chatroomSelection = document.getElementById('chatroom-selection');
  if (chatroomSelection) {
    chatroomSelection.onsubmit = changeChatRoom;
  } else {
    console.warn('no chatroom selection');
  }
};
