const leftBody = document.getElementById('leftBody');
const leftBodyDefault = document.getElementById('leftBodyDefault');
const rightHead = document.getElementById('rightHead');
const rightBody = document.getElementById('rightBody');
const configurationElem = document.getElementsByClassName('configurationElem')[0];
let preClickKeyElem;
let modifyStatus = false;

function newKeyElem(key) {
    const keyElem = document.createElement('div');
    keyElem.style.cssText = 'width: 100%; height: 50px; line-height: 50px; text-align: center; border-bottom: 1px solid #C0C4CC; cursor: pointer; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; background-color: #f0f9eb; position: relative;';
    keyElem.classList.add('keyElem');

    const delBtn = document.createElement('span');
    delBtn.style.cssText = 'color: #f89898; font-size: 20px; cursor: default; display: none; position: absolute; right:15px';
    delBtn.textContent = 'x';

    delBtn.addEventListener('click', (event) => {
        event.stopPropagation(); //阻止向上传播(父元素还有click事件)
        if (!confirm("通常只会在因意外关闭log-agent(目前已知有：删除log-agent容器)而导致无法注销服务时删除，确定删除吗？")) {
            return;
        }
        if (preClickKeyElem === keyElem) {
            if (!modifyInterceptor("你修改了配置项但还未提交，确定删除吗？")) {
                return;
            }
            rHeadAndBodyToDefault();
            keyElem.style.backgroundColor = '#f0f9eb';
        }
        axios.delete('/key/' + getKeyElemInnerText(keyElem))
            .then(res => res.data)
            .then(function (res) {
                if (res.code !== 0) {
                    alert(res.msg);
                    return;
                }
                keyElem.remove();
                if (leftBody.querySelectorAll('.keyElem').length === 0) {
                    leftBodyDefault.style.display = 'flex';
                }
                alert("删除成功！");
            })
            .catch(function (error) {
                alert(error);
            })
    });
    delBtn.addEventListener('mouseover', () => {
        delBtn.style.color = 'red';
    });
    delBtn.addEventListener('mouseout', () => {
        delBtn.style.color = '#f89898';
    });

    keyElem.addEventListener('mouseover', () => {
        keyElem.style.backgroundColor = '#d1edc4';
        delBtn.style.display = 'inline';
    });
    keyElem.addEventListener('mouseout', () => {
        if (keyElem !== preClickKeyElem) {
            keyElem.style.backgroundColor = '#f0f9eb';
        }
        delBtn.style.display = 'none';
    });
    keyElem.addEventListener("click", function () {
        if (preClickKeyElem === keyElem) {
            return;
        }
        if (!modifyInterceptor("你修改了配置项但还未提交，确定离开吗？")) {
            return;
        }

        if (preClickKeyElem !== undefined) {
            preClickKeyElem.style.backgroundColor = '#f0f9eb';
        }
        preClickKeyElem = keyElem;
        keyElem.style.backgroundColor = '#d1edc4';
        deleteAllConfigurationElem();
        rHeadAndBodyToNoData();

        axios.get('/configuration/' + getKeyElemInnerText(keyElem))
            .then(res => res.data)
            .then(function (res) {
                if (res.code !== 0) {
                    alert(res.msg);
                    return;
                }
                if (res.data.length > 0) {
                    rHeadAndBodyToData();
                    for (let i = 0; i < res.data.length; i++) {
                        newConfigurationElem(res.data[i].topic, res.data[i].path, false);
                    }
                }
            })
            .catch(function (error) {
                alert(error);
            })
    });

    keyElem.innerText = key;
    leftBody.appendChild(keyElem);
    keyElem.appendChild(delBtn);
}

// 因为在KeyElem中添加了一个text值为x的span作为删除按钮，所以在取innerText的时候会将x算进去
function getKeyElemInnerText(ke) {
    if (ke === undefined) {
        return "undefined"
    }
    return ke.innerText.slice(0, -1);
}

