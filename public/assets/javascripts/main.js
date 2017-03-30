function checkextension() {
  var file = document.querySelector("#fileInput");
  var error = document.querySelector("#error");
  if ( /\.(jpe?g)$/i.test(file.files[0].name) === false ) {
     error.classList.remove('hidden');
     file.value = '';
   }
}
