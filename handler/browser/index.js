const api = axios.create({
  baseURL: "/",
  headers: {
    Accept: "application/json",
  },
});

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
    loading: false,

    async handleClickBack() {
      const storeDir = Alpine.store("dir");
      if (!storeDir.canBack) {
        return;
      }

      try {
        this.loading = true;
        await storeDir.loadDir(storeDir.oneDirectoryBackward());
      } finally {
        this.loading = false;
      }
    },

    async handleClickFile(file) {
      const storeDir = Alpine.store("dir");
      const { name, isDir } = file;
      if (isDir) {
        try {
          this.loading = true;
          await storeDir.loadRelativeDir(name);
        } finally {
          this.loading = false;
        }
        return;
      }

      storeDir.openFile(name);
    },
  };
}

document.addEventListener("alpine:init", () => {
  Alpine.store("dir", {
    pathParts: [],
    fileList: null,
    loading: false,

    get path() {
      return this.pathPartsToString(this.pathParts);
    },

    get canBack() {
      return this.pathParts.length >= 3;
    },

    init() {
      this.curDir = this.getCurDir();
      this._loadDirPathPart(this.curDir);
    },

    pathPartsToString(pp) {
      return pp.join("/");
    },

    pathToPathParts(p) {
      const paths = p.split("/");
      return paths;
    },

    getCurDir() {
      const curUrl = new URL(window.location.toString());
      const dir = curUrl.searchParams.get("dir");
      return dir ? this.pathToPathParts(dir) : this.pathToPathParts("/");
    },

    oneDirectoryBackward() {
      if (!this.canBack) {
        return null;
      }

      this.pathParts.splice(this.pathParts.length - 2, 1);
      return this.pathPartsToString(this.pathParts);
    },

    async loadRelativeDir(dir) {
      const dirPp = dir.split("/");
      const newPathParts = [
        ...this.pathParts.slice(0, -1),
        ...dirPp,
        this.pathParts[this.pathParts.length - 1],
      ];
      this._loadDirPathPart(newPathParts);
    },

    openFile(name) {
      const newPathParts = this.pathParts.slice(0, -1).concat(name);
      const p = this.pathPartsToString(newPathParts);
      const destUrl = new URL(p, window.location.toString());
      window.location.href = destUrl.href;
    },

    async loadDir(dir) {
      if (!dir || !dir.startsWith("/") || !dir.endsWith("/")) {
        throw new Error(`InvalidArgumentException dir: ${dir}`);
      }
      this._loadDirPathPart(this.pathToPathParts(dir));
    },

    async _loadDirPathPart(pp) {
      // TODO: usar pushState?
      // https://stackoverflow.com/questions/824349/how-do-i-modify-the-url-without-reloading-the-page

      try {
        this.loading = true;
        const dir = this.pathPartsToString(pp);
        const apiRes = await api.get(dir);
        this.fileList = apiRes.data;
        this.pathParts = pp;
      } catch (error) {
        console.error("error", error);
        throw error;
      } finally {
        this.loading = false;
      }
    },
  });

  // Alpine.effect(() => {
  //   const pathParts = Alpine.store('dir').pathParts;
  //   const path = Alpine.store('dir').path;

  //   console.debug(pathParts, path);
  // })
});
