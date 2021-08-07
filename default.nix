{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_16;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "1fbw04dcih7716r833z7v3g9lg8cibixjfikc9r9qh4rhbzfqqsk";
  src = ./.;
}
