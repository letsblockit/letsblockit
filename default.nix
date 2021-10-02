{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_16;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "1qd56hv517sxbsd4yg3r590bn7lzm7hiviy3xazx15y2sd78nyn7";
  src = ./.;
}
