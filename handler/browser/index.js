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

/**
 * @param {number} mode
 * @returns string chmod em octal
 */
function formatChmod(mode) {
  // https://ss64.com/bash/chmod.html
  // chmod calculator
  return mode.toString(8);
}

function formatDate(date) {
  try {
    if (date) {
      return new Date(date).toLocaleString();
    }
  } catch (error) {
    console.error("date parse error:", error);
  }
  return "?";
}

function fileListDataFunc() {
  return {
    fileList: null,
    loading: false,
    curDir: null,

    init() {
        this.curDir = this.getCurDir();
      this.loadDir(this.curDir);
    },

    getCurDir() {
      const curUrl = new URL(window.location.toString());
      const dir = curUrl.searchParams.get("dir");
      return dir || "/";
    },

    async handleClickBack() {
        console.debug("click back")
        const paths = this.curDir.slice(0, this.curDir.length-1).split("/");
        paths.pop();
        newUrl = `${paths.join("/")}/`;
        await this.loadDir(newUrl);
        // TODO: usar um watch para curDir? ou centralizar mudanças?
        this.curDir = newUrl;
    },

    async handleClickFile(file) {
        const { name, isDir } = file;
        console.debug("file", file,name, isDir);
        // TODO: usar pushState?
        // https://stackoverflow.com/questions/824349/how-do-i-modify-the-url-without-reloading-the-page
      if (isDir) {
        const newUrl = `${this.curDir}${name}/`;
        await this.loadDir(newUrl);
        this.curDir = newUrl;
      } else {
        const destUrl = new URL(`${this.curDir}${name}`, window.location.toString());
        window.location.href = destUrl.href;
      }
    },

    async loadDir(dir) {
      const headers = { Accept: "application/json" };
      try {
        this.loading = true;
        const apiRes = await axios.get(dir, { headers });
        console.debug("apires", apiRes.data);
        this.fileList = apiRes.data;
      } catch (error) {
        console.error("error", error);
        throw error;
      } finally {
        this.loading = false;
      }
    },
  };
}
