document.addEventListener("DOMContentLoaded", function() {
    var btnRek    = document.getElementById('rek');
    var btnSay    = document.getElementById('say');
    var btnRekSay = document.getElementById('reksay');
    var w         = document.getElementById('w');

    var rek = new webkitSpeechRecognition(); // 音声認識APIの使用
    rek.lang = "ja"; // 言語を日本語に設定

    function insertTextArea(val) {
        w.value = val;
    }

    // ボタンクリックで認識開始
    btnRek.addEventListener('click', function () {
        rek.onresult = function (e) {
            if(e.results.length > 0){
                var word = e.results[0][0].transcript;
                insertTextArea(word);
            }
        };

        rek.start();
    });

    btnSay.addEventListener('click', function () {
        var wVal = w.value;
        // alert(wVal);
        var xhr = new XMLHttpRequest();
        xhr.onreadystatechange = function() {
            if (xhr.readyState === 4) {
                if (xhr.status == 200) {
                    // alert("OK"); //通信成功時
                } else {
                    // alert("NO"); //通信失敗時
                }
            }
        };
        xhr.onload = function() {
            // alert("complete"); //通信完了時
        };
        xhr.open("POST", "/say", true);
        xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        xhr.send("w="+wVal);
    });

    btnRekSay.addEventListener('click', function () {
        rek.onresult = function (e) {
            if(e.results.length > 0){
                var word = e.results[0][0].transcript;
                insertTextArea(word);
                btnSay.click();
            }
        };
        rek.start();
    });
});
