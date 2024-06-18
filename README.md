# Muse SCSS Analyser
![Screenshot 2024-06-18 at 12 13 03](https://github.com/IndependentIP/muse-scss-analyser/assets/83584112/96266861-3ea5-4eaf-a741-929073ace965)
![Screenshot 2024-06-18 at 12 17 25](https://github.com/IndependentIP/muse-scss-analyser/assets/83584112/4cae0989-4a84-4d59-bd3c-ea449a74ad2f)

The dream is to get rid of node-sass. For that we need to know what is being used and where, inside Muse.
This tool finds all the imports of SCSS files in Muse. Counts them, and shows the relationships between files in the visualisation tool.

## Prerequisites
- Golang 1.16
- The Muse project (it must be in your home directory) find it in the [Muse repository](https://github.com/IndependentIP/muse)

## How to run the analyser
1. Clone this repository

2. Go to the `analyser` directory
```bash
cd analyser
```

3. Run Go mod to install the dependencies
```bash
go mod tidy
```

4. Run the search command. Give the path to the src directory of your Muse project. Remember, the project must be in your home directory.
```bash
go run cmd/main.go search /Users/daniel/muse/front-end/src
```

5. The results are in the `generated` directory

## How to run the visualisation tool
1. copy the  `generated/d3data.json` file to the `webserver/static` directory
```bash
cp generated/d3data.json webserver/static/results.json
```

2. run the web server
```bash
go run cmd/main.go serve
```

3. Open your browser and go to `http://localhost:4321`

You can run the `search` command again to get new results and copy the new file to the `webserver/static` directory. On refresh, you'll see the new results.

## Building a binary
```bash
go build -o muse-analyser cmd/main.go
```

###  Running the binary
For searching
```bash
./muse-analyser search /Users/daniel/muse/front-end/src
```

For serving the visualisation tool
```bash
./muse-analyser serve
```

## About this repository
This repository contains two main parts:
- The Analyser tool, which searches for SCSS imports in the Muse project and shows the relationships between files in a series of JSON files
- tree-sitter-sass, a fork of the tree-sitter parser for SCSS, which is used by the Analyser tool to parse SCSS files, but with a go binding

## Folder structure
- `analyser` - The Go code for the Analyser tool
  - `cmd` - The main commands for the Analyser tool
  - `generated` - The results of the Analyser tool
  - `filesearch` - The code for searching for SCSS files in the Muse project
  - `webserver` - The code for the visualisation tool
- `tree-sitter-sass` - The tree-sitter parser for SCSS with a Go binding
  - `bindings` - Different bindings for the tree-sitter parser (auto-generated with treesitter-cli)
    - `go` - The Go binding















