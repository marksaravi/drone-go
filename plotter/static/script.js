const DATA_PER_SECOND = 3200
const CANVAS_TIME_SECONDS = 1
var TIME_SCALE = 1e9
const MAX_BUFFER_SIZE = DATA_PER_SECOND * 2

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

function lineTo(x, ay, gy, ry) {
    accCtx.lineTo(x, ay)
    gyroCtx.lineTo(x, gy)
    rotCtx.lineTo(x, ry)
}

function stroke() {
    accCtx.stroke();
    gyroCtx.stroke();
    rotCtx.stroke();
}

function plot(datalink) {
    let dataCounter = 0
    while (dataCounter < xyBuffer.length) {
        if (!datalink.data) {
            break
        }
        xyBuffer[dataCounter] = {
            aRol: datalink.data.a.r,
            gRol: datalink.data.g.r,
            rRol: datalink.data.r.r,
            t: datalink.data.t
        }
        dataCounter++
        datalink = datalink.prev
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
            lineTo(x, ay, gy, ry);
        }
    }
    stroke();
}
function createWebSocket() {
    // Create WebSocket connection.
    console.log("establishing connection");
    const socket = new WebSocket('ws://localhost:8081/conn');

    const firstlink = {
        next: null,
        prev: null,
        data: null,
    }
    let link = firstlink
    for (let i = 0; i < MAX_BUFFER_SIZE; i++) {
        newlink = {
            next: null,
            prev: null,
            data: null,
        }
        newlink.prev = link
        link.next = newlink
        link=link.next
    }
    link.next=firstlink
    firstlink.prev=link

    socket.addEventListener('message', function (event) {
        const packets = JSON.parse(JSON.parse(event.data));
        packets.forEach(p => {
            link=link.next
            link.data=p
            link.next.data=null
        })
        plot(link)
    });
}

setupPlotter()
createWebSocket()