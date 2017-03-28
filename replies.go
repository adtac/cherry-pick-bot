package main

var invalidEmail string = "You do not have a public email address. I need one to operate."

var cannotCherryPick string = "Uh-oh. I can't cherry-pick these commits. Any of the following could be the reason:\n\n- There are conflicts due to other commits being cherry-picked before.\n- Something has been merged into master and that's causing a conflict (in this case, ask the author of this commit to rebase to master and resolve all conflicts; nothing I can do here).\n- These commits have already been added for cherry-picking! If the commits have changed since, please close that PR and cherry-pick everything again."

var cannotRebase string = "Uh-oh. I couldn't rebase. This may happen because master has changed a lot and there are conflicts now. I can't really resolve conflicts, so you're going to have to do this one manually. Sorry!"
