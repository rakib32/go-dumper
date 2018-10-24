# go-dumper
This a data dump tool which will dump data from one database, then will upload it to GS or AWS bucket and finally will restore the data to destination server

- Creates backup form production database
- Loads backup sql file to staging database

## Options
- Usage:
  dumper [command]

- Available Commands:
  - `dump`      Creates backup form production database
  - `load`        Loads backup sql file to staging database

- Flags:
    - `--create-database`     Create new schema to target DB. Format will be like: dbname_20180101
    - `--delete-dump-files`   Delete Dump files .sql and .tar
    - dump 
        - `--skip-bucket-store`   It will ignore bucket storing
        - `--skip-restore`        Only backup will be created 
    - load 
        - `--dump-path`           It can be either local filepath or GS bucket path 
  
## Settings

- To run tool:

    - Run dumper via run.sh:
        - ./run.sh
        - ./run.sh prod [for production]
        
- To run dumper as a CronJob in kubernets cluster(Production)
    - Update `deploy/cronjob.yaml` file with your desired time to run 
    - Apply the above changes using following command: `kubectl apply -f cronjob.yaml`
    - Other userful commands 
        - `kubectl get pods`
        - `kubectl describe pod`
        - `kubectl delete jobs --all`



### Notes:
- *important* please use mysql client 8.0 or higher to avoid time out issues for large datbase
- used [viper](https://github.com/spf13/viper) to import config from [consul](https://www.consul.io) via [crypt](https://github.com/xordataexchange/crypt/tree/master/bin/crypt).
- used [logrus](github.com/sirupsen/logrus) to print logs
    - logrus is completely api-compatible with the stdlib logger, so you can replace your `log` imports everywhere with `log "github.com/sirupsen/logrus"`
