import WebSocket from "ws";

const webSocket = new WebSocket("ws://localhost:3080/ws/client/w3gv7?user_id=abc&lon=106.69380051915194&lat=10.78825445546148");
webSocket.onopen = function (e) {
  console.log("Web socket open");
  webSocket.send("Hello from client.ts");
  //webSocket.close();
}

webSocket.onmessage = function (e) {
  console.log("Client get message: ", e.data.toString())
  return true;
}

webSocket.onerror = function (e) {
  console.log("error: ", e.message);
}

webSocket.onclose = function (e) {
  console.log("Close");
}