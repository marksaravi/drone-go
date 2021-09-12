const DATA_PER_SECOND = 3200
const CANVAS_TIME_SECONDS = 10
var TIME_SCALE = 1e9
const MAX_BUFFER_SIZE = DATA_PER_SECOND * 120

const graphSettings = {
    width: 0,
    height: 0,
    xScale: 1,
    yScale: 1,
    timeSpan: 10,
    yMax: 90,
    axis: 'roll',
}

const ctx2D = {
    accelerometer: null,
    gyroscope: null,
    rotations: null,
}

let canvasWidth = 0;
let canvasHeight = 0;
var xScale = 1
var yScale = 1
var maxAngle = 90
const xyBuffer = new Array(DATA_PER_SECOND * CANVAS_TIME_SECONDS)

const grapgs = ['accelerometer', 'gyroscope', 'rotations']

function setupContainer() {
    const container = document.getElementById('canvas-container')
    const acc = document.getElementById('accelerometer')
    graphSettings.width = container.offsetWidth
    graphSettings.height = container.offsetHeight
    canvasWidth = container.offsetWidth
    canvasHeight = acc.offsetHeight
    xScale = canvasWidth / CANVAS_TIME_SECONDS / TIME_SCALE
    yScale = canvasHeight / 2 / maxAngle
    graphSettings.xScale = canvasWidth / CANVAS_TIME_SECONDS / TIME_SCALE
    graphSettings.yScale = canvasHeight / 2 / maxAngle
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
    ctx2D.accelerometer = getCanvasContext("accelerometer");
    ctx2D.gyroscope = getCanvasContext("gyroscope");
    ctx2D.rotations = getCanvasContext("rotations");

}

function getContextes(action) {
    grapgs.forEach(g => {
        action(ctx2D[g])
    })
}

function clearCanvases() {
    getContextes((ctx)=> ctx.clearRect(0, 0, canvasWidth, canvasHeight));
}

function beginPath() {
    getContextes((ctx)=> ctx.beginPath());
}

function lineTo(x, ay, gy, ry) {
    ctx2D['accelerometer'].lineTo(x, ay)
    ctx2D['gyroscope'].lineTo(x, gy)
    ctx2D['rotations'].lineTo(x, ry)
}

function stroke() {
    getContextes((ctx)=> ctx.stroke());
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
        link = link.next
    }
    link.next = firstlink
    firstlink.prev = link

    socket.addEventListener('message', function (event) {
        const packets = JSON.parse(JSON.parse(event.data));
        packets.forEach(p => {
            link = link.next
            link.data = p
            link.next.data = null
        })
        plot(link)
    });
}

setupPlotter()
createWebSocket()