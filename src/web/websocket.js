var socket = new WebSocket("ws://localhost:8080/ws");
socket.onmessage = function () {
  location.reload();
};