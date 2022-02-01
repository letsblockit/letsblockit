#! /usr/bin/env nix-shell
#! nix-shell -i bash -p go_1_17 -p reflex
#! nix-shell --pure --quiet

reflex -r "(cmd|src|data)/" -s -- go run -race ./cmd/server --reload
