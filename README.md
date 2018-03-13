# achelous
Sendmail replacement, which queues and uploads emails to a web service. We designed it for IoT devices and containers, which might send emails for administrative purposes.

**Naming:** The Achelous (Greek: Αχελώος, Ancient Greek: Ἀχελῷος Akhelôios), also Acheloos, is a river in western Greece. [Wikipedia](https://en.wikipedia.org/wiki/Achelous_River)

The project is split into two programs:

- `spring` is responsible for mail submission and implementes the `sendmail` command line interface. The `sendmail` binaries are provided by symlinks pointing to `spring`.
- `upstream` is responsible for mail upload and success reporting. It is implemented as a daemon process and can be managed via `systemd` or `rc.d` scripts.

## Features

- Sendmail compatible command line interface, designed as a drop-in replacement. Provides the virtual `mail-transport-agent` package.
- Compatible with ISC `cron` and GNU `mail`.
- A hierarchical queue structure, which allows embedding of container queues in sub-queues on the host system.
- URL and header variables of target web service are freely configurable.
- Reports successfully processed queues to the web service. Allows the implementation of a health check.

## Building

**Note:** Pre-built packages are available [here](https://github.com/core-process/achelous/releases). Have a look and skip this section, if you find the appropriate `deb` file for your architecture.

Make sure the following tools are available in your `PATH`:

- `make`
- `gcc`
- `go`
- `glide`
- `dpkg-deb` (for packaging)

Run `make` to build the binaries from source. The binaries will be placed in the `.build/bin` directoriy.

Run `make dist` to build the package file. The `deb` file will be placed in the `.build/dist` directory.

## Installation

Achelous conflicts with all packages providing the virtual `mail-transport-agent` package, e.g. `exim4-daemon-light`. Please remove them first.

Download or build the package file. Install achelous with `sudo dpkg -i achelous_*.deb`.

Create and adjust the configuration files `/etc/achelous/spring.json` and `/etc/achelous/upstream.json`.

## Configuration

**Note:** Usually specific reloading procedures are not required. `spring` reads its configuration file on every call. `upstream` reads its configuration file on every queue run.

### Examples

`/etc/achelous/spring.json`

```json
{
  "DefaultQueue": "",
  "PrettyJSON": false,
  "TriggerQueueRun": true
}
```

`/etc/achelous/upstream.json`

```json
{
  "PauseBetweenRuns": {
    "PreviousRunOK": "60m",
    "PreviousRunWithErrors": "1m"
  },
  "Target": {
    "Upload": {
      "URL": "https://myservice/api/email/upload",
      "Header": {
        "Authorization": [ "Bearer example" ]
      }
    },
    "Report": {
      "URL": "https://myservice/api/email/report",
      "Header": {
        "Authorization": [ "Bearer example" ]
      }
    },
    "RetriesPerRun": {
      "Attempts": 3,
      "PauseBetweenAttempts": "10s"
    }
  }
}
```

### Specification

`/etc/achelous/spring.json`

```json
{
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "DefaultQueue": {
      "type": "string",
      "description": "Relative queue path. The empty string represents the root or default queue."
    },
    "PrettyJSON": {
      "type": "boolean",
      "description": "Activates pretty printing of the email to be submitted to the queue."
    },
    "TriggerQueueRun": {
      "type": "boolean",
      "description": "Triggers a queue run on new mail submissions to reduce delays."
    }
  },
  "required": []
}
```

`/etc/achelous/upstream.json`

```json
{
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "PauseBetweenRuns": {
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "PreviousRunOK": {
          "type": "string",
          "pattern": "^([0-9]+[hms])+$",
          "description": "Pause till next queue run in case the previous queue run was successful."
        },
        "PreviousRunWithErrors": {
          "type": "string",
          "pattern": "^([0-9]+[hms])+$",
          "description": "Pause till next queue run in case the previous queue run failed partially or completely."
        }
      },
      "required": []
    },
    "Target": {
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "Upload": {
          "type": "object",
          "additionalProperties": false,
          "properties": {
            "URL": {
              "type": "string",
              "description": "The URL used for mail uploads via POST."
            },
            "Header": {
              "type": "object",
              "description": "The headers used for mail uploads.",
              "additionalProperties": false,
              "properties": {
                "Authorization": {
                  "type": "array",
                  "items": {
                    "type": "string"
                  }
                }
              },
              "required": []
            }
          },
          "required": []
        },
        "Report": {
          "type": "object",
          "additionalProperties": false,
          "properties": {
            "URL": {
              "type": "string",
              "description": "The URL used for success reporting via POST."
            },
            "Header": {
              "type": "object",
              "description": "The headers used for success reporting uploads.",
              "additionalProperties": false,
              "properties": {
                "Authorization": {
                  "type": "array",
                  "items": {
                    "type": "string"
                  }
                }
              },
              "required": []
            }
          },
          "required": []
        },
        "RetriesPerRun": {
          "type": "object",
          "additionalProperties": false,
          "properties": {
            "Attempts": {
              "type": "integer",
              "description": "Number of successive retries in case an upload failes, before giving up during the current queue run."
            },
            "PauseBetweenAttempts": {
              "type": "string",
              "pattern": "^([0-9]+[hms])+$",
              "description": "Pause between successive retries."
            }
          },
          "required": []
        }
      },
      "required": []
    }
  },
  "required": []
}
```

## Upload and Report Protocol

Uploading and reporting is done via a `POST` request. The URL as well as the headers to be sent are fully configurable via the configuration file `/etc/achelous/upstream.json`.

The following *JSON Schema* specifies the request body of the upload request:

```json
TODO
```

The reporting request is performed on successful queue runs only (success measure). The request body is empty.
