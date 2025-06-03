# Logger Environment Variables

The logger is now configured using a single environment variable for log level:

- `LOG_LEVEL`: Set to `debug`, `info`, `warn`, `error`, `fatal`, or `panic` to control the minimum log level. Default is `info`. (If not set, falls back to `TEAMSERVER_MODE` for legacy support.)
- `LOG_SHOW_TIME`: Set to `false` or `0` to hide timestamps in logs. Any other value (or unset) shows timestamps.
- `LOG_FILE`: If set, logs will be written to the specified file path instead of stdout.

Example usage in `.env`:

```
LOG_LEVEL=debug
LOG_SHOW_TIME=true
LOG_FILE=/var/log/meduza.log
```

Only `LOG_LEVEL` is used for log level control. `TEAMSERVER_MODE` is supported as a fallback for legacy setups.
