#!/usr/bin/env sh
curl --ftp-create-dirs -T bin/main -u $FTP_USER:$FTP_PASSWORD ftp://ease-62q56ueo.cloudapp.net
