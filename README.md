# Gator RSS Aggregator

Welcome to **Gator**, a command‑line RSS feed aggregator written in Go. This tool lets you subscribe to RSS feeds, follow/unfollow them, and browse recent posts pulled from the web. It stores data in a PostgreSQL database and helps you keep track of fresh content.

---

## Prerequisites

To run this project, you will need:

- **Go** (1.20+ recommended)
- **PostgreSQL** (running locally or accessible remotely)
- **Git** (to clone and track changes)

Verify your installations:
```sh
go version
psql --version
git --version
```

---

## Installation

Install the `gator` CLI using:
```sh
go install ./...
```

This builds the `gator` binary into your `$GOPATH/bin`. Make sure that directory is on your `$PATH`.

For example, you can add this to your shell configuration:
```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

Then verify installation:
```sh
gator --help
```

---

## Database Setup

Make sure PostgreSQL is running, then create a database:
```sh
createdb gator
```

Your project should include SQL migrations or schema files that create the required tables.
Apply them manually or using a migration tool.

---

## Configuration

Gator expects a configuration file that provides database connection information.

Create a file at:
```
~/.gatorconfig.json
```

Example:
```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable"
}
```

Replace with your actual database credentials.

---

## Usage

Run the CLI by calling:
```sh
blogag <command>
```

Example:
```sh
blogag register myuser
```

---

## Available Commands

Below are the currently implemented commands:

### Authentication
```
login       - Log in as a registered user
register    - Register a new user
reset       - Reset the database (dangerous!)
```

### Users & Feeds
```
users       - List all registered users
feeds       - List all feeds in the system
addfeed     - Add a new RSS feed for the logged‑in user
agg         - Fetch posts from feeds
```

### Following
```
follow      - Follow a feed
following   - Show feeds you're following
unfollow    - Unfollow a feed
```

### Browsing Content
```
browse      - Display recent aggregated posts for your user
```

---

## Typical Workflow

```sh
blogag register alice
blogag login alice
blogag addfeed https://example.com/rss.xml
blogag follow 1234-...-feed-id
blogag agg
blogag browse
```

---

## Git Tracking

This project should be tracked with Git. Commit early, commit often:
```sh
git init
git add .
git commit -m "initial commit"
```

---

## Troubleshooting
- Ensure PostgreSQL is running
- Check that your database URL is valid
- Ensure `$GOPATH/bin` is on `$PATH`
- Run `go mod tidy` if dependencies are missing

---

## License
This project is provided as‑is. Improve freely!

Happy hacking!

