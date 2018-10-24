MODE=$1
echo "MODE: "$MODE

if [ -z "$MODE" ]
then
        go build -o dumper -v main.go
       ./dumper dump --create-database=true --delete-dump-files=true  --skip-bucket-store=true 
else
        go build -o dumper -v main.go
       ./dumper dump --create-database=true --delete-dump-files=true  --skip-bucket-store=true >> log.txt 2>&1
fi

