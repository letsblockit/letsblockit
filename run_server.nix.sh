#! /usr/bin/env nix-shell
#! nix-shell -i bash -p go_1_16 -p reflex -p sqlite
#! nix-shell --pure --quiet

reflex -r "(main.go|(src|data)/)" -s -- go run -tags libsqlite3 . --reload --oryproject playground
