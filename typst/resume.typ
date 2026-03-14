#import "@preview/modern-cv:0.9.0": *

#show: resume.with(
  author: (
    firstname: "Yaroslav",
    lastname: "Antipov",
    email: "yarik.antipov2008@gmail.com",
    homepage: "https://wasdetchan.online",
    phone: "+79303535755",
    github: "WASDetchan",
    // twitter: "",
    // scholar: "",
    // orcid: "",
    // linkedin: "TODO",
    address: "Mersin, Turkey",
    positions: (
      "Software Engineer",
      "Developer",
    ),
    // custom: (
    //   (
    //     text: "Codeforces: WASDetchan",
    //     link: "https://codeforces.com/profile/WASDetchan"
    //   )
    // )
    custom: (
      (
        text: "Codeforces: WASDetchan",
        link: "https://codeforces.com/profile/WASDetchan",
        icon: "globe",
      ),
    ),
  ),
  keywords: ("Engineer", "Developer"),
  description: "resume",
  // profile-picture: image("profile.png"),
  profile-picture: none,

  date: datetime.today().display(),
  language: "en",
  colored-headers: true,
  show-footer: false,
  show-address-icon: true,
  paper-size: "us-letter",
)

= Experience

#resume-entry(
  title: "Software Engineer",
  location: "Remote",
  date: "2024 - 2025",
  description: "Nutrilinker, Dr. Elena Belyaeva",
  // title-link: "https://github.com/DeveloperPaul123",
)

#resume-item[
  - Full-stack development of Nutrilinker - a platform for meal planning with advanced doctor-to-client connection.
  - First iteration:
    - PHP Laravel backend
    - PostgreSQL database
    - React.js frontend using Laravel Inertia
  - Second iteration:
    - Rust Actix backend
    - PostgreSQL/SQLite database using Diesel ORM
    - React.js frontend
    - Available as both web and desktop app
  - Deployment on a server running Ubuntu with Docker, Nginx for HTTPS
]

#resume-entry(
  title: "Embedded Software Engineer",
  location: "Ivanovo, Russia",
  date: "2021 - 2022",
  description: "Scientist group in Ivanovo",
)

#resume-item[
  - Embedded software development and rapid hardware prototyping.
  - Tech stack:
    - Arduino Uno (for prototyping), ESP32 microcontrollers
    - Multiple devices: sensors, displays, etc.
    - C++ for embedded programming
  - Magaged complex asyncronous measurement workflows with feedback loops
  - Automated analyzis of the data, both on the microcontroller for quick result and on PC for better presizion and plotting
  - Made changes in the hardware on request for quicker development cycles, eleminating unneccesary meetings
]

= Projects

#resume-entry(
  title: "Bight - TUI spreadsheet",
  location: [#github-link("WASDetchan/bight.nvim") #github-link("WASDetchan/bight") ],
  date: "September 2025 - Present",
)

#resume-item[
  - Designed and implemented a spreadsheet engine in Rust
    - Generic Table API to allow using bight as a general-purpose table engine
    - Asynchronous parallel Lua formula evaluation using Tokio and mlua with dependency tracking for preventing deadlocks
    - Multiple file format export and import support
    - Simple, but powerfull plotting system, that uses plotters-rs to plot data from the tables
    - Extensive public API documentation
  - Implemented bight.nvim - a neovim plugin for bight
    - Written in Rust using nvim-oxi nvim API bindings
    - Powerful Lua API bindings for most of bight's features
    - Vim commands for quick preview plotting
]

#resume-entry(
  title: "Personal website",
  location: github-link("WASDetchan/wasdetchan-online"),
  date: "2026 - Present",
)

#resume-item[
  Tech stack:
  - Website:
    - Go + teml 
    - HTML, CSS
    - No database 
  - Deployment  
    - Ubuntu VPS
    - Docker with docker compose
    - Caddy for HTTPS
    - Xray for advanced routing and proxying
]

#resume-entry(
  title: "Meteor Crusher game",
  location: github-link("WASDetchan/space_shuttle"),
  date: "2020 - 2021",
)

#resume-item[
  - Custom game engine built up from pygame (SLD2 bindings for python)
  - Heavy usage of OOP patterns for easier development
]

#resume-entry(
  title: "Contributions to other open source projects",
)

#let github-link-left(github-path) = {
  set box(height: 11pt)

  align(left + horizon)[
    #fa-icon("github", fill: color-darkgray) #link(
      "https://github.com/" + github-path,
      github-path,
    )
  ]
}
#resume-item[
  - #github-link-left("XTLS/Xray-core")
  - #github-link-left("mlua-rs/mlua")
  - #github-link-left("matheus-git/logcast")
]
= Skills

#resume-skill-item(
  "Programming Languages",
  (
    strong("Rust"),
    strong("Go"),
    strong("C/C++"),
    "JavaScript/TypeScript",
    "Python",
    "PHP",
  ),
)
#resume-skill-item(
  "Related skills",
  (
    strong("Git"),
    strong("GitHub"),
    strong("Linux (server and desktop)"),
    strong("Docker"),
    strong("Competitive programming (Expert on Codeforces)"),
    "XTLS Xray",
    "Arduino/ESP32",
    "Vulkan",
  ),
)
#resume-skill-item("Spoken Languages", (strong("English (B2)"), "Russian", "Italian (A1)"))
#resume-skill-item(
  "Programs",
  (
    strong("Excel"),
    strong("Neovim"),
    "Word",
    "Powerpoint",
    "Visual Studio Code",
    "IntelliJ IDEA",
  ),
)

#resume-skill-item(
  "Other skills",
  (
    "Theoretical physics and data analyzis for physics",
    "Typst",
    "LaTeX",
    "Theoretical econommics",
  ),
)
// spacing fix, not needed if you use `resume-skill-grid`
#block(below: 0.65em)

// An alternative way of list out your resume skills
// #resume-skill-grid(
//   categories_with_values: (
//     "Programming Languages": (
//       strong("C++"),
//       strong("Python"),
//       "Rust",
//       "Java",
//       "C#",
//       "JavaScript",
//       "TypeScript",
//     ),
//     "Spoken Languages": (
//       strong("English"),
//       "Spanish",
//       "Greek",
//     ),
//     "Programs": (
//       strong("Excel"),
//       "Word",
//       "Powerpoint",
//       "Visual Studio",
//       "git",
//       "Zed"
//     ),
//     "Really Really Long Long Long Category": (
//       "Thing 1",
//       "Thing 2",
//       "Thing 3"
//     )
//   ),
// )

= Education

#resume-entry(
  title: "Phystech Lyceum named after P. L. Kapitsa",
  location: "Dolgoprudny, Russia",
  date: "2022 - 2026",
  description: "Physics and math profile",
)

#resume-item[
  - Advanced physics and math
  - Competitive programming
]
