# YouTube cookie runbook

viddl.me should not depend on Igor's personal Google account cookies. Treat YouTube as an unreliable upstream: cookies may reduce bot/account challenges, but they are not a guarantee.

## Recommended strategy

1. Keep cookieless mode working. The backend must omit `--cookies` when `YTDLP_COOKIES` is unset, missing, unreadable, or unwritable.
2. If YouTube starts requiring cookies, use a dedicated low-value Google/YouTube account created only for viddl.me.
3. Store cookies outside git and outside the web root if possible. A typical production path is `/var/www/viddl.me/backend/cookies.txt`.
4. Prefer manual refresh plus monitoring over automated Google login. Fully automated login is brittle, trips MFA/security flows, and is not worth it unless explicitly revisited.

## Cookie file requirements

- Netscape cookie format.
- Owned by the backend service user.
- Readable and writable by the backend service user because yt-dlp can update cookies on exit.
- Never committed, pasted into tickets, printed in logs, or sent in chat.

Example production setup:

```bash
sudo install -o www-data -g www-data -m 600 /tmp/cookies.txt /var/www/viddl.me/backend/cookies.txt
sudo systemctl restart viddl.service
```

Systemd environment:

```ini
Environment="YTDLP_COOKIES=/var/www/viddl.me/backend/cookies.txt"
```

## Manual refresh procedure

1. Open a browser profile used only for the dedicated viddl.me service Google account.
2. Log in to YouTube and confirm the account can view a normal public video.
3. Export YouTube cookies in Netscape format using a trusted local export tool.
4. Copy the file to the production cookie path with owner `www-data:www-data` and mode `0600`.
5. Restart the backend if the service does not pick up the changed cookie file.
6. Run the smoke probes below.

## Smoke probes

Use a tiny public YouTube video:

```bash
curl -sS --max-time 80 -D /tmp/viddl-info-headers.txt \
  -o /tmp/viddl-info-body.json \
  -H 'Content-Type: application/json' \
  -X POST https://viddl.me/api/info \
  --data '{"url":"https://www.youtube.com/watch?v=jNQXAC9IVRw"}'

curl -sS --max-time 180 -D /tmp/viddl-dl-headers.txt \
  -o /tmp/viddl-smoke.mp4 \
  -H 'Content-Type: application/json' \
  -X POST https://viddl.me/api/download \
  --data '{"url":"https://www.youtube.com/watch?v=jNQXAC9IVRw","format":"best"}'

file /tmp/viddl-smoke.mp4
```

Healthy results:

- `/api/info`: HTTP 200 and a JSON body containing a title.
- `/api/download`: HTTP 200, `Content-Type: video/mp4`, and `file` reports ISO Media / MP4.

Alert-worthy failures:

- `failed to fetch video information`
- `Sign in to confirm you're not a bot`
- cookie read/write errors
- non-MP4 response body from `/api/download`

## Monitoring recommendation

Run a scheduled smoke check against `/api/info` and `/api/download` with the short public YouTube URL above. Alert only on actionable failures; do not page for missing local cookies if cookieless mode succeeds.
