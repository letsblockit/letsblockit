#! /usr/bin/env nix-shell
#! nix-shell -i bash -p go_1_17 -p reflex
#! nix-shell --pure --quiet

reflex -r "(main.go|(src|data)/)" -s -- go run -race . --reload --oryurl https://inspiring-gates-6gd6xqshnz.projects.oryapis.com
