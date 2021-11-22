{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_17;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "0z1w7hzy2jd1j4km9anbpzz2946wjq5959x7b464yviblhbyikwk";
  src = ./.;
  doCheck = false;
}
