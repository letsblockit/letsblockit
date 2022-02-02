{
  description = "letsblock.it server and helpers";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-21.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachSystem [ "x86_64-linux" ] (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        overlay = (final: prev: {
          letsblockit = self.packages.${system}.server;
          ory = self.packages.${system}.ory;
        });

        defaultPackage = self.packages.${system}.server;
        packages.render = pkgs.callPackage ./nix/letsblockit.nix { cmd = "render"; };
        packages.server = pkgs.callPackage ./nix/letsblockit.nix { cmd = "server"; };
        packages.ory = pkgs.callPackage ./nix/ory.nix { };
        packages.sqlc = pkgs.callPackage ./nix/sqlc.nix { };

        packages.render-docker = pkgs.dockerTools.streamLayeredImage {
          name = "letsblockit-render";
          tag = "latest";
          created = "now";
          contents = self.packages.${system}.render;
          config = {
            Cmd = [ "render" "--help" ];
          };
        };

        defaultApp = self.apps.${system}.server;
        apps.render = flake-utils.lib.mkApp {
          drv = self.packages.${system}.render;
          exePath = "/bin/render";
        };
        apps.server = flake-utils.lib.mkApp {
          drv = self.packages.${system}.server;
          exePath = "/bin/server";
        };

        devShell = pkgs.mkShell {
          buildInputs = with self.packages.${system}; [
            ory
            sqlc
            pkgs.golangci-lint
            pkgs.mockgen
            pkgs.nix-prefetch
            pkgs.reflex
          ];
          inputsFrom = builtins.attrValues self.packages.${system};
        };
      });
}
