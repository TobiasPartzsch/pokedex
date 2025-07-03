# Pokedex CLI

A command-line interface for exploring and catching Pokemon using the [PokeAPI](https://pokeapi.co/). Features basic navigation, caching, and a simple REPL for user commands.

## Features

- Browse Pokemon locations (`map`, `mapb`)
- Explore areas for wild Pokemon (`explore <location>`)
- Attempt to catch Pokemon (`catch <pokemon>`)
- Inspect caught Pokemon (`inspect <pokemon>`)
- View your Pokedex (`pokedex`)
- Caching to reduce redundant API calls

## Getting Started

1. Clone this repository:
```
git clone https://github.com/TobiasPartzsch/pokedex
```

2. Build and run:
```
go build -o pokedex
./pokedex
```


## Usage

Type any of the following commands at the prompt:

- `help` – Show available commands
- `map` – Display the next 20 area locations
- `mapb` – Display the previous 20 area locations
- `explore <area>` – List Pokemon in a specific area
- `catch <pokemon>` – Try to catch a Pokemon by name
- `inspect <pokemon>` – Show detailed info on a caught Pokemon
- `pokedex` – List all your caught Pokemon
- `exit` – Quit the program

## Lessons Learned

- **Helper Functions:** Building reusable helpers for common actions led to cleaner, easier-to-test code and kept the main CLI loop focused on core logic.
- **File Organization:** I found it best to split out files (e.g., API client, caching, domain models) only when code complexity justified it, avoiding needless fragmentation and keeping things simple for as long as possible. Using an IDE certainly helped with not feeling the need to split.

## Contributing

Contributions and suggestions welcome! Please open an issue or pull request.
There will also be a thread on the Boots.dev Discord.

## License

MIT