{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_17;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "0z73jjl0i5w28xx7vp3x9wjzyjm2wdnyk39xcd2970prlf3y8p2h";
  src = ./.;
  doCheck = false;
}
