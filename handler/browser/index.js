console.debug("ok");

const fileForm = document.querySelector("#file-form");
fileForm.addEventListener("submit", (evt) => {
  evt.preventDefault();
  const formData = new FormData();
  const fileInput = document.querySelector("#file-input");
  formData.append("file", fileInput.files[0]);
  sendRequestUpload(formData);
});

function updateUploadProgress(text) {
  const uploadProgressDOM = document.querySelector("#upload-progress");
  uploadProgressDOM.innerText = text;
}

function sendRequestUpload(formData) {
  axios
    .post("./", formData, {
      onUploadProgress: (event) => {
        const progress = Math.round((event.loaded * 100) / event.total);
        const text =
          progress +
          "% (" +
          formatBytes(event.loaded) +
          " de " +
          formatBytes(event.total) +
          ")";
        updateUploadProgress(text);
      },
    })
    .then((response) => {
      console.info("O arquivo já foi enviada para o servidor");
      updateUploadProgress(
        "O arquivo já foi enviada para o servidor. Atualizando arquivos..."
      );
      fileForm.reset();
      setTimeout(() => location.reload(), 2000); // Reload the current page
    })
    .catch((err) => {
      console.error(
        "Houve um problema ao realizar o upload do arquivo no servidor",
        err
      );
      updateUploadProgress(
        "Houve um problema ao realizar o upload do arquivo no servidor"
      );
    });
}

function formatBytes(bytes, decimals = 2) {
  if (bytes === 0) return "0 Bytes";
  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
}
