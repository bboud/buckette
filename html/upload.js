$('#file_form').submit(function (e) {
    e.preventDefault();

    $('#file_input').attr('disabled', 'disabled');
    let fileInputs = $('#file_input').prop('files')

    promises = []

    for (let i = 0; i < fileInputs.length; i++) {
        promises.push(uploadFile(fileInputs[i], i))
    }
    Promise.all(promises).then((e) => {
        console.log(e)
        $('#file_input').val('');
        $('#file_input').removeAttr('disabled');
    }).catch((e) => {
        console.log(e)
    })
});

function uploadFile(file, index) {
    console.log(file)
    return new Promise((resolve, rejectt) => {
        var $progress_bar = $(`
        <div id="progress_bar_${index}" class="col-12 text-start">
            File name
            <div class="progress" role="progressbar" aria-label="File progress bar for ${file.name}" aria-valuenow="75" aria-valuemin="0" aria-valuemax="100">
            <div class="progress-bar" style="width: 0%"></div>
            </div>
        </div>
        `)

        let formData = new FormData()

        $('#progress_row').append($progress_bar)
        formData.append('file', file)

        $.ajax({
            xhr: ()=>{
                var xhr = new window.XMLHttpRequest();
                xhr.addEventListener("progress", (evt) => {
                    if (evt.lengthComputable) {
                        var percentComplete = (evt.loaded / evt.total) * 100;
                        // Place upload progress bar visibility code here
                    }
                }, false)
                return xhr
            },
            type: "POST",
            url: "/upl",
            data: formData,
            contentType: false,
            processData: false,
            headers: {
                'File-Name': file.name,
                'File-Type': file.type,
                'File-Size': file.size,
            },
            dataType: "json",
            success: function (response) {
                console.log(response)
            },
        });
    })
}