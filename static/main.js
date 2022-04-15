const baseURL = location.protocol + "//" + location.hostname + ":" + location.port

const getBlockURL = baseURL + "/getblocks/?ht=0";
const createWalletURL = baseURL + "/createwallet";
const getAddressURL = baseURL + "/getaddress?pub=";
const getBalanceURL = baseURL + "/getbalance?addr=";
const sendURL = baseURL + "/send";

const account = document.getElementById("account");
const def = document.getElementById("default");

const blocks = document.getElementById("blocks");

const showAlert = (msg, type) => {
    const alert = document.createElement("div");
    alert.className = `alert alert-${type} alert-dismissible fade show`;
    alert.setAttribute("role", "alert");

    const span = document.createElement("span");
    span.textContent = msg;

    const btn = document.createElement("button");
    btn.type = "button";
    btn.className = "btn-close";
    btn.setAttribute("data-bs-dismiss", "alert");
    btn.ariaLabel = "Close";

    alert.appendChild(span);
    alert.appendChild(btn);

    document.querySelector("body").prepend(alert);
}

const clearAllAlerts = () => {
    alerts = document.getElementsByClassName("alert");
    for (let i = 0; i < alerts.length; i ++) {
        alerts[i].remove();
    }
}

const parseUnixTime = (time) => {
    const date = new Date(time * 1000);
    const hours = date.getHours();
    const minutes = "0" + date.getMinutes();
    const seconds = "0" + date.getSeconds();

    const fmt = `${hours}:${minutes.substring(1, 4)}:${seconds.substring(1, 4)}`;
    return fmt;
};

const getTXs = (txs, div) => {
    let i = 1;
    txs.forEach((tx) => {
        const h = document.createElement("h2");
        h.textContent = `TX: ${i}`;
        i++;

        div.appendChild(h);

        const idiv = document.createElement("div");
        const odiv = document.createElement("div");

        const ih = document.createElement("h4");
        const oh = document.createElement("h4");

        ih.textContent = "Inputs";
        oh.textContent = "Outputs";

        idiv.appendChild(ih);
        odiv.appendChild(oh);

        const inputs = tx["inputs"];
        const outputs = tx["outputs"];

        inputs.forEach((i) => {
            const txid = document.createElement("li");
            const out = document.createElement("li");
            const pubkey = document.createElement("li");
            const sig = document.createElement("li");

            txid.textContent = `TXID: ${i["txid"]}`;
            out.textContent = `Out: ${i["out"]}`;
            pubkey.textContent = `Public Key: ${i["public_key"]}`;
            sig.textContent = `Signature: ${i["signature"]}`;

            const ul = document.createElement("ul");
            ul.appendChild(txid);
            ul.appendChild(out);
            ul.appendChild(pubkey);
            ul.appendChild(sig);

            idiv.appendChild(ul);
        });

        outputs.forEach((o) => {
            const value = document.createElement("li");
            const pubkeyHash = document.createElement("li");

            value.textContent = `Value: ${o["value"]}`;
            pubkeyHash.textContent = `Public Key Hash: ${o["public_key_hash"]}`;

            const ul = document.createElement("ul");
            ul.appendChild(value);
            ul.appendChild(pubkeyHash);

            odiv.appendChild(ul);
        });

        div.appendChild(idiv);
        div.appendChild(odiv);
    });
};

const fill = (b, div) => {
    nonce = document.createElement("div");
    nonce.textContent = `nonce: ${b["nonce"]}`;

    hash = document.createElement("div");
    hash.textContent = `hash: ${b["hash"]}`;

    prevBlockHash = document.createElement("div");
    prevBlockHash.textContent = `prev block's hash: ${b["prev_block_hash"]}`;

    div.appendChild(nonce);
    div.appendChild(hash);
    div.appendChild(prevBlockHash);

    getTXs(b["transactions"], div);
};

const cookBlocks = (b) => {
    parseUnixTime(b["timestamp"]);

    const item = document.createElement("div");
    item.className = "accordion-item";

    const head = document.createElement("h2");
    head.className = "accordion-header";
    head.id = `heading-${b["nonce"]}`;

    const btn = document.createElement("button");
    btn.className = "accordion-button collapsed";
    btn.type = "button";
    btn.ariaExpanded = "false";
    btn.setAttribute("data-bs-toggle", "collapse");
    btn.setAttribute("data-bs-target", `#collapse-${b["nonce"]}`);
    btn.setAttribute("aria-controls", `collapse-${b["nonce"]}`);

    btn.textContent = parseUnixTime(b["timestamp"]);

    head.appendChild(btn);
    item.appendChild(head);

    const body = document.createElement("div");
    body.id = `collapse-${b["nonce"]}`;
    body.className = "accordion-collapse collapse";
    body.ariaLabel = `heading-${b["nonce"]}`;
    body.setAttribute("data-bs-parent", "#blocks");

    const content = document.createElement("div");
    content.className = "accordion-body";

    fill(b, content);

    body.appendChild(content);
    item.appendChild(body);

    blocks.appendChild(item);
};

const createBtn = document.getElementById("create");
const loadBtn = document.getElementById("load");
const inputFile = document.getElementById("wallet-file");

const download = (pub, priv) => {
    const text = JSON.stringify({
        "public_key": pub,
        "private_key": priv
    })

    const element = document.createElement("a");
    element.style.display = "none";
    element.setAttribute(
        "href",
        "data:text/plain;charset=utf-8," + encodeURIComponent(text)
    );
    element.setAttribute("download", "wallet.json");
    document.body.append(element);

    element.click();
    document.body.removeChild(element);
}

const save = (pub, priv) => {
    localStorage.setItem("pub", pub);
    localStorage.setItem("priv", priv);
}

