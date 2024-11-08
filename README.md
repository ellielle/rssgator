# RSS Gator

This project aggregates RSS feeds provided by the user. Once a user registers, they can log in and:

- [x] Add feeds to the aggregator
- [x] Subscribe/Unsubscribe from feeds
- [x] View which feeds they're subscribed to
- [x] Can be run in the background and will refresh feeds on interval

## Installation

Clone this repo:

```bash
git clone git@github.com:ellielle/rssgator
```

This project uses a `.gatorconfig` file that is created when run for the first time.
The db_url field needs to be populated with your postgres connection string, which should look something like:

```bash
postgres://username:password@host:port/dbname
```

## Goose

[Goose](https://github.com/pressly/goose) is used to handle migrations, and [sqlc](https://docs.sqlc.dev/en/latest/index.html) is used to generate Go code.

```bash
goose postgres postgres://username:password@localhost:5432/dbname up
```

## sqlc

After migrating, bavigate to the `sql/schema` directory and run `sqlc generate` to generate Go code to interact with the database.

# Usage

```bash
go run . <command> [args...]
```

RSSGator accepts the following commands:

- `login [username]`
- `register [username]`
- `users` - lists all users and indicates who is logged in currently
- `feeds` - lists all feeds available to subscribe to

The following commands require you to be logged in:

- `addfeed [name] [url]` - adds a feed to the feed list, and follows it
- `follow [url]` - adds a feed to your follow list
- `following` - shows all feeds the user is following
- `unfollow [url]` - removes a feed from your follow list
- `agg` - starts the server to fetch feeds on a given interval
- `browse` - browse aggregated posts
