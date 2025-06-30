# jobd

a simple cicd daemon for executing jobs on a schedule

# configuration

default configuration

```
datadir: /srv/jobd
shell_defaults:
  workingdir: ""
  env: []
  immediateexit: false
  timeout: 3600
```

shell_defaults

- `workingdir`
- `env`
- `immediatexit` like adding `set -e` to scripts
- `timeout`


# Job spec

a job or a JobSpec defines how to run a job, which is a set of steps

jobs are define under the `jobs:` section of the configuration file. Each entry
is an array.

a Job definition looks like this


Minimal job spec:

```
jobs:
- jobid: someid
  schedule: '@every 1h'
  steps:
  - shell:
       script: df -h /data > example.txt
```


Full job spec

- `disabled` bool - true if job is disabled. default false.
- `immediate` bool - true if job runs immediately at startup, default false.
- `stackable` bool - true if job can run multiple instances, default false.
- `keep` int - number of previous runs to keep, 0 means keep all. default 0.
- `directory` - string - the starting directory of the job.
   default is {datadir}/jobs/{jobid}/run/{runid}/workspace
- `schedule` - string - cron schedule, see github.com/robfig/cron for expressions
- steps - array - array of job steps

# step spec

- `id` - string - optional step id
- `shell` - string - shell spec definition
- `vars` - array - variables to capture

```
jobs:
- jobid: someid
  disabled: true
  immediate: true
  stackable: true
  keep: 3
  directory: /some/dir
  schedule: '@every 1h'

  steps:
  - id: step1
    shell:
       script: |
            df -h /data > example.txt

    vars:
      - from_command:
        - foo: cat example.json | jq -r .foo
        - bar: cat example.json | jq -r .foo

  - id: step2
    shell:
       immediateexit: true
       script: |
            ls -1 stuff/*.json > todo.txt


```

# shell spec


- `script` - string - script
- `workingdir`
- `env`
- `immediatexit`
- `timeout`


# variables

A special `$ENV` variable is set to carry environment variable between steps

To pass a variable between steps

    echo FOO=BAR >> $ENV

to access the variable in a different step

    echo $FOO
