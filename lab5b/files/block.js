// svg adrress
const svgurl = "http://www.w3.org/2000/svg";

// parse and insert svg elements in <block> tags 
function parseBlocks() {
    const blocks = document.getElementsByTagName("block")
    for (const block of blocks) {
        const values = block.getAttribute("values")
        const selected = block.getAttribute("selected")
        const log = block.hasAttribute("log")
        createBlock(block, values, selected, log)
    }
}

// create SVG block
function createBlock(place, values, selected, log) {    
    const svg = document.createElementNS(svgurl, "svg")
    insertDefs(svg)
    const val = values.split(',').map(Number)
    const sel = selected.split(',').map(Number)
    insertValues(svg, val, log)
    const text = document.createElementNS(svgurl, "text")
    text.setAttribute("font-size", "16pt")
    text.setAttribute("text-anchor", "middle")
    text.setAttribute("x", 200)
    text.setAttribute("y", 20)
    text.setAttribute("font-weight", "bold")
    text.appendChild(document.createTextNode(place.getAttribute("name")));
    svg.appendChild(text)
    insertInterface(svg, val, sel, log)
    place.replaceChildren(svg)
}

// insert values to block
function insertValues(svg, val, log) {
    let min = val[0]
    let max = val[val.length - 1]
    for (let v of val) {
        let pos = (v - min) / (max - min)
        if (log && pos > 0.0) pos = Math.log10(1.0 + pos * 9.00)
        let posx = 26 + 360*pos
        let posy = 140
        const line = document.createElementNS(svgurl, "line")
        line.setAttribute("x1", posx - 6)
        line.setAttribute("x2", posx - 6)
        line.setAttribute("y1", 40)
        line.setAttribute("y2", 140)
        line.setAttribute("stroke", "silver")
        line.setAttribute("stroke-width", "0.5px")
        line.setAttribute("opacity", "0.3")
        svg.appendChild(line)
        const text = document.createElementNS(svgurl, "text")
        text.setAttribute("font-size", "16pt")
        text.setAttribute("text-anchor", "end")
        text.setAttribute("transform", "translate("+posx+","+posy+") rotate(-90)")
        text.appendChild(document.createTextNode(v));
        svg.appendChild(text)
    }
}

// insert interface
function insertInterface(svg, val, bez, log) {
    // create circles
    if (bez.length != 4) {
        console.log("selected attribute need 4 numeric values separated by commas", bez)
        return
    }
    let min = val[0]
    let max = val[val.length - 1]
    let circles = []
    for (let i = 0; i < bez.length; i++) {
        let pos = (bez[i] - min) / (max - min)
        if (log && pos > 0.0) pos = Math.log10(1.0 + pos * 9.00)
        let posx = 20 + 360*pos
        let posy = i == 1 || i == 2 ? 40 : 130
        const circle = document.createElementNS(svgurl, "circle")
        circle.setAttribute("cx", posx)
        circle.setAttribute("cy", posy)
        circle.setAttribute("r", 8)
        circle.updates = []
        circles.push(circle)
    }
    // link circles
    for (let i = 1; i < circles.length; i++) {
        circles[i-1].next = circles[i]
        circles[i].prev = circles[i-1]
    }
    // move circles
    let move = null
    for (let i = 0; i < circles.length; i++) {
        circles[i].onmousedown = (event) => {
            move = circles[i]
        }
    }
    // move on svg
    svg.onmousemove = (event) => {
        if (move == null) return
        let px = event.clientX - svg.getBoundingClientRect().left
        if (move.next != null && (px >= parseInt(move.next.getAttribute("cx")))) px = parseInt(move.next.getAttribute("cx"))-1
        if (move.prev != null && (px <= parseInt(move.prev.getAttribute("cx")))) px = parseInt(move.prev.getAttribute("cx"))+1
        move.setAttribute("cx", px)
        for (let uf of move.updates) uf()
    }
    svg.onmouseup = (event) => {
        move = null
    }
    // beziers 
    createBezier(svg, circles[0], circles[1], "url(#redblue)")
    createBezier(svg, circles[2], circles[3], "url(#bluered)")
    // lines
    createLine(svg, null, circles[0], "red")
    createLine(svg, circles[3], null, "red")
    createLine(svg, circles[1], circles[2], "#4169E1")
    // add circles on top
    circles[0].setAttribute("fill", "#7B1113")
    circles[1].setAttribute("fill", "#4169E1")
    circles[2].setAttribute("fill", "#4169E1")
    circles[3].setAttribute("fill", "#7B1113")
    for (let v of circles) svg.appendChild(v)
    // function to return circle values
    svg.getvalues = () => {
        let pos = []
        for (let i = 0; i < circles.length; i++) {
            let cx = parseInt(10000 * (circles[i].getAttribute("cx") - 20) / 360)
            pos.push(cx)
        }
        return "0,"+svg.parentNode.getAttribute("id")+","+min+","+max+","+(log?"1":"0")+","+pos.join(",")
    }
}

