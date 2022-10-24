const api = axios.create({
  baseURL: "/",
  headers: {
    Accept: "application/json",
  },
});

function formatBytes(bytes, decimals = 2, unit = "kB") {
  if (bytes === 0) return "0 Bytes";
  const k = unit === "Kib" ? 1000 : 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes =
    unit === "Kib"
      ? ["Bytes", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"]
      : ["Bytes", "kB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  const n = bytes / Math.pow(k, i);
  const nLocalText = n.toLocaleString(undefined, {
    maximumFractionDigits: dm,
  });
  return `${nLocalText} ${sizes[i]}`;
}

/**
 * @param {number} mode
 * @returns string chmod em octal
 */
function formatChmod(mode) {
  // https://ss64.com/bash/chmod.html
  // https://chmod-calculator.com/
  // https://unix.stackexchange.com/questions/39710/how-to-get-permission-number-by-string-rw-r-r

  // TODO: mode 20000000775
  // ex: mode = 436 (0o664)
  // (0o664).toString(2) -> '110110100'
  // (0o664 & 0b111).toString(2) -> '100'
  // (0o664 >> 3 & 0b111).toString(2) -> '110'
  // (0o664 >> 6 & 0b111).toString(2) -> '110'
  // return mode.toString(8);

  const fileModeDir = 1 << (32 - 1);
  const rBit = 0b100;
  const wBit = 0b010;
  const xBit = 0b001;

  const chmodText = (n) => {
    const r = (n & rBit) !== 0 ? "r" : "-";
    const w = (n & wBit) !== 0 ? "w" : "-";
    const x = (n & xBit) !== 0 ? "x" : "-";
    return `${r}${w}${x}`;
  };
  
  const dir = mode & fileModeDir !== 0 ? "d" : "-";
  const ogpText = {
    owner: chmodText((mode >> 6) & 0b111),
    group: chmodText((mode >> 3) & 0b111),
    public: chmodText(mode & 0b111),
  };

  return `${dir}${ogpText.owner}${ogpText.group}${ogpText.public}`;
}

const fileFormData = () => ({
  uploading: false,
  uploadProgress: {
    loaded: 0,
    total: null,
    progress: 0,
    rate: null,
  },
  selectedFiles: [],
  isDragover: false,

  get uploadProgressText() {
    return `${(this.uploadProgress.progress * 100).toLocaleString()} %`;
  },

  /** @param {FileList} fileList */
  handleFileList(fileList) {
    // console.debug("fileList", fileList);

    if (!fileList || fileList.length == 0) {
      // TODO: Alert aqui
      console.warn("!fileList");
      return;
    }

    /** @type {Array<File>}  */
    const fileArr = [];
    for (let i = 0; i < fileList.length; i++) {
      const file = fileList[i];
      //   const file = files.item(i);
      // if (file.type !== "application/pdf") {
      //   // TODO: Alert aqui
      //   console.warn("!arquivo não pdf", file, file.type);
      //   continue;
      // }

      fileArr.push(file);
    }

    // Não permite arquivos duplicados
    const newSelectedFiles = fileArr.concat(this.selectedFiles);
    const newUniqueSelectedFiles = [];
    for (let i = 0; i < newSelectedFiles.length; i++) {
      const file = newSelectedFiles[i];
      const fileIdx = newSelectedFiles.findIndex(
        (sf) => sf.name === file.name && sf.size === file.size
      );
      const fileExists = fileIdx !== i;
      if (fileExists) {
        // TODO: Alert aqui
        console.warn("!arquivo já carregado ignorado", file, file.name);
        continue;
      }

      newUniqueSelectedFiles.push(file);
    }

    this.selectedFiles = newUniqueSelectedFiles;
  },

  /** @param {DragEvent} event */
  handleDropFiles(event) {
    const { files } = event.dataTransfer;
    this.handleFileList(files);
  },

  handleFileInputChange(e) {
    /** @type {{ target: HTMLInputElement }} */
    const { target: fileInput } = e;
    const { files } = fileInput;
    this.handleFileList(files);

    // Limpa o fileInput para que o onchage seja chamado quando selecionar o mesmo arquivo
    // https://stackoverflow.com/questions/1703228/how-can-i-clear-an-html-file-input-with-javascript
    fileInput.value = "";
  },

  handleFormSubmit() {
    if (!this.selectedFiles.length) {
      return;
    }

    const formData = new FormData();
    for (const f of this.selectedFiles) {
      formData.append("file", f);
    }
    this.sendRequestUpload(formData);
  },

  async sendRequestUpload(formData) {
    const onUploadProgress = (e) => {
      this.uploadProgress.loaded = e.loaded;
      this.uploadProgress.progress = e.progress;
      this.uploadProgress.rate = e.rate;
      this.uploadProgress.total = e.total;
    };

    const headers = {
      "Content-Type": "multipart/form-data",
    };

    try {
      this.uploading = true;
      const serverRes = await axios.post(this.$store.dir.path, formData, {
        headers,
        onUploadProgress,
      });
      console.debug(serverRes);
      console.info("O arquivo já foi enviada para o servidor");
      // TODO:
      // updateUploadProgress(
      //   "O arquivo já foi enviada para o servidor. Atualizando arquivos..."
      // );
      // fileForm.reset();
    } catch (error) {
      console.error(
        "Houve um problema ao realizar o upload do arquivo no servidor",
        err
      );
      // TODO:
      // updateUploadProgress(
      //   "Houve um problema ao realizar o upload do arquivo no servidor"
      // );
    } finally {
      this.uploading = false;
    }
  },
});

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

    formatDate(date) {
      try {
        if (date) {
          return new Date(date).toLocaleString(undefined, {
            day: "numeric",
            month: "short",
            year: "numeric",
            weekday: "short",
            hour: "2-digit",
            minute: "2-digit",
            second: "2-digit"
          });
        }
      } catch (error) {
        console.error("date parse error:", error);
      }
      return "?";
    }
  };
}

document.addEventListener("alpine:init", () => {
  Alpine.store("dir", {
    pathParts: [],
    fileList: [],
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

  Alpine.directive(
    "fbytes",
    (el, { value, modifiers, expression }, { evaluate }) => {
      const unit = value === "kib" ? "Kib" : undefined;
      const dec = parseInt(modifiers?.[0]) || undefined;
      const num = evaluate(expression);
      // console.debug(value,modifiers, num, unit,dec);
      const res = formatBytes(num, dec, unit);
      // console.debug(num, dec, unit, res);
      el.textContent = res;
    }
  );
});
