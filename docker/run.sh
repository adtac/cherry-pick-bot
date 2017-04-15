#!/usr/bin/env bash

export GITHUB_PRIVATE_KEY=/tmp/ssh_key

function generate_keys {
    ssh-keygen -b 4096 -N '' -f $GITHUB_PRIVATE_KEY

    echo "New keys has been generated, below is the public key"
    cat "$GITHUB_PRIVATE_KEY.pub"

    sleep 5s
}

eval $(ssh-agent -s)

[ ! -f $GITHUB_PRIVATE_KEY ] && generate_keys

ssh-add $GITHUB_PRIVATE_KEY

exec cherry-pick-bot