const get = () => {
    pub = localStorage.getItem("pub");
    priv = localStorage.getItem("priv");
    return { pub, priv };
}

createBtn.addEventListener("click", () => {
    fetch(createWalletURL, {
        method: "GET",
    })
        .then(resp => {
            if (!resp.ok) {
                return;
            }

            return resp.json();
        })
        .then(data => {
            if (data === undefined) {
                return;
            }

            if (!data.successful) {
                console.log(data["error"]);
                return;
            }

            pub = data["public_key"];
            priv = data["private_key"];
            download(pub, priv);
        })
        .catch(err => console.log(err));
});

loadBtn.addEventListener("click", () => {
    let pub, priv;

    if (inputFile.files.length === 0) {
        showAlert("no file selected", "info");
        return;
    }

    const file = inputFile.files[0];
    const reader = new FileReader()
    reader.onload = () => {
        try {
            const wallet = JSON.parse(reader.result)
            pub = wallet["public_key"];
            priv = wallet["private_key"];
        } catch(e) {
            success = false;
            showAlert("not a wallet file", "danger");
            return;
        }

        if (pub === undefined || priv === undefined) {
            showAlert("not a wallet file", "danger");
            return;
        }

        save(pub, priv);
        inputFile.value = null;
        loadAccountPage(pub);
    }

    reader.readAsText(file)
});

const loadDefaultPage = () => {
    clearAllAlerts();
    account.style.display = "none";
    def.style.display = "block";

    try {
        document.getElementById("login").remove();
    } catch(e) {}

    const { pub, priv } = get();
    if (pub !== null && priv !== null) {
        const nav = document.querySelector("#default nav");

        const btn = document.createElement("button");
        btn.type = "button";
        btn.className = "btn btn-secondary";
        btn.id = "login";
        btn.textContent = "Login"

        nav.appendChild(btn);

        btn.addEventListener("click", () => {
            loadAccountPage(pub);
        });
    }

    fetch(getBlockURL, {
        method: "GET",
    })
        .then((resp) => {
            if (!resp.ok) {
                return;
            }

            return resp.json();
        })
        .then((data) => {
            if (data === undefined) {
                return;
            }

            if (!data.successful) {
                console.log(data["error"]);
                return;
            }

            blocks.innerHTML = null;

            data["blocks"].blocks.forEach((b) => {
                cookBlocks(b);
            });
        })
        .catch((err) => console.log(err));
}

const loadBalance = (addr, input) => {
    fetch(getBalanceURL + addr, {
        method: "GET",
    })
        .then(resp => {
            if (!resp.ok) {
                return;
            }

            return resp.json();
        })
        .then(data => {
            if (data === undefined) {
                return;
            }

            if (!data.successful) {
                console.log(data["error"]);
                return;
            }

            input.value = data["balance"];
        })
        .catch(err => console.log(err))
}

const cookAccount = addr => {
    document.getElementById("input-address").value = addr;
    loadBalance(addr, document.getElementById("input-balance"));
}

const loadAccountPage = pub => {
    clearAllAlerts();
    account.style.display = "block";
    def.style.display = "none";

    fetch(getAddressURL + pub, {
        method: "GET",
    })
        .then(resp => {
            if (!resp.ok) {
                loadDefaultPage();
                return;
            }

            return resp.json();
        })
        .then(data => {
            if (data === undefined) {
                return;
            }

            if (!data.successful) {
                console.log(data["error"]);
                loadDefaultPage();
                return;
            }

            const addr = data["address"];
            cookAccount(addr);
        })
        .catch(err => {
            console.log(err);
            loadDefaultPage();
        });
}

function main() {
    const { pub, priv } = get();
    if (pub === null || priv === null) {
        loadDefaultPage();
        return;
    }

    loadAccountPage(pub);
}

main();

const isNum = (num) => {
    if(/^\d+$/.test(num)) {
        return true;
    }

    return false;
}

const send = (amount, recv) => {
    const { pub, priv } = get();
    if (pub === null || priv === null) {
        showAlert("failed to get public and private keys from local storage", "danger");
        return;
    }

    const body = {
        "recv_addr": recv,
        "amount": parseInt(amount, 10),
        "sender_pub": pub,
        "sender_priv": priv
    }

    const headers = {
        'Content-Type': 'application/json;charset=utf-8'
    }

    const options = {
        method: "POST",
        headers: headers,
        body: JSON.stringify(body)
    }

    fetch(sendURL, options)
        .then(resp => resp.json())
        .then(data => {
            if (data === undefined) {
                return;
            }

            if (!data.successful) {
                err = data["error"];
                if (err === "invalid TX") {
                    showAlert("Invalid TX or UTXO part of this TX is already present in the mempool", "danger");
                    return;
                }

                showAlert(err, "danger");
                return;
            }

            showAlert(data["msg"], "success");
        })
        .catch(err => {
            showAlert(err, "danger");
        })
}

(function () {
    'use strict'

    // fetch all the forms we want to apply custom bootstrap validation styles to
    const form = document.getElementById("send");

    // loop over them and prevent submission
    form.addEventListener('submit', e => {
        e.preventDefault()

        if (!form.checkValidity()) {
            e.stopPropagation()
            form.classList.add('was-validated')
            return;
        }

        const amount = document.getElementById("amount").value;
        const recv = document.getElementById("recv-address").value;

        if (!isNum(amount) || amount <= 0) {
            showAlert("Amount must be a positive integer", "info");
            form.reset();
            return;
        }

        send(amount, recv);

        form.reset();
    }, false)
})()

document.getElementById("logout").addEventListener("click", () => {
    localStorage.clear();
    loadDefaultPage();
});

document.getElementById("home").addEventListener("click", () => {
    loadDefaultPage();
});
