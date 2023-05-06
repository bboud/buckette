var inputForm = document.getElementById('file_form');
inputForm.addEventListener("submit", function (e) {
    e.preventDefault();
    var fileInput = document.getElementById('file_input');
    fileInput.classList.add('disabled')

    for (var i = 0; i < fileInput.files.length; i++) {
        var fileForm = new FormData();
        fileForm.append('file', fileInput.files[i]);
        fetch('/upl', {
            method: 'POST',
            headers: {
                'File-Name': fileInput.files[i].name,
                'File-Size': fileInput.files[i].size,
                'File-Type': fileInput.files[i].type,
            },
            body: fileForm
        }).then(function (response) {
            console.log(response);
        });
    }
    fileInput.classList.remove('disabled')
});
