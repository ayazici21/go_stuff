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
}

function createCardFromTask(task: string) {
  let div = document.createElement("div")
  div.setAttribute("id", task["_id"])
  let card = document.createElement("div")
  card.classList.add("uk-card", "uk-card-default", "uk-card-hover", "uk-padding-small")
  //card.setAttribute("uk-grid", "")

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
  body.appendChild(status)

  let completeTaskButton: HTMLButtonElement = document.createElement("button")

  completeTaskButton.setAttribute("uk-icon", "check")
  completeTaskButton.setAttribute("ratio", "2")
  completeTaskButton.setAttribute("id", task["_id"])
  completeTaskButton.onclick = (event) => completeTaskEvent(div)

  body.appendChild(completeTaskButton)
  
  let deleteTaskButton = document.createElement("button")
  deleteTaskButton.setAttribute("uk-icon", "trash")
  deleteTaskButton.setAttribute("ratio", "2")
  deleteTaskButton.setAttribute("id", task["_id"])
  deleteTaskButton.onclick = deleteTaskEvent

  body.appendChild(deleteTaskButton)

  card.appendChild(title)
  card.appendChild(body)
  div.appendChild(card)
  cards.appendChild(div)

}


async function completeTaskEvent(caller: HTMLDivElement) {
    let task_id: string = caller.getAttribute("id")!

    await fetch("http://127.0.0.1:3131/task/" + task_id, {method: "PUT"})
    .then(response => {
        if (response.ok) {
            console.log("OK")
        } else {
            console.log(response)
        }
    })
    return
}


function deleteTaskEvent(event: MouseEvent) {
    return
}