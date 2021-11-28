{
  description = "letsblock.it server";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/release-21.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachSystem [ "x86_64-linux" ] (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        pinnedGo = pkgs.go_1_17;
        buildGoModule = pkgs.buildGoModule.override {
          go = pinnedGo;
        };
      in
      {
        defaultPackage = self.packages.${system}.letsblockit;
        packages.letsblockit = buildGoModule {
          pname = "letsblockit";
          version = "1.0";
          vendorSha256 = "1zvm97x0rm65d1gv4hrksdb3gk18rmba5d5z1ciksxakxw0dj2gz";
          src = ./.;
          doCheck = false;
        };

        packages.ory-cli = buildGoModule {
          pname = "ory-cli";
          version = "0.1.0";
          src = pkgs.fetchFromGitHub {
            owner = "ory";
            repo = "cli";
            rev = "v0.1.0";
            sha256 = "1fg069gzjsvz933pz867ghy0wizvmaf99x17r5vw6hc7a0s2nvqs";
          };
          vendorSha256 = "0dkfis7h1il8xyj1vl55agfrpr94qc9v4ml5w2i8rrpxg7fdxhpk";
          doCheck = false;
          installPhase = ''
            install -D $GOPATH/bin/cli $out/bin/ory-cli
          '';
        };

        defaultApp = self.apps.${system}.letsblockit;
        apps.letsblockit = flake-utils.lib.mkApp {
          drv = self.packages.${system}.letsblockit;
        };
        apps.ory-proxy = flake-utils.lib.mkApp {
          drv = self.packages.${system}.ory-cli;
        };

        devShell = pkgs.mkShell {
          buildInputs = with pkgs; [
            self.packages.${system}.ory-cli
          ];
          inputsFrom = builtins.attrValues self.packages.${system};
        };
      });
}
