# go_movies_api_rest
Build, manage &amp; deploy a RESTful JSON API (a movies web application)in GO.


## Database Setup and Configuration
You'll need to install PostgreSQL.
```sh
brew install postgresql@15  # mac

sudo apt install postgresql # ubuntu
```

```sh
psql --version
```

### Connecting to the PostgreSQL Interactive Terminal
```sh
sudo -u postgres psql

# now the prompt should become like: postgres=#
```
Check out the current user
```sql
SELECT current_user;

-- you can also use meta command to check the users:
\du
-- the following lists all the databases
\l
-- list the tables
\dt
```

Create the project databae & connect to it.
```sql
CREATE DATABASE moviesdb;

\c moviesdb -- connect to the database created
/*
Now the prompt should look like: moviesdb=#
*/

CREATE ROLE moviesdb_user WITH LOGIN PASSWORD 'Ch@ng3M3!';

/*
citext extension adds a case-insensitive character string to PostgreSQL.
We'll use this to store user email adresses.
*/
CREATE EXTENSION IF NOT EXISTS citext;

exit
```

Connect as the new user
```sh
psql --host=localhost --dbname=moviesdb --username=moviesdb_user
# You'll be prompted to enter the password for the user.
```

```sql
-- Now you should be connected to the `moviesdb` & the prompt looks like: moviesdb=>
SELECT current_user;

exit
```

### Connecting to PostgreSQL from Go
We'll be needing a **database driver** to act as the *middleman* between Go and the database itself. We'll be using the `pq` package.

```sh
go get github.com/lib/pq@v1
```

We also want to decouple the connection string, aka DSN, from our Go code.
So, in our go application we'll pass the DSN like so:
```go
flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("MOVIESDB_DSN"), "PostgreSQL DSN")
```
And to store the DSN as an environment variable on the local machine, open the `.profile` file in your HOME directory, and add the DSN like so

```md
export MOVIESDB_DSN='postgres://moviesdb_user:Ch@ng3M3!@localhost/moviesdb'
```

For our environmental variable to be recognized in the terminal either reboot your system or in the terminal run `source $HOME/.profile`.
```sh
# Make sure the variable is recognized in your terminal.
echo MOVIESDB_DSN
```

As a bonus you can connect to the database via the DSN saved as the environmental variable
```sh
psql $MOVIESDB_DSN
```

## SQL Migrations
To manage SQL migrations we'll be using the [migrate](https://github.com/golang-migrate/migrate/releases) command-line tool (which itself is written in Go).
To install:
```sh
brew install golang-migrate # macOS

# On Linux, simply download the pre-built binary and move it to a location on your system path.
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.2/migrate.linux-amd64.tar.gz | tar xvz -C /tmp
mv /tmp/migrate ~/go/bin/
## you may need to add to the path
export PATH=$HOME/go/bin:$PATH

# make sure it's working
migrate -version
```

Having installed the `migrate` tool, we can generate our **migration files**.
```sh
migrate create -seq -ext=.sql -dir=./migrations create_movies_table
# -seq flag indicates to use sequential numbering like 0001, 0002 ...
# now we should have two files in the `migrations` directory:
# 001_create_movies_table.up.sql
# 001_create_movies_table.down.sql
```

