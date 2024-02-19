## Entain BE Technical Test Solution

This is the solution of Entain BE Technical Test.

### Directory Structure

- `api`: A basic REST gateway, forwarding requests onto service(s).
- `racing`: A very bare-bones racing service.
- `sport`: A very bare-bones sport service.

```
entain/
├─ api/
│  ├─ proto/
│  ├─ main.go
├─ racing/
│  ├─ db/
│  ├─ proto/
│  ├─ service/
│  ├─ main.go
├─ sport/
│  ├─ db/
│  ├─ proto/
│  ├─ service/
│  ├─ main.go
├─ README_OG.md
├─ README.md
```

### Resource

- `README_OG.md`: File contains the original task. 
- [init](https://github.com/jie1311/Entain/commit/d147197769b2312b5a6c32f86a85ed686536e3f5) contians the original codes.

### Getting Started

1. Install Go (latest).

[see here](https://golang.org/doc/install).

2. Install `protoc`

[see here](https://grpc.io/docs/protoc-installation/).

3. In a terminal window, cd to the entain directory, start the service...

```bash
make start
```

4. In the same terminal window, still within the entain directory, stop the service...

```bash
make stop_all
```

### Data structure
```ptoto
message Race {
  // ID represents a unique identifier for the race.
  int64 id = 1;
  // MeetingID represents a unique identifier for the races meeting.
  int64 meeting_id = 2;
  // Name is the official name given to the race.
  string name = 3;
  // Number represents the number of the race.
  int64 number = 4;
  // Visible represents whether or not the race is visible.
  bool visible = 5;
  // AdvertisedStartTime is the time the race is advertised to run.
  google.protobuf.Timestamp advertised_start_time = 6;
  // Satus represents whether the race is open or closed
  Status status = 7;
}
```

```ptoto
message Event {
  // ID represents a unique identifier for the event.
  int64 id = 1;
  // Name is the official name given to the event.
  string name = 2;
  // Type is the type of the event.
  string type = 3;
  // Location is the location where the event happens.
  string location = 4;
  // Visible represents whether or not the event is visible.
  bool visible = 5;
  // AdvertisedStartTime is the time the event is advertised to run.
  google.protobuf.Timestamp advertised_start_time = 6;
  // Satus represents whether the event is open or closed
  Status status = 7;
}
```

### APIs 

1. `ListRaces`

In a terminal window, call `ListRaces` for simply listing all races... 

```bash
curl -X "POST" "http://localhost:8000/v1/list-races" \
     -H 'Content-Type: application/json' \
     -d $'{
  "filter": {}
}'
```

It can be filtered by `meething_ids` and/or `visibleOnly`... 

```bash
curl -X "POST" "http://localhost:8000/v1/list-races" \
     -H 'Content-Type: application/json' \
     -d $'{
  "filter": {"meetingIds": [10, 8], "visibleOnly": true}
}'
```

- `meething_ids` takes an array of integers 
- `visibleOnly` takes a boolean value, `false` would not filter the results, `true` would only return races whose `visible` is `true`

It can be sorted by one of `ID`, `MEETING_ID` , `NAME` , `NUMBER` , `VISIBLE` and `ADVERTISED_START_TIME` (all caps)... 
```bash
curl -X "POST" "http://localhost:8000/v1/list-races" \
     -H 'Content-Type: application/json' \
     -d $'{
   "sorting": {"sort_by": "NAME", "descend": true}
}'
```

- if there's no `sorting` object, or `sort_by` is set to `UNSPECIFIED`, results would by default sorted by `ADVERTISED_START_TIME` in ascending order
- bt default, results are ordered in ascending order, however, they can be ordered in decending order by adding `"descend": true` in the `sorting` object

2. `GetRaces`

In a terminal window, call `GetRaces` for simply getting all race macted by `id`, `96` for example... 

```bash
curl -X "GET" "http://localhost:8000/v1/get-race/96"
```
- if there's no race with requested id, it would return `null`

3. `ListEvents`

In a terminal window, call `ListEvents` for simply listing all events... 

```bash
curl -X "POST" "http://localhost:8000/v1/list-events" \
     -H 'Content-Type: application/json' \
     -d $'{
  "filter": {}
}'
```

It can be filtered by `visibleOnly`... 

```bash
curl -X "POST" "http://localhost:8000/v1/list-events" \
     -H 'Content-Type: application/json' \
     -d $'{
  "filter": {"visibleOnly": true}
}'
```

- `visibleOnly` takes a boolean value, `false` would not filter the results, `true` would only return events whose `visible` is `true`

It can be sorted by one of `ID`, `NAME` , `TYPE` , `LOCATION` , `VISIBLE` and `ADVERTISED_START_TIME` (all caps)... 
```bash
curl -X "POST" "http://localhost:8000/v1/list-events" \
     -H 'Content-Type: application/json' \
     -d $'{
   "sorting":{"sort_by":"NAME"}
}'
```

- if there's no `sorting` object, or `sort_by` is set to `UNSPECIFIED`, results would by default sorted by `ADVERTISED_START_TIME` in ascending order
- bt default, results are ordered in ascending order, however, they can be ordered in decending order by adding `"descend": true` in the `sorting` object
