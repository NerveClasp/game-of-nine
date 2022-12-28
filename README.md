# [WIP] Game of Nine

## Idea

"Server" is started on one of the computers in a network and anyone
can open a local url, create a Player, create or join a Game, play it.
Perfect for cases when there's no electricity and/or Internet and you are
a person short for a true card game (Game of Nine requires 3-6 players)

## Issues

These's an issue in FireFox right now, some console errors and WebSocket
does not seem to be working. I'll get to that later, just FYI. Chrome(ium)
works fine, so that's an issue for future Romka))

## Local development

`yarn` is needed (cause npm is slow and pnpm is buggy with SMUI) and `go` (Golang)

```sh
yarn install && cd frontend && yarn install && cd ..
```

```sh
yarn start
```

## Details

This project is using Golang and SvelteKit. SvelteKit is being built
to `static` folder and used by `main.go` to serve the frontend. Go
is used as a backend and a WebSocket server. Everything is using `7331` port
by default.

## Development

You can start Golang part via either

```sh
go run .
```

or

```sh
yarn go
```

- that's a backend part and the game logic. Golang simply serves what it finds
  in `static` folder (you may want to create this folder manually if something goes
  wrong, I did not test it, sorry)

In parallel to that you might want to run

```sh
yarn watch
```

- that will do `yarn build` in `frontend` folder and `nodemon` will watch any changes
  to any `svelte`, `ts`, or `css` files (check `./nodemon.json` for details). Manual
  `Ctrl/Cmd+R` is required to see changes. If you do any changes to `.go` files -
  restart go backend manually.
