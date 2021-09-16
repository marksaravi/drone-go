const MAX_DATA_PER_SECOND = 4000
const MAX_TIME_SPAN = 60
const TIME_SCALE = 1e-9
const MAX_BUFFER_SIZE = MAX_DATA_PER_SECOND * 120
const FONT_SIZE = 10

const X_PADDING = 32
const Y_PADDING = 20

const graphSettings = {
    width: 0,
    height: 0,
    graphWidth: 0,
    graphHeight: 0,
    xScale: 1,
    yScale: 1,
    timeSpan: 10,
    yMax: 105,
    axis: 'roll',
    yGrid: 15,
    firstTime: 0,
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
    graphSettings.graphWidth = graphSettings.width - X_PADDING
    graphSettings.graphHeight = graphSettings.height - Y_PADDING
    graphSettings.xScale = graphSettings.graphWidth / graphSettings.timeSpan
    graphSettings.yScale = graphSettings.graphHeight / 2 / graphSettings.yMax
}

function getCanvasContext(id) {
    const canvas = document.getElementById(id);
    canvas.width = graphSettings.width;
    canvas.height = graphSettings.height;
    const ctx = canvas.getContext("2d");
    ctx.font = `${FONT_SIZE}px Comic Sans MS`;
    ctx.fillStyle = "red";
    ctx.textAlign = "center";
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
    getContextes((ctx) => ctx.clearRect(0, 0, graphSettings.width, graphSettings.height));
}

function Y(y) {
    return graphSettings.graphHeight / 2 - y * graphSettings.yScale
}

function X(sec) {
    return sec * graphSettings.xScale + X_PADDING
}


function drawYGrids() {
    getContextes((ctx) => {
        ctx.lineWidth = 0.5
        ctx.strokeStyle = 'darkgray';
    });
    let y = -graphSettings.yMax
    while (y <= graphSettings.yMax) {
        const yg = Y(y)
        const xs = X(0)
        const xe = X(graphSettings.timeSpan)
        getContextes((ctx) => {
            ctx.textAlign = "right";
            ctx.strokeStyle = 'darkgray';
            if (y > -graphSettings.yMax && y < graphSettings.yMax) {
                ctx.fillText(`${Math.floor(y * 10) / 10}`, xs - 8, yg + FONT_SIZE * 0.33)
            }
            ctx.lineWidth = 0.5
            ctx.beginPath()
            ctx.lineTo(xs, yg)
            ctx.lineTo(xe, yg)
            ctx.stroke()
        });
        y += graphSettings.yGrid
    }

}

function drawXGrids(datalink, secOffset) {
    let dl = datalink
    while (dl.prev && dl.prev.data) {
        const sec = dl.data.sec - secOffset
        if (sec <= 0) {
            break
        }
        if (dl.data.secMarker || dl.data.decSecMarker) {
            const x = X(sec)
            const yt = Y(graphSettings.yMax)
            const yb = Y(-graphSettings.yMax)
            getContextes((ctx) => {
                ctx.textAlign = "center";
                if (dl.data.secMarker) {
                    ctx.strokeStyle = 'darkgray';
                    ctx.lineWidth = 0.5
                    ctx.fillText(`${Math.floor(dl.data.sec)}`, x, yb + FONT_SIZE * 1.5)
                } else {
                    ctx.strokeStyle = 'lightgray';
                    ctx.lineWidth = 0.25
                }
                ctx.beginPath()
                ctx.lineTo(x, yt)
                ctx.lineTo(x, yb)
                ctx.stroke()
            });
        }

        dl = dl.prev
    }
}

function plot(datalink) {
    let dl = datalink
    const secOffset = dl.data.sec <= graphSettings.timeSpan ? 0 : dl.data.sec - graphSettings.timeSpan
    clearCanvases()
    drawXGrids(dl, secOffset)
    drawYGrids()
    getContextes((ctx) => {
        ctx.strokeStyle = 'darkgreen';
        ctx.lineWidth = 1
        ctx.beginPath()
    });
    let prevX = -1000
    while (dl.prev && dl.prev.data) {
        const sec = dl.data.sec - secOffset
        if (sec <= 0) {
            break
        }
        const x = X(sec)
        const fx = Math.floor(x)
        if (fx !== prevX) {
            ctx2D.accelerometer.lineTo(x, Y(dl.data.a[graphSettings.axis]))
            ctx2D.gyroscope.lineTo(x, Y(dl.data.g[graphSettings.axis]))
            ctx2D.rotations.lineTo(x, Y(dl.data.r[graphSettings.axis]))
            prevX = fx
        }
        dl = dl.prev
    }
    getContextes((ctx) => {
        ctx.stroke()
    });
}

function setSecondsAndMarkers(l) {
    l.data.sec = (l.data.t - graphSettings.firstTime) * TIME_SCALE
    if (l.prev && l.prev.data) {
        const currS = l.data.sec
        const prevS = l.prev.data.sec
        if (Math.floor(currS) - Math.floor(prevS) === 1) {
            l.prev.data.secMarker = false
            l.data.secMarker = true
        } else {
            if (Math.floor(currS * 10) - Math.floor(prevS * 10) === 1) {
                l.prev.data.decSecMarker = false
                l.data.decSecMarker = true
            }
        }
    }
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
        const packets = JSON.parse(JSON.parse(event.data))
        packets.forEach(p => {
            if (graphSettings.firstTime === 0) {
                graphSettings.firstTime = p.t
            }
            link = link.next
            link.data = p
            setSecondsAndMarkers(link)
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