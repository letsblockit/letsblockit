{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_17;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "1947h4x7clig6z00fb0qnxxzhclxqjgikdq6i8gmq5bzmk33fiaf";
  src = ./.;
  doCheck = false;
}
