# uit-hub

## Fake server

### Run

```bash
cd server
go run main.go
```

Or run from repo root:

```bash
SEED=default go -C server run .
```

Server runs at `http://localhost:3000` by default.

### Auth flow

1. Login to get a token.

```bash
curl -X POST http://localhost:3000/login \
	-H 'Content-Type: application/json' \
	-d '{"student_id":"22520001","password":"pass123"}'
```

Response:

```json
{ "message": "Successfully!", "data": { "token": "mock-22520001" } }
```

2. Use the token for endpoints under `/student`, `/rooms/availability`, and `/contact`.

````bash
curl http://localhost:3000/student/profile \
	-H 'Authorization: Bearer mock-22520001'

### Reset mock dataset

Reset the in-memory dataset back to the loaded fixture.

```bash
curl -X POST http://localhost:3000/admin/reset
````

Optional protection: if you set `ADMIN_TOKEN`, then you must provide header `X-Admin-Token`.

```bash
ADMIN_TOKEN=secret SEED=default go -C server run .
curl -X POST http://localhost:3000/admin/reset -H 'X-Admin-Token: secret'
```

```

```
