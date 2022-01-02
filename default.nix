{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_17;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "sha256-ZR/vfl7uDEnapYI3VBcWHMhZC+8vhh4OYzX+KTM80gI=";
  src = ./.;
  doCheck = false;
}
