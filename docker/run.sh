#!/usr/bin/env bash

export GITHUB_PRIVATE_KEY=/tmp/ssh_key

function generate_keys {
    ssh-keygen -b 4096 -N '' -f $GITHUB_PRIVATE_KEY

    echo "New keys has been generated, below is the public key"
    cat "$GITHUB_PRIVATE_KEY.pub"

    sleep 5s
}

# start the ssh-agent in the background
eval $(ssh-agent -s)

# if no private keys are provided, generate them (and output the public key to console)
[ ! -f $GITHUB_PRIVATE_KEY ] && generate_keys

# add the SSH key
ssh-add $GITHUB_PRIVATE_KEY

# set the GIT_SSH_COMMAND
export GIT_SSH_COMMAND="ssh -i $GITHUB_PRIVATE_KEY"

# add github.com to the list of known hosts
mkdir -p $HOME/.ssh
ssh-keyscan github.com >> $HOME/.ssh/known_hosts

# start cherry-pick-bot
exec cherry-pick-bot
