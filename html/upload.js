$('#file_form').submit(function (e) {
    e.preventDefault();

    //$('#file_input').attr('disabled', 'disabled');
    $('#file_input').fadeOut()

    let fileInputs = $('#file_input').prop('files')

    promises = []

    for (let i = 0; i < fileInputs.length; i++) {
        promises.push(uploadFile(fileInputs[i], i))
    }
    Promise.all(promises).then((values) => {
        console.log(values)
        $('#file_input').val('');
        $('#file_input').removeAttr('disabled');
    }).catch((e) => {
        $('#file_input').val('');
        $('#file_input').removeAttr('disabled');
        console.log(e)
    })
});

var uploadFile = function (file, index) {
    var $progress_bar = $(`
        <div class="col-12 text-start upload_removeable">
            ${file.name}
            <div class="progress" role="progressbar" aria-label="File progress bar for ${file.name}" aria-valuenow="75" aria-valuemin="0" aria-valuemax="100">
            <div id="progress_bar_${index}_bar" class="progress-bar" style="width: 0%"></div>
            </div>
        </div>
        `)

    let formData = new FormData()

    $('#progress_row').append($progress_bar)
    formData.append('file', file)

    return $.ajax({
        xhr: () => {
            var xhr = new XMLHttpRequest();
            xhr.onprogress = (p) => {
                if (p.lengthComputable) {
                    console.log(p)
                    var percentComplete = (p.loaded / p.total) * 100;
                    $(`#progress_bar_${index}_bar`).width(`${percentComplete}%`);
                }
            }
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
    });
}

$('#modal').on('hide.bs.modal', ()=>{
    console.log('this is called')
    $('.upload_removeable').remove()
    $('#file_input').show()
});