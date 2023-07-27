import WebSocket from "ws"; 

const webSocket = new WebSocket("ws://localhost:3080/ws/driver/Secrete_trip_id");

webSocket.onopen = function(e){
  console.log("Web socket open");
  webSocket.send("Hello from driver.ts");
  //webSocket.close();
}

webSocket.onmessage = function(e){
  console.log("Driver get message: ", e.data.toString())
  return true;
}

webSocket.onerror = function(e){
  console.log("error: ", e.message);
}

webSocket.onclose = function(e){
  console.log("Close");
}