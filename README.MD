# game service

Game Service is designed to serve game app(s) on multiple devices via `REST` calls.<br/> All Internal `APIs` & the `Database` are designed to support multiple games in future, with multiple supporting modes.<br/>
The Service is bootstrapped to work with a game named `FAMOUS SHOOTER` and supports following `3 modes`:
><ol>
>	<li>Battle Royal</li>
>	<li>Team Deathmatch</li>
>	<li>Capture the flag</li>
></ol>

## assumptions

>
><ul>
>    <li>gameid=1 is default game id for all endpoints</li>
>    <li>A player can play only one game at a time</li>
>    <li>Area is always 3 letter code</li>
>    <li>Game names are Unique & can support multiple modes</li>
>    <li>Multiple games can have same modes</li>
>    <li>A Game can support any number of players</li>
>    <li>JSON is the supported format for now</i>
></ul>

## endpoints

>
> ### Retrieving game's top modes in a given area
>
> - GET `/top-modes?area=blr`
>   - returns a json array of all active modes in the given `area` in `descending` order of popularility i.e. the `mode` with highest number of active players will be first followed by sencond highest and so on
>
> ### Add new player to the game in a given area
>
> - POST `/player?modeid=1&area=blr`
>   - returns the `playerid` assigned to the new player, this `playerid` can be used to remove the player from the game
>
> ### Remove a player from the game 
>
> - DELETE `/player?playerid=1`


## prerequisites
`Game Service` at bare minimim needs `docker` installed/running and a `.env` at the project's root with following details
```
LILA_DB_HOST=postgres
LILA_DB_PORT=5432
LILA_DB_USER=postgres
LILA_DB_PASSWORD=postgres
LILA_DB_DATABASE=postgres
LILA_DB_DRIVER=postgres
LILA_DB_EXPLORER_EMAIL=someone@postgres.com
LILA_DB_EXPLORER_PASSWORD=postgres
LILA_DB_RUN_MIGRATION=true
LILA_SERVER_PORT=:1234
LILA_TIME_TO_WAIT_FOR_SUPPORTING_SERVICES_TO_COME_UP=s30s
``` 

## execution

>`make up` will bring up all services.

After all services have successfully come up visit
[localhost](http://localhost:1234).
To explore the postgres data visit [postgres](http://localhost:5050) and login with `LILA_DB_EXPLORER_EMAIL` and `LILA_DB_EXPLORER_PASSWORD` and configure the database in there with the said configuration from `.env` file.

>`make down` will teardown all services


## load testing

```
$ bombardier -c 125 -n 100000 http://localhost:1234/top-modes\?area\=blr

Bombarding http://localhost:1234/top-modes?area=blr with 100000 request(s) using 125 connection(s)
 100000 / 100000 [========================================================================================================================] 100.00% 7453/s 13s
Done!
Statistics        Avg      Stdev        Max
  Reqs/sec      7543.95    1912.33   12149.01
  Latency       16.57ms    23.13ms   529.39ms
  HTTP codes:
    1xx - 0, 2xx - 100000, 3xx - 0, 4xx - 0, 5xx - 0
    others - 0
  Throughput:     1.54MB/s
```

## architecture & performance

The Service is designed following the `go` conventions, `executables` are kept inside `cmd` directory, re-usable packages are kept under `pkg` directory &  critical business logic is kept under `internal` directory.<br/> 
All database tables have been properly indexed to suit the service's need. <br/>
`Server` has dependency on `Service` <br/>
`Service` has denepdency on `storage layer` or the `Database` which can be swapped with `in-memory` or `elastic` should there be a need pertaining to performance (in future).<br/>
We can also opt `gRPC` based communications using `protobuf` for performance & to cap bandwidth - keeping in view mobile devices.<br/>
Based on `load-testing` above the `service` is able to serve an average of `7.5K requets/second`.

> server -> service -> database
>
> Game Service is composed up of three microservices
> - postgres
>   - is the database used     
> - game_app
>   - actual game service serving the `REST` calls  
> - postgres_explorer
>   - handy UI for exploring the postgres database




## problem solving time

> It took ~6 hours in total to build the service.
>
> Following are the bifurcation
>   - Architecting the service, code structure and database structure took approximately 2 hours.
>   - Coding & Dockerizing took 3 hours.
>   - Testing took an 1 hour.