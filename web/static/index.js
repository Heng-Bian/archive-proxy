(() => {
    const remoteUrl = document.getElementById("url-input")
    const charset = document.getElementById("encoding-input")
    const entries = document.getElementById("download-entries")

    const readButton = document.getElementById("read")
    const downloadButton = document.getElementById("download")
    const urlInput = document.getElementById("url-input")

    function read() {
        try {
            new URL(remoteUrl.value)
            let queryParams = new URLSearchParams()
            queryParams.set("charset", charset.value)
            queryParams.set("url", remoteUrl.value)
            fetch("/list" + "?" + queryParams.toString())
                .then(res => {
                    if (res.status == 500) {
                        res.text().then(message => {
                            alert(message)
                        })
                    } else {
                        res.json().then(
                            json => {
                                files = json["Files"]
                                let entries = document.getElementById("download-entries");
                                entries.innerHTML = ""
                                for (let i in files) {
                                    let option = document.createElement("option")
                                    option.appendChild(document.createTextNode(files[i]));
                                    option.setAttribute("value", files[i]);
                                    entries.appendChild(option);
                                }
                            }
                        )
                    }
                })
        } catch (error) {
            alert(errovaluer)
        }
    }

    function download() {
        let files = []
        let options = entries.options
        for (let i in options) {
            if (options[i].selected) {
                files.push(options[i].value)
            }
        }
        if (files.length == 0) {
            alert("no seleted files!")
            return
        }
        let queryParams = new URLSearchParams()
        queryParams.set("charset", charset.value)
        queryParams.set("url", remoteUrl.value)
        fetch("/pack" + "?" + queryParams.toString(), {
            method: 'POST',
            body: JSON.stringify(files),
        }).then(res => {
            if (res.status == 500) {
                res.text().then(message => {
                    alert(message)
                })
            } else {
                const fileStream = streamSaver.createWriteStream("package.zip")
                const readableStream = res.body;
                if (window.WritableStream && readableStream.pipeTo) {
                    return readableStream.pipeTo(fileStream).then(() => {

                    })
                }
                window.writer = fileStream.getWriter();
                const reader = res.body.getReader();
                const pump = () => reader.read().then(res => res.done ? window.writer.close() : window.writer.write(res.value).then(pump))
                pump()
            }
        })
    }

    function clearFiles() {
        entries.innerHTML = ""
    }

    readButton.addEventListener("click", read)
    downloadButton.addEventListener("click", download)
    urlInput.addEventListener("change", clearFiles)
})()

