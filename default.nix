{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_16;
}
{
  pname = "letsblockit";
  version = "1.0";
  buildInputs = [ pkgs.sqlite ];
  buildFlags = "-tags libsqlite3";
  vendorSha256 = "00qmnvsjn5r5sdja3lm9w3xhyvxnsidfyp5nji8gahlif35g6xx8";
  src = ./.;
  doCheck = false;
}
