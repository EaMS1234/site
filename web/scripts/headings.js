// Creates a link for each heading across the whole page

document.querySelectorAll("h1, h2, h3, h4, h5, h6").forEach((heading) => {
  // Ignores the heading which contains the title of the page
  if (heading.id != "content-title") {

    // On enter, create and display the link on the right of the heading
    heading.addEventListener("mouseenter", () => {
      let link = document.createElement("a")

      link.innerHTML = " #"
      link.setAttribute("href", `#${heading.id}`)
      link.setAttribute("title", `${document.getElementById("content-title").innerHTML}#${heading.id}`)

      heading.appendChild(link)
    })

    // Remove the link once the heading is not on focus
    heading.addEventListener("mouseleave", () => {
      // Checks if the node is an anchor, and then removes it
      heading.removeChild(heading.lastChild.nodeName === "A" ? heading.lastChild : null)
    })
  }
})
