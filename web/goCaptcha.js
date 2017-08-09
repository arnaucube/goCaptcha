
var selection = [0, 0, 0, 0, 0, 0];
var captchaid = "";

function httpGet(url) {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open("GET", url, false); // false for synchronous request
    xmlHttp.send(null);
    return xmlHttp.responseText;
}

function httpPost(url, data) {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open("POST", url, false);
    xmlHttp.setRequestHeader("Access-Control-Allow-Origin", "*");
    xmlHttp.send(data);
    return xmlHttp.responseText;
}

function getCaptcha() {
    data = httpGet(goCaptchaURL + "/captcha")
    captcha = JSON.parse(data);
    captchaid = captcha.id;
    showCaptcha(captcha);
}

function showCaptcha(captcha) {
    var html;
    html = "";
    html += "<h2>Select all " + captcha.question + "s</h2>";
    for (k in captcha.imgs) {
        html += "<img id='" + k + "' onclick='selectCaptchaImg(this)' src='" + goCaptchaURL + "/image/" + captcha.imgs[k] + "' style='width:150px;cursor:pointer;' />";
    }
    html += "<div onclick='validateCaptcha()' class='g_button c_blue300 g_floatRight'>Validate</div>";
    document.getElementById("goCaptcha").innerHTML = html;
}

function selectCaptchaImg(elem) {
    if (selection[elem.id] == 0) {
        selection[elem.id] = 1;
        document.getElementById(elem.id).className = "g_selected";
    } else {
        selection[elem.id] = 0;
        document.getElementById(elem.id).className = "g_unselected";
    }
}

function validateCaptcha() {
    var answer = {
        selection: selection,
        captchaid: captcha.id
    };
    data = httpPost(goCaptchaURL + "/answer", JSON.stringify(answer));
    resp = JSON.parse(data);
    var html = "";
    if (resp) {
        html += "<h2>goCaptcha validated</h2>";
    } else {
        selection = [0, 0, 0, 0, 0, 0];
        html += "<h2>goCaptcha failed</h2>";
        html += "<div onclick='getCaptcha()' class='g_button c_red300 g_floatRight'>Reload Captcha</div>";
    }
    document.getElementById("goCaptcha").innerHTML = html;
}

if (document.getElementById("goCaptcha")) {
    getCaptcha();
}
