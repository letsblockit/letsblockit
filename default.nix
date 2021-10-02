{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_16;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "0rhk3rg3lkn6z047yy6lndrshibf0mrcg3mqy8vhkhv7b98bz200";
  src = ./.;
}