function newConfigurationElem(topic, path, focus) {
    const clonedElement = configurationElem.cloneNode(true);
    rightBody.appendChild(clonedElement);

    const topicElem = clonedElement.getElementsByClassName("topic")[0];
    topicElem.value = topic;
    topicElem.addEventListener("focus", function () {
        topicElem.style.border = "1px solid #b88230";
    });
    topicElem.addEventListener("input", function () {
        modifyStatus = true;
    });
    if (focus) {
        topicElem.focus();
    }

    const pathElem = clonedElement.getElementsByClassName("path")[0];
    pathElem.value = path;
    pathElem.addEventListener("focus", function () {
        pathElem.style.border = "1px solid #b88230";
    });
    pathElem.addEventListener("input", function () {
        modifyStatus = true;
    });

    const deleteElem = clonedElement.getElementsByClassName("delete")[0];
    deleteElem.addEventListener("click", function () {
        if (document.querySelectorAll('.configurationElem').length === 1) {
            if (!confirm("这是最后一个配置项，删除后会自动提交，你确定吗？")) {
                return;
            }
            submitAPI(getKeyElemInnerText(preClickKeyElem), "")
                .then(function (res) {
                    if (res) {
                        clonedElement.remove();
                        rHeadAndBodyToNoData();
                    }
                })
            return;
        }
        clonedElement.remove();
        modifyStatus = true;
    });
}

function rHeadAndBodyToNoData() {
    rightHead.style.display = 'none';
    rightBody.style.height = 'calc(100% - 50px)';
    rightBodyDefault.style.display = 'flex';
    document.getElementById("noDataAddButton").style.display = "inline-block";
    document.getElementById("notice").innerText = "无配置项";
    modifyStatus = false;
}

function rHeadAndBodyToData() {
    rightBodyDefault.style.display = 'none';
    rightBody.style.height = 'calc(100% - 100px)';
    rightHead.style.display = 'flex';
}

function rHeadAndBodyToDefault() {
    preClickKeyElem = undefined;
    deleteAllConfigurationElem();
    rightHead.style.display = 'none';
    rightBody.style.height = 'calc(100% - 50px)';
    rightBodyDefault.style.display = 'flex';
    document.getElementById("noDataAddButton").style.display = "none";
    document.getElementById("notice").innerText = "请选择Key";
    modifyStatus = false;
}

function deleteAllConfigurationElem() {
    const configurationElems = document.querySelectorAll('.configurationElem');
    configurationElems.forEach(function (elem) {
        elem.remove();
    });
}

function deleteAll() {
    if (!confirm("清空后会自动提交，你确定吗？")) {
        return;
    }
    submitAPI(getKeyElemInnerText(preClickKeyElem), "")
        .then(function (res) {
            if (res) {
                deleteAllConfigurationElem();
                rHeadAndBodyToNoData();
            }
        })
}

function modifyInterceptor(msg) {
    if (!modifyStatus) {
        return true;
    }
    return confirm(msg);
}

function submit() {
    let ok = true;
    let confArr = [];

    const inputs = document.querySelectorAll('input');
    for (let i = 0; i < inputs.length; i++) {
        let elem = inputs[i];
        elem.value = elem.value.trim();
        if (elem.value.length === 0) {
            elem.style.border = "2px solid red";
            ok = false;
        } else if (i % 2 === 1) {
            confArr.push({ "topic": inputs[i - 1].value, "path": elem.value })
        }
    }

    if (!ok) {
        alert("提交失败，有配置项为空！");
        return false;
    }

    submitAPI(getKeyElemInnerText(preClickKeyElem), JSON.stringify(confArr));
}

function submitAPI(key, data) {
    return new Promise(function (resolve, reject) {
        axios.put('/configuration/' + key, {
            data: data,
        })
            .then(function (res) {
                if (res.data.code !== 0) {
                    alert(res.data.msg);
                    resolve(false);
                } else {
                    alert("提交成功！");
                    modifyStatus = false;
                    resolve(true);
                }
            })
            .catch(function (error) {
                alert(error);
                reject(error);
            });
    });
}

document.getElementById("noDataAddButton").addEventListener("click", function () {
    rHeadAndBodyToData();
    newConfigurationElem("", "", true);
});
document.getElementById("addButton").addEventListener("click", function () {
    newConfigurationElem("", "", true);
});
document.getElementById("deleteAllButton").addEventListener("click", deleteAll);
document.getElementById("submitButton").addEventListener("click", submit);

configurationElem.remove();

axios.get('/keys')
    .then(res => res.data)
    .then(function (res) {
        if (res.code !== 0) {
            alert(res.msg);
            return;
        }
        if (res.data.length > 0) {
            leftBodyDefault.style.display = 'none';
            for (let i = 0; i < res.data.length; i++) {
                newKeyElem(res.data[i].key);
            }
        }
    })
    .catch(function (error) {
        alert(error);
    })
