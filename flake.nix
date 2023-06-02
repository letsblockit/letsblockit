{
  description = "letsblock.it server and helpers";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-23.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        pinnedGo = pkgs.go_1_19;
        commonImageLabels = import ./nix/labels.nix;

        # Scripts to wrap, with their dependencies, available via `nix run .#script-name`
        scripts = with pkgs; {
          add-migration = [ self.packages.${system}.migrate ];
          run-migrate = [ self.packages.${system}.migrate ];
          run-server = [ pinnedGo reflex self.packages.${system}.ory ];
          run-tests = [ pinnedGo golangci-lint ];
          update-assets = [ pinnedGo nodejs-slim-18_x nodePackages.npm ];
          update-contributors = [ pinnedGo curl imagemagick git ];
          update-codegen = [ mockgen self.packages.${system}.sqlc ];
          update-labels = [ coreutils ];
          update-vendorsha = [ nix-prefetch gnused ];
          upgrade-deps = [ nodejs-slim-18_x nodePackages.npm pinnedGo nix-prefetch git ];
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
          vector = pkgs.callPackage ./nix/vector.nix { };

          render-container = pkgs.dockerTools.streamLayeredImage {
            name = "ghcr.io/letsblockit/render";
            tag = "latest";
            created = builtins.substring 0 8 self.lastModifiedDate;
            contents = self.packages.${system}.render;
            config = {
              Cmd = [ "render" ];
              Labels = {
                "org.opencontainers.image.title" = "letsblock.it render CLI";
                "org.opencontainers.image.documentation" = "https://github.com/letsblockit/letsblockit/blob/main/cmd/render/README.md";
              } // commonImageLabels;
            };
          };
          server-container = pkgs.dockerTools.streamLayeredImage {
            name = "ghcr.io/letsblockit/server";
            tag = "latest";
            created = builtins.substring 0 8 self.lastModifiedDate;
            contents = [ pkgs.cacert self.packages.${system}.vector self.packages.${system}.server ];
            config = {
              Cmd = [ "/bin/server" ];
              Env = [ "LETSBLOCKIT_ADDRESS=:8765" "PATH=/bin" ];
              ExposedPorts."8765/tcp" = { };
              Labels = {
                "org.opencontainers.image.title" = "letsblock.it server";
                "org.opencontainers.image.documentation" = "https://github.com/letsblockit/letsblockit/blob/main/cmd/server/README.md";
              } // commonImageLabels;
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
