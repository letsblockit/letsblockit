{
  description = "letsblock.it server and helpers";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-21.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        pinnedGo = pkgs.go_1_17;

        # Scripts to wrap, with their dependencies, available via `nix run .#script-name`
        scripts = with pkgs; {
          add-migration = [ self.packages.${system}.migrate ];
          run-migrate = [ self.packages.${system}.migrate ];
          run-server = [ pinnedGo reflex self.packages.${system}.ory ];
          run-tests = [ pinnedGo golangci-lint ];
          update-assets = [ nodejs-17_x nodePackages.npm ];
          update-codegen = [ mockgen self.packages.${system}.sqlc ];
          update-vendorsha = [ pkgs.nix-prefetch ];
        };
      in
      {
        defaultPackage = self.packages.${system}.run-server;
        packages = {
          render = pkgs.callPackage ./nix/letsblockit.nix { cmd = "render"; };
          server = pkgs.callPackage ./nix/letsblockit.nix { cmd = "server"; };
          migrate = pkgs.callPackage ./nix/migrate.nix { };
          ory = pkgs.callPackage ./nix/ory.nix { };
          sqlc = pkgs.callPackage ./nix/sqlc.nix { };

          render-docker = pkgs.dockerTools.streamLayeredImage {
            name = "letsblockit-render";
            tag = "latest";
            created = "now";
            contents = self.packages.${system}.render;
            config = {
              Cmd = [ "render" "--help" ];
            };
          };
        } // (builtins.mapAttrs
          (name: deps: pkgs.writeShellApplication {
            name = name;
            runtimeInputs = deps;
            text = ''
              # Make nix-prefetch use nixpkgs from the flake lock
              export NIX_PATH="nixpkgs=${nixpkgs.sourceInfo.outPath}"
              ./scripts/${name}.sh "$@"
            '';
          })
          scripts);

        apps = {
          render = flake-utils.lib.mkApp {
            drv = self.packages.${system}.render;
            exePath = "/bin/render";
          };
          server = flake-utils.lib.mkApp {
            drv = self.packages.${system}.server;
            exePath = "/bin/server";
          };
        };

        devShell = pkgs.mkShell {
          # Build inputs from the packages
          inputsFrom = builtins.attrValues self.packages.${system};
          # Runtime inputs from the scripts
          buildInputs = builtins.concatLists (builtins.attrValues scripts);
        };

        overlay = final: prev: {
          letsblockit = self.packages.${system}.server;
          migrate = self.packages.${system}.migrate;
          ory = self.packages.${system}.ory;
        };
      });
}
