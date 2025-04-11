var selectedChat = "general";
var connection;
class Event {
  constructor(type, payload) {
    this.type = type;
    this.payload = payload;
  }
}

class SendMessageEvent {
  constructor(message, from) {
    this.message = message;
    this.from = from;
  }
}

class NewMessageEvent {
  constructor(message, from, sent) {
    this.message = message;
    this.from = from;
    this.sent = sent;
  }
}

function routeEvent(event) {
  if (event.type === undefined) {
    window.alert("No type field in the event.");
  }
  const { type, _ } = event;
  switch (type) {
    case "new_message":
      const messageEvent = Object.assign(new NewMessageEvent(), event.payload);
      appendChatMessage(messageEvent);
      break;
    default:
      alert("Unsupported Message Type.");
      break;
  }
}

function appendChatMessage(messageEvent) {
  var date = new Date(messageEvent.sent);
  const formattedMsg = `${date.toLocaleString()}: ${messageEvent.message}`;
  var textarea = document.getElementById("chat-messages");
  textarea.innerHTML = textarea.innerHTML + "\n" + formattedMsg;
  textarea.scrollTop = textarea.scrollHeight;
}

function sendEvent(eventName, payload) {
  const event = new Event(eventName, payload);
  connection.send(JSON.stringify(event));
}

function changeChatRoom() {
  var newChat = document.getElementById("chatroom");
  if (newChat != null && newChat.value != selectedChat) {
    console.log(newChat);
  }
  return false;
}

function sendMessage() {
  const newMessage = document.getElementById("message");
  if (newMessage != null) {
    let outgoingEvent = new SendMessageEvent(newMessage.value, "adnan");
    sendEvent("send_message", outgoingEvent);
  }
  return false;
}

function connectWebsocket(otp) {
  if (window["WebSocket"]) {
    console.log("Websockets supported.");
    // new Websocket(url, protocols)
    connection = new WebSocket(
      "ws://" + document.location.host + "/ws?otp=" + otp
    );

    connection.onopen = function (event) {
      document.getElementById(
        "connection-header"
      ).innerHTML = `Connected to websocket : True`;
    };

    connection.onclose = function (event) {
      document.getElementById(
        "connection-header"
      ).innerHTML = `Connected to websocket : False`;
    };

    // Fires on receiving the message.
    connection.onmessage = function (event) {
      const eventData = JSON.parse(event.data);
      const evt = Object.assign(new Event(), eventData);
      routeEvent(evt);
    };
  } else {
    window.alert("Websockets not supported.");
  }
}

function login() {
  let formData = {
    username: document.getElementById("username").value,
    password: document.getElementById("password").value,
  };
  fetch("login", {
    method: "psot",
    body: JSON.stringify(formData),
    mode: "cors",
  })
    .then((res) => {
      if (res.ok) {
        return res.json();
      } else {
        throw new Error("Unauthorized!");
      }
    })
    .then((data) => {
      connectWebsocket(data.otp);
    })
    .catch((err) => {
      alert(err);
    });
  return false;
}

window.onload = function () {
  document.getElementById("chatroom-selection").onsubmit = changeChatRoom;
  document.getElementById("chatroom-message").onsubmit = sendMessage;
  document.getElementById("login-form").onsubmit = login;
};
