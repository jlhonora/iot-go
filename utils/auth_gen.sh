#!/bin/bash

user=$1
pass=$2

pass_hash=`printf $pass | shasum -a 1 | awk -F ' ' '{print $1}'`

echo $user:$pass_hash | base64
