{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_17;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "1gmxhzn7qmxjmjbk9mskqbgsnyfwrbbllgcs86n6mkad321i8i4d";
  src = ./.;
  doCheck = false;
}
