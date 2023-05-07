const inputForm : HTMLElement | null = document.getElementById('file_form');
console.log('this ran')
if(inputForm){
    console.log('this ran')
    inputForm.addEventListener("submit", (e)=>{
        e.preventDefault()
        let fileInput : FileList | null = document.getElementById('file_input')
        if (!fileInput){return}
        
    
        for(let i = 0; i < fileInput.length; i++){
            const fileForm = new FormData()

            fileForm.append('file', fileInput.item[i])
            fetch('/upl', {
                method: 'POST',
                headers: {
                    'File-Name': fileInput.item[i],
                    'File-Size': fileInput.item[i],
                    'File-Type': fileInput.item[i],
                },
                body: fileForm
            }).then((response)=>{
                console.log(response)
            })
        }
    });
}
