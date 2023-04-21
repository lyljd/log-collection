const leftBody = document.getElementById('leftBody');
const leftBodyDefault = document.getElementById('leftBodyDefault');
const rightHead = document.getElementById('rightHead');
const rightBody = document.getElementById('rightBody');
const configurationElem = document.getElementsByClassName('configurationElem')[0];
let preClickKeyElem;
let modifyStatus = false;

function newKeyElem(key) {
    const keyElem = document.createElement('div');
    keyElem.style.cssText = 'width: 100%; height: 50px; line-height: 50px; text-align: center; border-bottom: 1px solid #C0C4CC; cursor: pointer; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; background-color: #f0f9eb;';
    keyElem.addEventListener('mouseover', () => {
        keyElem.style.backgroundColor = '#e1f3d8';
    });
    keyElem.addEventListener('mouseout', () => {
        if (keyElem !== preClickKeyElem) {
            keyElem.style.backgroundColor = '#f0f9eb';
        }
    });
    keyElem.addEventListener("click", function () {
        if (!modifyInterceptor("你修改了配置项但还未提交，确定离开吗？")) {
            return;
        }

        if (preClickKeyElem !== undefined) {
            preClickKeyElem.style.backgroundColor = '#f0f9eb';
        }
        preClickKeyElem = keyElem;
        keyElem.style.backgroundColor = '#e1f3d8';
        deleteAllConfigurationElem();

        // TODO axios请求选择key的value
        //alert(keyElem.innerText);
    });
    keyElem.innerText = key;
    leftBody.appendChild(keyElem);
}

function newConfigurationElem(topic, path) {
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
    topicElem.focus();

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
            if (!modifyInterceptor("这是最后一个配置项，删除后会自动提交，你确定吗？")) {
                return;
            }
            if (!submitAPI("")) {
                return;
            }
            clonedElement.remove();
            rHeadAndBodyToNoData();
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
    modifyStatus = false;
}

function rHeadAndBodyToData() {
    rightBodyDefault.style.display = 'none';
    rightBody.style.height = 'calc(100% - 100px)';
    rightHead.style.display = 'flex';
}

function deleteAllConfigurationElem() {
    const configurationElems = document.querySelectorAll('.configurationElem');
    configurationElems.forEach(function (elem) {
        elem.remove();
    });
    rHeadAndBodyToNoData();
}

function deleteAll() {
    if (!modifyInterceptor("清空后会自动提交，你确定吗？")) {
        return;
    }
    if (!submitAPI("")) {
        return;
    }
    deleteAllConfigurationElem();
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

    return submitAPI(JSON.stringify(confArr));
}

function submitAPI(data) {
    // TODO axios

    alert("提交成功！");
    return true;
}

document.getElementById("noDataAddButton").addEventListener("click", function () {
    rHeadAndBodyToData();
    newConfigurationElem("", "");
});
document.getElementById("addButton").addEventListener("click", function () {
    newConfigurationElem("", "");
});
document.getElementById("deleteAllButton").addEventListener("click", deleteAll);
document.getElementById("submitButton").addEventListener("click", submit);

configurationElem.remove();

//rHeadAndBodyToNoData();

// leftBodyDefault.style.display = 'none';
// for (let i = 0; i < 10; i++) {
//     newKeyElem(i);
// }

/* TODO
先发送axios请求所有的key，若没有key就啥都不干，若请求失败就仅提示请求失败，然后也啥都不干，
若请求成功且有key则调用newKeyElem(key)，就拿着第一个key去请求value(etcd中查出配置项)，
若没有value就调用rHeadAndBodyToNoData()，若请求失败就仅提示请求失败，然后也啥都不干，
若请求成功且有value，就调用，然后将第一个key置于选中状态，然后把value传入一个统一处理json的函数，
该函数在上面切换key时调用axios后也会用到，这个函数就是将json转化成一个一个的配置项，
会调用newConfigurationElem(topic, path)
*/
