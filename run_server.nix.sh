#! /usr/bin/env nix-shell
#! nix-shell -i bash -p go_1_16 -p reflex
#! nix-shell --quiet

xdg-open http://localhost:8765
reflex -r "(main.go|(src|data)/)" -s go run .
