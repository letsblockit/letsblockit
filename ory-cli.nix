{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_17;
}
{
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
    mkdir -p $out/bin
    cp $GOPATH/bin/cli $out/bin/ory
  '';
}
