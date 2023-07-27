const client_link_input = document.getElementById("client_link_input") as HTMLInputElement;
const driver_link_input = document.getElementById("driver_link_input") as HTMLInputElement;

const client_msg_input = document.getElementById("client_message_input") as HTMLInputElement;
const driver_msg_input = document.getElementById("driver_message_input") as HTMLInputElement;

const client_button = document.getElementById("client_button") as HTMLButtonElement;
const driver_button = document.getElementById("diver_button") as HTMLButtonElement;

function Client_Connect(){
  const ws = new WebSocket(client_link_input.value);
  ws.onclose = function(e){
    client_button.onclick = Client_Connect;

    client_button.innerHTML = "Connect to socket"

    console.log("Client socket closed: ",e.code);
  }

  ws.onopen = function(e){
    client_button.innerHTML = "Send text";
    console.log("Socket open");
  }
}


client_button.onclick = Client_Connect;