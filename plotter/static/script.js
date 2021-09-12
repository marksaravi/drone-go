const DATA_PER_SECOND = 3200
const CANVAS_TIME_SECONDS = 10
var TIME_SCALE = 1e9
const MAX_BUFFER_SIZE = DATA_PER_SECOND * 60

let grpah = null;
let canvasWidth = 0;
let canvasHeight = 0;
var gyroCanvas = null;
var gyroCtx = null;
var xScale = 1
var yScale = 1
var maxAngle = 90
const xyBuffer = new Array(DATA_PER_SECOND * CANVAS_TIME_SECONDS)


function setupCanvas() {
    grpah = document.querySelector('#gyro');
    canvasWidth = grpah.offsetWidth;
    canvasHeight = grpah.offsetHeight;
    gyroCanvas = document.getElementById("gyro");
    xScale = canvasWidth / CANVAS_TIME_SECONDS / TIME_SCALE
    yScale = canvasHeight / 2 / maxAngle
    gyroCtx = gyroCanvas.getContext("2d");
}

function plot(buffer, lastIndex, numOfPackets) {
    let i = lastIndex;
    let packetCounter = 0
    let dataCounter = 0
    while (packetCounter < numOfPackets && dataCounter < xyBuffer.length) {
        i--;
        if (i < 0) {
            i = MAX_BUFFER_SIZE
        }
        packetCounter++
        xyBuffer[dataCounter] = {
            rRol: buffer[i].r.r,
            t: buffer[i].t
        }
        dataCounter++
    }
    gyroCtx.clearRect(0, 0, 800, 100);
    gyroCtx.stroke();
    let first = true;
    const startTime = xyBuffer[dataCounter - 1].t;
    let prevX = 0
    let strokeCounter = 0
    for (let i = 0; i < dataCounter; i++) {
        const x = canvasWidth - (xyBuffer[i].t - startTime) * xScale
        const y = canvasHeight / 2 - xyBuffer[i].rRol * yScale
        if (first) {
            first = false;
            strokeCounter++;
            gyroCtx.moveTo(x, y)
            prevX = x
        } else {
            if (Math.floor(prevX) !== Math.floor(x)) {
                gyroCtx.lineTo(x, y)
                prevX = x;
                strokeCounter++;
            }
        }
    }
    console.log(strokeCounter)
    gyroCtx.stroke();
}
function createWebSocket() {
    // Create WebSocket connection.
    console.log("establishing connection");
    const socket = new WebSocket('ws://localhost:8081/conn');

    window.plotterBuffer = new Array(MAX_BUFFER_SIZE)
    window.plotterBufferLastIndex = 0
    window.plotterBufferCounter = 0
    // Listen for messages
    socket.addEventListener('message', function (event) {
        const packets = JSON.parse(JSON.parse(event.data));
        packets.forEach(p => {
            window.plotterBuffer[window.plotterBufferLastIndex] = p
            window.plotterBufferLastIndex++
            if (window.plotterBufferLastIndex == MAX_BUFFER_SIZE) {
                window.plotterBufferLastIndex = 0
            }

            if (window.plotterBufferCounter < MAX_BUFFER_SIZE) {
                window.plotterBufferCounter++
            }
        })
        // console.log(window.plotterBufferLastIndex)
        plot(window.plotterBuffer, window.plotterBufferLastIndex, window.plotterBufferCounter)
    });
}

setupCanvas()
createWebSocket()