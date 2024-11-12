function alerting(){
    alert("Hello World!");
}


function setAlerts(alerts) {
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

function setComputers(computers) {
    let computersTable = document.getElementById("computersTable");
    let newElement = document.createElement("div");
    newElement.className = "computersBlockContentRow";
    let computer = document.createElement("div");
    computer.className = "computersCell";
    computer.innerHTML = computers.Name;
    let address = document.createElement("div");
    address.className = "computersCell";
    address.innerHTML = computers.Address;
    newElement.appendChild(computer);
    newElement.appendChild(address);
    computersTable.appendChild(newElement);
}

function setRules(rules) {
    let rulesTable = document.getElementById("rulesTable");
    let newElement = document.createElement("div");
    newElement.className = "rulesBlockContentRow";
    let deleteBut = document.createElement("div");
    deleteBut.className = "deleteButton";
    deleteBut.innerHTML = "-";
    newElement.appendChild(deleteBut);
    let layer = document.createElement("div");
    layer.className = "rulesCell";
    layer.innerHTML = rules.Netlayer.Name;
    let srcip = document.createElement("div");
    srcip.className = "rulesCell";
    srcip.innerHTML = rules.SrcIp;
    let dstip = document.createElement("div");
    dstip.className = "rulesCell";
    dstip.innerHTML = rules.DstIp;
    let contains = document.createElement("div");
    contains.className = "rulesCell";
    contains.innerHTML = rules.PayloadContains;
    let srcport = document.createElement("div");
    srcport.className = "rulesCell";
    srcport.innerHTML = rules.SrcPort;
    let dstport = document.createElement("div");
    dstport.className = "rulesCell";
    dstport.innerHTML = rules.DstPort;
    let checksum = document.createElement("div");
    checksum.className = "rulesCell";
    checksum.innerHTML = rules.Checksum;
    newElement.appendChild(layer);
    newElement.appendChild(srcip);
    newElement.appendChild(dstip);
    newElement.appendChild(contains);
    newElement.appendChild(srcport);
    newElement.appendChild(dstport);
    newElement.appendChild(checksum);
    rulesTable.appendChild(newElement);
}
function newAlert(alerts) {
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
    alertsTable.prepend(newElement);
}
function onMessage(event) {
    const response = JSON.parse(event.data);
    if (response.table_name == "alerts") {
        setAlerts(response.data);
    }
    if (response.table_name == "computers") {
        setComputers(response.data);
    }
    if (response.table_name == "rules") {
        setRules(response.data);
    }
    if (response.table_name == "new_computers") {
        setComputers(response.data);
    }
    if (response.table_name == "new_rules") {
        setRules(response.data);
    }
    if (response.table_name == "new_alerts") {
        newAlert(response.data);
    }
}
