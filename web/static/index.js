function read() {
    let remoteUrl = document.getElementById("url-input").value
    try {
        let url = new URL(remoteUrl)
    } catch(error) {
        alert(error)
    }

}
function download() {
    alert("download")
}