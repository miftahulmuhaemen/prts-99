# Arknights Wikia Chatbot AI-powered

This project is a web scraper for extracting data from the Arknights Wiki. It uses the Colly library to scrape information about operators, including their skills, talents, promotions, and other attributes.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)

## Installation

To install the necessary dependencies, run:

```sh
go mod download
```


```
.
├── .air.toml
├── .gitignore
├── cache/
├── go.mod
├── go.sum
├── internal/
│   ├── cache/
│   ├── embed_model/
│   ├── model/
│   │   └── character.go
│   ├── scrapper/
│   │   ├── scrapper.go
│   │   └── scrapper_test.go
├── LICENSE
├── main.go
├── README.md
└── tmp/
```