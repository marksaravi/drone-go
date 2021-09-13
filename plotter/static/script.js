const MAX_DATA_PER_SECOND = 4000
const MAX_TIME_SPAN = 60
var TIME_SCALE = 1e9
const MAX_BUFFER_SIZE = MAX_DATA_PER_SECOND * 120

const X_PADDING = 24
const Y_PADDING = 24

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
const xyBuffer = new Array(MAX_DATA_PER_SECOND * MAX_TIME_SPAN)
const grapgs = ['accelerometer', 'gyroscope', 'rotations']

function setupContainer() {
    const container = document.getElementById('canvas-container')
    const acc = document.getElementById('accelerometer')
    graphSettings.width = container.offsetWidth
    graphSettings.height = acc.offsetHeight
    graphSettings.xScale = graphSettings.width / graphSettings.timeSpan / TIME_SCALE
    graphSettings.yScale = graphSettings.height / 2 / graphSettings.yMax
}

function getCanvasContext(id) {
    const canvas = document.getElementById(id);
    canvas.width = graphSettings.width;
    canvas.height = graphSettings.height;
    const ctx = canvas.getContext("2d");
    return ctx
}

function setupPlotter() {
    setupContainer()
    ctx2D.accelerometer = getCanvasContext("accelerometer");
    ctx2D.gyroscope = getCanvasContext("gyroscope");
    ctx2D.rotations = getCanvasContext("rotations");
    getContextes((ctx) => {
        ctx.lineWidth = 0.5
        ctx.strokeStyle = '#006400';
    });
}

function getContextes(action) {
    grapgs.forEach(g => {
        action(ctx2D[g])
    })
}


function clearCanvases() {
    getContextes((ctx) => ctx.clearRect(0, 0, graphSettings.width, graphSettings.height));
}

function beginPath() {
    getContextes((ctx) => ctx.beginPath());
}

function lineTo(x, ay, gy, ry) {
    ctx2D['accelerometer'].lineTo(x, ay)
    ctx2D['gyroscope'].lineTo(x, gy)
    ctx2D['rotations'].lineTo(x, ry)
}

function stroke() {
    getContextes((ctx) => ctx.stroke());
}

function Y(y) {
    return graphSettings.height / 2 - y * graphSettings.yScale
}

function X(t) {
    return graphSettings.width - t * graphSettings.xScale
}

function plot(datalink) {
    let dataCounter = 0
    const startTime = datalink.data.t
    while (dataCounter < xyBuffer.length) {
        if (!datalink.data) {
            break
        }
        const x = X(startTime - datalink.data.t)
        const ay = Y(datalink.data.a[graphSettings.axis])
        const gy = Y(datalink.data.g[graphSettings.axis])
        const ry = Y(datalink.data.r[graphSettings.axis])
        xyBuffer[dataCounter] = {
            x,
            ay,
            gy,
            ry
        }

        dataCounter++
        if (x < 0) {
            break
        }
        datalink = datalink.prev
    }
    clearCanvases()
    beginPath();
    let prevX = 0
    for (let i = dataCounter - 1; i >= 0; i--) {
        if (Math.floor(prevX) != Math.floor(xyBuffer[i].x)) {
            prevX = xyBuffer[i].x
            lineTo(xyBuffer[i].x, xyBuffer[i].ay, xyBuffer[i].gy, xyBuffer[i].ry);
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

function setAxis(axis) {
    graphSettings.axis = axis
    const label = document.getElementById('axis')
    label.innerHTML = axis
}

setupPlotter()
createWebSocket()