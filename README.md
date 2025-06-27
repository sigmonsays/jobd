# jobd

a simple cicd daemon for executing jobs on a schedule

# Job Spec

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

```
jobs:
- jobid: someid
  disable: true
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
