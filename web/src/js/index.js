import 'bulma';
import '../style/index.scss';
import Vue from 'vue/dist/vue.esm';
import axios from 'axios';
import '@fortawesome/fontawesome-free/js/all';

import copyToClipboard from './clipboard';
import wangwangAudio from '../audio/wangwang.mp3';

let app;
function upload(fileInput) {
    let fd = new FormData();
    fd.append("file", fileInput.files[0]);
    axios({
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
    }).then(res => {
        if (res.data.code != 0) {
            alert(res.data.result);
        } else {
            // 添加上传的文件信息
            app.fileList.unshift(res.data.result)
        }
        app.isUploading = false;
    }).catch(err => {
        console.log(err);
    });
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
            syncFilename(event) {
                //console.log(event);
                let target = event.target;
                let files;
                if (target) {
                    files = target.files;
                }
                let file;
                if (files.length > 0) {
                    file = files[0];
                }
                if (file) {
                    this.fileName = file.name;
                    this.fileSize = file.size;
                    this.isUploading = true;
                    upload(target);
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
            }
        }
    });
}



window.addEventListener('load', init);