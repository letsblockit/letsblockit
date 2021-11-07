{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_17;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "15virha8bah1jr24srhwfhvxmnjxwaqsmxfsz02xrxnq3nnmnr5q";
  src = ./.;
  doCheck = false;
}
