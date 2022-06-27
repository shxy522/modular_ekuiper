#!/bin/sh

cur=$(dirname "$0")
echo "Base path $cur"
pip install -r $cur/requirements.txt
echo "Done"