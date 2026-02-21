// In Go, I get a list of files and sort them by their timestamp, so latest comes first.
// However, the dates are saved in a Map inside the struct passed to the front-end.
// This allows me to "link" an object to another, but Maps in Go cannot be sorted easily.
// The simples solution seemed to sort everything on the front end: that's what this script does.

const list = document.getElementById("list")

// Matches all headers with four digits
const years = Array.from(document.querySelectorAll('[id^="year-"]'))

// Sorts from higher (latest) to lower (earliest)
years.sort((a, b) => {
  if (a.id > b.id) {
    return -1
  }

  if (a.id < b.id) {
    return 1
  }

  return 0
})

// Readds each div, so it comes on top
years.forEach((y) => {
  list.appendChild(y)
})