// create bezier
function createBezier(svg, ca, cb, stroke) {
    const path = document.createElementNS(svgurl, "path")
    let px1 = parseInt(ca.getAttribute("cx"))
    let py1 = parseInt(ca.getAttribute("cy"))
    let px2 = parseInt(cb.getAttribute("cx"))
    let py2 = parseInt(cb.getAttribute("cy"))
    path.setAttribute("d", bezierPathFromPoints(px1, py1, px2, py2))
    path.setAttribute("stroke", stroke)
    svg.appendChild(path)
    ca.updates.push(() => {
        let px1u = parseInt(ca.getAttribute("cx"))
        let py1u = parseInt(ca.getAttribute("cy"))
        let px2u = parseInt(cb.getAttribute("cx"))
        let py2u = parseInt(cb.getAttribute("cy"))
        path.setAttribute("d", bezierPathFromPoints(px1u, py1u, px2u, py2u))
    })
    cb.updates.push(() => {
        let px1u = parseInt(ca.getAttribute("cx"))
        let py1u = parseInt(ca.getAttribute("cy"))
        let px2u = parseInt(cb.getAttribute("cx"))
        let py2u = parseInt(cb.getAttribute("cy"))
        path.setAttribute("d", bezierPathFromPoints(px1u, py1u, px2u, py2u))
    })
}

// calculate bezier path
function bezierPathFromPoints(px1, py1, px2, py2) {
    let df = parseInt(0.3 * (px2 - px1))
    return "M "+px1+" "+py1+" C "+(px2-df)+" "+py1+", "+(px1+df)+" "+py2+", "+px2+" "+py2+""
}

// create line
function createLine(svg, ca, cb, stroke) {
    const line = document.createElementNS(svgurl, "line")
    let px1 = ca != null ? parseInt(ca.getAttribute("cx")) : 0
    let py1 = ca != null ? parseInt(ca.getAttribute("cy")) : 130
    let px2 = cb != null ? parseInt(cb.getAttribute("cx")) : 400
    let py2 = cb != null ? parseInt(cb.getAttribute("cy")) : 130
    line.setAttribute("x1", px1)
    line.setAttribute("y1", py1)
    line.setAttribute("x2", px2)
    line.setAttribute("y2", py2)
    line.setAttribute("stroke", stroke)
    svg.appendChild(line)
    if (ca != null) {
        ca.updates.push(() => {
            let px1u = parseInt(ca.getAttribute("cx"))
            let py1u = parseInt(ca.getAttribute("cy"))
            line.setAttribute("x1", px1u)
            line.setAttribute("y1", py1u)
        })
    }
    if (cb != null) {
        cb.updates.push(() => {
            let px2u = parseInt(cb.getAttribute("cx"))
            let py2u = parseInt(cb.getAttribute("cy"))
            line.setAttribute("x2", px2u)
            line.setAttribute("y2", py2u)
        })
    }
}

// insert gradient definitions
function insertDefs(svg) {
    const defs = document.createElementNS(svgurl, "defs")
    const g1 = document.createElementNS(svgurl, "linearGradient")
    g1.setAttribute("id", "redblue")
    g1.setAttribute("x1", "0%")
    g1.setAttribute("y1", "0%")
    g1.setAttribute("x2", "100%")
    g1.setAttribute("y2", "0%")
    const s1a = document.createElementNS(svgurl, "stop")
    s1a.setAttribute("offset", "0%")
    s1a.setAttribute("stop-color", "red")
    g1.appendChild(s1a)
    const s1b = document.createElementNS(svgurl, "stop")
    s1b.setAttribute("offset", "100%")
    s1b.setAttribute("stop-color", "royalblue")
    g1.appendChild(s1b)
    const g2 = document.createElementNS(svgurl, "linearGradient")
    g2.setAttribute("id", "bluered")
    g2.setAttribute("x1", "0%")
    g2.setAttribute("y1", "0%")
    g2.setAttribute("x2", "100%")
    g2.setAttribute("y2", "0%")
    const s2a = document.createElementNS(svgurl, "stop")
    s2a.setAttribute("offset", "0%")
    s2a.setAttribute("stop-color", "royalblue")
    g2.appendChild(s2a)
    const s2b = document.createElementNS(svgurl, "stop")
    s2b.setAttribute("offset", "100%")
    s2b.setAttribute("stop-color", "red")
    g2.appendChild(s2b)
    defs.appendChild(g1)
    defs.appendChild(g2)
    svg.appendChild(defs)
}

// call main function
parseBlocks()

// search
function search() {
    let links = ""
    let values = document.querySelectorAll("svg")
    for (let value of values) {
        links += value.getvalues() + ";"
    }
    links += document.getElementById("sortby").value
    fetch("/search/" + links).then(response => {
        if (!response.ok) {
            throw new Error("Błąd sieci: " + response.status)
        }
        return response.json()
    }).then(data => {
        show(data)
    }).catch(error => {
        console.error("Wystąpił błąd:", error)
    })
}

// field formats
const formats = {
    "Screen": (value) => { return value.toFixed(2) },
    "CamFront":(value) => { return value.toFixed(1) },
    "CamBack":(value) => { return value.toFixed(1) },
    "Value": (value) => { return value.toFixed(4) },
    "Weight": (value) => { return value.toFixed(0) },
}

// show data
function show(data) {
    const results = document.getElementById("results")
    if (data == null || data.length == 0) {
        results.innerHTML = "endpoint error"
        return
    } 
    const table = document.createElement("table")
    const tr = document.createElement("tr")
    const first = data[0]
    console.log(first)
    for (let key of Object.keys(data[0])) {
        const th = document.createElement("th")
        th.innerHTML = key
        tr.appendChild(th)       
    }
    table.appendChild(tr)
    for (let phone of data) {
        const tr = document.createElement("tr")
        for (let [key, value] of Object.entries(phone)) {
            const td = document.createElement("td")
            let format = formats[key]
            if (format == null) format = (value) => { return value }
            td.innerHTML = format(value)
            tr.appendChild(td)       
        }
        table.appendChild(tr)
    }
    results.replaceChildren(table)
}
