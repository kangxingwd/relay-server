#!/bin/bash

find apisrv -type f -exec dos2unix {} \;
find nginx -type f -exec dos2unix {} \;
find server -type f -exec dos2unix {} \;
find webui -type f -exec dos2unix {} \;


