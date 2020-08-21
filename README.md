Query cancellation problem
==========================

pgx/v4 does not work as expected with pgbouncer due to
query cancellation ([issue #679](https://github.com/jackc/pgx/issues/679)). In order to replicate:

1. `docker-compose up`
2. `echo "show pools;" | psql -h localhost -p 6432 -U postgres pgbouncer`

Output:

 database  |   user    | cl_active | cl_waiting | sv_active | sv_idle | sv_used | sv_tested | sv_login | maxwait | maxwait_us |  pool_mode
-----------|-----------|-----------|------------|-----------|---------|---------|-----------|----------|---------|------------|-------------
 pgbouncer | pgbouncer |         1 |          0 |         0 |       0 |       0 |         0 |        0 |       0 |          0 | statement
 postgres  | postgres  |         2 |         97 |         0 |       0 |       0 |         0 |        1 |       0 |     748640 | transaction
 
The test runs 100 workers each sending queries that are always canceled.
Pgbouncer's pool of client connections is limited to 100, service's pool is limited to 10. One would expect the total number of  connections equal 10, but here we have 100.
