{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_17;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "sha256-ZMGTL58tCbyG1QvmTFzTd+mkJ1t2/403e5podNr3aOY=";
  src = ./.;
  doCheck = false;
}
