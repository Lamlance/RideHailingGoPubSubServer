const client_link_input = document.getElementById("client_link_input") as HTMLInputElement;
const driver_link_input = document.getElementById("driver_link_input") as HTMLInputElement;

const client_msg_input = document.getElementById("client_message_input") as HTMLInputElement;
const driver_msg_input = document.getElementById("driver_message_input") as HTMLInputElement;

const client_button = document.getElementById("client_button") as HTMLButtonElement;
const driver_button = document.getElementById("driver_button") as HTMLButtonElement;

const log_textarea = document.getElementById("log_textarea") as HTMLTextAreaElement;


let client_socket:WebSocket;
let driver_socket:WebSocket;

function Client_Connect(){
  client_socket = new WebSocket(client_link_input.value);
  client_socket.onclose = function(){
    client_button.onclick = Client_Connect;
    client_button.innerHTML = "Connect"
    log_textarea.value += `Client socket closed \n`;
  }

  client_socket.onopen = function(){
    client_button.innerHTML = "Send text";
    client_button.onclick = ()=>{
      client_socket.send(client_msg_input.value);
      log_textarea.value += `Client send: ${client_msg_input.value} \n`;
    }

    client_socket.onmessage = function(e){
      log_textarea.value += `Client get message: ${e.data.toString()}`
    }    

    log_textarea.value += `Client socket opened \n`;
  }

  client_socket.onerror = function(e){
    log_textarea.value += `Client socket error: ${e}`
  }

}

function DriverConnect(){
  driver_socket = new WebSocket(driver_link_input.value);
  driver_socket.onclose = function(){
    driver_button.onclick = DriverConnect;
    driver_button.innerHTML = "Connect";
    log_textarea.value += "Driver socket closed \n";
  }

  driver_socket.onopen = function(){
    driver_button.innerHTML = "Send text";
    driver_button.onclick = ()=>{
      driver_socket.send(driver_msg_input.value);
      log_textarea.value += `Driver send: ${driver_msg_input.value} \n`;
    }
    
    driver_socket.onmessage = function(e){
      log_textarea.value += `Driver get message: ${e.data.toString()} \n`;
    }

    log_textarea.value += "Driver socket opened \n"
  }

  driver_socket.onerror = function(e){
    log_textarea.value += `Driver socket error: ${e}`
  }
}


client_button.onclick = Client_Connect;
driver_button.onclick = DriverConnect;
console.log(driver_button.onclick);
console.log(client_button.onclick);
