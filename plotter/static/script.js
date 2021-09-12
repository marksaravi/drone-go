function plot() {
    var c = document.getElementById("gyro");
    var ctx = c.getContext("2d");
    ctx.moveTo(0, 0);
    ctx.lineTo(100, 100);
    ctx.lineTo(200, 0);
    ctx.stroke();
}
function createWebSocket() {
    // Create WebSocket connection.
    console.log("establishing connection");
    const socket = new WebSocket('ws://localhost:8081/conn');

    // Listen for messages
    socket.addEventListener('message', function (event) {
        const json = JSON.parse(event.data);
    });
}
plot()
createWebSocket()