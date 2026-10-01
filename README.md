# keltas
## Running

`make run` builds the app, generates a self-signed `cert.pem`/`key.pem` (once) and
serves HTTPS only on `localhost:8080` (override with `NERINGA_LISTEN_ADDR`).

To restrict access, create `allow.txt` with one IP or CIDR per line (`#` comments
allowed). Connections from other addresses are dropped. Without the file, all
addresses are allowed.
