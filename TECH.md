# Hegemonie: The Technical Stack

## Architecture

The Hegemonie platform consist in a set of microservices.

* **region server** is a tuple of 3 ``grpc`` services responsible for the game
  logic within a single region. The service itself is not authenticated. Please
  refer to the ``grpc`` description of the
  [City API](https://github.com/hegemonie-rpg/engine/blob/master/pkg/region/city.proto)
  ,
  [Army API](https://github.com/hegemonie-rpg/engine/blob/master/pkg/region/army.proto)
  and
  [Admin API](https://github.com/hegemonie-rpg/engine/blob/master/pkg/region/admin.proto)
  for more information.

* **events server** is a ``grpc`` service to subscribe to events related to
  in-game characters. EAch character's log collects structured events destined
  to be rendered with the proper localization.

* **maps server** is a ``grpc`` service managing maps as directed graphs that
  implement the maps in Hegemonie. It proposes path computations and paginated
  listing of the graph elements (vertices, nodes).

In addition, side services act as technology enablers.

* The services of the [ORY suite](https://ory.sh) provide OpenAPI interfaces to
  manage the general authentication and authorization needs.

* Services implementing the [OpenTelemetry](https://opentelemetry.io/) suite
  provide the collection, the aggregation, the storage and the display of events
  traces. Those services help troubleshooting problems in the whole solution,
  despite its largely distributed character.

* An **API gateway** ensuring the required authentication of the calls to the
  backend, doing the load balancing among the target backend, rate limiting on a
  per-user basis, etc. This is currently achieved by an
  [HAProxy](https://github.com/haproxy/haproxy) that will soon delegate the
  authentication and authorization decisions to an
  [ORY OathKeeper](https://github.com/ory/oathkeeper).

* A [Prometheus](https://prometheus.io) service to keep the in-game statistics
  that serve the Game Mastery purpose.

![Hegemonie Architecture](https://raw.githubusercontent.com/hegemonie-rpg/engine/master/docs/system-architecture.png)

## A Single administration CLI

Everything is controlled by a single CLI tool, ``hege`` that allows starting the
Hegemonie services, doing the daily operations, and a giving a shortcut access
to the backend services (without any authentication).

## Hegemonie is pure golang

Written in 100% in [Go](https://golang.org): for the sake of Simplicity and
Portability. The code mostly depends
on [Golang grpc](https://github.com/grpc/grpc-go) and
Golang [standard library](https://golang.org/pkg). At the moment no special
attention has been paid to the performance of the whole thing: this will happen
after the release of a very first MVP.

## Reliability

The current effort tend to make every workdload stateless with a proven
persistence backend solution.

## Scalability

Scalability is not a concern yet. Large communities are not a target for
Hegemonie. However, there are already a few opportunities to let the game
scaling oppotunities:

* The **API gateway**, whatever nginx or haproxy, acts as a stateless ingress
  proxy and can ensure HA in an active/active fashion.

* The **event server** is currently stateful because it relies on a local
  storage. Further scaling plans exist, based on a stateless service in front of
  a relatively scalabale KV backend (TiKV), plus a partitioning/sharding of the
  users if necessary. ``TiKV`` services have their own scalability model.

* The **region server** is stateful: it manages all the entities in-game.
  Distinct region services (i.e. processes) will host distinct datasets. So
  there is a limit in size for a region, but a de facto natural sharding
  opportunity of the users among the regions (i.e. services).

* The **map server** keeps a cache (loaded from a read-only reference) but
  serves stateless requests on read-only content. It can be multiplied _ad lib_.

* The services of the **ORY suite** are stateless and can be multiplied as much
  as required. Their underlying ``PostgreSQL`` instances have their own
  scalability model.

* The services of the **OpenTelemetry stack** are mostly stateless, and the
  storage solutions have their own scalability measures.

## Performance

Further than the design choice, in the region server, to keep each region "live
in RAM", the performance is not a topic yet.

We roughly target a system that can manage a game instance for a small community
of less than 50 players, that would be lightweight enough to run on a ARM-based
single board computer (e.g. a RaspberryPi v3).

## Security

Hegemonie makes several choices related to the security.

First, the security is enforced at the gate. The API gateway is a TLS
termination endpoint but no other SSL connection is used internally. There is
nothing forbidding an internal TLS usage, this is just the default chosen.

The API gateway is also a Zero Trust proxy for RPC calls that requires a valid
authentication and authorizations. So that gRPC calls are identified by JWT at
the gate but not in each service's implementation.

Second, Hegemonie made the choice of the ORY suite of tools. Those tools give
clues of "security well done", in other words all the garantys of best-in-class
implementation.

## Deployment

### Deploy with Docker

TL;DR:

```shell
set -ex
git clone https://github.com/hegemonie-rpg/engine.git
cd hegemonie
./bin/hege-docker up
```

More information to come soon.
