{ pkgs ? import <nixpkgs> { } }:
pkgs.buildGoModule.override
{
  go = pkgs.go_1_17;
}
{
  pname = "letsblockit";
  version = "1.0";
  vendorSha256 = "0c2rcck6fvn77p4c1pa9s9ygjwq6lvxnynkmvnz6mwmwi0mshi10";
  src = ./.;
  doCheck = false;
}
