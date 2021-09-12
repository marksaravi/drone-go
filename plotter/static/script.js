const CANVAS_TIME_SECONDS = 3
var TIME_SCALE = 1e9

let startTime = 0
let grpah = null;
let canvasWidth = 0;
let canvasHeight = 0;
var gyroCanvas = null;
var gyroCtx = null;
var xScale = 1
var yScale = 1
var maxAngle = 90


function setupCanvas() {
    grpah = document.querySelector('#gyro');
    canvasWidth = grpah.offsetWidth;
    canvasHeight = grpah.offsetHeight;
    gyroCanvas = document.getElementById("gyro");
    xScale = canvasWidth / CANVAS_TIME_SECONDS / TIME_SCALE
    yScale = canvasHeight / 2 / maxAngle
}

function plot(buffer) {
    const gyroCtx = gyroCanvas.getContext("2d");
    gyroCtx.lineWidth = 1;
    console.log('plot', buffer.length);
    const lastTime = buffer[buffer.length - 1].t
    const totalGraphTime = (lastTime - startTime)
    const xEnd = totalGraphTime * xScale
    gyroCtx.fillStyle = "red";
    gyroCtx.fillRect(0, 0, 800, 100);
    gyroCtx.stroke();
    let xOffset = 0
    if (xEnd > canvasWidth) {
        xOffset = xEnd - canvasWidth
    }
    let first = true
    let prevX = 0
    gyroCtx.clearRect(0, 0, 800, 100)
    gyroCtx.stroke();
    for (let i = buffer.length - 1; i >= 0; i--) {
        const d = buffer[i]
        const x = (d.t - startTime) * xScale - xOffset
        const y = canvasHeight / 2 - d.r.r * yScale
        if (x<0) {
            console.log(i)
            break
        }
        if (first) {
            gyroCtx.moveTo(x, y)
            prevX = x;
            first = false
        } else {
            if (Math.floor(prevX)!==Math.floor(x)) {
                prevX = x
                gyroCtx.lineTo(x, y)
            }
        }
    }
    gyroCtx.stroke();
    
}
function createWebSocket() {
    // Create WebSocket connection.
    console.log("establishing connection");
    const socket = new WebSocket('ws://localhost:8081/conn');

    window.plotterBuffer = []
    // Listen for messages
    socket.addEventListener('message', function (event) {
        const packets = JSON.parse(JSON.parse(event.data));
        if (startTime === 0) {
            startTime = packets[0].t
        }
        window.plotterBuffer = [...window.plotterBuffer, ...packets]
        // packet.map(element => {
        //     buffer.push(element)
        // });
        const maxlen = 5000
        const l1 = window.plotterBuffer.length
        if (window.plotterBuffer.length > maxlen) {
            window.plotterBuffer = window.plotterBuffer.slice(window.plotterBuffer.length-maxlen)
            console.log("slicing", window.plotterBuffer.length, l1)
        }
        plot(window.plotterBuffer)
    });
}

setupCanvas()
createWebSocket()