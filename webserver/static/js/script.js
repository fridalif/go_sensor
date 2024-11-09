function alerting(){
    alert("Hello World!");
}


function onMessage(event) {
    const alerts = JSON.parse(event.data);
    console.log(alerts);
    let alertsTable = document.getElementById("alertsTable");
    let newElement = document.createElement("div");
    newElement.className = "rulesBlockContentRow";
    let computer = document.createElement("div");
    computer.className = "rulesCell";
    computer.innerHTML = alerts.Computer.Name;
    let time = document.createElement("div");
    time.className = "rulesCell";
    time.innerHTML = alerts.Timestamp;
    let layer = document.createElement("div");
    layer.className = "rulesCell";
    layer.innerHTML = alerts.Rule.Netlayer.Name;
    let srcip = document.createElement("div");
    srcip.className = "rulesCell";
    srcip.innerHTML = alerts.Rule.SrcIp;
    let dstip = document.createElement("div");
    dstip.className = "rulesCell";
    dstip.innerHTML = alerts.Rule.DstIp;
    let contains = document.createElement("div");
    contains.className = "rulesCell";
    contains.innerHTML = alerts.Rule.PayloadContains;
    let srcport = document.createElement("div");
    srcport.className = "rulesCell";
    srcport.innerHTML = alerts.Rule.SrcPort;
    let dstport = document.createElement("div");
    dstport.className = "rulesCell";
    dstport.innerHTML = alerts.Rule.DstPort;
    let checksum = document.createElement("div");
    checksum.className = "rulesCell";
    checksum.innerHTML = alerts.Rule.Checksum;
    let ttl = document.createElement("div");
    ttl.className = "rulesCell";
    ttl.innerHTML = alerts.Rule.TTL;
    newElement.appendChild(computer);
    newElement.appendChild(time);
    newElement.appendChild(layer);
    newElement.appendChild(srcip);
    newElement.appendChild(dstip);
    newElement.appendChild(contains);
    newElement.appendChild(srcport);
    newElement.appendChild(dstport);
    newElement.appendChild(checksum);
    newElement.appendChild(ttl);
    alertsTable.appendChild(newElement);
}
