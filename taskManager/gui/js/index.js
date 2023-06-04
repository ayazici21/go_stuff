var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (g && (g = 0, op[0] && (_ = 0)), _) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
var taskIDs = new Set();
var cards;
window.onload = function () {
    cards = document.getElementById("cards");
    viewTasks();
};
function viewTasks() {
    return __awaiter(this, void 0, void 0, function () {
        var tasks;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, fetch("http://127.0.0.1:3131/tasks")
                        .then(function (response) {
                        if (response.ok) {
                            return response.json();
                        }
                        return [];
                    })
                        .then(function (tasks) {
                        tasks.forEach(function (task) {
                            console.log(task);
                            if (task["_id"] in taskIDs) {
                                return;
                            }
                            taskIDs.add(task["_id"]);
                            createCardFromTask(task);
                        });
                    })];
                case 1:
                    tasks = _a.sent();
                    return [4 /*yield*/, fetch("http://127.0.0.1:3131/filter") // i should probably set filter to 1 every time
                            .then(function (response) {
                            if (response.ok) {
                                return response.json();
                            }
                            else {
                                return 1;
                            }
                        }).then(function (data) {
                            [1, 2, 3].forEach(function (num) { var _a; return (_a = document.getElementById("filter" + num)) === null || _a === void 0 ? void 0 : _a.classList.remove("uk-active"); });
                            document.getElementById("filter" + data["filter"]).classList.add("uk-active");
                        })];
                case 2:
                    _a.sent();
                    return [2 /*return*/];
            }
        });
    });
}
function createCardFromTask(task) {
    var div = document.createElement("div");
    div.setAttribute("id", task["_id"]);
    var card = document.createElement("div");
    card.classList.add("uk-card", "uk-card-default", "uk-card-hover", "uk-padding");
    var title = document.createElement("h3");
    title.classList.add("uk-card-title");
    title.appendChild(document.createTextNode(task["title"]));
    var body = document.createElement("div");
    body.classList.add("uk-card-body");
    var desc = document.createElement("p");
    desc.appendChild(document.createTextNode(task["description"]));
    body.appendChild(desc);
    var status = document.createElement("p");
    status.appendChild(document.createTextNode("Status: " + (task["status"] ? "Completed" : "Pending")));
    status.setAttribute("id", "status_" + task["_id"]);
    status.setAttribute("data-status", task["status"]);
    body.appendChild(status);
    var completeTaskButton = document.createElement("button");
    completeTaskButton.setAttribute("uk-icon", "check");
    completeTaskButton.setAttribute("ratio", "2");
    completeTaskButton.onclick = function (event) { return completeTaskEvent(status); };
    body.appendChild(completeTaskButton);
    var deleteTaskButton = document.createElement("button");
    deleteTaskButton.setAttribute("uk-icon", "trash");
    deleteTaskButton.setAttribute("ratio", "2");
    deleteTaskButton.onclick = function (event) { return deleteTaskEvent(div); };
    body.appendChild(deleteTaskButton);
    card.appendChild(title);
    card.appendChild(body);
    div.appendChild(card);
    cards.appendChild(div);
}
function completeTaskEvent(caller) {
    return __awaiter(this, void 0, void 0, function () {
        var task_id;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    task_id = caller.getAttribute("id").substring(7);
                    if (caller.getAttribute("data-status") == "true") {
                        // instead of making complete task button disabled, just do nothing for now
                        console.log("Task already completed");
                        return [2 /*return*/];
                    }
                    return [4 /*yield*/, fetch("http://127.0.0.1:3131/task/" + task_id, { method: "PUT" })
                            .then(function (response) {
                            if (response.ok) {
                                // change the status in the UI
                                caller.textContent = "Status: Completed";
                            }
                            else {
                                // TODO: show error
                            }
                        })];
                case 1:
                    _a.sent();
                    return [2 /*return*/];
            }
        });
    });
}
function deleteTaskEvent(caller) {
    return __awaiter(this, void 0, void 0, function () {
        var task_id;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    task_id = caller.getAttribute("id");
                    return [4 /*yield*/, fetch("http://127.0.0.1:3131/task/" + task_id, { method: "DELETE" })
                            .then(function (response) {
                            if (response.ok) {
                                caller.remove();
                            }
                            else {
                                // TODO: show error
                            }
                        })];
                case 1:
                    _a.sent();
                    return [2 /*return*/];
            }
        });
    });
}
function addTask() {
    return __awaiter(this, void 0, void 0, function () {
        var titleField, descField, title, description;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    titleField = document.getElementById("task_title");
                    descField = document.getElementById("task_desc");
                    title = titleField.value;
                    description = descField.value;
                    if (title.length == 0 || description.length == 0) {
                        // handle this
                        return [2 /*return*/];
                    }
                    console.log(title);
                    console.log(description);
                    console.log(JSON.stringify({
                        title: title,
                        description: description
                    }));
                    return [4 /*yield*/, fetch("http://127.0.0.1:3131/task", {
                            method: "POST",
                            body: JSON.stringify({
                                title: title,
                                description: description
                            }),
                            headers: { "Content-Type": "application/json" }
                        }).then(function (response) {
                            if (response.ok) {
                                return response.json();
                            }
                            else {
                                // handle error
                                return null;
                            }
                        }).then(function (data) {
                            if (data == null) {
                                return;
                            }
                            taskIDs.add(data["_id"]);
                            createCardFromTask(data);
                        })];
                case 1:
                    _a.sent();
                    return [2 /*return*/];
            }
        });
    });
}
function switchFilter(filter) {
    return __awaiter(this, void 0, void 0, function () {
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, fetch("http://127.0.0.1:3131/filter/" + filter, { method: "PUT" })
                        .then(function (response) {
                        if (response.ok) {
                            cards.replaceChildren();
                            viewTasks();
                        }
                        else {
                            // TODO: show error
                        }
                    })];
                case 1:
                    _a.sent();
                    return [2 /*return*/];
            }
        });
    });
}
