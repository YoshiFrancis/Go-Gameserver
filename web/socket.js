let socket;

const connect = async (username = "guest", url="ws://localhost:8080/ws") => {
  console.log(username, url);

  socket = new WebSocket(url);

  socket.addEventListener('open', function (event) {
      console.log('Connected to the WebSocket server');
      socket.send(username);
  });

  socket.addEventListener('message', function (event) {
      console.log('Message from server', event.data);
  });

  socket.addEventListener('close', function (event) {
      console.log('Disconnected from the WebSocket server');
  });

  socket.addEventListener('error', function (event) {
      console.error('WebSocket error', event);
  });
}

document.getElementById("connectForm").addEventListener('submit', function(event) {
  event.preventDefault();
  const username = document.getElementById('username').value;
  const url = document.getElementById('url').value;
  connect(username, url);
})

document.getElementById("sendMessage").addEventListener('submit', function(event) {
  event.preventDefault();
  const message = document.getElementById('message').value;
  if (socket) {
    console.log("emmitting ", message);
    socket.send(message);
  }
})