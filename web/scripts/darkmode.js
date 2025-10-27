function toggle() {
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches)
        document.getElementById("style").setAttribute("href", "styles/dark.css")
    else
        document.getElementById("style").setAttribute("href", "styles/dark.css")
}

toggle()
