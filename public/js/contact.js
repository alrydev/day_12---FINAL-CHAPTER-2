function submitData () {
    
    let name = document.getElementById("input-name").value
    let email = document.getElementById("input-email").value
    let number = document.getElementById("input-number").value
    let subject = document.getElementById("input-subject").value
    let message = document.getElementById("input-message").value


    let a = document.createElement('a')
    let emailReceiver = "alriydev@gmail.com"

    a.href = `mailto:${emailReceiver}?subject=${subject}&body=Nama saya ${name}, ${message}. Tolong hubungi kembali di ${email} ${number} Terima Kasih` 
    a.click()


}
