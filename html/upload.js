$('#file_form').submit(function (e) { 
    e.preventDefault();
    $('#file_input').attr('disabled', 'disabled');
    let fileInputs = $('#file_input').prop('files')

    let formData = new FormData()
    promises = []

    for(let i=0; i < fileInputs.length; i++){
        formData.append('file', fileInputs[i])
        promises.push(fetch('/upl', {
            method: "post",
            headers: {
                'File-Name': fileInputs[i].name,
                'File-Size': fileInputs[i].size,
                'File-Type': fileInputs[i].type,
            },
            body: formData,
        }));
    }
    Promise.all(promises).then(()=>{
        $('#file_input').removeAttr('disabled');
    })
});