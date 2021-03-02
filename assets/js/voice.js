const URL = '/upload';
let div = document.createElement('div');
div.id = 'messages';
let start = document.createElement('button');
start.id = 'start';
start.innerHTML = 'Start';
let stop = document.createElement('button');
stop.id = 'stop';
stop.innerHTML = 'Stop';
document.body.appendChild(div);
document.body.appendChild(start);
document.body.appendChild(stop);

var options = {
    tag: 'audio',
    type: 'audio/ogg',
    ext: '.ogg',
    gUM: {audio: true}
}



navigator.mediaDevices.getUserMedia({ audio: true})
    .then(stream => {
        const mediaRecorder = new MediaRecorder(stream,options);

        document.querySelector('#start').addEventListener('click', function(){
            mediaRecorder.start();
        });
        let audioChunks = [];
        mediaRecorder.addEventListener("dataavailable",function(event) {
            audioChunks.push(event.data);
        });

        document.querySelector('#stop').addEventListener('click', function(){
            mediaRecorder.stop();
        });

        mediaRecorder.addEventListener("stop", function() {
            const audioBlob = new Blob(audioChunks, {
                type: 'audio/ogg; codecs=0'
            });

            let fd = new FormData();
           // fd.append('fname', 'test1.wav');
            fd.append('voice', audioBlob);

            sendVoice(fd);
            audioChunks = [];
        });
    });

async function sendVoice(form) {
    let promise = await fetch(URL, {
        method: 'POST',
        body: form});
    if (promise.ok) {
        let response =  await promise.json();
       // alert(response);
        console.log(response.data);

        var desc = document.getElementById('desc');
        var textarea = document.createElement('textarea');

        desc.textContent = '';
        desc.appendChild(textarea);
        var message = document.createTextNode(response.data);
        textarea.appendChild(message);
        let audio = document.createElement('audio');
        audio.src = response.data;
        audio.controls = true;
        audio.autoplay = true;
        document.querySelector('#messages').appendChild(audio);
    }
}