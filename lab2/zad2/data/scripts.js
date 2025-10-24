const rates = []
let ranked = 0
let getlink = ""

function prepare() {
    let brates = document.getElementsByTagName("brate")
    for (let brate of brates) {
        let bid = brate.getAttribute("id")
        for (let i = 1; i < 6; i++) {
            brate.innerHTML += "<img src='/files/bempty.png' onmouseover='bover(" + i + "," + bid + ")' onmouseout='bshow(" + bid + ")' onclick='selrate(" + i + "," + bid + ")'>"
        }
    }
    generatestats()
}

function bover(selid, bid) {
    let brate = document.getElementById(bid)
    let images = brate.querySelectorAll("img")
    for (let i = 1; i <= selid; i++) {
        images[i - 1].src = "/files/bfull.png"
    }
    for (let i = selid + 1; i < 6; i++) {
        images[i - 1].src = "/files/bempty.png"
    }
} 

function bshow(bid) {
    let selid = 0
    if (rates[bid] != undefined) {
        selid = rates[bid]
    }
    bover(selid, bid)
}

function selrate(selid, bid) {
    rates[bid] = selid
    generatestats()
}

function loadreco() {
    if (ranked < 3) {
        document.getElementById("stats").innerHTML = "Oceń więcej piw!"
        return
    }
    const xhttp = new XMLHttpRequest();
    xhttp.onload = function() {
        document.getElementById("reco").innerHTML = this.responseText
        var arr = document.getElementById("reco").getElementsByTagName('script')
        for (var n = 0; n < arr.length; n++) eval(arr[n].innerHTML)
    }
    xhttp.open("GET", "/reco/" + getlink, true)
    xhttp.send()
}

function generatestats() {
    ranked = 0
    getlink = ""
    let idlist = ""
    for (let i in rates) {
        if (rates[i] != undefined) {
            getlink += i + ";" + rates[i] + ";"
            idlist += i + "=" + rates[i] + ", "
            ranked++
        }        
    }
    let stats = document.getElementById("stats")
    stats.innerHTML = "Ocenionych: " + ranked + "<br>" + idlist
    loadreco()
}
