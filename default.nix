{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_17;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "1zvm97x0rm65d1gv4hrksdb3gk18rmba5d5z1ciksxakxw0dj2gz";
  src = ./.;
  doCheck = false;
}
