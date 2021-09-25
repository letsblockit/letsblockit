{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_16;
}
{
  pname = "ory-cli";
  version = "0.0.99";

  src = pkgs.fetchFromGitHub {
    owner = "xvello";
    repo = "ory-cli";
    rev = "e0c51d4ef27d30e245bc6046540b5aaec5ded4a2";
    sha256 = "0fk42y0517ga2mwg7zhma44r4wqvnvxmj7yxd6d3p7iypydwzlpx";
  };
  vendorSha256 = "06l6l04zdzziwjc5y1ji9k2ahyrhv36f2filq1f1sjxrik9wz6a4";

  doCheck = false;
  installPhase = ''
    mkdir -p $out/bin
    cp $GOPATH/bin/cli $out/bin/ory
  '';
}
