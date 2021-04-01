import 'bulma';
import '../style/index.scss';
import Vue from 'vue/dist/vue.esm';
import axios from 'axios';
import '@fortawesome/fontawesome-free/js/all';

import copyToClipboard from './clipboard';
import wangwangAudio from '../audio/wangwang.mp3';

let app;
async function upload(files) {
    app.isUploading = true;
    for (let i = 0; i < files.length; i++) {
        let file = files[i];
        let fd = new FormData();
        fd.append("file", file);
        app.fileSize = file.size;
        app.fileUploadSize = 0;
        app.fileName = file.name + " (" + (i + 1) + "/" + files.length + ")"
        try {
            let res = await axios({
                method: "post",
                url: "file",
                data: fd,
                onUploadProgress: function (event) {
                    //console.log(event);
                    if (event.lengthComputable) {
                        app.fileUploadSize = event.loaded;
                    }
                    else {
                        alert('unable to compute');
                    }
                }
            });
            if (res.data.code != 0) {
                alert(res.data.result);
            } else {
                // 添加上传的文件信息
                app.fileList.unshift(res.data.result)
            }
        } catch (err) {
            alert(err);
        }
    }
    app.isUploading = false;
}

function updateFileList() {
    axios.get('file')
        .then((res) => {
            if (res.data.code == 0) {
                app.fileList = res.data.result;
            } else {
                alert(res.data.result);
            }
        })
        .catch(err => {
            console.log(err);
        });
}

function init() {
    app = new Vue({
        el: "#app",
        data: {
            fileName: "",
            fileSize: 0,
            fileUploadSize: 0,
            isUploading: false,
            fileUploadPercent: "0%",
            fileList: []
        },
        created() {
            updateFileList();
            // 刷新倒计时
            setInterval(() => {
                app.fileList = app.fileList.filter((item) => item.expireTimestamp > new Date());
            }, 1000);
        },
        computed: {
            fileSizeHuman() {
                let fileSize = this.fileSize;
                return this.getHumanFileSize(fileSize);
            }
        },
        watch: {
            fileUploadSize(newVal, oldVal) {
                if (this.fileSize == 0) {
                    this.fileUploadPercent = 0 + "%";
                    return
                }
                let percent = parseInt((this.fileUploadSize / this.fileSize) * 100);
                if (percent >= 100) {
                    percent = 100;
                }
                this.fileUploadPercent = percent + "%";
            }
        },
        methods: {
            wangwang() {
                // 由于file-loader打包限制，这里需要手动解决路径
                let audio = new Audio(location.href + "assets" + wangwangAudio.substr(1));
                audio.play();
                audio.onplaying = () => {
                    alert("🐶wangwang!!");
                };
            },
            uploadFile(event) {
                //console.log(event);
                let target = event.target;
                let files;
                if (target) {
                    files = target.files;
                }
                if (files && files.length > 0) {
                    upload(files);
                } else {
                    this.fileName = "";
                    this.fileSize = 0;
                    this.fileUploadSize = 0;
                    this.isUploading = false;
                }
            },
            copyURL(token) {
                let url = location.href + "file/" + token;
                copyToClipboard(url);
            },
            copyMD(name,token) {
                let md = `[${name}](${location.href}file/${token})`;
                copyToClipboard(md);
            },
            openURL(token) {
                let url = location.href + "file/" + token;
                open(url);
            },
            getHumanFileSize(fileSize) {
                if (fileSize < 1024) {
                    return fileSize + 'B';
                }
                if (fileSize < 1024 * 1024) {
                    return `${(fileSize / 1024).toFixed(2)}KB`;
                }
                if (fileSize < 1024 * 1024 * 1204) {
                    return `${(fileSize / 1024 / 1024).toFixed(2)}MB`;
                }
                if (fileSize < 1024 * 1024 * 1204 * 1024) {
                    return `${(fileSize / 1024 / 1024 / 1024).toFixed(2)}GB`;
                }
            },
            getCountDown(expireTimestamp) {
                let t = parseInt((expireTimestamp - (new Date())) / 1000);
                let h = String(parseInt(t / 3600));
                if (h.length == 1) {
                    h = `0${h}`;
                }
                let m = `0${parseInt((t % 3600) / 60)}`.substr(-2);
                let s = `0${(t % 3600) % 60}`.substr(-2);
                return `${h}:${m}:${s}`;
            },
            fileDragEnter(event) {
                event.preventDefault();
                event.stopPropagation();
                if(this.isUploading){
                    alert("You can't drag files when uploding...");
                }
            },
            fileDragOver(event) {
                event.preventDefault();
                event.stopPropagation();
            },
            fileDrop(event) {
                event.preventDefault();
                event.stopPropagation();
                if(this.isUploading){
                    return;
                }
                // 去掉目录
                (async (files) => {
                    let targets = [];
                    let isFile = false;
                    for (let file of files) {
                        try {
                            isFile = await new Promise((resolve, reject) => {
                                let reader = new FileReader();
                                reader.onload = () => {
                                    resolve(true);
                                };
                                reader.onerror = () => {
                                    resolve(false);
                                };
                                reader.readAsArrayBuffer(file.slice(0, file.size <= 10 ? file.size : 10));
                            });
                        } catch (error) {
                            isFile = false;
                        }
                        if (isFile) {
                            targets.push(file);
                        }
                    }
                    return targets;
                })(event.dataTransfer.files)
                    .then(files => upload(files));
            }
        }
    });
}



window.addEventListener('load', init);