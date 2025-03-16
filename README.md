# README

## DBMX

A modern and productivity focused database management tool. It is under active development and is NOT ready for production use.

## Early Preview
#### Sidebar and Tabs
<img width="1440" alt="Screenshot 2025-03-16 at 9 53 23 AM" src="https://github.com/user-attachments/assets/a5b962ca-3330-4c81-a4d4-550ebcd6025e" />

#### Connect to a Database Server
<img width="1440" alt="Screenshot 2025-03-16 at 9 53 34 AM" src="https://github.com/user-attachments/assets/c221522d-3b2f-40a6-ab5e-9f496dd7bc31" />

#### Resizable Editor View and Autocomplete Suggestions
<img width="1440" alt="Screenshot 2025-03-16 at 9 54 40 AM" src="https://github.com/user-attachments/assets/a7406e8a-a76a-4648-bc9b-4777fd79d6fb" />

## Tech Stack
- Wails Go
- Svelte 5
- Shadcn UI
- SQLite3db

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.
