---
outline: deep
---

# SCIM Provisioning

<span class="idpzero-text">idpzero</span> can push its configured users to an external
[SCIM 2.0](https://datatracker.ietf.org/doc/html/rfc7644) service. In this exchange
<span class="idpzero-text">idpzero</span> acts as the **SCIM client** (the provisioning
side) and the external system is the **service provider** (the target). This lets you
seed or keep a downstream, SCIM-enabled system in sync with the same identities your
local IdP issues tokens for.

::: tip SCOPE
Only the **User** resource is supported, and provisioning is **outbound only** —
<span class="idpzero-text">idpzero</span> pushes users out; it does not expose a SCIM
server for others to write into. Groups are not provisioned.
:::

## Configuration

Add an optional `scim` block to `.idpzero/server.yaml`:

```yaml
scim:
  endpoint: https://your-target.example.com/scim/v2   # base URL of the SCIM service
  bearer_token: ""                                     # optional; see below
```

| Field | Description |
| --- | --- |
| `endpoint` | Base URL of the target SCIM service. <span class="idpzero-text">idpzero</span> appends `/Users` (and `/Users/{id}`) to this. |
| `bearer_token` | Optional. Sent as `Authorization: Bearer <token>` on every request. Omit it if the target needs no auth. |

The users that get pushed are the same ones under the top-level `users:` key in
`server.yaml` — the ones shown on the **Users** page.

### Providing the token via environment variable

So you don't have to commit a secret to source control, the token can also be supplied
through the `IDPZERO_SCIM_TOKEN` environment variable, which **takes precedence** over
`scim.bearer_token`:

```sh
IDPZERO_SCIM_TOKEN="my-secret-token" idpzero serve
```

::: warning KEEP SECRETS OUT OF SOURCE CONTROL
Prefer `IDPZERO_SCIM_TOKEN` for real tokens. Only put `bearer_token` directly in
`server.yaml` for local/throwaway targets, and never commit a production token.
:::

## Running a sync

Provisioning is triggered manually from the dashboard so it stays predictable and
observable — there is no automatic background push.

1. Start the server with `idpzero serve`.
2. Open the dashboard and go to the **Provisioning** page (default
   [http://localhost:4379/scim](http://localhost:4379/scim)).
3. Click **Sync now**.

Each configured user is reconciled against the target and the result is shown per user
(`created`, `updated`, or `failed` with the error detail). The page shows the outcome of
the most recent run only.

### Reset history

**Reset history** clears the last-run status so the page returns to its initial state.

::: tip
Reset only clears <span class="idpzero-text">idpzero</span>'s **local** view of the last
sync. It does **not** delete users already provisioned on the target.
:::

## How reconciliation works

<span class="idpzero-text">idpzero</span> keeps **no local record** of what it has
provisioned — each sync reconciles from scratch. For every configured user it:

1. Looks the user up on the target with `GET /Users?filter=externalId eq "<subject>"`.
2. If **no match** exists → `POST /Users` (**create**).
3. If a **match** exists → `PUT /Users/{id}` (**replace**).

The user's `subject` is used as the SCIM `externalId`, which is what makes repeated syncs
idempotent — re-running a sync updates existing users rather than duplicating them.

### Attribute mapping

| idpzero (`server.yaml`) | SCIM User attribute |
| --- | --- |
| `subject` | `externalId` |
| `claims.preferred_username` (falls back to `subject`) | `userName` |
| `login_display` (falls back to `claims.name`) | `displayName` |
| `claims.name` / `given_name` / `middle_name` / `family_name` | `name.formatted` / `givenName` / `middleName` / `familyName` |
| `claims.email` | `emails[0]` (primary) |
| `claims.phone` | `phoneNumbers[0]` (primary) |
| — | `active` (always `true`) |

## Requirements on the target

For provisioning to work correctly, the SCIM service must:

- **Honour the `externalId` filter** on `GET /Users` (RFC 7644 §3.4.2.2). This is how
  <span class="idpzero-text">idpzero</span> matches existing users. If the target ignores
  the filter and returns unrelated records, <span class="idpzero-text">idpzero</span>
  detects the mismatch and fails the user rather than modifying the wrong record.
- **Accept the bearer token** (when configured) and return **JSON** responses, including
  JSON SCIM error envelopes on failure.
- Not sit behind an interactive-login gateway — see below.

## Troubleshooting

**`unexpected redirect (302) to .../auth/login`**
The target bounced the request to a browser login page instead of authenticating the API
call. This usually means the SCIM route is behind an interactive-login / session gateway,
or the bearer token isn't being accepted. A correctly configured SCIM endpoint returns a
JSON `401`, not a redirect. <span class="idpzero-text">idpzero</span> deliberately does not
follow the redirect and reports it clearly.

**`... does not appear to support filtering by externalId`**
The target returned users that don't match the requested `externalId`, meaning it isn't
applying the `filter`. Fix filtering on the service provider; until then existing users
can't be matched and updated safely.

**`expected a JSON response but got Content-Type ...`**
The endpoint returned a non-JSON (e.g. HTML) body with a success status. Double-check that
`endpoint` points at the SCIM API base URL and not a web page.
