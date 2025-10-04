# Notes on testing scenarios

## Basic Stress-testing

Enter PostgreSQL shell.

```sh
docker exec -it timeseries-timescaledb-1 bash
```

Initialize database with scale factor `625` with table name `example`.
For the record Each scale factor creates approximately 16MB, so scale factor of `625` is roughly 10GB.

```sh
pgbench -i -s 625 -U postgres example
```

Run stress test with `10` concurrent clients. Use `2` threads for executing transactions and run
`10_000` transactions per client.

```sh
pgbench -c 10 -j 2 -t 10000 example
```

