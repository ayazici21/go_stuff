var taskIDs = new Set<String>()
var cards: HTMLElement;

window.onload = () => {
    cards = document.getElementById("cards")!;
    viewTasks()
}

async function viewTasks() {
    let tasks = await fetch("http://127.0.0.1:3131/tasks")
        .then(response => {
            if (response.ok) {
                return response.json()
            }
            return []
        })
        .then(tasks => {
            tasks.forEach(task => {
                console.log(task)
                if (task["_id"] in taskIDs) {
                    return
                }
                taskIDs.add(task["_id"])
                createCardFromTask(task)
            });
        })

        await fetch("http://127.0.0.1:3131/filter") // i should probably set filter to 1 every time
            .then(response => {                     // the app is restarted, but not now.
                if (response.ok) {
                    return response.json()
                } else {
                    return 1
                }
            }).then(data => {
                [1,2,3].forEach(num => document.getElementById("filter"+num)?.classList.remove("uk-active"));
                document.getElementById("filter" + data["filter"])!.classList.add("uk-active")
            })
}

function createCardFromTask(task: string) {
    let div = document.createElement("div")
    div.setAttribute("id", task["_id"])

    let card = document.createElement("div")
    card.classList.add("uk-card", "uk-card-default", "uk-card-hover", "uk-padding")

    let title = document.createElement("h3")
    title.classList.add("uk-card-title")
    title.appendChild(document.createTextNode(task["title"]))

    let body = document.createElement("div")
    body.classList.add("uk-card-body")

    let desc = document.createElement("p")
    desc.appendChild(document.createTextNode(task["description"]))
    body.appendChild(desc)

    let status = document.createElement("p")
    status.appendChild(document.createTextNode("Status: " + (task["status"] ? "Completed" : "Pending")))
    status.setAttribute("id", "status_" + task["_id"])
    status.setAttribute("data-status", task["status"])
    body.appendChild(status)

    let completeTaskButton: HTMLButtonElement = document.createElement("button")

    completeTaskButton.setAttribute("uk-icon", "check")
    completeTaskButton.setAttribute("ratio", "2")
    completeTaskButton.onclick = (event) => completeTaskEvent(status)

    body.appendChild(completeTaskButton)

    let deleteTaskButton = document.createElement("button")
    deleteTaskButton.setAttribute("uk-icon", "trash")
    deleteTaskButton.setAttribute("ratio", "2")
    deleteTaskButton.onclick = (event) => deleteTaskEvent(div)

    body.appendChild(deleteTaskButton)

    card.appendChild(title)
    card.appendChild(body)
    div.appendChild(card)
    cards.appendChild(div)
}


async function completeTaskEvent(caller: HTMLParagraphElement) {

    let task_id: string = caller.getAttribute("id")!.substring(7)


    if (caller.getAttribute("data-status") == "true") {
        // instead of making complete task button disabled, just do nothing for now
        console.log("Task already completed")
        return
    }


    await fetch("http://127.0.0.1:3131/task/" + task_id, { method: "PUT" })
        .then(response => {
            if (response.ok) {
                // change the status in the UI
                caller.textContent = "Status: Completed"
            } else {
                // TODO: show error
            }
        })


}


async function deleteTaskEvent(caller: HTMLDivElement) {
    let task_id: string = caller.getAttribute("id")!

    await fetch("http://127.0.0.1:3131/task/" + task_id, { method: "DELETE" })
        .then(response => {
            if (response.ok) {
                caller.remove()
            } else {
                // TODO: show error
            }
        })
}

async function addTask() {
    let titleField = document.getElementById("task_title")! as HTMLInputElement
    let descField = document.getElementById("task_desc")! as HTMLInputElement
    let title: string = titleField.value!
    let description: string = descField.value!

    if(title.length == 0 || description.length == 0) {
        // handle this
        return
    }
    console.log(title)
    console.log(description)
    console.log(JSON.stringify({
        title: title,
        description: description
    }))
    await fetch("http://127.0.0.1:3131/task", {
        method: "POST",
        body: JSON.stringify({
            title: title,
            description: description
        }),
        headers: {"Content-Type": "application/json"}
    }).then(response => {
        if (response.ok) {
            return response.json()
        } else {
            // handle error
            return null
        }
    }).then(data => {
        if (data == null) {
            return
        }
        taskIDs.add(data["_id"])
        createCardFromTask(data)
    })

    

}

async function switchFilter(filter: number) {
    await fetch("http://127.0.0.1:3131/filter/" + filter, { method: "PUT" })
        .then(response => {
            if (response.ok) {
                cards.replaceChildren()
                viewTasks()
            } else {
                // TODO: show error
            }
        })
}