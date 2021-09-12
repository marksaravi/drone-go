const DATA_PER_SECOND = 3200
const CANVAS_TIME_SECONDS = 10
var TIME_SCALE = 1e9
const MAX_BUFFER_SIZE = DATA_PER_SECOND * 60

let grpah = null;
let canvasWidth = 0;
let canvasHeight = 0;
var gyroCanvas = null;
var gyroCtx = null;
var accCanvas = null;
var accCtx = null;
var rotCanvas = null;
var rotCtx = null;
var xScale = 1
var yScale = 1
var maxAngle = 90
const xyBuffer = new Array(DATA_PER_SECOND * CANVAS_TIME_SECONDS)


function setupContainer() {
    const container = document.getElementById('canvas-container');
    const acc = document.getElementById('acc');
    canvasWidth = container.offsetWidth;
    canvasHeight = acc.offsetHeight;
    xScale = canvasWidth / CANVAS_TIME_SECONDS / TIME_SCALE;
    yScale = canvasHeight / 2 / maxAngle;
}

function getCanvasContext(id) {
    const canvas = document.getElementById(id);
    canvas.width = canvasWidth;
    canvas.height = canvasHeight;
    const ctx = canvas.getContext("2d");
    return ctx
}

function setupPlotter() {
    setupContainer()
    accCtx = getCanvasContext("acc");
    gyroCtx = getCanvasContext("gyro");
    rotCtx = getCanvasContext("rotations");
}

function clearCanvases() {
    accCtx.clearRect(0, 0, canvasWidth, canvasHeight);
    gyroCtx.clearRect(0, 0, canvasWidth, canvasHeight);
    rotCtx.clearRect(0, 0, canvasWidth, canvasHeight);
}

function beginPath() {
    accCtx.beginPath();
    gyroCtx.beginPath();
    rotCtx.beginPath();
}

function lineTo(x, ay,gy,ry) {
    accCtx.lineTo(x, ay)
    gyroCtx.lineTo(x, gy)
    rotCtx.lineTo(x, ry)
}

function stroke() {
    accCtx.stroke();
    gyroCtx.stroke();
    rotCtx.stroke();
}

function plot(buffer, lastIndex, numOfPackets) {
    let i = lastIndex;
    let packetCounter = 0
    let dataCounter = 0
    while (packetCounter < numOfPackets && dataCounter < xyBuffer.length) {
        i--;
        if (i < 0) {
            i = MAX_BUFFER_SIZE - 1
        }
        packetCounter++
        xyBuffer[dataCounter] = {
            aRol: buffer[i].a.r,
            gRol: buffer[i].g.r,
            rRol: buffer[i].r.r,
            t: buffer[i].t
        }
        dataCounter++
    }
    const startTime = xyBuffer[dataCounter - 1].t;
    clearCanvases()
    beginPath();
    let prevX = 0
    for (let i = dataCounter - 1; i >= 0; i--) {
        const x = canvasWidth - (xyBuffer[i].t - startTime) * xScale
        const ay = canvasHeight / 2 - xyBuffer[i].aRol * yScale
        const gy = canvasHeight / 2 - xyBuffer[i].gRol * yScale
        const ry = canvasHeight / 2 - xyBuffer[i].rRol * yScale
        if (Math.floor(prevX) != Math.floor(x)) {
            prevX = x
            lineTo(x, ay,gy,ry);
        }
    }
    stroke();
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

setupPlotter()
createWebSocket()