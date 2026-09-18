# Bruno collection: go-micro

Open this folder (`bruno/`) in Bruno. URLs, `tenant` and `addressId` are collection defaults in `collection.bru`. Environments only override ports. Select `local` for `configs/service_local.yaml`.

| Environment | Config | HTTP | HTTPS |
| --- | --- | --- | --- |
| `local` | `configs/service_local.yaml` | 9480 | 9443 |
| `testdata` | `testdata/service_local.yaml` | 9000 | 9443 |
| `minimal` | `testdata/service_local_minimal.yaml` | 8000 | 8543 |

With TLS enabled, the HTTP port only serves `/livez`, `/readyz`, `/` and `/metrics`. Address CRUD runs against `apiUrl` (HTTPS).

The service uses a generated or file-based certificate. Disable TLS verification in Bruno, or run:

```
bru run --insecure --env local
```

JWT is off in the default local config. If you enable `auth.type: jwt`, set a Bearer token on the collection.
