FROM centos:latest
RUN yum -y update && yum install -y ca-certificates && yum -y localinstall https://dev.mysql.com/get/mysql80-community-release-el7-1.noarch.rpm && yum install -y mysql-community-client 

COPY dumper /bin/

ENTRYPOINT ["/bin/dumper"]